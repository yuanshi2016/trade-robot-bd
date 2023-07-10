package goex

import (
	"encoding/csv"
	"fmt"
	"github.com/Jeffail/tunny"
	"github.com/shopspring/decimal"
	"github.com/tealeg/xlsx"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
	"trade-robot-bd/app/grid-strategy-svc/util/goex"
	"trade-robot-bd/app/grid-strategy-svc/util/utils"
	"trade-robot-bd/libs/helper"
)

const (
	NewOrder_CM   = -1 + iota //币本位
	NewOrder_SPOT             //现货
	NewOrder_UM               //U本位
)
const (
	NewOrder_Sell = -1 + iota //币本位
	NewOrder_Buy  = iota      //U本位
)

type MockOrder struct {
	Open        float64 `json:"open"`        //开仓价
	Close       float64 `json:"close"`       //平仓价
	HighClose   float64 `json:"highClose"`   //平仓成本价
	Liquidation float64 `json:"liquidation"` //爆仓价
	Direction   int64   `json:"direction"`   //-1 空 1多
	Type        int64   `json:"type"`        // 1 U本位 -1 币本位 0 现货
	Lever       int64   `json:"lever"`       //杠杆
	Quantity    float64 `json:"quantity"`    //开仓数量/张数
	CpUsd       int64   `json:"cpUsd"`       //面值 仅币本位使用
	FeeRate     float64 `json:"feeRate"`     //手续费比例
	Fee         float64 `json:"fee"`         //手续费
	FeeUsd      float64 `json:"feeUsd"`      //手续费USD - 仅币本位或现货
	FeeDiscount float64 `json:"feeDiscount"` //手续费折扣
	Bail        float64 `json:"bail"`        //保证金
	Gain        float64 `json:"gain"`        //收益
	NetGain     float64 `json:"netGain"`     //净收益
	NetGainUSd  float64 `json:"netGainUSd"`  //收益转USD
	Rate        float64 `json:"rate"`        //收益率
	NetRate     float64 `json:"netRate"`     //净收益率
	Usd         float64 //Usd快照
	BidTime     int64
	AskTime     int64
}

type MockResults []MockResult
type MockResult struct {
	OldUsd        float64 `excel:"column:A;desc:初始余额;width:30"`
	Usd           float64 `excel:"column:B;desc:当前余额;width:30"`
	RsiMin        float64 `excel:"column:C;desc:Rsi买入;width:30"`
	RsiMax        float64 `excel:"column:D;desc:Rsi卖出;width:30"`
	RsiLength     int     `excel:"column:E;desc:Rsi长度;width:30"`
	RviMin        float64 `excel:"column:F;desc:Rvi买入;width:30"`
	RviMax        float64 `excel:"column:G;desc:Rvi买入;width:30"`
	RviLength     int     `excel:"column:H;desc:Rvi长度;width:30"`
	MfiMin        float64 `excel:"column:I;desc:Mfi买入;width:30"`
	MfiMax        float64 `excel:"column:J;desc:Mfi买入;width:30"`
	MfiLength     int     `excel:"column:K;desc:Mfi长度;width:30"`
	AtrLength     int     `excel:"column:L;desc:Atr长度;width:30"`
	TradeNum      int64   `excel:"column:M;desc:交易次数;width:30"`
	ProfitRate    float64 `excel:"column:N;desc:收益率;width:30"`
	ProfitType    string  `excel:"column:O;desc:平仓条件;width:30"`
	StopLossRate  float64 `excel:"column:P;desc:止损率;width:30"`
	IsLiquidation bool    `json:"isLiquidation" excel:"column:Q;desc:是否爆仓;width:30"` //是否爆仓
	Order         []*MockOrder
}

type MockDetail struct { //仅回测下使用 实盘状态下无用
	Pair         *goex.CurrencyPair
	Usd          float64      //模拟余额
	OldUsd       float64      //模拟余额
	Balance      float64      //币余额
	TradeNum     int64        //下单次数
	ProfitRate   int64        //累计收益率
	BuyOrder     *MockOrder   //买入或做多
	SellOrder    *MockOrder   //做空订单
	HistoryOrder []*MockOrder //历史订单
	Mutex        sync.Mutex
	Direction    int64   `json:"direction"` //-1 空 1多
	Type         int64   `json:"type"`      // 1 U本位 -1 币本位 0 现货
	Lever        int64   `json:"lever"`     //杠杆
	FeeRate      float64 `json:"feeRate"`   //手续费比例
	CpUsd        int64   `json:"cpUsd"`     //面值 仅币本位使用
}

