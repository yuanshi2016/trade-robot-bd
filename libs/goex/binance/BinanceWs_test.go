package binance

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
	"trade-robot-bd/libs/goex"
	"unsafe"
)

var bnWs = NewBinanceWs(&goex.APIConfig{})

func init() {
	//bnWs.proxyUrl = "socks5://10.10.1.3:12528"
	bnWs.SetBaseUrl("wss://dstream.binance.com/ws")
	//bnWs.SetCombinedBaseURL("wss://fstream.binancezh.com/stream?streams=boz16htWU3gjF7hkoRUG3EZponFxRpbQ3XDDUywmKO80krZAzoY13tzAuKIQ")
	//bnWs.SetCallbacks(printfTicker, printfDepth, printfTrade, printfKline)
}

func printfTicker(ticker *goex.Ticker) {
	log.Println("ticker:", ticker)
}

func printfDepth(depth *goex.Depth) {
	log.Println("depth:", depth)
}

func printfTrade(trade *goex.Trade) {
	log.Println("trade:", trade)
	log.Println("trade:", (*RawTrade)(unsafe.Pointer(trade)))
}

func printfAggTrade(aggTrade *goex.Trade) {
	log.Println("trade:", (*AggTrade)(unsafe.Pointer(aggTrade)))
}
func printfKline(kline *goex.Kline, period int) {
	log.Println("kline:", kline)
}

// go test -run TestBinanceWs_SubscribeTicker
func TestBinanceWs_SubscribeTicker(t *testing.T) {
	err := bnWs.SubscribeTicker(goex.BNB_USD)
	fmt.Println(err)
	time.Sleep(time.Second * 5)
}

func TestBinanceWs_GetDepthWithWs(t *testing.T) {
	bnWs.SubscribeDepth(goex.BTC_USDT, 5)
	time.Sleep(time.Second * 10)
}

func TestBinanceWs_GetKLineWithWs(t *testing.T) {
	return
	bnWs.SubscribeKline(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN)
	time.Sleep(time.Second * 10)
}

func TestBinanceWs_GetTradesWithWs(t *testing.T) {
	bnWs.SubscribeTrade(goex.BTC_USDT)
	time.Sleep(time.Second * 5)
}

// go test -run TestBinanceWs_SubscribeAggTrade

func TestBinanceWs_SubscribeAggTrade(t *testing.T) {
	var bg = new(sync.WaitGroup)
	bg.Add(1)
	err := bnWs.SubscribeAggTrade(goex.BNB_USD, printfAggTrade)
	if err != nil {
		log.Fatal("运行失败", err)
	}
	bg.Wait()
	time.Sleep(time.Second * 5)
}

func TestBinanceWs_SubscribeDiffDepth(t *testing.T) {
	bnWs.SubscribeDiffDepth(goex.BTC_USDT, printfDepth)
	time.Sleep(time.Second * 10)
}

func TestBinanceWs_SubscribeDepth(t *testing.T) {
	bnWs.SubscribeDepth(goex.BTC_USDT, 5)
	bnWs.SubscribeDepth(goex.LTC_USDT, 5)
	bnWs.SubscribeDepth(goex.ETC_USDT, 5)
	time.Sleep(time.Second * 60)
}
