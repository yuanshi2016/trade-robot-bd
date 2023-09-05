package mockers

import (
	"github.com/shopspring/decimal"
	"log"
	"math"
	"time"
	"trade-robot-bd/app/grid-strategy-svc/util/utils"
	"trade-robot-bd/libs/goex"
	"trade-robot-bd/libs/goex/binance"
	"trade-robot-bd/libs/helper"
)

func (m *WhereCycleOne) Res(whereS *MockCyCle) (r *goex.MockResult) {
	for i := 0; i < len(whereS.kLineData)-1; i++ {
		m.CalcKline(whereS.kLineData[i])
		r = m.Run()
	}
	whereS.MockResults = append(whereS.MockResults, r)
	return
}
func (m *WhereCycleOne) Run() *goex.MockResult {
	// -- 是否有实盘订单
	if m.TradeType == TradeTypeOline && m.MockDetail.SellOrderOnline != nil {
		m.actionSell(m.MockDetail.SellOrderOnline.Price, time.Now().UnixNano())
	}
	if m.TradeType == TradeTypeOline && m.MockDetail.BuyOrderOnline != nil {
		m.actionBuy(m.MockDetail.BuyOrderOnline.Price, time.Now().UnixNano())
	}
	if m.MockResult == nil {
		m.MockResult = &goex.MockResult{
			AtrLength:     m.AtrLength,
			RenKoMoveI:    m.RenKoMoveI,
			RenKoMoveType: m.RenKoMoveType,
			ProfitRate:    m.ProfitRate,
			OldUsd:        m.MockDetail.OldUsd,
			StopLossRate:  m.StopLossRate,
			ProfitType:    m.ProfitType.String(),
		}
	}
	m.GetSignal()
	if len(m.Renko) < 1 {
		return m.MockResult
	}
	m.closeOrder(goex.BUY, m.KlineLast.Close, m.KlineLast.Timestamp)
	m.closeOrder(goex.SELL, m.KlineLast.Close, m.KlineLast.CloseTime)
	if m.Signal == goex.SELL && m.MockDetail.SellOrder == nil {
		m.actionSell(m.KlineLast.Close, m.KlineLast.CloseTime)
	}
	if m.Signal == goex.BUY && m.MockDetail.BuyOrder == nil {
		m.actionBuy(m.KlineLast.Close, m.KlineLast.CloseTime)
	}
	m.MockResult.Usd = m.MockDetail.Usd
	m.MockResult.TradeNum = m.MockDetail.TradeNum
	m.MockResult.Order = m.MockDetail.HistoryOrder
	return m.MockResult
}

// calcQuantity 计算开仓数量
func (m *WhereCycleOne) calcQuantity() float64 {
	var balance float64
	var onlineBalance goex.SubAccount
	if m.TradeType == TradeTypeOline {
		switch m.MockDetail.Type {
		case goex.NewOrder_UM:
			Account, err := m.BnSwap.GetAccount()
			if err != nil {
				return 0
			}
			onlineBalance = Account.SubAccounts[m.Symbol.CurrencyB]
			break
		case goex.NewOrder_CM:
			Account, err := m.BnSwap.GetAccount()
			if err != nil {
				return 0
			}
			onlineBalance = Account.SubAccounts[m.Symbol.CurrencyB]
			break
		case goex.NewOrder_SPOT:
			Account, err := m.Bn.GetAccount()
			if err != nil {
				return 0
			}
			onlineBalance = Account.SubAccounts[m.Symbol.CurrencyB]
			break
		}
		balance = onlineBalance.Balance
		m.MockDetail.OnlineUsd = onlineBalance.MarginBalance
	} else {
		// 非实盘模式下 采用模拟余额
		balance = m.MockDetail.Usd
	}

	quantity := helper.IfThen(m.MockDetail.Type == goex.NewOrder_UM, balance*m.TotalBalanceRatio*float64(m.MockDetail.Lever), m.MockDetail.Usd)
	// 如果是双向持仓 且存在多/空单，则 /2
	if m.IsTowWay && m.MockDetail.SellOrderOnline == nil && m.MockDetail.BuyOrderOnline == nil {
		quantity = math.Floor(quantity / 2)
	}
	//限制最大持仓量
	if quantity > m.MaxHold {
		quantity = m.MaxHold
	}
	return quantity //如果为单向持仓 则返回全部
}

