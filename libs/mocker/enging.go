package mocker

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/shopspring/decimal"
	"github.com/xjieinfo/xjgo/xjcore/xjexcel"
	"log"
	"sync"
	"time"
	"trade-robot-bd/app/grid-strategy-svc/util/utils"
	"trade-robot-bd/libs/goex"
	"trade-robot-bd/libs/goex/binance"
	"trade-robot-bd/libs/helper"
)

func (m *MockCyCle) RunCycleWhere() (r []*WhereCycleOne) {
	var mc = WhereCycleOne{
		Comm: Comm{Symbol: m.Symbol},
		Bn:   m.Bn,
		MockDetail: &goex.MockDetail{
			OldUsd:  300,
			Usd:     300,
			Balance: 0,
			Type:    goex.NewOrder_UM,
			CpUsd:   10,
			FeeRate: 0.04,
		},
	}
	if len(m.StopLossRate) < 1 {
		m.StopLossRate = []float64{0}
	}
	if len(m.AtrLength) < 1 {
		m.AtrLength = []int{0}
	}
	if len(m.ProfitType) < 1 {
		m.ProfitType = []CloseSignal{ProfitRate, ProfitSignal}
	}
	if len(m.Lever) < 1 {
		m.Lever = []int{5} //默认5倍合约
	}

	for rateI, _ := range m.ProfitRate {
		for StopLossI, _ := range m.StopLossRate {
			for AtrLengthI, _ := range m.AtrLength {
				for ProfitTypeI, _ := range m.ProfitType {
					for LeverI, _ := range m.Lever {
						icopy := mc
						icopy.ProfitRate = m.ProfitRate[rateI]
						icopy.StopLossRate = m.StopLossRate[StopLossI]
						icopy.AtrLength = m.AtrLength[AtrLengthI]
						icopy.ProfitType = m.ProfitType[ProfitTypeI]
						icopy.MockDetail.Lever = int64(m.Lever[LeverI])
						r = append(r, &icopy)
					}
				}
			}
		}
	}
	return r
}