// ExportOrders 导出最优订单
func (c MockResult) ExportOrders(Cycle int, start, end string) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("student_list")
	if err != nil {
		fmt.Printf(err.Error())
	}
	ce := sheet.AddRow().AddCell()
	ce.HMerge = 11
	ce.SetValue(fmt.Sprintf("%v分钟[%s-%s]回测数据: 初始余额:%v 当前余额:%v 交易次数:%v Atr:%v 止损率:%v 爆仓:%v \n",
		Cycle, start, end, c.OldUsd, decimal.NewFromFloat(c.Usd).StringFixed(2), c.TradeNum,
		c.AtrLength,
		c.StopLossRate,
		helper.IfThen(c.IsLiquidation, "是", "否")))
	row := sheet.AddRow()
	row.AddCell().SetValue("订单方向")
	row.AddCell().SetValue("持仓数量")
	row.AddCell().SetValue("买入价格")
	row.AddCell().SetValue("买入时间")
	row.AddCell().SetValue("卖出价格")
	row.AddCell().SetValue("卖出时间")
	row.AddCell().SetValue("强平价格")
	row.AddCell().SetValue("保证金")
	row.AddCell().SetValue("净收益")
	row.AddCell().SetValue("余额快照")
	for i := 0; i < len(c.Order); i++ {
		order := c.Order[i]
		row := sheet.AddRow()
		row.AddCell().SetValue(helper.IfThen(order.Direction < 0, "做空", "做多"))
		row.AddCell().SetValue(order.Quantity)
		row.AddCell().SetValue(order.Open)
		row.AddCell().SetValue(time.UnixMilli(order.BidTime).Format(helper.TimeFormatYmdHis))
		row.AddCell().SetValue(order.Close)
		row.AddCell().SetValue(time.UnixMilli(order.AskTime).Format(helper.TimeFormatYmdHis))
		row.AddCell().SetValue(order.Liquidation)
		row.AddCell().SetValue(order.Bail)
		row.AddCell().SetValue(order.NetGain)
		row.AddCell().SetValue(decimal.NewFromFloat(order.Usd).StringFixed(2))
	}
	err = file.Save(fmt.Sprintf("./mocktest/%vm回测数据最优数据.xlsx", Cycle))
	if err != nil {
		fmt.Printf(err.Error())
	}
}

// SplitToFlotList
//
//	@Description:
//	@receiver c
//	@param actions排序类型
//	@param showLiquidation 是否显示爆仓数据
//	@return i32List
func (c MockResults) SplitToFlotList(actions string, showLiquidation bool) (i32List MockResults) {
	if len(c) == 0 {
		return
	}
	for _, item := range c {
		if !item.IsLiquidation || showLiquidation {
			i32List = append(i32List, item)
		}
	}
	switch actions {
	case "asc":
		sort.Slice(i32List, func(i, j int) bool {
			return i32List[i].Usd < i32List[j].Usd
		})
		break
	case "desc":
		sort.Slice(i32List, func(i, j int) bool {
			return i32List[i].Usd > i32List[j].Usd
		})
		break
	case "Reverse":
		for i, j := 0, len(i32List)-1; i < j; i, j = i+1, j-1 {
			i32List[i].Usd, i32List[i].Usd = i32List[j].Usd, i32List[j].Usd
		}
		break
	default:
		break
	}
	return i32List
}

// CalcCpUp 计算盈亏详情
func (conf *MockOrder) CalcCpUp() *MockOrder {
	// U本位实际盈亏 = ((1 / 开仓价) - (1 / 平仓价)) * (开仓数量 * 平仓价)
	// U本位开仓保证金 = 开仓数量  * (1 / 杠杆)
	// U本位盈亏比例 = (实际盈亏 / 保证金) * 100
	// 币本位实际盈亏 = ((1 / 开仓价) - (1 / 平仓价)) * (开仓数量 * 面值)
	// 币本位开仓保证金 = 初始保证金= 面值 * 开仓张数 / (平仓价* 杠杆倍数)
	// 币本位盈亏比例 = (实际盈亏 / 保证金) * 100
	if conf.Type == NewOrder_CM && conf.CpUsd <= 1 {
		log.Fatalf("币本位下，面值为必填")
	}
	var rate float64
	switch conf.Type {
	case NewOrder_SPOT:
		conf.Gain = ((1 / conf.Open) - (1 / conf.Close)) * (conf.Quantity * conf.Close)
		conf.NetGain = conf.Gain - conf.FeeUsd
		conf.NetGainUSd = conf.NetGain
		rate = (conf.Gain / conf.Quantity) * 100
		conf.NetRate = (conf.NetGain / conf.Quantity) * 100
		break
	case NewOrder_CM:
		conf.Bail = conf.Quantity * float64(conf.CpUsd) / (conf.Close * float64(conf.Lever))
		conf.Gain = ((1 / conf.Open) - (1 / conf.Close)) * (helper.IfThen(conf.Direction == NewOrder_Buy, conf.Quantity, -conf.Quantity) * float64(conf.CpUsd))
		conf.NetGain = conf.Gain - conf.Fee
		conf.NetGainUSd = conf.NetGain * conf.Close
		rate = (conf.Gain / conf.Bail) * 100
		conf.NetRate = (conf.NetGain / conf.Bail) * 100
		break
	case NewOrder_UM:
		conf.Bail = conf.Quantity / float64(conf.Lever)
		conf.Gain = ((1 / conf.Open) - (1 / conf.Close)) * (helper.IfThen(conf.Direction == NewOrder_Buy, conf.Quantity, -conf.Quantity) * conf.Close)
		rate = (conf.Gain / conf.Bail) * 100
		conf.NetGain = conf.Gain - conf.FeeUsd
		conf.NetGainUSd = conf.NetGain
		conf.NetRate = (conf.NetGain / conf.Bail) * 100
		break
	}
	conf.Rate = rate
	return conf
}