// newOrder 创建订单
func (m *WhereCycleOne) newOrder(direction goex.TradeSide, price float64, actionTime int64) *goex.MockOrder {
	quantity := m.calcQuantity()
	if quantity == 0 || m.KlineLast == nil {
		return nil
	}
	//最新价格为0  返回空
	if m.KlineLast.Close <= 0 && m.TradeType == TradeTypeOline {
		return nil
	}
	if direction == goex.BUY && (m.MockDetail.BuyOrder != nil || m.MockDetail.BuyOrderOnline != nil) {
		return nil
	}
	if direction == goex.SELL && (m.MockDetail.SellOrder != nil || m.MockDetail.SellOrderOnline != nil) {
		return nil
	}
	o := &goex.MockOrder{
		Direction: direction,
		Type:      m.MockDetail.Type,
		Lever:     m.MockDetail.Lever,
		Quantity:  quantity,
		FeeRate:   m.MockDetail.FeeRate,
		CpUsd:     m.MockDetail.CpUsd,
		BidTime:   actionTime,
	}
	o.OpenMarket(price)
	if m.TradeType == TradeTypeOline { //实盘模式
		exinfo, err := m.BnSwap.GetTradeSymbol(m.Symbol)
		var onlineOrder *goex.Order
		if err == nil && (m.MockDetail.BuyOrderOnline == nil || m.MockDetail.SellOrderOnline == nil) {
			amount := decimal.NewFromFloat(o.Quantity)
			_price := decimal.NewFromFloat(helper.IfThen(m.TradeType == TradeTypeOline, m.KlineLast.Close, price)).Round(int32(exinfo.PricePrecision))
			amount = amount.Div(_price).Round(int32(exinfo.QuantityPrecision))

			switch direction {
			case goex.BUY, goex.BUY_MARKET:
				if m.MockDetail.BuyOrderOnline == nil {
					onlineOrder, err = m.BnSwap.MarketBuy(amount, _price, m.Symbol, false)
					m.MockDetail.BuyOrderOnline = onlineOrder
				}
				break
			case goex.SELL, goex.SELL_MARKET:
				if m.MockDetail.SellOrderOnline == nil {
					onlineOrder, err = m.BnSwap.MarketSell(amount, _price, m.Symbol, false)
					m.MockDetail.SellOrderOnline = onlineOrder
				}
				break
			}
			if err != nil {
				log.Printf("[%v]下单失败[%v],数量:%v 价格:%v | %v", m.Symbol.String(), direction.String(), amount, _price, err.Error())
			} else if onlineOrder != nil {
				log.Printf("[%v]下单成功[%v],数量:%v 价格:%v 订单ID:%v", m.Symbol.String(), direction.String(), amount, _price, onlineOrder.OrderID2)
			}
		}
	}
	// 计算爆仓价
	o.CalcLiquidation(helper.IfThen(m.TradeType == TradeTypeOline, m.MockDetail.OnlineUsd, m.MockDetail.Usd), m.Brackets[m.Symbol.ToSymbol("")]) //计算爆仓价 传参格式 余额 保证金维持比例
	return o
}

