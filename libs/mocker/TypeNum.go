package mocker

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

const (
	WhereAll = 0 + iota
	wheresplit
)

type Brackets map[string]*goex.BracketsList
type WhereMul struct {
	Pair     []*goex.CurrencyPair
	Usd      float64 //模拟余额
	OldUsd   float64 //模拟余额
	Balance  float64 //币余额
	TradeNum int64   //下单次数
	Cone     []*WhereCycleOne
}

// WhereCycleOne 高精度测试
type WhereCycleOne struct {
	Comm
	Bn             *binance.Binance
	Brackets       Brackets
	Signal         goex.TradeSide //信号
	OlineType      int            //0=本地 1=实时数据 模拟交易 2=真实交易
	AtrLength      int
	MaxHold        float64                 //最大持仓量
	ProfitType     CloseSignal             //平仓条件
	StopLossRate   float64                 //止损百分比 - 当前亏损率
	IsLiquidation  bool                    //是否爆仓
	klinsLikeTrade map[int64][]*goex.Trade //K线周期对应的交易记录
	Gain           float64                 //收益
	NetGain        float64                 //净收益
	ProfitRate     float64                 //盈利平仓比率
	MockDetail     *goex.MockDetail
}
type Comm struct {
	wg        sync.WaitGroup
	mutex     sync.Mutex
	Symbol    goex.CurrencyPair // 参与币种
	KlineLast *goex.Kline
	kLineData []*goex.Kline
	Renko     []*goex.Kline
	TradeData []*goex.Trade
}

// MockCyCle 数据回测
type MockCyCle struct {
	Comm
	Cycle          int                     //测试周期
	KlinsLikeTrade map[int64][]*goex.Trade //K线周期对应的交易记录
	Lever          []int
	AtrLength      []int
	Bn             *binance.Binance
	MaxHold        float64   //最大持仓量
	StopLossRate   []float64 //止损值
	ProfitRate     []float64
	ProfitType     []CloseSignal
	StartDay       string
	EndDay         string
	Brackets       Brackets
}

func (m *MockCyCle) KlineLinkTrade() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	m.KlinsLikeTrade = make(map[int64][]*goex.Trade)
	fl := helper.MakeMathQuantity(60/m.Cycle, 0, m.Cycle)
	pool := tunny.NewFunc(runtime.NumCPU()*1024, func(i interface{}) interface{} {
		item := m.TradeData[i.(int)]
		t := time.UnixMilli(item.Date)
		level, _ := strconv.Atoi(fmt.Sprintf("%02d", fl[int(t.Minute()/m.Cycle)]+m.Cycle-1))
		t1 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), level, 59, 999999999, t.Location())
		mu.Lock()
		m.KlinsLikeTrade[t1.UnixMilli()] = append(m.KlinsLikeTrade[t1.UnixMilli()], item)
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