func (m *MockCyCle) MakeCycleWhere(bn *binance.Binance, actionType int) {
	var list goex.MockResults
	log.Printf("数据回测 - 数据整理中....")
	if actionType == WhereAll {
		list = m.MakeCycleWhereAll(bn)
	} else {
		list = m.MakeCycleWhereSplit(bn)
	}
	var f *excelize.File
	log.Printf("[%vm]回测完毕,数据: %v 组\n", m.Cycle, len(list))
	pl := list.SplitToFlotList("desc", true)
	for i := 0; i < len(pl); i++ {
		itemMock := pl[i]
		//if itemMock.IsLiquidation {
		//	continue
		//}
		itemMock.ExportOrders(m.Cycle, m.StartDay, m.EndDay)
		log.Printf("%v分钟[%s-%s]回测数据: 初始余额:%v 当前余额:%v 交易次数:%v Atr:%v 平仓条件:%v 收益率:%v 止损率:%v %%  爆仓:%v \n",
			m.Cycle, m.StartDay, m.EndDay, itemMock.OldUsd, decimal.NewFromFloat(itemMock.Usd).StringFixed(2), itemMock.TradeNum,
			itemMock.AtrLength,
			itemMock.ProfitType,
			itemMock.ProfitRate,
			itemMock.StopLossRate,
			helper.IfThen(itemMock.IsLiquidation, "是", "否"))
		break
	}

	f = xjexcel.ListToExcel(list, "", fmt.Sprintf("%vm回测数据", m.Cycle))
	helper.Exists("mocktest", true, true)
	err := f.SaveAs(fmt.Sprintf("mocktest/%vm回测数据.xlsx", m.Cycle))
	if err != nil {
		log.Fatal("保存失败")
		return
	}
}
func (m *MockCyCle) MakeCycleWhereAll(bn *binance.Binance) (list goex.MockResults) {
	m.Brackets = LoadBrackets(bn) //加载杠杆信息
	whereOne := m.RunCycleWhere()
	log.Printf("条件重组完成:共%v组数据", len(whereOne))

	kline := binance.GetKLines("FUTURES_UM", "klines", "daily", m.StartDay, m.EndDay,
		[]string{fmt.Sprintf("%vm", m.Cycle)}, []string{m.Symbol.ToSymbol("")}, m.Symbol)
	log.Printf("K线数据读取完毕,共%v条 等待排序...", len(kline))
	m.kLineData = goex.KlineSort(kline, "asc")
	log.Println("K线数据排序完成")
	////-----tradeReArr - Start
	//trades := binance.GetTrade("FUTURES_UM", "trades", "daily", m.StartDay, m.EndDay, []string{}, []string{m.Symbol.ToSymbol("")}, m.Symbol)
	//log.Printf("交易数据读取完毕,共%v条 等待排序...", len(trades))
	//m.TradeData = goex.BubbleSortGeneric1(trades, runtime.NumCPU())
	//log.Printf("数据加载完毕，等待整合,[K线:%v-交易记录:%v]\n", len(m.kLineData), len(m.TradeData))
	//m.KlineLinkTrade() //关联K线与交易记录
	//log.Printf("[%v-%v]日数据关联完成[%v]", m.StartDay, m.EndDay, len(m.KlinsLikeTrade))
	////-----tradeReArr - End
	var wg = new(sync.WaitGroup)

	for iR := 0; iR < len(whereOne); iR++ {
		wg.Add(1)
		iR := iR
		go func() {
			var r = whereOne[iR].RsiOrder(m)
			list = append(list, *r)
			wg.Done()
		}()
		//go pool.Process(iR)
	}
	log.Println("任务分发完成，正在执行...")
	wg.Wait()
	log.Println("任务执行完成，正在将结果写入文件...")
	return list
}

// MakeCycleWhereSplit 分割回测 防止数据过大 造成溢出
func (m *MockCyCle) MakeCycleWhereSplit(bn *binance.Binance) (list goex.MockResults) {
	// 将K线数据 按天分割后 逐一回测
	var klinData = binance.DownloadData("FUTURES_UM", "klines", "daily", m.StartDay, m.EndDay, []string{fmt.Sprintf("%vm", m.Cycle)}, []string{m.Symbol.ToSymbol("")})
	var tradeData = binance.DownloadData("FUTURES_UM", "trades", "daily", m.StartDay, m.EndDay, []string{}, []string{m.Symbol.ToSymbol("")})
	log.Printf("数据加载完毕，等待整合,[%v-%v]\n", len(klinData), len(tradeData))
	m.Brackets = LoadBrackets(bn) //加载杠杆信息

	whereOne := m.RunCycleWhere()
	log.Printf("条件重组完成:共%v组数据", len(whereOne))
	var wg = new(sync.WaitGroup)
	for i := 0; i < len(klinData); i++ {
		m.kLineData = []*goex.Kline{}
		m.TradeData = []*goex.Trade{}
		// 按天分割后的K线数据与交易记录
		m.kLineData = goex.LocalKlineCsv(klinData[i], m.Symbol)
		m.TradeData = goex.LocalTradeCsv(tradeData[i], m.Symbol)
		m.KlineLinkTrade() //关联K线与交易记录
		log.Printf("[%v]日数据关联完成[%v]", time.UnixMilli(m.kLineData[0].CloseTime).Format(helper.TimeFormatYmdHis), len(m.KlinsLikeTrade))
		for iR := 0; iR < len(whereOne); iR++ {
			_iR := iR
			go func() {
				wg.Add(1)
				_one := whereOne[_iR]
				var r = _one.RsiOrder(m)
				wg.Done()
				if i >= len(klinData)-1 {
					list = append(list, *r)
				}
			}()
		}
		wg.Wait()
		log.Printf("处理完毕")
	}
	return list
}

