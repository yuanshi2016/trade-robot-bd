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
	"trade-robot-bd/libs/mockers"
)

var ApiKey = "ppJiDo1sA8jIPtJtD37a7r7slZHWKnpOWgzTJNgDbRA8zbZyxFcLU500uItwXAdZ"
var ApiSecretKey = "hNcNer35AEvnR1WxyxTkNJ8iGAZhVC8oIhVEQsyoo8BAKMj4xggrLhkW2BdPvcVR"
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
	loadAssKey(ApiKey, ApiSecretKey)
}
func loadAssKey(ApiKey, ApiSecretKey string) {
	config.ApiSecretKey = ApiSecretKey
	config.ApiKey = ApiKey
	config.Endpoint = ""
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
	klinData, _ := bnHttp.GetKlineRecords(goex.ETH_USDT, goex.KLINE_PERIOD_5MIN, 167, 0)
	Renko := goex.KlineToRenKo(klinData, 20, goex.InClose, 2, 5, goex.RenKoMoveTypeAMi)
	for i := 0; i < len(Renko); i++ {
		item := Renko[i]
		fmt.Printf("%v. 开:%v 高:%v 	收:%v 低:%v Diff:%v --%v \r\n", i, item.Open, item.High, item.Close, item.Low, item.Vol, time.UnixMilli(item.Timestamp).Format(helper.TimeFormatYmdHis))
	}
}
func rsi() {
	var klinData = binance.GetKLines("FUTURES_UM", "klines", "daily",
		"2023-06-27", "2023-06-31",
		[]string{fmt.Sprintf("%vm", 5)}, []string{goex.ETH_USDT.ToSymbol("")}, goex.ETH_USDT)
	klinData = goex.KlineSort(klinData, "asc")
	c := goex.CalcRsi(klinData, 14)
	log.Fatalln(len(klinData), len(c), c, time.UnixMilli(klinData[len(klinData)-1].CloseTime).Format(helper.TimeFormatYmdHis))
}
func catAccount() {
	loadAssKey("Dc289rn6Os0F2G26950igEQQOYKm3LelvaaSyS081hGEBkYUMNYj3MFJoTOlQtYP", "3GgSS5Vdigtn41TfK3Bp2X27PgXEQesGsDIRw102XwfYW29hY9TGZu4OFjK3bJss")
	Account, err := bnHttp.GetAccount()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Fatalf("%#v", Account.SubAccounts[goex.USDT])
}
func main() {
	var bg = new(sync.WaitGroup)
	bg.Add(1)
	//mains()
	T_Oline()
	//T_Oline_N([]goex.CurrencyPair{goex.ETH_USDT, goex.BNB_USDT, goex.BTC_USDT, goex.LTC_USDT})
	//go helper.ForTime(500*time.Millisecond, looketBalanceT)
	for _, lev := range []int{5} {
		m := mockers.MockCyCle{
			Comm: mockers.Comm{
				Symbol:    goex.ETH_USDT,
				MaxHold:   100000,
				TradeType: mockers.TradeTypeLocal, //测试模型 本地回测
			},
			Cycle:             lev,
			Lever:             helper.MakeMathQuantity(1, 5, 1),
			ProfitRate:        helper.MakeMathQuantity(1, 5.0, 2.0),
			StopLossRate:      helper.MakeMathQuantity(1, -30.0, -2.0),     //止损率
			ProfitType:        []mockers.CloseSignal{mockers.ProfitSignal}, // mockers.ProfitSignal mockers.ProfitRate
			FeeRate:           0.04,
			Usd:               1000,
			IsTowWay:          true,
			TotalBalanceRatio: 1,
			StartDay:          "2023-07-01",
			EndDay:            "2023-07-31",
		}
		m.Bn = bnHttpWith
		m.AtrLength = helper.MakeMathQuantity(1, 18, 1) //-
		/**
		18 - 5:
		2023-6. 456%
		2023-5. 324%
		2023-4. 32%
		2023-1 - 2023 -8 11272%

		18 - 9:
		2023-1 - 2023 -8 11288%
		8月份行情 ATRLen=18 MoveI=9 收益为正
		*/
		m.RenKoMoveI = helper.MakeMathQuantity(1, 5, 1) //-
		m.RenKoMoveType = []goex.RenKoMoveType{goex.RenKoMoveTypeAMM}
		//m.RenKoMoveType = []goex.RenKoMoveType{goex.RenKoMoveTypeAX, goex.RenKoMoveTypeAX, goex.RenKoMoveTypeAMi, goex.RenKoMoveTypeMM, goex.RenKoMoveTypeAMM}
		m.RunCycleWhere(mockers.WhereTypeAll)
	}
	bg.Wait()

}

