/**
 * @Notes:
 * @class MockCyCle
 * @package
 * @author: 原始
 * @Time: 2023/8/2   00:55
 */
package mockers

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/shopspring/decimal"
	"github.com/xjieinfo/xjgo/xjcore/xjexcel"
	"log"
	"sync"
	"time"
	"trade-robot-bd/libs/goex"
	"trade-robot-bd/libs/goex/binance"
	"trade-robot-bd/libs/helper"
)

func (m *MockCyCle) MakeCycleWhere() (r []*WhereCycleOne) {
	var mc = WhereCycleOne{
		Comm: Comm{Symbol: m.Symbol, MaxHold: m.MaxHold},
		MockDetail: goex.MockDetail{
			OldUsd:  m.Usd,
			Usd:     m.Usd,
			Balance: m.Balance,
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
	if len(m.RenKoMoveType) < 1 {
		m.RenKoMoveType = []goex.RenKoMoveType{goex.RenKoMoveTypeAMi} //默认avg+min
	}
	if len(m.RenKoMoveI) < 1 {
		m.RenKoMoveI = []int{5} //默认avg+min
	}

	for rateI, _ := range m.ProfitRate {
		for StopLossI, _ := range m.StopLossRate {
			for AtrLengthI, _ := range m.AtrLength {
				for ProfitTypeI, _ := range m.ProfitType {
					for LeverI, _ := range m.Lever {
						for AtrMoveTypeI, _ := range m.RenKoMoveType {
							for AtrMoveII, _ := range m.RenKoMoveI {
								icopy := mc
								icopy.ProfitRate = m.ProfitRate[rateI]
								icopy.StopLossRate = m.StopLossRate[StopLossI]
								icopy.AtrLength = m.AtrLength[AtrLengthI]
								icopy.ProfitType = m.ProfitType[ProfitTypeI]
								icopy.RenKoMoveType = m.RenKoMoveType[AtrMoveTypeI]
								icopy.RenKoMoveI = m.RenKoMoveI[AtrMoveII]
								icopy.MockDetail.Lever = int64(m.Lever[LeverI])
								icopy.Bn = m.Bn
								icopy.Brackets = m.Brackets
								icopy.IsTowWay = m.IsTowWay
								icopy.TotalBalanceRatio = m.TotalBalanceRatio
								r = append(r, &icopy)
							}
						}
					}
				}
			}
		}
	}
	m.WhereCycleOnes = r
	return r
}
func (m *MockCyCle) RunCycleWhere(actionType WhereType) {
	var list goex.MockResults
	log.Printf("数据回测 - 数据整理中....")
	if actionType == WhereTypeAll {
		list = m.MakeCycleWhereAll()
	} else {
		list = m.MakeCycleWhereSplit()
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
		log.Printf("%v分钟[%s-%s]回测数据: 初始余额:%v 当前余额:%v 交易次数:%v Atr:%v Atr移动:%v Atr方式:%v 平仓条件:%v 收益率:%v%% 止损率:%v %%  爆仓:%v \n",
			m.Cycle, m.StartDay, m.EndDay, itemMock.OldUsd, decimal.NewFromFloat(itemMock.Usd).StringFixed(2), itemMock.TradeNum,
			itemMock.AtrLength,
			itemMock.RenKoMoveI,
			itemMock.RenKoMoveType,
			itemMock.ProfitType,
			helper.Round((itemMock.Usd/itemMock.OldUsd*100)-100, 4),
			itemMock.StopLossRate,
			helper.IfThen(itemMock.IsLiquidation, "是", "否"))
		break
	}
	saveGroupList := func() {
		//将带指针的测试结果 转为不带指针的测试结果 否则xsl会报错
		var l []goex.MockResult
		for i := 0; i < len(list)-1; i++ {
			l = append(l, *list[i])
		}
		f = xjexcel.ListToExcel(l, "", fmt.Sprintf("%vm回测数据", m.Cycle))
		helper.Exists("mocktest", true, true)
		err := f.SaveAs(fmt.Sprintf("mocktest/%vm回测数据.xlsx", m.Cycle))
		if err != nil {
			log.Fatal("保存失败")
			return
		}
	}
	saveGroupList()
}

// MakeCycleWhereSplit 分割回测 防止数据过大 造成溢出
func (m *MockCyCle) MakeCycleWhereSplit() (list goex.MockResults) {
	// 将K线数据 按天分割后 逐一回测
	var klinData = binance.DownloadData("FUTURES_UM", "klines", "daily", m.StartDay, m.EndDay, []string{fmt.Sprintf("%vm", m.Cycle)}, []string{m.Symbol.ToSymbol("")})
	log.Printf("数据加载完毕，等待整合,[%v]\n", len(klinData))
	m.Brackets = LoadBrackets(m.Bn) //加载杠杆信息
	// 将多条件编译组合为多个单条件
	m.MakeCycleWhere()
	log.Printf("条件重组完成:共%v组数据", len(m.WhereCycleOnes))
	var wg = new(sync.WaitGroup)
	for i := 0; i < len(klinData); i++ {
		m.kLineData = []*goex.Kline{}
		// 按天分割后的K线数据与交易记录
		m.kLineData = goex.LocalKlineCsv(klinData[i], m.Symbol)
		log.Printf("[%v]日数据关联完成[%v]", time.UnixMilli(m.kLineData[0].CloseTime).Format(helper.TimeFormatYmdHis), len(m.klineLikeTrade))
		list = append(list, m.Run()...)
		wg.Wait()
		log.Printf("处理完毕")
	}
	return list
}
func (m *MockCyCle) MakeCycleWhereAll() goex.MockResults {
	m.Brackets = LoadBrackets(m.Bn) //加载杠杆信息
	log.Printf("加载杠杆信息:共%v组数据", len(m.Brackets))
	m.MakeCycleWhere()
	log.Printf("条件重组完成:共%v组数据", len(m.WhereCycleOnes))

	kline := binance.GetKLines("FUTURES_UM", "klines", "daily", m.StartDay, m.EndDay,
		[]string{fmt.Sprintf("%vm", m.Cycle)}, []string{m.Symbol.ToSymbol("")}, m.Symbol)
	log.Printf("K线数据读取完毕,共%v条 等待排序...", len(kline))
	m.kLineData = goex.KlineSort(kline, "asc")
	log.Println("K线数据排序完成")
	return m.Run()
}
func (m *MockCyCle) Run() (list goex.MockResults) {
	var wgA = new(sync.WaitGroup)
	// 将任务组 并行运行至单个任务组 - 此应是独立环境
	for i, _ := range m.WhereCycleOnes {
		wgA.Add(1)
		i := i
		go func() {
			aa := m.WhereCycleOnes[i].Res(m)
			list = append(list, aa)
			wgA.Done()
		}()
	}
	log.Println("任务分发完成，正在执行...")
	wgA.Wait()
	return list
}