// RsiOrder Rsi量化大法
func (m *WhereCycleOne) RsiOrder(mockData *MockCyCle) (r *goex.MockResult) {
	_mockResult := &goex.MockResult{
		AtrLength:    m.AtrLength,
		ProfitRate:   m.ProfitRate,
		OldUsd:       m.MockDetail.OldUsd,
		StopLossRate: m.StopLossRate,
		ProfitType:   m.ProfitType.String(),
	}
	//创建订单
	newOrder := func(direction, closeTime int64) *goex.MockOrder {
		o := &goex.MockOrder{
			Direction: direction,
			Type:      m.MockDetail.Type,
			Lever:     m.MockDetail.Lever,
			Quantity:  helper.IfThen(m.MockDetail.Type == goex.NewOrder_UM, m.MockDetail.Usd*float64(m.MockDetail.Lever), m.MockDetail.Usd),
			FeeRate:   m.MockDetail.FeeRate,
			CpUsd:     m.MockDetail.CpUsd,
			BidTime:   closeTime,
		}
		//限制最大持仓量
		if o.Quantity > mockData.MaxHold {
			o.Quantity = mockData.MaxHold
		}
		return o
	}
	openBuy := func(actionPrice float64, closeTime int64) bool {
		if m.Signal == goex.BUY && m.MockDetail.BuyOrder == nil { //达到买多条件 - 判断有无货币购买(包含资金切割) 注意:非本地回测下，数据具有实时性
			var order = newOrder(goex.NewOrder_Buy, closeTime)                                //生成做多订单
			order.Buy(actionPrice)                                                            //非本地数据回测下 数据具有实时性 本地模拟 暂时取最低价与收盘价折中值
			order.CalcLiquidation(m.MockDetail.Usd, mockData.Brackets[m.Symbol.ToSymbol("")]) //计算爆仓价 传参格式 余额 保证金维持比例
			m.MockDetail.BuyOrder = order
			m.MockDetail.TradeNum++
			//因为本地数据回测  是根据K线收盘价操作  所以下买单后  需要在下一个K线判断是否卖出/平多
			return true
		}
		return false
	}
	// 打开空单
	openSell := func(actionPrice float64, closeTime int64) bool {
		//检测当前是否有空单进行做空
		if m.MockDetail.SellOrder == nil && m.Signal == goex.SELL {
			var order = newOrder(goex.NewOrder_Sell, closeTime)                               //生成做空订单
			order.Buy(actionPrice)                                                            //非本地数据回测下 数据具有实时性 本地模拟 暂时取最高与开盘折中值
			order.CalcLiquidation(m.MockDetail.Usd, mockData.Brackets[m.Symbol.ToSymbol("")]) //计算爆仓价 传参格式 余额 保证金维持比例
			m.MockDetail.SellOrder = order
			m.MockDetail.TradeNum++
			return true
		}
		return false
	}
	// 关闭空单
	closeSell := func(actionPrice float64, closeTime int64) {
		order := m.MockDetail.SellOrder
		if order == nil || m.Signal != goex.BUY {
			return
		}
		spot := &goex.MockOrder{
			Direction: order.Direction,
			Type:      order.Type,
			Lever:     order.Lever,
			Quantity:  order.Quantity,
			FeeRate:   order.FeeRate,
			CpUsd:     order.CpUsd,
		}
		spot.Buy(order.Open)
		spot.Sell(actionPrice)
		spot.CalcCpUp()
		action := func() {
			m.MockDetail.SellOrder.Sell(actionPrice) //非本地数据回测下 数据具有实时性 本地模拟 暂时取最低价与收盘价折中值
			m.MockDetail.SellOrder.AskTime = closeTime
			m.MockDetail.SellOrder.CalcCpUp()
			m.MockDetail.Usd += m.MockDetail.SellOrder.NetGainUSd
			m.MockDetail.SellOrder.Usd = m.MockDetail.Usd //快照 模拟余额
			m.MockDetail.HistoryOrder = append(m.MockDetail.HistoryOrder, m.MockDetail.SellOrder)
			m.MockDetail.SellOrder = nil
		}
		//触发止损
		if spot.NetRate <= m.StopLossRate || (spot.NetRate >= m.ProfitRate && m.ProfitType == ProfitRate) {
			action()
			return
		}
		if m.Signal == goex.BUY && m.ProfitType == ProfitSignal {
			if spot.NetRate > m.StopLossRate && spot.NetRate < 0 { //出现亏损 但未达到止损
				return
			}
			action()
		}
	}
	closeBuy := func(actionPrice float64, closeTime int64) *goex.MockResult {
		if m.MockDetail.BuyOrder == nil || _mockResult.IsLiquidation {
			return nil
		}
		action := func() {
			m.MockDetail.BuyOrder.Sell(actionPrice)
			m.MockDetail.BuyOrder.AskTime = closeTime
			m.MockDetail.BuyOrder.CalcCpUp()
			m.MockDetail.Usd += m.MockDetail.BuyOrder.NetGainUSd
			m.MockDetail.BuyOrder.Usd = m.MockDetail.Usd //快照 模拟余额
			m.MockDetail.HistoryOrder = append(m.MockDetail.HistoryOrder, m.MockDetail.BuyOrder)
			m.MockDetail.BuyOrder = nil
		}
		order := m.MockDetail.BuyOrder
		if helper.IfThen(order.Type == goex.NewOrder_Buy, actionPrice <= order.Liquidation, actionPrice >= order.Liquidation) { //爆仓啦
			_mockResult.IsLiquidation = true
			m.IsLiquidation = true
			_mockResult.Usd = m.MockDetail.Usd
			_mockResult.TradeNum = m.MockDetail.TradeNum
			m.MockDetail.BuyOrder.Sell(actionPrice)
			m.MockDetail.BuyOrder.AskTime = closeTime
			m.MockDetail.BuyOrder.CalcCpUp()
			m.MockDetail.HistoryOrder = append(m.MockDetail.HistoryOrder, m.MockDetail.BuyOrder)
			_mockResult.Order = m.MockDetail.HistoryOrder
			return _mockResult
		}
		spot := &goex.MockOrder{
			Direction: order.Direction,
			Type:      order.Type,
			Lever:     order.Lever,
			Quantity:  order.Quantity,
			FeeRate:   order.FeeRate,
			CpUsd:     order.CpUsd,
		}
		spot.Buy(order.Open)
		spot.Sell(actionPrice)
		spot.CalcCpUp()
		//触发止损
		if spot.NetRate <= m.StopLossRate || (spot.NetRate >= m.ProfitRate && m.ProfitType == ProfitRate) {
			action()
			return nil
		}
		if m.Signal == goex.SELL && m.ProfitType == ProfitSignal {
			if spot.NetRate > m.StopLossRate && spot.NetRate < 0 { //出现亏损 但未达到止损
				return nil
			}
			action()
			return nil
		}
		return nil
	}
	// 测试数据爆仓后 不再进行操作
	if m.IsLiquidation {
		return
	}
	for i := 0; i < len(mockData.kLineData); i++ {
		if i >= len(mockData.kLineData)-2 {
			_mockResult.Usd = m.MockDetail.Usd
			_mockResult.TradeNum = m.MockDetail.TradeNum
			_mockResult.Order = m.MockDetail.HistoryOrder
		}
		var item = mockData.kLineData[i]
		m.CalcKline(item, 0)
		m.GetSignal()
		if len(m.Renko) < 1 {
			continue
		}
		_, _ = openSell, closeSell
		_, _ = openBuy, closeBuy
		if ok := openSell(m.Renko[len(m.Renko)-1].High, item.CloseTime); ok {
			continue
		}
		closeSell(m.Renko[len(m.Renko)-1].Low, item.CloseTime) //查看是否有做空订单进行平空
		if ok := openBuy(m.Renko[len(m.Renko)-1].Low, item.CloseTime); ok {
			continue
		}
		//如果结果不为nil  则返回爆仓数据 该条测试组意味着暂停
		if r := closeBuy(m.Renko[len(m.Renko)-1].High, item.CloseTime); r != nil {
			return r
		}
	}

	return _mockResult
}

