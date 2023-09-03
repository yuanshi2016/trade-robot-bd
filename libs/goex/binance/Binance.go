package binance

import (
	"errors"
	"fmt"
	"github.com/Jeffail/tunny"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	. "trade-robot-bd/libs/goex"
	"trade-robot-bd/libs/helper"
)

const (
	GLOBAL_API_BASE_URL = "https://api.binance.com"
	WWW_API_BASE_URL    = "https://www.binance.com"
	F_API_BASE_URL      = "https://fapi.binance.com"
	US_API_BASE_URL     = "https://api.binance.us"
	JE_API_BASE_URL     = "https://api.binance.je"
	//API_V1       = API_BASE_URL + "api/v1/"
	//API_V3       = API_BASE_URL + "api/v3/"

	TICKER_URI             = "ticker/24hr?symbol=%s"
	TICKER_ALL_URI         = "ticker/24hr"
	TICKERS_URI            = "ticker/allBookTickers"
	DEPTH_URI              = "depth?symbol=%s&limit=%d"
	ACCOUNT_URI            = "account?"
	ORDER_URI              = "order"
	LEVERAGE_URI           = "leverage"
	UNFINISHED_ORDERS_INFO = "openOrders?"
	KLINE_URI              = "klines"
	SERVER_TIME_URL        = "time"
	CREATE_SUB_ACC_URL     = ""
	UserDataStream         = "userDataStream?"
)

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:   "1m",
	KLINE_PERIOD_3MIN:   "3m",
	KLINE_PERIOD_5MIN:   "5m",
	KLINE_PERIOD_15MIN:  "15m",
	KLINE_PERIOD_30MIN:  "30m",
	KLINE_PERIOD_60MIN:  "1h",
	KLINE_PERIOD_1H:     "1h",
	KLINE_PERIOD_2H:     "2h",
	KLINE_PERIOD_4H:     "4h",
	KLINE_PERIOD_6H:     "6h",
	KLINE_PERIOD_8H:     "8h",
	KLINE_PERIOD_12H:    "12h",
	KLINE_PERIOD_1DAY:   "1d",
	KLINE_PERIOD_3DAY:   "3d",
	KLINE_PERIOD_1WEEK:  "1w",
	KLINE_PERIOD_1MONTH: "1M",
}

type Filter struct {
	FilterType          string  `json:"filterType"`
	MaxPrice            float64 `json:"maxPrice,string"`
	MinPrice            float64 `json:"minPrice,string"`
	TickSize            float64 `json:"tickSize,string"`
	MultiplierUp        float64 `json:"multiplierUp,string"`
	MultiplierDown      float64 `json:"multiplierDown,string"`
	AvgPriceMins        int     `json:"avgPriceMins"`
	MinQty              float64 `json:"minQty,string"`
	MaxQty              float64 `json:"maxQty,string"`
	StepSize            float64 `json:"stepSize,string"`
	MinNotional         float64 `json:"minNotional,string"`
	ApplyToMarket       bool    `json:"applyToMarket"`
	Limit               int     `json:"limit"`
	MaxNumAlgoOrders    int     `json:"maxNumAlgoOrders"`
	MaxNumIcebergOrders int     `json:"maxNumIcebergOrders"`
	MaxNumOrders        int     `json:"maxNumOrders"`
}

type RateLimit struct {
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
	RateLimitType string `json:"rateLimitType"`
}

type TradeSymbol struct {
	Symbol                     string   `json:"symbol"`
	Status                     string   `json:"status"`
	BaseAsset                  string   `json:"baseAsset"`
	BaseAssetPrecision         int      `json:"baseAssetPrecision"`
	PricePrecision             int      `json:"pricePrecision"`
	QuantityPrecision          int      `json:"quantityPrecision"`
	QuoteAsset                 string   `json:"quoteAsset"`
	QuotePrecision             int      `json:"quotePrecision"`
	BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
	Filters                    []Filter `json:"filters"`
	IcebergAllowed             bool     `json:"icebergAllowed"`
	IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
	IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
	OcoAllowed                 bool     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
	OrderTypes                 []string `json:"orderTypes"`
}

