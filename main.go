package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"sync"
	"time"
	"trade-robot-bd/libs/goex"
	"trade-robot-bd/libs/goex/binance"
	"trade-robot-bd/libs/helper"
	"trade-robot-bd/libs/mocker"
)

var ApiKey = "HsOBLroCIbCxexEAuC1S45dXEzKq1CwxwPOFfBQq2cPxKur3z7TvJibnQFE164Rg"
var ApiSecretKey = "8JEoF1voA4lmff0cAXouBEPILr9W5V8bvry9b1k88fSQkDMzMz2CeHQTjQL5c0GZ"
var config = goex.APIConfig{
	HttpClient: &http.Client{
		Timeout: 10 * time.Second,
	},
	ClientType:   "f",
	ApiKey:       ApiKey,
	ApiSecretKey: ApiSecretKey,
}
var bnWs = &binance.BinanceWs{}
var bnHttp = &binance.BinanceSwap{}
var bnHttpWith = &binance.Binance{}

func init() {
	bnWs = binance.NewBinanceWs(&config) //Websocket
	config.Endpoint = ""
	bnHttp = binance.NewBinanceSwap(&config) //合约
	config.Endpoint = ""
	bnHttpWith = binance.NewWithConfig(&config) //现货
}
func asit() {
	klinData, _ := bnHttp.GetKlineRecords(goex.BNB_USDT, goex.KLINE_PERIOD_1H, 23, 0)
	//klinData = goex.KlineSort(klinData, "asc")
	c := goex.CalcRvgi(klinData, 10)
	log.Fatalln(len(klinData), len(c), c, klinData[len(klinData)-1].Close, klinData[len(klinData)-2].Close)
}
func atr() {
	klinData, _ := bnHttp.GetKlineRecords(goex.BNB_USDT, goex.KLINE_PERIOD_5MIN, 167, 0)
	Renko := goex.KlineToRenKo(klinData, 11, goex.InClose, 6)

	for i := 0; i < len(Renko); i++ {
		item := Renko[i]
		fmt.Printf("%v. 开:%v 高:%v 	收:%v 低:%v Diff:%v --%v \r\n", i, item.Open, item.High, item.Close, item.Low, item.Vol, time.UnixMilli(item.Timestamp).Format(helper.TimeFormatYmdHis))
	}

	log.Fatalln("")
	//item := openRenko[len(openRenko)-1]
	//log.Fatalf("开:%v 收:%v 高:%v 低:%v Diff:%v --%v \r\n\n", item.Open, item.Close, item.High, item.Low, item.Vol, time.UnixMilli(item.Timestamp).Format(helper.TimeFormatYmdHis))

	//c := goex.CalcAtr(klinData, 20)
	//log.Fatalln(len(klinData), len(c), c, klinData[len(klinData)-1].Close)
}
func rsi() {
	var klinData = binance.GetKLines("FUTURES_UM", "klines", "daily",
		"2023-06-27", "2023-06-31",
		[]string{fmt.Sprintf("%vm", 3)}, []string{goex.BNB_USDT.ToSymbol("")}, goex.BNB_USDT)
	klinData = goex.KlineSort(klinData, "asc")
	c := goex.CalcRsi(klinData, 14)
	log.Fatalln(len(klinData), len(c), c, time.UnixMilli(klinData[len(klinData)-1].CloseTime).Format(helper.TimeFormatYmdHis))
}
func main() {
	//var bg = new(sync.WaitGroup)
	//bg.Add(1)
	//bg.Wait()
	//runtime.GOMAXPROCS(24)
	T_Oline_N([]goex.CurrencyPair{goex.BCH_USDT, goex.MTL_USDT, goex.FIL_USDT,
		goex.ETC_USDT, goex.SOL_USDT, goex.PEPE_USDT, goex.CTSI_USDT, goex.ONT_USDT,
		goex.ETH_USDT, goex.BNB_USDT, goex.BTC_USDT, goex.LTC_USDT})
	//atr()
	//rsi()
	//asit()

	//MockTest()
	//onLineKline()
	//go helper.ForTime(500*time.Millisecond, looketBalanceT)
	for _, lev := range []int{5} {
		m := mocker.MockCyCle{
			Comm:           mocker.Comm{Symbol: goex.CTSI_USDT},
			Cycle:          lev,
			KlinsLikeTrade: make(map[int64][]*goex.Trade),
			Lever:          helper.MakeMathQuantity(1, 5, 1),
			ProfitRate:     helper.MakeMathQuantity(1, 1.0, 1.0),
			StopLossRate:   helper.MakeMathQuantity(1, -5.0, -1.0),
			ProfitType:     []mocker.CloseSignal{mocker.ProfitSignal},
			MaxHold:        100000,
			StartDay:       "2023-01-01",
			EndDay:         "2023-07-31",
		}
		m.Bn = bnHttpWith
		m.AtrLength = helper.MakeMathQuantity(1, 11, 1) //-
		m.MakeCycleWhere(bnHttpWith, mocker.WhereAll)
	}
	//bg.Wait()
}