// LoadBrackets 加载币种杠杆与保证金信息
func LoadBrackets(bn *binance.Binance) Brackets {
	brk := make(Brackets)
	lp := bn.Brackets()
	for i := range lp {
		item := lp[i]
		brk[item.Symbol] = &item
	}
	log.Printf("杠杆与保证金列表加载完毕:%v", len(brk))
	return brk
}

// CalcTrade 通过传入交易记录 进行后置处理
func CalcTrade(trade *goex.Trade) {

}

// CalcKline 通过K线 进行数据后置处理
func (m *WhereCycleOne) CalcKline(kline *goex.Kline, i int) {
	m.KlineLast = kline
	if len(m.kLineData) <= 0 {
		m.kLineData = append(m.kLineData, kline)
	} else if kline.Timestamp == m.kLineData[len(m.kLineData)-1].Timestamp {
		m.kLineData[len(m.kLineData)-1] = kline
	} else {
		m.kLineData = append(m.kLineData, kline)
	}
	if len(m.kLineData) > 167 {
		m.kLineData = m.kLineData[1:]
	}
	m.GetSignal()
}
func (m *WhereCycleOne) GetSignal() goex.TradeSide {
	m.Signal = m.SignalRenko()
	switch m.OlineType {
	case 1:
		m.oline_N()
		break
	}
	return m.Signal
}
func (m *WhereCycleOne) oline_N() {
	//创建订单
	newOrder := func(direction, closeTime int64) *goex.MockOrder {
		o := &goex.MockOrder{
			Direction: direction,
			Type:      m.MockDetail.Type,
			Lever:     m.MockDetail.Lever,
			Quantity:  helper.IfThen(m.MockDetail.Type == goex.NewOrder_UM, m.MockDetail.Usd*float64(m.MockDetail.Lever), m.MockDetail.Usd),
			FeeRate:   m.MockDetail.FeeRate,
			CpUsd:     m.MockDetail.CpUsd,
			BidTime:   closeTime,
		}
		//限制最大持仓量
		if o.Quantity > m.MaxHold {
			o.Quantity = m.MaxHold
		}
		return o
	}
	openBuy := func(actionPrice float64, closeTime int64) bool {
		if m.Signal == goex.BUY && m.MockDetail.BuyOrder == nil { //达到买多条件 - 判断有无货币购买(包含资金切割) 注意:非本地回测下，数据具有实时性
			m.MockDetail.Mutex.Lock()
			var order = newOrder(goex.NewOrder_Buy, closeTime)                         //生成做多订单
			order.Buy(actionPrice)                                                     //非本地数据回测下 数据具有实时性 本地模拟 暂时取最低价与收盘价折中值
			order.CalcLiquidation(m.MockDetail.Usd, m.Brackets[m.Symbol.ToSymbol("")]) //计算爆仓价 传参格式 余额 保证金维持比例
			m.MockDetail.BuyOrder = order
			m.MockDetail.TradeNum++
			//因为本地数据回测  是根据K线收盘价操作  所以下买单后  需要在下一个K线判断是否卖出/平多
			m.MockDetail.Mutex.Unlock()
			log.Printf("[%v] 方向:做多 开仓:%v  当前余额:%v",
				m.Symbol.String(), m.MockDetail.BuyOrder.Open, m.MockDetail.Usd)
			return true
		}
		return false
	}
	closeBuy := func(actionPrice float64, closeTime int64) *goex.MockResult {
		if m.MockDetail.BuyOrder == nil || m.IsLiquidation {
			return nil
		}
		action := func() {
			m.MockDetail.BuyOrder.Sell(actionPrice)
			m.MockDetail.BuyOrder.AskTime = closeTime
			m.MockDetail.BuyOrder.CalcCpUp()
			m.MockDetail.Usd += m.MockDetail.BuyOrder.NetGainUSd
			m.MockDetail.BuyOrder.Usd = m.MockDetail.Usd //快照 模拟余额
			m.MockDetail.HistoryOrder = append(m.MockDetail.HistoryOrder, m.MockDetail.BuyOrder)
			log.Printf("[%v] 方向:做多 开仓:%v 闭仓:%v 利润:%v 当前余额:%v",
				m.Symbol.String(), m.MockDetail.BuyOrder.Open, m.MockDetail.BuyOrder.Close, m.MockDetail.BuyOrder.NetGainUSd, m.MockDetail.Usd)
			m.MockDetail.BuyOrder = nil
		}
		order := m.MockDetail.BuyOrder
		if helper.IfThen(order.Type == goex.NewOrder_Buy, actionPrice <= order.Liquidation, actionPrice >= order.Liquidation) { //爆仓啦
			m.IsLiquidation = true
			log.Fatalf("[%v]已爆仓 开仓价:%v 当前价:%v 爆仓价:%v", m.Symbol.String(), order.Open, actionPrice, order.Liquidation)
		}
		spot := &goex.MockOrder{
			Direction: order.Direction,
			Type:      order.Type,
			Lever:     order.Lever,
			Quantity:  order.Quantity,
			FeeRate:   order.FeeRate,
			CpUsd:     order.CpUsd,
		}
		spot.Buy(order.Open)
		spot.Sell(actionPrice)
		spot.CalcCpUp()
		//触发止损
		if spot.NetRate <= m.StopLossRate || (spot.NetRate >= m.ProfitRate && m.ProfitType == ProfitRate) {
			action()
			return nil
		}
		if m.Signal == goex.SELL && m.ProfitType == ProfitSignal {
			if spot.NetRate > m.StopLossRate && spot.NetRate < 0 { //出现亏损 但未达到止损
				return nil
			}
			action()
			return nil
		}
		return nil
	}
	// 打开空单
	openSell := func(actionPrice float64, closeTime int64) bool {
		//检测当前是否有空单进行做空
		if m.MockDetail.SellOrder == nil && m.Signal == goex.SELL {
			var order = newOrder(goex.NewOrder_Sell, closeTime)                        //生成做空订单
			order.Buy(actionPrice)                                                     //非本地数据回测下 数据具有实时性 本地模拟 暂时取最高与开盘折中值
			order.CalcLiquidation(m.MockDetail.Usd, m.Brackets[m.Symbol.ToSymbol("")]) //计算爆仓价 传参格式 余额 保证金维持比例
			m.MockDetail.SellOrder = order
			m.MockDetail.TradeNum++
			return true
		}
		return false
	}
	// 关闭空单
	closeSell := func(actionPrice float64, closeTime int64) {
		if m.MockDetail.SellOrder == nil || m.Signal == goex.SELL {
			return
		}
		spot := &goex.MockOrder{
			Direction: m.MockDetail.SellOrder.Direction,
			Type:      m.MockDetail.SellOrder.Type,
			Lever:     m.MockDetail.SellOrder.Lever,
			Quantity:  m.MockDetail.SellOrder.Quantity,
			FeeRate:   m.MockDetail.SellOrder.FeeRate,
			CpUsd:     m.MockDetail.SellOrder.CpUsd,
		}
		spot.Buy(m.MockDetail.SellOrder.Open)
		spot.Sell(actionPrice)
		spot.CalcCpUp()
		action := func() {
			m.MockDetail.SellOrder.Sell(actionPrice) //非本地数据回测下 数据具有实时性 本地模拟 暂时取最低价与收盘价折中值
			m.MockDetail.SellOrder.AskTime = closeTime
			m.MockDetail.SellOrder.CalcCpUp()
			m.MockDetail.Usd += m.MockDetail.SellOrder.NetGainUSd
			m.MockDetail.SellOrder.Usd = m.MockDetail.Usd //快照 模拟余额
			m.MockDetail.HistoryOrder = append(m.MockDetail.HistoryOrder, m.MockDetail.SellOrder)
			m.MockDetail.SellOrder = nil
		}
		//触发止损
		if spot.NetRate <= m.StopLossRate || (spot.NetRate >= m.ProfitRate && m.ProfitType == ProfitRate) {
			action()
			return
		}
		if m.Signal == goex.BUY && m.ProfitType == ProfitSignal {
			if spot.NetRate > m.StopLossRate && spot.NetRate < 0 { //出现亏损 但未达到止损
				return
			}
			action()
		}
	}
	openBuy(m.KlineLast.Close, m.KlineLast.Timestamp)
	closeBuy(m.KlineLast.Close, m.KlineLast.Timestamp)
	openSell(m.KlineLast.Close, m.KlineLast.Timestamp)
	closeSell(m.KlineLast.Close, m.KlineLast.Timestamp)
}

