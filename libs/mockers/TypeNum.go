package mockers

import (
	"fmt"
	"github.com/Jeffail/tunny"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
	"trade-robot-bd/libs/goex"
	"trade-robot-bd/libs/goex/binance"
	"trade-robot-bd/libs/helper"
)

// -- 模拟类型 start
type WhereType int

const (
	WhereTypeAll WhereType = 0 + iota
	WhereTypeSplit
)

func (w WhereType) String() string {
	switch w {
	case WhereTypeAll:
		return "全部数据"
	case WhereTypeSplit:
		return "分段模拟"
	}
	return "--"
}
func (w WhereType) Int() int {
	switch w {
	case WhereTypeAll:
		return int(WhereTypeAll)
	case WhereTypeSplit:
		return int(WhereTypeSplit)
	}
	return 0
}

// -- 模拟类型 end
// -- 测试类型 start
type TradeType int

const (
	TradeTypeLocal      TradeType = 0 + iota //回测交易
	TradeTypeOnlineData                      //仅实时数据 交易为模拟
	TradeTypeOline                           //完全在线数据 :数据与交易都为线上 会发生真实交易
)

func (w TradeType) String() string {
	switch w {
	case TradeTypeLocal:
		return "回测模拟"
	case TradeTypeOnlineData:
		return "实时回测"
	case TradeTypeOline:
		return "实时交易"
	}
	return "--"
}

// -- 测试类型 end

// WhereCycleOne 单条件测试
type WhereCycleOne struct {
	Comm
	KlineLast         *goex.Kline        //最新一条K线
	Ticker            *goex.Ticker       //最新一条K线
	Signal            goex.TradeSide     //交易信号信号
	AtrLength         int                //Atr指标长度
	RenKoMoveType     goex.RenKoMoveType //Renko指标移动类型
	RenKoMoveI        int                //Renko指标平均价移动取值
	ProfitType        CloseSignal        //平仓条件
	StopLossRate      float64            //止损百分比 - 当前亏损率
	IsLiquidation     bool               //是否爆仓
	Gain              float64            //收益
	NetGain           float64            //净收益
	ProfitRate        float64            //盈利平仓比率
	IsTowWay          bool               //是否双向持仓
	TotalBalanceRatio float64            // 总资金分配占比 - 仅实盘
	MockDetail        goex.MockDetail
	MockResult        *goex.MockResult
}
type Comm struct {
	TradeType      TradeType
	Bn             *binance.Binance
	BnWs           *binance.BinanceWs
	BnSwap         *binance.BinanceSwap
	Symbol         goex.CurrencyPair       // 参与币种
	kLineData      []*goex.Kline           //K线数据
	Renko          []*goex.Kline           //Renko指标
	TradeData      []*goex.Trade           //交易记录
	klineLikeTrade map[int64][]*goex.Trade //K线周期对应的交易记录
	Brackets       map[string]*goex.BracketsList
	MaxHold        float64 //最大持仓量
}

// MockCyCle 数据回测
type MockCyCle struct {
	Comm
	wg                sync.WaitGroup
	Cycle             int //测试周期
	Lever             []int
	AtrLength         []int
	RenKoMoveType     []goex.RenKoMoveType
	RenKoMoveI        []int
	Bn                *binance.Binance
	StopLossRate      []float64 //止损值
	ProfitRate        []float64
	ProfitType        []CloseSignal
	WhereCycleOnes    []*WhereCycleOne
	Usd               float64 //模拟余额
	Balance           float64 //币余额
	FeeRate           float64 //手续费比例
	StartDay          string
	EndDay            string
	IsTowWay          bool //是否单向持仓
	Brackets          map[string]*goex.BracketsList
	MockResults       goex.MockResults
	TotalBalanceRatio float64
}

func (m *MockCyCle) KlineLinkTrade() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	m.klineLikeTrade = make(map[int64][]*goex.Trade)
	fl := helper.MakeMathQuantity(60/m.Cycle, 0, m.Cycle)
	pool := tunny.NewFunc(runtime.NumCPU()*1024, func(i interface{}) interface{} {
		item := m.TradeData[i.(int)]
		t := time.UnixMilli(item.Date)
		level, _ := strconv.Atoi(fmt.Sprintf("%02d", fl[int(t.Minute()/m.Cycle)]+m.Cycle-1))
		t1 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), level, 59, 999999999, t.Location())
		mu.Lock()
		m.klineLikeTrade[t1.UnixMilli()] = append(m.klineLikeTrade[t1.UnixMilli()], item)
		mu.Unlock()
		wg.Done()
		return nil
	})
	wg.Add(len(m.TradeData))
	for i := 0; i < len(m.TradeData); i++ {
		go pool.Process(i)
	}
	wg.Wait()
}
func (m *MockCyCle) Sort() (mapValues []helper.MapValue) {
	sVal := reflect.ValueOf(m)
	sType := reflect.TypeOf(m)
	if sType.Kind() == reflect.Ptr {
		//用Elem()获得实际的value
		sVal = sVal.Elem()
		sType = sType.Elem()
	}
	num := sVal.NumField()
	for i := 0; i < num; i++ {
		f := sType.Field(i)
		if f.Type == reflect.TypeOf([]int{}) {
			val := sVal.Field(i).Interface().([]int)
			mapValues = append(mapValues, helper.MapValue{Key: f.Name, Value: val})
		}
	}
	sort.Slice(mapValues, func(i, j int) bool {
		return len(mapValues[i].Value) > len(mapValues[j].Value)
	})
	fmt.Println(mapValues)
	return mapValues
}

// CloseSignal -----------平仓条件 start
type CloseSignal int

const (
	ProfitRate   CloseSignal = -1
	ProfitSignal CloseSignal = 1
)

func (s CloseSignal) String() string {
	switch s {
	case ProfitRate:
		return "收益率"
	case ProfitSignal:
		return "平仓信号"
	default:
		return "Unknown"
	}
}

// CloseSignal -----------平仓条件 end
