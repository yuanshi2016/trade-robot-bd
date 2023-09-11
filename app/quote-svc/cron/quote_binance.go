package cron

import (
	"fmt"
	"strings"
	"time"
	"trade-robot-bd/libs/cache"
	"trade-robot-bd/libs/exchangeclient"
	"trade-robot-bd/libs/goex"
	"trade-robot-bd/libs/logger"
)

var TickBinanceAll = "tick:binance:all"
var BinanceTickArrayAll = make([]Ticker, 0)
var BinanceTickMapAll = make(map[string]interface{}) //all保存所有品种
var BinaceTickArrayBtc = make([]Ticker, 0)
var BinanceKlineAll Klines

func StoreBinanceTick() {
	d := time.Second * 8
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		storeBinanceTick()
		storeBinanceKline()
		if len(BinanceTickMapAll) == 0 {
			continue
		}
		if err := cache.Redis().HMSet(TickBinanceAll, BinanceTickMapAll).Err(); err != nil {
			logger.Errorf("binance将行情存到redis失败 %v", err)
		}
	}
}
func storeBinanceKline() {
	client := exchangeclient.InitBinance("", "")
	_kline, err := client.ApiClient.GetKlineRecords(goex.BNB_USDT, goex.KLINE_PERIOD_1MIN, 1, 0)
	if err != nil {
		logger.Infof("storeBinanceTick GetTickers has err %v", err)
		return
	}
	var klinesData []Kline
	for _, item := range _kline {
		klinesData = append(klinesData, Kline{
			Open:      item.Open,
			Close:     item.Close,
			High:      item.High,
			Low:       item.Low,
			Vol:       item.Vol,
			CloseTime: uint64(item.Timestamp),
			QuoteVol:  item.Vol,
		})
	}
	BinanceKlineAll = *new(Klines)
	BinanceKlineAll = Klines{
		Type:   "kline",
		Symbol: goex.BNB_USDT.String(),
		Fin:    0,
		Data:   klinesData[0],
	}
}
func storeBinanceTick() {
	client := exchangeclient.InitBinance("", "")
	tickers, err := client.ApiClient.GetTickers()
	if err != nil {
		logger.Infof("storeBinanceTick GetTickers has err %v", err)
		return
	}
	//重置为空
	BinanceTickArrayAll = BinanceTickArrayAll[:0]
	BinaceTickArrayBtc = BinaceTickArrayBtc[:0]
	for _, v := range tickers {
		if v.Open <= 0 {
			continue
		}
		change := (v.Last - v.Open) / v.Open * 100
		tick := Ticker{
			Symbol: v.Symbol,
			Last:   v.Last,
			Buy:    v.Buy,
			Open:   v.Open,
			Sell:   v.Sell,
			High:   v.High,
			Low:    v.Low,
			Vol:    v.Vol,
			Change: fmt.Sprintf("%.2f%v", change, "%"),
			Date:   v.Date,
		}
		if change > 0.0 {
			tick.Change = "+" + tick.Change
		}
		if strings.HasSuffix(v.Symbol, "USDT") {
			tick.Symbol = strings.ReplaceAll(v.Symbol, "USDT", "-USDT")
			BinanceTickMapAll[tick.Symbol] = tick
			BinanceTickArrayAll = append(BinanceTickArrayAll, tick)
		}
		if len(v.Symbol) >= 4 && strings.HasSuffix(v.Symbol, "BTC") {
			tick.Symbol = strings.ReplaceAll(v.Symbol, "BTC", "-BTC")
			BinaceTickArrayBtc = append(BinaceTickArrayBtc, tick)
		}
	}
	logger.Infof("获取USDT行情数据数量: %d", len(BinanceTickArrayAll))
}