func (ts TradeSymbol) GetMinAmount() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "LOT_SIZE" {
			return v.MinQty
		}
	}
	return 0
}

func (ts TradeSymbol) GetAmountPrecision() int {
	for _, v := range ts.Filters {
		if v.FilterType == "LOT_SIZE" {
			step := strconv.FormatFloat(v.StepSize, 'f', -1, 64)
			pres := strings.Split(step, ".")
			if len(pres) == 1 {
				return 0
			}
			return len(pres[1])
		}
	}
	return 0
}

func (ts TradeSymbol) GetMinPrice() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "PRICE_FILTER" {
			return v.MinPrice
		}
	}
	return 0
}

func (ts TradeSymbol) GetMinValue() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "MIN_NOTIONAL" {
			return v.MinNotional
		}
	}
	return 0
}

func (ts TradeSymbol) GetPricePrecision() int {
	for _, v := range ts.Filters {
		if v.FilterType == "PRICE_FILTER" {
			step := strconv.FormatFloat(v.TickSize, 'f', -1, 64)
			pres := strings.Split(step, ".")
			if len(pres) == 1 {
				return 0
			}
			return len(pres[1])
		}
	}
	return 0
}

type ExchangeInfo struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int           `json:"serverTime"`
	ExchangeFilters []interface{} `json:"exchangeFilters,omitempty"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	Symbols         []TradeSymbol `json:"symbols"`
}

type Binance struct {
	accessKey  string
	secretKey  string
	baseUrl    string
	apiV1      string
	apiV3      string
	wapiV3     string
	sapiV1     string
	fapiV1     string
	dapiV1     string
	bapiv1     string
	clientType string
	httpClient *http.Client
	timeOffset int64 //nanosecond
	*ExchangeInfo
}

func (bn *Binance) buildParamsSigned(postForm *url.Values) error {
	postForm.Set("recvWindow", "60000")
	tonce := strconv.FormatInt(time.Now().UnixNano()+bn.timeOffset, 10)[0:13]
	postForm.Set("timestamp", tonce)
	payload := postForm.Encode()
	sign, _ := GetParamHmacSHA256Sign(bn.secretKey, payload)
	postForm.Set("signature", sign)
	return nil
}

func New(client *http.Client, api_key, secret_key string) *Binance {
	return NewWithConfig(&APIConfig{
		HttpClient:   client,
		Endpoint:     GLOBAL_API_BASE_URL,
		ApiKey:       api_key,
		ApiSecretKey: secret_key})
}

func NewWithConfig(config *APIConfig) *Binance {
	if config.Endpoint == "" {
		config.Endpoint = GLOBAL_API_BASE_URL
	}
	bn := &Binance{
		baseUrl:    config.Endpoint,
		apiV1:      config.Endpoint + "/api/v1/",
		apiV3:      config.Endpoint + "/api/v3/",
		sapiV1:     config.Endpoint + "/sapi/v1/",
		wapiV3:     config.Endpoint + "/wapi/v3/",
		fapiV1:     F_API_BASE_URL + "/fapi/v1/",
		bapiv1:     WWW_API_BASE_URL + "/bapi/",
		accessKey:  config.ApiKey,
		secretKey:  config.ApiSecretKey,
		httpClient: config.HttpClient}
	bn.setTimeOffset()
	return bn
}

func (bn *Binance) GetExchangeName() string {
	return BINANCE
}

func (bn *Binance) Ping() bool {
	_, err := HttpGet(bn.httpClient, bn.apiV3+"ping")
	if err != nil {
		return false
	}
	return true
}

func (bn *Binance) setTimeOffset() error {
	respmap, err := HttpGet(bn.httpClient, bn.apiV3+SERVER_TIME_URL)
	if err != nil {
		return err
	}

	stime := int64(ToInt(respmap["serverTime"]))
	st := time.Unix(stime/1000, 1000000*(stime%1000))
	lt := time.Now()
	offset := st.Sub(lt).Nanoseconds()
	bn.timeOffset = int64(offset)
	return nil
}

func (bn *Binance) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := bn.apiV3 + fmt.Sprintf(TICKER_URI, currency.ToSymbol(""))
	tickerMap, err := HttpGet(bn.httpClient, tickerUri)

	if err != nil {
		return nil, err
	}

	var ticker Ticker
	ticker.Pair = currency
	t, _ := tickerMap["closeTime"].(float64)
	ticker.Date = uint64(t / 1000)
	ticker.Open = ToFloat64(tickerMap["openPrice"])
	ticker.Last = ToFloat64(tickerMap["lastPrice"])
	ticker.Buy = ToFloat64(tickerMap["bidPrice"])
	ticker.Sell = ToFloat64(tickerMap["askPrice"])
	ticker.Low = ToFloat64(tickerMap["lowPrice"])
	ticker.High = ToFloat64(tickerMap["highPrice"])
	ticker.Vol = ToFloat64(tickerMap["volume"])
	return &ticker, nil
}

func (bn *Binance) GetTickers() ([]*Ticker, error) {
	tickerUri := bn.apiV3 + TICKER_ALL_URI
	var tickersArray []map[string]interface{}
	err := HttpGet4(bn.httpClient, tickerUri, nil, &tickersArray)

	if err != nil {
		return nil, err
	}
	var tickers []*Ticker

	for _, tickerMap := range tickersArray {
		ticker := &Ticker{}

		if symbol, ok := tickerMap["symbol"]; ok {
			ticker.Symbol = symbol.(string)
		}
		t, _ := tickerMap["closeTime"].(float64)
		ticker.Date = uint64(t / 1000)
		ticker.Open = ToFloat64(tickerMap["openPrice"])
		ticker.Last = ToFloat64(tickerMap["lastPrice"])
		ticker.Buy = ToFloat64(tickerMap["bidPrice"])
		ticker.Sell = ToFloat64(tickerMap["askPrice"])
		ticker.Low = ToFloat64(tickerMap["lowPrice"])
		ticker.High = ToFloat64(tickerMap["highPrice"])
		ticker.Vol = ToFloat64(tickerMap["volume"])
		tickers = append(tickers, ticker)
	}

	return tickers, nil
}

func (bn *Binance) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	if size <= 5 {
		size = 5
	} else if size <= 10 {
		size = 10
	} else if size <= 20 {
		size = 20
	} else if size <= 50 {
		size = 50
	} else if size <= 100 {
		size = 100
	} else if size <= 500 {
		size = 500
	} else {
		size = 1000
	}

	apiUrl := fmt.Sprintf(bn.apiV3+DEPTH_URI, currencyPair.ToSymbol(""), size)
	resp, err := HttpGet(bn.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	if _, isok := resp["code"]; isok {
		return nil, errors.New(resp["msg"].(string))
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	depth := new(Depth)
	depth.Pair = currencyPair
	depth.UTime = time.Now()
	n := 0
	for _, bid := range bids {
		_bid := bid.([]interface{})
		amount := ToFloat64(_bid[1])
		price := ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
		n++
		if n == size {
			break
		}
	}

	n = 0
	for _, ask := range asks {
		_ask := ask.([]interface{})
		amount := ToFloat64(_ask[1])
		price := ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
		n++
		if n == size {
			break
		}
	}

	return depth, nil
}

func (bn *Binance) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	path := bn.apiV3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("side", orderSide)
	params.Set("type", orderType)
	params.Set("newOrderRespType", "ACK")
	params.Set("quantity", amount)

	switch orderType {
	case "LIMIT":
		params.Set("timeInForce", "GTC")
		params.Set("price", price)
	case "MARKET":
		params.Set("newOrderRespType", "RESULT")
	}

	bn.buildParamsSigned(&params)

	resp, err := HttpPostForm2(bn.httpClient, path, params,
		map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return nil, err
	}

	orderId := ToInt(respmap["orderId"])
	if orderId <= 0 {
		return nil, errors.New(string(resp))
	}

	side := BUY
	if orderSide == "SELL" {
		side = SELL
	}

	dealAmount := ToFloat64(respmap["executedQty"])
	cummulativeQuoteQty := ToFloat64(respmap["cummulativeQuoteQty"])
	avgPrice := 0.0
	if cummulativeQuoteQty > 0 && dealAmount > 0 {
		avgPrice = cummulativeQuoteQty / dealAmount
	}

	return &Order{
		Currency:   pair,
		OrderID:    orderId,
		OrderID2:   strconv.Itoa(orderId),
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: dealAmount,
		AvgPrice:   avgPrice,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  ToInt(respmap["transactTime"])}, nil
}

func (bn *Binance) GetAccount() (*Account, error) {
	params := url.Values{}
	bn.buildParamsSigned(&params)
	path := bn.apiV3 + ACCOUNT_URI + params.Encode()
	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}
	if _, isok := respmap["code"]; isok == true {
		return nil, errors.New(respmap["msg"].(string))
	}
	acc := Account{}
	acc.Exchange = bn.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	balances := respmap["balances"].([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["asset"].(string), "").AdaptBccToBch()
		balance, _ := decimal.NewFromFloat(ToFloat64(vv["free"])).Add(decimal.NewFromFloat(ToFloat64(vv["locked"]))).Float64()
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["free"]),
			ForzenAmount: ToFloat64(vv["locked"]),
			Balance:      balance,
		}
	}

	return &acc, nil
}

func (bn *Binance) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "BUY")
}

func (bn *Binance) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "SELL")
}

func (bn *Binance) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "BUY")
}

func (bn *Binance) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "SELL")
}

func (bn *Binance) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	path := bn.apiV3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)

	resp, err := HttpDeleteForm(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})

	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	orderIdCanceled := ToInt(respmap["orderId"])
	if orderIdCanceled <= 0 {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (bn *Binance) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	if orderId != "" {
		params.Set("orderId", orderId)
	}
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + ORDER_URI + "?" + params.Encode()

	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	status := respmap["status"].(string)
	side := respmap["side"].(string)

	ord := Order{}
	ord.Currency = currencyPair
	ord.OrderID = ToInt(orderId)
	ord.OrderID2 = orderId
	ord.Cid, _ = respmap["clientOrderId"].(string)
	ord.Type = respmap["type"].(string)

	if side == "SELL" {
		ord.Side = SELL
	} else {
		ord.Side = BUY
	}

	switch status {
	case "NEW":
		ord.Status = ORDER_UNFINISH
	case "FILLED":
		ord.Status = ORDER_FINISH
	case "PARTIALLY_FILLED":
		ord.Status = ORDER_PART_FINISH
	case "CANCELED":
		ord.Status = ORDER_CANCEL
	case "PENDING_CANCEL":
		ord.Status = ORDER_CANCEL_ING
	case "REJECTED":
		ord.Status = ORDER_REJECT
	}

	ord.Amount = ToFloat64(respmap["origQty"].(string))
	ord.Price = ToFloat64(respmap["price"].(string))
	ord.DealAmount = ToFloat64(respmap["executedQty"])
	ord.AvgPrice = ord.Price // response no avg price ， fill price
	ord.OrderTime = ToInt(respmap["time"])

	cummulativeQuoteQty := ToFloat64(respmap["cummulativeQuoteQty"])
	if cummulativeQuoteQty > 0 {
		ord.AvgPrice = cummulativeQuoteQty / ord.DealAmount
	}

	return &ord, nil
}

func (bn *Binance) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + UNFINISHED_ORDERS_INFO + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		side := ord["side"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}
		ordId := ToInt(ord["orderId"])
		orders = append(orders, Order{
			OrderID:   ordId,
			OrderID2:  strconv.Itoa(ordId),
			Currency:  currencyPair,
			Price:     ToFloat64(ord["price"]),
			Amount:    ToFloat64(ord["origQty"]),
			Side:      TradeSide(orderSide),
			Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["time"])})
	}
	return orders, nil
}

func (bn *Binance) GetAllUnfinishOrders() ([]Order, error) {
	params := url.Values{}

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + UNFINISHED_ORDERS_INFO + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		side := ord["side"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}

		ordId := ToInt(ord["orderId"])
		orders = append(orders, Order{
			OrderID:   ToInt(ord["orderId"]),
			OrderID2:  strconv.Itoa(ordId),
			Currency:  bn.toCurrencyPair(ord["symbol"].(string)),
			Price:     ToFloat64(ord["price"]),
			Amount:    ToFloat64(ord["origQty"]),
			Side:      TradeSide(orderSide),
			Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["time"])})
	}
	return orders, nil
}

func (bn *Binance) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]*Kline, error) {
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))
	params.Set("interval", _INERNAL_KLINE_PERIOD_CONVERTER[period])
	if since > 0 {
		params.Set("startTime", strconv.Itoa(since))
	}
	//params.Set("endTime", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("limit", fmt.Sprintf("%d", size))

	klineUrl := bn.apiV3 + KLINE_URI + "?" + params.Encode()

	klines, err := HttpGet3(bn.httpClient, klineUrl, nil)
	if err != nil {
		return nil, err
	}
	var klineRecords []*Kline

	for _, _record := range klines {
		r := &Kline{Pair: currency}
		record := _record.([]interface{})
		r.Timestamp = int64(record[0].(float64)) / 1000 //to unix timestramp
		r.Open = ToFloat64(record[1])
		r.High = ToFloat64(record[2])
		r.Low = ToFloat64(record[3])
		r.Close = ToFloat64(record[4])
		r.Vol = ToFloat64(record[5])

		klineRecords = append(klineRecords, r)
	}

	return klineRecords, nil

}

// 非个人，整个交易所的交易记录
// 注意：since is fromId
func (bn *Binance) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	param := url.Values{}
	param.Set("symbol", currencyPair.ToSymbol(""))
	param.Set("limit", "500")
	if since > 0 {
		param.Set("fromId", strconv.Itoa(int(since)))
	}
	apiUrl := bn.apiV3 + "historicalTrades?" + param.Encode()
	resp, err := HttpGet3(bn.httpClient, apiUrl, map[string]string{
		"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, v := range resp {
		m := v.(map[string]interface{})
		ty := SELL
		if m["isBuyerMaker"].(bool) {
			ty = BUY
		}
		trades = append(trades, Trade{
			Tid:    ToInt64(m["id"]),
			Type:   ty,
			Amount: ToFloat64(m["qty"]),
			Price:  ToFloat64(m["price"]),
			Date:   ToInt64(m["time"]),
			Pair:   currencyPair,
		})
	}

	return trades, nil
}

func (bn *Binance) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + "allOrders?" + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		log.Printf("%+v", ord)
		side := ord["side"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}
		ordId := ToInt(ord["orderId"])
		orders = append(orders, Order{
			OrderID:   ToInt(ord["orderId"]),
			OrderID2:  strconv.Itoa(ordId),
			Currency:  currency,
			Price:     ToFloat64(ord["price"]),
			Amount:    ToFloat64(ord["origQty"]),
			Side:      TradeSide(orderSide),
			Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["time"])})
	}
	return orders, nil

}

func (bn *Binance) GetUserDataStream() (string, error) {
	params := url.Values{}
	_ = bn.buildParamsSigned(&params)
	path := bn.apiV3 + "userDataStream"
	log.Println(path)
	respmap, err := HttpPostForm2(bn.httpClient, path, nil, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return "", err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(respmap, &resp)
	if _, isok := resp["code"]; isok == true {
		return "", errors.New(resp["msg"].(string))
	}
	return resp["listenKey"].(string), nil
}

func (bn *Binance) PutUserDataStream(listenKey string) (bool, error) {
	path := bn.apiV3 + "userDataStream?" + "listenKey=" + listenKey
	respmap, err := HttpPutData(bn.httpClient, path, "", map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return false, err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(respmap, &resp)
	if _, isok := resp["code"]; isok == true {
		return false, errors.New(resp["msg"].(string))
	}
	return true, nil
}

func (bn *Binance) toCurrencyPair(symbol string) CurrencyPair {
	if bn.ExchangeInfo == nil {
		var err error
		bn.ExchangeInfo, err = bn.GetExchangeInfo()
		if err != nil {
			return CurrencyPair{}
		}
	}
	for _, v := range bn.ExchangeInfo.Symbols {
		if v.Symbol == symbol {
			return NewCurrencyPair2(v.BaseAsset + "_" + v.QuoteAsset)
		}
	}
	return CurrencyPair{}
}

func (bn *Binance) GetExchangeInfo() (*ExchangeInfo, error) {
	resp, err := HttpGet5(bn.httpClient, bn.fapiV1+"exchangeInfo", nil)
	if err != nil {
		return nil, err
	}
	info := &ExchangeInfo{}
	err = json.Unmarshal(resp, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (bn *Binance) GetTradeSymbol(currencyPair CurrencyPair) (*TradeSymbol, error) {
	if bn.ExchangeInfo == nil {
		var err error
		bn.ExchangeInfo, err = bn.GetExchangeInfo()
		if err != nil {
			return nil, err
		}
	}
	for k, v := range bn.ExchangeInfo.Symbols {
		if v.Symbol == currencyPair.ToSymbol("") {
			return &bn.ExchangeInfo.Symbols[k], nil
		}
	}
	return nil, errors.New("symbol not found")
}

func (bn *Binance) CancelAllOrders(currencyPair CurrencyPair) (bool, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	_ = bn.buildParamsSigned(&params)

	path := bn.apiV3 + "openOrders"
	respmap, err := HttpDeleteForm(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return false, err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(respmap, &resp)
	if _, isok := resp["code"]; isok == true {
		return false, errors.New(resp["msg"].(string))
	}
	log.Printf("%+v", resp)
	return true, nil
}

func (bn *Binance) CreateSubAccount() (subAccountId string, err error) {
	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	_ = bn.buildParamsSigned(&params)
	path := bn.sapiV1 + "broker/subAccount"
	form, err := HttpPostForm2(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return "", err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(form, &resp)
	if _, isok := resp["code"]; isok == true {
		return "false", errors.New(resp["msg"].(string))
	}
	log.Printf("%+v", resp)
	return resp["subaccountId"].(string), nil
}

func (bn *Binance) EnableSubAccountMargin(subAccount string) (err error) {
	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	params.Set("subAccountId", subAccount)
	params.Set("margin", "true")
	_ = bn.buildParamsSigned(&params)
	path := bn.sapiV1 + "broker/subAccount/margin"
	form, err := HttpPostForm2(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(form, &resp)
	if _, isok := resp["code"]; isok == true {
		return errors.New(resp["msg"].(string))
	}
	log.Printf("%+v", resp)
	return nil
}

type TransFerResp struct {
	txnId        string
	tranId       string
	clientTranId string
}

func (bn *Binance) SubAccountTransfer(fromId, toId, clientTranId, asset, amount string) (data TransFerResp, err error) {
	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	params.Set("fromId", fromId)
	params.Set("toId", toId)
	params.Set("clientTranId", clientTranId)
	params.Set("asset", asset)
	params.Set("amount", amount)
	_ = bn.buildParamsSigned(&params)
	path := bn.sapiV1 + "broker/transfer"
	form, err := HttpPostForm2(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return TransFerResp{}, err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(form, &resp)
	if _, isok := resp["code"]; isok == true {
		return TransFerResp{}, errors.New(resp["msg"].(string))
	}
	txnID := decimal.NewFromFloat(resp["txnId"].(float64)).String()
	return TransFerResp{txnId: txnID, clientTranId: resp["clientTranId"].(string)}, nil
}

// 万向划转
func (bn *Binance) UserUniversalTransfer(Type AssetTransferType, asset, amount string) (data TransFerResp, err error) {
	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	params.Set("type", Type.String())
	params.Set("asset", asset)
	params.Set("amount", amount)
	_ = bn.buildParamsSigned(&params)
	path := bn.sapiV1 + "asset/transfer"
	form, err := HttpPostForm2(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return TransFerResp{}, err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(form, &resp)
	if _, isok := resp["code"]; isok == true {
		return TransFerResp{}, errors.New(resp["msg"].(string))
	}
	tranId := decimal.NewFromFloat(resp["tranId"].(float64)).String()
	return TransFerResp{tranId: tranId}, nil
}

type DepositAddrResp struct {
	Address string
	Coin    string
	Tag     string
	Url     string
}

func (bn *Binance) GetSubAccountDepositAddress(email, symbol string) (resp DepositAddrResp, err error) {
	params := url.Values{}
	params.Set("email", email)
	params.Set("coin", symbol)
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	_ = bn.buildParamsSigned(&params)
	path := bn.sapiV1 + "capital/deposit/subAddress?" + params.Encode()
	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return DepositAddrResp{}, err
	}
	return DepositAddrResp{Address: respmap["address"].(string), Coin: respmap["coin"].(string), Tag: respmap["tag"].(string),
		Url: respmap["url"].(string),
	}, nil
}

func (bn *Binance) GetAccountDepositAddress(symbol string) (resp DepositAddrResp, err error) {
	params := url.Values{}
	params.Set("coin", symbol)
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	_ = bn.buildParamsSigned(&params)
	path := bn.sapiV1 + "capital/deposit/address?" + params.Encode()
	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return DepositAddrResp{}, err
	}
	return DepositAddrResp{Address: respmap["address"].(string), Coin: respmap["coin"].(string), Tag: respmap["tag"].(string),
		Url: respmap["url"].(string),
	}, nil
}

//
//{
//"subaccountId": "1",
//"apikey":"vmPUZE6mv9SD5VNHk4HlWFsOr6aKE2zvsw0MuIgwCIPy6utIco14y7Ju91duEh8A",
//"secretkey":"NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0",
//"canTrade": true,
//"marginTrade": false,
//"futuresTrade": false
//}

type SubAccountApiResp struct {
	SubAccountID string
	ApiKey       string
	SecretKey    string
	CanTrade     bool
	MarginTrade  bool
	FuturesTrade bool
}

func (bn *Binance) CreateSubAccountApi(subAccountId, canTrade string) (data SubAccountApiResp, err error) {
	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	params.Set("subAccountId", subAccountId)
	params.Set("canTrade", canTrade)
	_ = bn.buildParamsSigned(&params)
	path := bn.sapiV1 + "broker/subAccountApi"
	form, err := HttpPostForm2(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return SubAccountApiResp{}, err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(form, &resp)
	if _, isok := resp["code"]; isok == true {
		return SubAccountApiResp{}, errors.New(resp["msg"].(string))
	}
	log.Printf("%+v", resp)

	return SubAccountApiResp{
		SubAccountID: resp["subaccountId"].(string),
		ApiKey:       resp["apiKey"].(string),
		SecretKey:    resp["secretKey"].(string),
		CanTrade:     resp["canTrade"].(bool),
		MarginTrade:  resp["marginTrade"].(bool),
		FuturesTrade: resp["futuresTrade"].(bool),
	}, nil
}

// ListDownloadData
//
//	@Description: 下载币安K线历史数据
//	@receiver bn
//	@param bizType	下载类型 FUTURES_UM-U本位	FUTURES_CM-币本位	SPOT-现货
//	@param productName	类型	trades-最新成交
//	@param interval	monthly月 | daily日
//	@param startDay	开始时间 2023-01-01
//	@param endDay	结束时间 2023-01-02
//	@param granularityList	周期时间 [1s,1m]
//	@param symbolList	币种(多个)["BNBUSDT"]
//	@return data
//	@return err
func (bn *Binance) ListDownloadData(bizType, productName, interval, startDay, endDay string, granularityList, symbolList []string) (data DownloadItemLists, err error) {
	params := make(map[string]interface{})
	params["bizType"] = bizType
	params["productName"] = productName
	params["interval"] = interval
	params["startDay"] = startDay
	params["endDay"] = endDay
	params["granularityList"] = granularityList
	params["symbolList"] = symbolList
	path := bn.bapiv1 + "bigdata/v1/public/bigdata/finance/exchange/listDownloadData"
	form, err := HttpPostForm4(bn.httpClient, path, params, map[string]string{})
	if err != nil {
		return data, err
	}
	var resp HttpResult
	_ = jsoniter.Unmarshal(form, &resp)
	resp.Data.DownloadItemList.Sort("asc")
	return resp.Data.DownloadItemList, nil
}
func (bn *Binance) Brackets() (data []BracketsList) {
	_type := []string{"delivery", "future"}
	for _, s := range _type {
		params := make(map[string]interface{})
		var resp *HttpResult
		path := bn.bapiv1 + fmt.Sprintf("futures/v1/friendly/%s/common/brackets", s)
		// futures/v1/friendly/delivery/common/brackets
		form, err := HttpPostForm4(bn.httpClient, path, params, map[string]string{})
		if err != nil {
			return data
		}
		err = json.Unmarshal(form, &resp)
		data = append(data, resp.Data.Brackets...)
	}

	return data
}

// DownloadData 下载K线数据并整合到K线数组
func DownloadData(bizType, productName, interval, startDay, endDay string, granularityList, symbolList []string) (k []string) {
	bnHttpWith := NewWithConfig(&APIConfig{HttpClient: &http.Client{
		Timeout: 120 * time.Second,
	}})
	var wg sync.WaitGroup
	da, err := bnHttpWith.ListDownloadData(bizType, productName, interval, startDay, endDay, granularityList, symbolList)
	if err != nil {
		log.Fatalf("数据下载失败%v", err.Error())
	}
	var mu = new(sync.Mutex)
	pool := tunny.NewFunc(10, func(i interface{}) interface{} {
		var item = da[i.(int)]
		rootp, _ := os.UserHomeDir() // helper.GetCurrentPath()
		path := fmt.Sprintf("%s/.binanceData/resource/%s/%s/%s/%s", rootp, item.ProductName, item.BizType, item.Interval, item.Symbol)
		downFile := fmt.Sprintf("%s/%s", path, item.Filename)
		unZipFile := fmt.Sprintf("%s/%s.csv", path, helper.GetFileName(item.Filename, false))
		helper.Exists(path, true, true) //检测对应目录是否存在
		// 检测解压后文件是否存在
		if !helper.Exists(unZipFile, false, false) {
			if !helper.Exists(downFile, false, false) {
				helper.Download(item.Url, downFile)
			}
			_, _ = helper.Unzip(downFile, path)
			_ = os.Remove(downFile)
		}
		mu.Lock()
		k = append(k, unZipFile)
		mu.Unlock()
		wg.Done()
		return nil
	})
	defer pool.Close()
	for i := 0; i < len(da); i++ {
		wg.Add(1)
		go pool.Process(i)
	}
	wg.Wait()
	return k
}

func GetKLines(bizType, productName, interval, startDay, endDay string, granularityList, symbolList []string, pair CurrencyPair) (klins []*Kline) {
	kData := DownloadData(bizType, productName, interval, startDay, endDay, granularityList, symbolList)
	var wg sync.WaitGroup
	var mu sync.Mutex
	pool := tunny.NewFunc(runtime.NumCPU()*10, func(i interface{}) interface{} {
		k := LocalKlineCsv(kData[i.(int)], pair)
		mu.Lock()
		klins = append(klins, k...)
		mu.Unlock()
		wg.Done()
		return nil
	})
	defer pool.Close()
	for i := 0; i < len(kData); i++ {
		wg.Add(1)
		go pool.Process(i)
		//klins = append(klins, LocalKlineCsv(kData[i], pair)...)
	}

	wg.Wait()
	return klins
}
func GetTrade(bizType, productName, interval, startDay, endDay string, granularityList, symbolList []string, pair CurrencyPair) (trades []*Trade) {
	kData := DownloadData(bizType, productName, interval, startDay, endDay, granularityList, symbolList)
	var wg sync.WaitGroup
	var mu sync.Mutex
	pool1 := tunny.NewFunc(runtime.NumCPU(), func(i interface{}) interface{} {
		t := LocalTradeCsvSpili(kData[i.(int)], pair)
		mu.Lock()
		trades = append(trades, t...)
		//log.Println("交易记录数量:", len(kData), "已完成数量:", len(t), "时间:", time.UnixMilli(t[0].Date).Format(helper.TimeFormatYmd))

		mu.Unlock()
		wg.Done()
		return nil
	})
	for i := 0; i < len(kData); i++ {
		wg.Add(1)
		go pool1.Process(i)
	}

	defer func() {
		pool1.Close()
	}()

	wg.Wait()
	return trades
}