// closeOrder 根据传入订单方向 获取对应订单 卖出/平仓
func (m *WhereCycleOne) closeOrder(direction goex.TradeSide, actionPrice float64, actionTime int64) bool {
	var order = m.getWhereOrder(direction)
	var onlineOrder *goex.Order
	var onlineOrderErr error
	// 如果当前方向订单为空 或者当前模拟已爆仓 则直接返回false
	if order == nil || m.IsLiquidation {
		return false
	}
	spot := m.getSpotOrder(direction, actionPrice)

	action := func() bool {
		if m.TradeType == TradeTypeOline {
			log.Printf("[%v]平仓指令[%v] ", m.Symbol.String(), direction.String())
			if m.MockDetail.BuyOrderOnline != nil && direction == goex.BUY {
				onlineOrder, onlineOrderErr = m.BnSwap.MarketSell(decimal.NewFromFloat(m.MockDetail.BuyOrderOnline.Amount).Abs(), decimal.NewFromFloat(m.MockDetail.BuyOrderOnline.Price), m.Symbol, true)
			}
			if m.MockDetail.SellOrderOnline != nil && direction == goex.SELL {
				onlineOrder, onlineOrderErr = m.BnSwap.MarketBuy(decimal.NewFromFloat(m.MockDetail.SellOrderOnline.Amount).Abs(), decimal.NewFromFloat(m.MockDetail.SellOrderOnline.Price), m.Symbol, true)
			}
			// 如果实盘订单没有平仓成功则 模拟订单不会平仓
			if onlineOrderErr != nil {
				oOrder := helper.IfThen(direction == goex.BUY, m.MockDetail.BuyOrderOnline, m.MockDetail.SellOrderOnline)
				log.Printf("[%v]平仓失败[%v],数量:%v 价格:%v | %v", m.Symbol.String(), direction.String(), oOrder.Amount, oOrder.Price, onlineOrderErr.Error())
			} else if onlineOrder != nil {
				switch direction {
				case goex.BUY, goex.BUY_MARKET:
					m.MockDetail.BuyOrderOnline = nil
					break
				case goex.SELL, goex.SELL_MARKET:
					m.MockDetail.SellOrderOnline = nil
					break
				}
				log.Printf("[%v]平仓成功[%v],数量:%v 价格:%v  订单ID:%v", m.Symbol.String(), direction.String(), onlineOrder.Amount, onlineOrder.AvgPrice, onlineOrder.OrderID2)
			}
		}
		order.CloseMarket(actionPrice)
		order.AskTime = actionTime
		order.CalcCpUp()
		log.Printf("关闭订单 方向:%v 开仓价:%v 当前价:%v 爆仓价:%v 模拟收益:%v 持仓数量:%v 收益率:%v",
			order.Direction, order.Open, actionPrice, order.Liquidation, order.NetGainUSd, order.Quantity, spot.NetRate)
		if direction == goex.BUY {
			m.MockDetail.BuyOrder = nil
		}
		if direction == goex.SELL {
			m.MockDetail.SellOrder = nil
		}
		m.MockDetail.Usd = m.MockDetail.Usd + order.NetGainUSd
		order.Usd = m.MockDetail.Usd //快照 模拟余额
		m.MockDetail.HistoryOrder = append(m.MockDetail.HistoryOrder, order)
		order = nil
		return true
	}

	//触发止损
	if spot.NetRate <= m.StopLossRate || (spot.NetRate >= m.ProfitRate && m.ProfitType == ProfitRate) {
		return action()
	}
	// ifThen 的意义在于 获取当前操作的反向  也就是:当前是买入 要获取到卖出的信号 才可以卖出
	if m.Signal == helper.IfThen(direction == goex.BUY, goex.SELL, goex.BUY) && m.ProfitType == ProfitSignal {
		if spot.NetRate >= m.StopLossRate && spot.NetRate < 0 { //出现亏损 但未达到止损
			return false
		}
		//log.Printf("关闭订单 方向:%v 开仓价:%v 当前价:%v 爆仓价:%v 模拟收益:%v", order.Direction, order.Open, actionPrice, order.Liquidation, spot.NetGainUSd)
		return action()
	}

	return false
}

// getWhereOrder 根据方向 返回当前方向订单 如果订单为空 返回nil
func (m *WhereCycleOne) getWhereOrder(direction goex.TradeSide) *goex.MockOrder {
	var order *goex.MockOrder
	switch direction {
	case goex.BUY:
		order = m.MockDetail.BuyOrder
		break
	case goex.SELL:
		order = m.MockDetail.SellOrder
		break
	}
	return helper.IfThen(order == nil || order.Open == 0, nil, order)
}

// isLiquidation 计算是否爆仓
func (m *WhereCycleOne) isLiquidation(direction goex.TradeSide, actionPrice float64, closeTime int64) bool {
	var order = m.getWhereOrder(direction)
	if order == nil {
		return false
	}
	if order.CheckLiquidation(actionPrice) { //爆仓啦
		m.IsLiquidation = true
		order.CloseMarket(actionPrice)
		order.AskTime = closeTime
		order.CalcCpUp()
		m.MockDetail.HistoryOrder = append(m.MockDetail.HistoryOrder, order)
		order = nil
		return true
	}
	return false
}