// SignalRenko Renko信号
func (m *WhereCycleOne) SignalRenko() goex.TradeSide {
	if len(m.kLineData) < m.AtrLength+1 {
		return -1
	}
	places, err := m.Bn.GetTradeSymbol(m.Symbol)
	if err != nil || places.PricePrecision <= 0 {
		places = &binance.TradeSymbol{}
		places.PricePrecision = 6
	}
	m.Renko = goex.KlineToRenKo(m.kLineData, m.AtrLength, goex.InClose, int32(places.BaseAssetPrecision))
	openRenko := utils.SplitToFlotList(goex.RealData(m.Renko, goex.InOpen), "Reverse")
	closeRenko := utils.SplitToFlotList(goex.RealData(m.Renko, goex.InClose), "Reverse")
	buySignal := goex.Crossunder(openRenko, closeRenko)
	sellSignal := goex.Crossover(openRenko, closeRenko)
	if buySignal {
		return goex.BUY
	}
	if sellSignal {
		return goex.SELL
	}
	return -1 // -1为未识别信号
}

// SubTrade 订阅交易记录
func (m *WhereCycleOne) SubTrade(bn *binance.BinanceWs) {
	err := bn.SubscribeAggTrade(m.Symbol, func(trade *goex.Trade) {
		CalcTrade(trade)
	})
	if err != nil {
		fmt.Println("运行失败", err)
	}
}

// SubKline 读取并订阅K线
func (m *WhereCycleOne) SubKline(bn *binance.BinanceWs, bnt *binance.BinanceSwap) {
	var err error
	m.kLineData, err = bnt.GetKlineRecords(m.Symbol, goex.KLINE_PERIOD_5MIN, 167, 0)
	if err != nil {
		log.Fatalf("运行失败,获取K线报错-%v", err.Error())
	}
	bn.KlineCallback = m.CalcKline
	err = bn.SubscribeKline(m.Symbol, goex.KLINE_PERIOD_5MIN)
	if err != nil {
		fmt.Println("运行失败", err)
	}
}

// OnLineKline 读取在线K线数据
func (m *WhereCycleOne) OnLineKline(bn *binance.BinanceWs, bnt *binance.BinanceSwap) {
	var err error
	m.kLineData, err = bnt.GetKlineRecords(m.Symbol, goex.KLINE_PERIOD_5MIN, 167, 0)
	if err == nil && len(m.kLineData) > 0 {
		m.KlineLast = m.kLineData[len(m.kLineData)-1]
		m.GetSignal()
	}
}
