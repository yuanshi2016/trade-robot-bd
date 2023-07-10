package goex

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"trade-robot-bd/libs/helper"
)

type Order struct {
	Price        float64
	Amount       float64
	AvgPrice     float64
	DealAmount   float64
	Fee          float64
	Cid          string //客户端自定义ID
	OrderID2     string
	OrderID      int //deprecated
	Status       TradeStatus
	Currency     CurrencyPair
	Side         TradeSide
	Type         string //limit / market
	OrderType    int    //0:default,1:maker,2:fok,3:ioc
	OrderTime    int    // create  timestamp
	FinishedTime int64  //finished timestamp
	Symbol       string
}

type Trade struct {
	Tid          int64        `json:"tid"`
	Type         TradeSide    `json:"type"`
	Amount       float64      `json:"amount,string"`
	Price        float64      `json:"price,string"`
	Date         int64        `json:"date_ms"`
	Pair         CurrencyPair `json:"omitempty"`
	Side         string       `json:"side"`
	Fee          string       `json:"fee"`
	InstrumentId string       `json:"instrument_id"`
	DateTime     time.Time    `json:"date_time"`
}

func BubbleSortGeneric1(x []*Trade, f int) []*Trade {
	var r []*Trade
	// 将切片切分排序重组 按照 (X.len - (X.len % Cpu.Runing)) / Cpu.Runing
	y := len(x) % f        //按照cpu数量拆分后的余数
	f_ := (len(x) - y) / f //整拆分后的数量
	leftNum := 0           //截止到上次
	wg := &sync.WaitGroup{}
	var mu sync.Mutex
	var dd = make(map[int][]*Trade)
	wg.Add(f)
	for si := 1; si <= f; si++ {
		currI := si * f_
		if si == f {
			currI += y
		}
		_si := si
		i1 := x[leftNum:currI]
		leftNum = currI
		go func([]*Trade) {
			it := TradeSort(i1, "asc")
			mu.Lock()
			dd[_si-1] = it
			mu.Unlock()
			//log.Printf("第%v次排序完成:Start:%v end：%v count:%v", _si, it[0].Date, it[len(it)-1].Date, len(dd))
			wg.Done()
		}(r)
	}
	wg.Wait()
	for i := f; i >= 0; i-- {
		r = append(r, dd[i]...)
	}
	return r
}
func TradeSort(c []*Trade, actions string) []*Trade {
	switch actions {
	case "asc":
		sort.Slice(c, func(i, j int) bool {
			return c[i].Date < c[j].Date
		})
		break
	case "desc":
		sort.Slice(c, func(i, j int) bool {
			return c[i].Date > c[j].Date
		})
		break
	case "Reverse":
		for i, j := 0, len(c)-1; i < j; i, j = i+1, j-1 {
			c[i].Date, c[i].Date = c[j].Date, c[j].Date
		}
		break
	default:
		break
	}
	return c
}

type SubAccount struct {
	Currency     Currency
	Amount       float64
	ForzenAmount float64
	LoanAmount   float64
	Balance      float64
}

type MarginSubAccount struct {
	Balance     float64
	Frozen      float64
	Available   float64
	CanWithdraw float64
	Loan        float64
	LendingFee  float64
}

type MarginAccount struct {
	Sub              map[Currency]MarginSubAccount
	LiquidationPrice float64
	RiskRate         float64
	MarginRatio      float64
}

type Account struct {
	Exchange    string
	Asset       float64 //总资产
	NetAsset    float64 //净资产
	SubAccounts map[Currency]SubAccount
}

type AccountTrade struct {
	Currency     string               `json:"currency"`
	Details      *AccountTradeDetails `json:"details"`
	InstrumentID string               `json:"instrument_id"`
	Amount       string               `json:"amount"`
	Balance      string               `json:"balance"`
	LedgerID     string               `json:"ledger_id"`
	Timestamp    time.Time            `json:"timestamp"`
	Fee          string               `json:"fee"`
	Type         string               `json:"type"`
}