// getSpotOrder 根据传入价格 计算盈亏情况
func (m *WhereCycleOne) getSpotOrder(direction goex.TradeSide, actionPrice float64) *goex.MockOrder {
	var order = m.getWhereOrder(direction)
	spot := &goex.MockOrder{
		Direction: order.Direction,
		Type:      order.Type,
		Lever:     order.Lever,
		Quantity:  order.Quantity,
		FeeRate:   order.FeeRate,
		CpUsd:     order.CpUsd,
	}
	spot.OpenMarket(order.Open)
	spot.CloseMarket(actionPrice)
	spot.CalcCpUp()
	return spot
}

// openBuy 买入 如果爆仓 则返回nil
func (m *WhereCycleOne) actionBuy(actionPrice float64, actionTime int64) bool {
	// 已爆仓
	if m.isLiquidation(goex.BUY, actionPrice, actionTime) {
		return false
	}
	var order = m.newOrder(goex.BUY, actionPrice, actionTime) //生成做多订单
	if order != nil {
		m.MockDetail.BuyOrder = order
		m.MockDetail.TradeNum++
	}
	return true

	return false
}

// actionSell 卖出 如果爆仓 则返回nil
func (m *WhereCycleOne) actionSell(actionPrice float64, actionTime int64) bool {
	// 已爆仓
	if m.isLiquidation(goex.SELL, actionPrice, actionTime) {
		return false
	}
	if m.MockDetail.SellOrder == nil { //达到买多条件 - 判断有无货币购买(包含资金切割) 注意:非本地回测下，数据具有实时性
		var order = m.newOrder(goex.SELL, actionPrice, actionTime) //生成做多订单
		if order != nil {
			m.MockDetail.SellOrder = order
			m.MockDetail.TradeNum++
		}
		//因为本地数据回测  是根据K线收盘价操作  所以下买单后  需要在下一个K线判断是否卖出/平多
		return true
	}
	return false
}

// CalcKline 通过K线 进行数据后置处理
func (m *WhereCycleOne) CalcKline(kline *goex.Kline) {
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
}
func (m *WhereCycleOne) GetSignal() goex.TradeSide {
	m.Signal = m.SignalRenko()
	switch m.TradeType {
	case TradeTypeLocal:
		//m.oline_N()
		break
	}
	return m.Signal
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
	m.Renko = goex.KlineToRenKo(m.kLineData, m.AtrLength, goex.InClose, int32(places.BaseAssetPrecision), m.RenKoMoveI, m.RenKoMoveType)
	// openRenko 和 closeRenko 反着来了
	openRenko := utils.SplitToFlotList(goex.RealData(m.Renko, goex.InClose), "Reverse")
	closeRenko := utils.SplitToFlotList(goex.RealData(m.Renko, goex.InOpen), "Reverse")
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

// OnLineKline 读取在线K线数据
func (m *WhereCycleOne) OnLineKline(bnt *binance.BinanceSwap) {
	var err error
	m.kLineData, err = bnt.GetKlineRecords(m.Symbol, goex.KLINE_PERIOD_5MIN, 167, 0)
	// 创建最新K线订阅 如果K线订阅基址为空 或者 启动类型不是本地回测
	if m.BnWs.KlineCallback == nil && m.TradeType != TradeTypeLocal {
		go func() {
			m.BnWs.KlineCallback = func(kline *goex.Kline, period int) {
				m.KlineLast = kline
			}
			err := m.BnWs.SubscribeKline(m.Symbol, 5)
			if err != nil {
				log.Fatalf("订阅K线失败:%v", err.Error())
			}
		}()
	}
	if err == nil && len(m.kLineData) > 0 {
		if m.KlineLast == nil {
			m.KlineLast = m.kLineData[len(m.kLineData)-1]
		}
		//m.Run()
	}
}

// LoadBrackets 加载币种杠杆与保证金信息
func LoadBrackets(bn *binance.Binance) map[string]*goex.BracketsList {
	brk := make(map[string]*goex.BracketsList)
	lp := bn.Brackets()
	for i := range lp {
		item := lp[i]
		brk[item.Symbol] = &item
	}
	return brk
}