// T_Oline_N 在线 但非交易
func T_Oline_N(list []goex.CurrencyPair) {
	var wg sync.WaitGroup
	var view []*mocker.WhereCycleOne
	Brackets := mocker.LoadBrackets(bnHttpWith)
	var port int
	var stop float64
	var lever int64
	flag.IntVar(&port, "P", 8081, "端口号,默认为空")
	flag.Int64Var(&lever, "L", 5, "倍数")
	flag.Float64Var(&stop, "S", 5, "止损")
	flag.Parse()
	for _, pair := range list {
		wg.Add(1)
		o := &mocker.WhereCycleOne{
			Comm: mocker.Comm{
				Symbol: pair,
			},
			Brackets:     Brackets,
			MaxHold:      100000,
			AtrLength:    11,
			ProfitType:   mocker.ProfitSignal,
			StopLossRate: -stop,
			OlineType:    1, //模拟交易
			MockDetail: &goex.MockDetail{
				Usd:     1000,
				OldUsd:  1000,
				Type:    goex.NewOrder_UM,
				Lever:   lever,
				FeeRate: 0.04,
				CpUsd:   10,
			},
			Bn: bnHttpWith,
		}
		view = append(view, o)
		c := config
		c.Endpoint = ""
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			// 循环接收ticker的触发事件
			for range ticker.C {
				o.OnLineKline(bnHttp)
			}
		}()
	}
	viewRun := func() {

		//gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.LoadHTMLGlob("resource/html/*")
		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{"data": view})
		})
		_ = r.Run(fmt.Sprintf(":%v", port))
		log.Println("可视化网关已运行")
	}
	go viewRun()
	wg.Wait()
}
func TradesF(maxTime int64, trades *[]*goex.Trade) (r []*goex.Trade) {
	for _, item := range *trades {
		if item.Date > maxTime-9999 && item.Date < maxTime {
			r = append(r, item)
		}
	}
	return r
}

var isLog = false

// 通过监控现货余额，实时转入到币本位
func looketBalanceT() {
	Account, err := bnHttpWith.GetAccount()
	if err != nil {
		log.Println(err.Error())
	}
	this := Account.SubAccounts[goex.BNB]
	if !isLog {
		isLog = true
		log.Printf("[%s]余额监控-转至币本位", this.Currency.String())
	}
	Balance := decimal.NewFromFloat(this.Balance)
	if this.Balance > 0.0001 {
		Id, err := bnHttpWith.UserUniversalTransfer(goex.Transfer_MAIN_CMFUTURE, this.Currency.String(), Balance.String())
		if err != nil {
			log.Println(err.Error())
		}
		log.Printf("[%s]划转ID%v,划转金额:%v", this.Currency.String(), Id, Balance)
	} else {
		//log.Printf("[%s]当前余额:%v", this.Currency.String(), Balance.String())
	}
}

// MockTest 下单数据测试
func MockTest() {
	Brackets := mocker.LoadBrackets(bnHttpWith)
	//spot := &goex.MockOrder{
	//	Direction: goex.NewOrder_Buy,
	//	Type:      goex.NewOrder_SPOT,
	//	Quantity:  3000,
	//	FeeRate:   0.075,
	//}
	//spot.Buy(235.6)
	//spot.Sell(241)
	//fmt.Printf("现货测试%#v\r\n", spot.CalcCpUp())
	//cm := &goex.MockOrder{
	//	Direction: goex.NewOrder_Buy,
	//	Type:      goex.NewOrder_CM,
	//	Lever:     5,
	//	Quantity:  50,
	//	FeeRate:   0.05,
	//	CpUsd:     10,
	//}
	//cm.Buy(250.6)
	//cm.Sell(251.6)
	//fmt.Printf("币本位测试%#v\r\n", cm.CalcCpUp())
	um := &goex.MockOrder{
		Direction: goex.NewOrder_Buy,
		Type:      goex.NewOrder_UM,
		Lever:     5,
		Quantity:  1500,
		FeeRate:   0.05,
	}
	um.Buy(246)
	um.CalcLiquidation(300, Brackets[goex.BNB_USDT.ToSymbol("")])
	um.Sell(246.5)
	log.Fatalf("U本位测试%#v\r\n", um.CalcCpUp())
}