// T_Oline_N 在线 但非交易
func T_Oline_N(list []goex.CurrencyPair) {
	var wg sync.WaitGroup
	var view []*mockers.WhereCycleOne
	Brackets := mockers.LoadBrackets(bnHttpWith)
	var port int
	var stop float64
	var lever int64
	flag.IntVar(&port, "P", 8088, "端口号,默认为空")
	flag.Int64Var(&lever, "L", 5, "倍数")
	flag.Float64Var(&stop, "S", 5, "止损")
	flag.Parse()
	for _, pair := range list {
		wg.Add(1)
		o := &mockers.WhereCycleOne{
			Comm: mockers.Comm{
				Symbol:    pair,
				Brackets:  Brackets,
				MaxHold:   100000,
				Bn:        bnHttpWith,
				BnSwap:    bnHttp,
				TradeType: mockers.TradeTypeOnlineData, //模拟交易
			},
			ProfitType:    mockers.ProfitSignal,
			StopLossRate:  -stop,
			AtrLength:     18,
			RenKoMoveI:    5,
			IsTowWay:      true,
			RenKoMoveType: goex.RenKoMoveTypeAMM,
			MockDetail: goex.MockDetail{
				Usd:     1000,
				OldUsd:  1000,
				Type:    goex.NewOrder_UM,
				Lever:   lever,
				FeeRate: 0.04,
				CpUsd:   10,
			},
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
func T_Oline() {
	var wg sync.WaitGroup
	var view []*mockers.WhereCycleOne
	Brackets := mockers.LoadBrackets(bnHttpWith)
	var port int
	var stop float64
	var lever int64
	var pair goex.CurrencyPair
	var ASK string
	var AK string
	flag.IntVar(&port, "P", 8083, "端口号,默认为空")
	flag.Int64Var(&lever, "L", 20, "倍数")
	flag.Float64Var(&stop, "S", 60, "止损")
	flag.StringVar(&pair.CurrencyA.Symbol, "SymA", "ltc", "止损")
	flag.StringVar(&pair.CurrencyB.Symbol, "SymB", "usdt", "止损")
	flag.StringVar(&ASK, "ASK", ApiSecretKey, "ApiSecretKey")
	flag.StringVar(&AK, "AK", ApiKey, "ApiKey")
	flag.Parse()
	loadAssKey(AK, ASK)
	pair = pair.ToUpper()
	Account, err := bnHttp.GetFutureUserinfo()
	if err != nil {
		log.Fatalf("账户信息加载失败:%v", err.Error())
	}
	wg.Add(1)
	_, err = bnHttp.SetLeverage(pair, int(lever))
	if err != nil {
		log.Fatalf("[%v]修改杠杆倍率失败:%v", pair.String(), err.Error())
	}
	o := &mockers.WhereCycleOne{
		Comm: mockers.Comm{
			Symbol:    pair,
			Brackets:  Brackets,
			MaxHold:   100000,
			Bn:        bnHttpWith,
			BnSwap:    bnHttp,
			TradeType: mockers.TradeTypeOline, //模拟交易
		},
		ProfitType:        mockers.ProfitSignal,
		StopLossRate:      -stop,
		AtrLength:         18,
		RenKoMoveI:        5,
		IsTowWay:          true,
		TotalBalanceRatio: 0.95,
		RenKoMoveType:     goex.RenKoMoveTypeAMM,
		MockDetail: goex.MockDetail{
			Usd:     Account.FutureSubAccounts[goex.USDT].AccountRights,
			OldUsd:  Account.FutureSubAccounts[goex.USDT].AccountRights,
			Type:    goex.NewOrder_UM,
			Lever:   lever,
			FeeRate: 0.04,
			CpUsd:   10,
		},
	}
	/// 加载目前已有订单
	a, err := bnHttp.GetFuturePosition(o.Symbol)
	if err == nil {
		for _, position := range a {
			if position.SellAmount < 0 {
				o.MockDetail.SellOrderOnline = &goex.Order{
					Price:    position.SellPriceAvg,
					Amount:   position.SellAmount,
					AvgPrice: position.SellPriceAvg,
					Symbol:   o.Symbol.String(),
					Currency: pair,
					Side:     goex.SELL_MARKET,
				}
			}
			if position.BuyAmount > 0 {
				o.MockDetail.BuyOrderOnline = &goex.Order{
					Price:    position.BuyPriceAvg,
					Amount:   position.BuyAmount,
					AvgPrice: position.BuyPriceAvg,
					Symbol:   o.Symbol.String(),
					Currency: pair,
					Side:     goex.BUY_MARKET,
				}
			}
		}
	}
	view = append(view, o)
	c := config
	c.Endpoint = ""
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		// 循环接收ticker的触发事件
		for range ticker.C {
			o.OnLineKline(bnHttp)
		}
	}()

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

// 测试下单
func newOrderTest() {
	//Account, err := bnHttp.GetFutureUserinfo()
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//bnHttp.GetPositionSideDual()
	pair := goex.ETH_USDT
	exinfo, err := bnHttp.GetTradeSymbol(pair)
	amount := decimal.NewFromFloat(10)
	price := decimal.NewFromFloat(1588)
	//log.Fatalf("%#v", exinfo)
	order, err := bnHttp.MarketBuy(amount.Div(price).Round(int32(exinfo.QuantityPrecision)), price, goex.ETH_USDT, false)
	if err != nil {
		log.Printf("下单失败1:%v", err.Error())
	}
	order1, err := bnHttp.MarketSell(amount.Div(price).Round(int32(exinfo.QuantityPrecision)), price, goex.ETH_USDT, true)
	if err != nil {
		log.Printf("下单失败2:%v", err.Error())
	}
	log.Println(order)
	log.Println(order1)
}

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

func mains() {
	bnHttpWith := binance.NewWithConfig(&goex.APIConfig{
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ClientType: "f",
	})
	Brackets := mockers.LoadBrackets(bnHttpWith)
	//spot := &goex.MockOrder{
	//	Direction: goex.BUY,
	//	Type:      goex.NewOrder_SPOT,
	//	Lever:     1,
	//	Quantity:  230,
	//	FeeRate:   0.075,
	//	CpUsd:     1,
	//}
	//spot.Buy(230)
	//spot.Sell(231)
	//log.Printf("现货测试%#v\r\n", spot.CalcCpUp())
	//cm := &goex.MockOrder{
	//	Direction: goex.BUY,
	//	Type:      goex.NewOrder_CM,
	//	Lever:     5,
	//	Quantity:  50,
	//	FeeRate:   0.05,
	//	CpUsd:     10,
	//}
	//cm.Buy(250.6)
	//cm.Sell(251.6)
	//log.Printf("币本位测试%#v\r\n", cm.CalcCpUp())
	//方向:SELL 开仓价:1836.62832162 当前价:1739.76763514 爆仓价:4791.401682700546 模拟收益:5193.831691464065

	um := &goex.MockOrder{
		Direction: goex.SELL,
		Type:      goex.NewOrder_UM,
		Lever:     5,
		Quantity:  1000,
		FeeRate:   0.04,
	}
	um.OpenMarket(1836.62)
	um.CalcLiquidation(5000, Brackets[goex.ETH_USDT.ToSymbol("")])
	um.CloseMarket(1739.76)
	log.Printf("U本位测试%#v\r\n", um.CalcCpUp())
	um1 := &goex.MockOrder{
		Direction: goex.BUY,
		Type:      goex.NewOrder_UM,
		Lever:     5,
		Quantity:  1000,
		FeeRate:   0.04,
	}
	um1.OpenMarket(1739.76)
	um1.CalcLiquidation(5000, Brackets[goex.ETH_USDT.ToSymbol("")])
	um1.CloseMarket(1836.62)
	log.Fatalf("U本位测试1%#v\r\n", um1.CalcCpUp())
}