type AccountTradeDetails struct {
	InstrumentID string `json:"instrument_id"`
	OrderID      string `json:"order_id"`
}

type Ticker struct {
	Symbol string       `json:"symbol"`
	Pair   CurrencyPair `json:"omitempty"`
	Last   float64      `json:"last,string"`
	Buy    float64      `json:"buy,string"`
	Open   float64      `json:"open,string"`
	Sell   float64      `json:"sell,string"`
	High   float64      `json:"high,string"`
	Low    float64      `json:"low,string"`
	Vol    float64      `json:"vol,string"`
	Date   uint64       `json:"date"` // 单位:ms
}

type FutureTicker struct {
	*Ticker
	ContractType string  `json:"omitempty"`
	ContractId   int     `json:"contractId"`
	LimitHigh    float64 `json:"limitHigh,string"`
	LimitLow     float64 `json:"limitLow,string"`
	HoldAmount   float64 `json:"hold_amount,string"`
	UnitAmount   float64 `json:"unitAmount,string"`
}

type DepthRecord struct {
	Price  float64
	Amount float64
}

type DepthRecords []DepthRecord

func (dr DepthRecords) Len() int {
	return len(dr)
}

func (dr DepthRecords) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DepthRecords) Less(i, j int) bool {
	return dr[i].Price < dr[j].Price
}

type Depth struct {
	ContractType string //for future
	Pair         CurrencyPair
	UTime        time.Time
	AskList      DepthRecords // Descending order
	BidList      DepthRecords // Descending order
}

type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string //for okex.com v3 api
	ClientId      string //for bitstamp.net , huobi.pro
	ClientType    string
	Lever         int //杠杆倍数 , for future
}

type Kline struct {
	Pair        CurrencyPair
	Timestamp   int64
	Open        float64
	Close       float64
	High        float64
	Low         float64
	Vol         float64 //成交额
	CloseTime   int64
	QuoteVolume float64 //主动买入成交额
	Count       int64   //成交笔数
}

func KlineSort(c []*Kline, actions string) (d []*Kline) {
	switch actions {
	case "asc":
		sort.Slice(c, func(i, j int) bool {
			return c[i].CloseTime < c[j].CloseTime
		})
		break
	case "desc":
		sort.Slice(c, func(i, j int) bool {
			return c[i].CloseTime > c[j].CloseTime
		})
		break
	case "Reverse":
		for i, j := 0, len(c)-1; i < j; i, j = i+1, j-1 {
			c[i].CloseTime, c[i].CloseTime = c[j].CloseTime, c[j].CloseTime
		}
		break
	default:
		break
	}
	return c
}

type FutureKline struct {
	*Kline
	Vol2 float64 //个数
}

type FutureSubAccount struct {
	Currency          Currency
	AccountRights     float64 //账户权益
	KeepDeposit       float64 //保证金
	ProfitReal        float64 //已实现盈亏
	ProfitUnreal      float64
	RiskRate          float64 //保证金率
	TotalAvailBalance float64
	MarginFrozen      float64
	Symbol            string
}

type FutureAccount struct {
	FutureSubAccounts map[Currency]FutureSubAccount
}

type FutureOrder struct {
	ClientOid      string //自定义ID，GoEx内部自动生成
	OrderID2       string //请尽量用这个字段替代OrderID字段
	Price          float64
	Amount         float64
	AvgPrice       float64
	DealAmount     float64
	OrderID        int64 //deprecated
	OrderTime      int64
	Status         TradeStatus
	State          string
	Currency       CurrencyPair
	OrderType      int       //ORDINARY=0 POST_ONLY=1 FOK= 2 IOC= 3
	OType          TradeSide //1：开多 2：开空 3：平多 4： 平空
	LeverRate      int       //倍数
	Fee            float64   //手续费
	ContractName   string
	FinishedTime   int64 // finished timestamp
	LastFillId     string
	LastFillPrice  float64 //最新成交价格（如果没有，推0
	LastFillAmount float64 //最新成交数量（如果没有，推0）
	LastFillTime   time.Time
}