// CalcHithClose 计算成本价
func (conf *MockOrder) CalcHithClose() {
	switch conf.Type {
	case NewOrder_SPOT:
		conf.HighClose = (((conf.FeeRate * 2.03) * conf.Open) / 100) + conf.Open
		//conf.HighClose = ((conf.Fee * 2) * conf.Open) + conf.Open
		break
	case NewOrder_CM:
		// 币本位成本价计算 (手续费 * 开仓价) + 开仓价
		conf.HighClose = (((conf.FeeRate * 2.03) * conf.Open) / 100) + conf.Open
		//conf.HighClose = ((conf.Fee * 2) * conf.Open) + conf.Open
		break
	case NewOrder_UM:
		conf.HighClose = (((conf.FeeRate * 2.03) * conf.Open) / 100) + conf.Open
		break
	}
}

// CalcLiquidation 计算爆仓价
func (conf *MockOrder) CalcLiquidation(WB float64, info *BracketsList) {
	var itemInfo RiskBrackets
	for i := range info.RiskBrackets {
		if conf.Quantity <= info.RiskBrackets[i].BracketNotionalCap {
			itemInfo = info.RiskBrackets[i]
		}
		break
	}
	//var xUsdt = conf.Open*conf.Quantity - 58045
	var TMM1 float64 = 0                             //其他合约下的全部保证金(除合约1外)
	var UPNL1 float64 = 0                            // 全部其他合约的未实现盈亏(除合约1外)
	var cumB = itemInfo.CumFastMaintenanceAmount     //单向模式下 合约1的维持保证金速算额
	var cumL float64 = 0                             //开多合约1下的维持保证金速算额(单向持仓模式)
	var cumS float64 = 0                             //开空合约1下的维持保证金速算额(双向持仓模式)
	var Side1BOTH float64 = 1                        // 合约1的方向(单向持仓模式) 1=开多 -1=开空
	var Position1BOTH = conf.Quantity / conf.Open    //合约1 的持仓大小 (单向持仓模式) 无论开多或开空 取绝对值
	var EP1BOTH = conf.Open                          // 合约1的头寸价格 (单向持仓模式)
	var Position1LONG float64 = 0                    //开多仓位大小（双向持仓模式）； 无论开多或开空 取绝对值
	var EP1LONG float64 = 0                          //开多持仓头寸（双向持仓模式）； 无论开多或开空 取绝对值
	var Position1SHORT float64 = 0                   //开空仓位大小 (双向持仓模式) 无论开多或开空 取绝对值
	var EP1SHORT float64 = 0                         //开空持仓的头寸(双向持仓模式) 无论开多或开空 取绝对值
	var MMRB = itemInfo.BracketMaintenanceMarginRate // 单向持仓模式合约的维持保证金费率
	var MMRL float64 = 0                             //开多合约的维持保证金费率(双向持仓模式)
	var MMRS float64 = 0                             //开空合约的维持保证金费率(双向持仓模式)
	lpt := WB - TMM1 + UPNL1 + cumB + cumL + cumS - Side1BOTH*Position1BOTH*EP1BOTH - Position1LONG*EP1LONG + Position1SHORT*EP1SHORT
	lpb := Position1BOTH*MMRB + Position1LONG*MMRL + Position1SHORT*MMRS - Side1BOTH*Position1BOTH - Position1LONG + Position1SHORT
	conf.Liquidation = lpt / lpb
	conf.Liquidation = helper.IfThen(conf.Liquidation < 0, 0, conf.Liquidation)
}