type FuturePosition struct {
	BuyAmount      float64
	BuyAvailable   float64
	BuyPriceAvg    float64
	BuyPriceCost   float64
	BuyProfitReal  float64
	CreateDate     int64
	LeverRate      int
	SellAmount     float64
	SellAvailable  float64
	SellPriceAvg   float64
	SellPriceCost  float64
	SellProfitReal float64
	Symbol         CurrencyPair //btc_usd:比特币,ltc_usd:莱特币
	ContractType   string
	ContractId     int64
	ForceLiquPrice float64 //预估爆仓价
}

type HistoricalFunding struct {
	InstrumentId string    `json:"instrument_id"`
	RealizedRate float64   `json:"realized_rate,string"`
	FundingTime  time.Time `json:"funding_time"`
}

type TickSize struct {
	InstrumentID    string
	UnderlyingIndex string
	QuoteCurrency   string
	PriceTickSize   float64 //下单价格精度
	AmountTickSize  float64 //数量精度
}

type FuturesContractInfo struct {
	*TickSize
	ContractVal  float64 //合约面值(美元)
	Delivery     string  //交割日期
	ContractType string  //	本周 this_week 次周 next_week 季度 quarter
}

//api parameter struct

type BorrowParameter struct {
	CurrencyPair CurrencyPair
	Currency     Currency
	Amount       float64
}

type RepaymentParameter struct {
	BorrowParameter
	BorrowId string
}

// HttpResult --- Http 币安接口结构体
type HttpResult struct {
	Code          string         `json:"code"`
	Message       string         `json:"message"`
	MessageDetail string         `json:"messageDetail"`
	Data          HttpResultData `json:"data"`
	Success       bool           `json:"success"`
}
type HttpResultData struct {
	DownloadItemList DownloadItemLists `json:"downloadItemList"`
	Brackets         []BracketsList    `json:"brackets"`
}
type BracketsList struct {
	Symbol        string         `json:"symbol"`
	UpdateTime    int64          `json:"updateTime"`
	NotionalLimit int64          `json:"notionalLimit"`
	RiskBrackets  []RiskBrackets `json:"riskBrackets"`
}
type RiskBrackets struct {
	BracketSeq                   float64 `json:"bracketSeq"`
	BracketNotionalFloor         float64 `json:"bracketNotionalFloor"`
	BracketNotionalCap           float64 `json:"bracketNotionalCap"`
	BracketMaintenanceMarginRate float64 `json:"bracketMaintenanceMarginRate"`
	CumFastMaintenanceAmount     float64 `json:"cumFastMaintenanceAmount"`
	MinOpenPosLeverage           float64 `json:"minOpenPosLeverage"`
	MaxOpenPosLeverage           float64 `json:"maxOpenPosLeverage"`
}
type DownloadItemLists []*DownloadItemList
type DownloadItemList struct {
	ProductName string `json:"productName"`
	BizType     string `json:"bizType"`
	Granularity string `json:"granularity"`
	Interval    string `json:"interval"`
	Symbol      string `json:"symbol"`
	Day         string `json:"day"`
	Url         string `json:"url"`
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
}

func (d DownloadItemLists) Sort(actions string) {

	switch actions {
	case "asc":
		sort.Slice(d, func(i, j int) bool {
			l, _ := strconv.Atoi(helper.StringJoinString(strings.Split(d[i].Day, "-")...))
			r, _ := strconv.Atoi(helper.StringJoinString(strings.Split(d[j].Day, "-")...))
			return l < r
		})
		break
	case "desc":
		sort.Slice(d, func(i, j int) bool {
			l, _ := strconv.Atoi(helper.StringJoinString(strings.Split(d[i].Day, "-")...))
			r, _ := strconv.Atoi(helper.StringJoinString(strings.Split(d[j].Day, "-")...))
			return l > r
		})
		break
	default:
		break
	}
}