// Buy 做多/平空
func (conf *MockOrder) Buy(price float64) *MockOrder {
	conf.Open = price
	switch conf.Type {
	case NewOrder_SPOT:
		conf.Fee += (conf.Quantity / conf.Open) * conf.FeeRate / 100
		conf.FeeUsd += conf.Fee * conf.Open
		break
	case NewOrder_CM:
		conf.Fee += (conf.Quantity * float64(conf.CpUsd) / conf.Open) * conf.FeeRate / 100
		conf.FeeUsd = conf.Fee * conf.Open
		break
	case NewOrder_UM:
		conf.Fee += conf.Quantity * (conf.FeeRate / 100)
		conf.FeeUsd += conf.Fee
		break
	}
	conf.CalcHithClose()
	return conf
}

// Sell 做空/平多
func (conf *MockOrder) Sell(price float64) *MockOrder {
	conf.Close = price
	switch conf.Type {
	case NewOrder_SPOT:
		var sellFee = (((conf.Quantity / conf.Open) * conf.Close) * conf.FeeRate / 100) / conf.Close
		conf.Fee += sellFee
		conf.FeeUsd += sellFee * conf.Close
		break
	case NewOrder_CM:
		var sellFee = (conf.Quantity * float64(conf.CpUsd) / conf.Close) * conf.FeeRate / 100
		conf.Fee += sellFee
		conf.FeeUsd = conf.Fee * conf.Close
		break
	case NewOrder_UM:
		conf.Fee += conf.Quantity * (conf.FeeRate / 100)
		conf.FeeUsd = conf.Fee
		break
	}
	return conf
}

// LocalTradeCsvSpili 读取本地交易数据
func LocalTradeCsvSpili(file string, curr CurrencyPair) (tr []*Trade) {
	csvFile, err := os.Open(file)
	csvReader := csv.NewReader(csvFile)
	if err != nil {
		log.Fatal("读取失败")
	}
	defer csvFile.Close()
	//csvReader.Comma = ','             // 分隔符
	//csvReader.FieldsPerRecord = 3     // 每行有几个字段
	//csvReader.TrimLeadingSpace = true // 是否去除前导空格

	chunkSize := 1024 * 200 // 每个分片的大小
	chunk := make([][]string, chunkSize)
	kk := 0
	for {
		for i := 0; i < chunkSize; i++ {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			chunk[i] = record
		}
		// 消费分片
		for _, record := range chunk {
			if kk == 0 {
				kk = 1
				continue
			}
			if len(record) <= 0 {
				goto loop
			}
			t := &Trade{
				Pair:   curr,
				Tid:    int64(ToUint64(record[0])),
				Type:   TradeSide(utils.IfThen(record[5] == "true", 1, 2)),
				Amount: ToFloat64(record[2]),
				Price:  ToFloat64(record[1]),
				Date:   int64(ToUint64(record[4])),
			}
			tr = append(tr, t)
		}
		// 清空分片，准备下一轮读取
		chunk = make([][]string, chunkSize)
	}
loop:
	return tr
}
func LocalTradeCsv(file string, curr CurrencyPair) []*Trade {
	csvFile, err := os.Open(file)
	csvReader := csv.NewReader(csvFile)
	if err != nil {
		log.Fatal("读取失败")
	}
	defer csvFile.Close()
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("%v读取失败:%v", file, err.Error())
	}
	var tr []*Trade

	for i := range rows {
		if i == 0 {
			continue
		}
		row := rows[i]
		t := &Trade{
			Pair:   curr,
			Tid:    int64(ToUint64(row[0])),
			Type:   TradeSide(utils.IfThen(row[5] == "true", 1, 2)),
			Amount: ToFloat64(row[2]),
			Price:  ToFloat64(row[1]),
			Date:   int64(ToUint64(row[4])),
		}
		tr = append(tr, t)
	}
	//wg.Wait()
	return tr
}

func LocalKlineCsv(file string, curr CurrencyPair) (k []*Kline) {
	csvFile, err := os.Open(file)
	csvReader := csv.NewReader(csvFile)
	if err != nil {
		log.Fatal("读取失败")
	}

	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("%v读取失败:%v", file, err.Error())
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	pool := tunny.NewFunc(runtime.NumCPU()*10, func(i interface{}) interface{} {
		row := rows[i.(int)]
		t := &Kline{
			Pair:        curr,
			Timestamp:   int64(ToUint64(row[0])),
			Open:        ToFloat64(row[1]),
			Close:       ToFloat64(row[4]),
			High:        ToFloat64(row[2]),
			Low:         ToFloat64(row[3]),
			Vol:         ToFloat64(row[5]),
			CloseTime:   int64(ToUint64(row[6])),
			QuoteVolume: float64(ToUint64(row[7])),
			Count:       int64(ToUint64(row[8])),
		}
		mu.Lock()
		k = append(k, t)
		mu.Unlock()
		wg.Done()
		return nil
	})
	defer func() {
		csvFile.Close()
		pool.Close()
	}()
	for i := range rows {
		if i == 0 {
			continue
		}
		wg.Add(1)
		go pool.Process(i)
	}
	wg.Wait()
	return k
}
