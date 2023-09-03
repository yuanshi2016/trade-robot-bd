/**
 * @Notes:
 * @class mockOrder_test
 * @package
 * @author: 原始
 * @Time: 2023/6/13   02:20
 */
package goex

import (
	"log"
	"net/http"
	"time"
	"trade-robot-bd/app/grid-strategy-svc/util/goex"
	"trade-robot-bd/libs/goex/binance"
)

func mains() {
	bnHttpWith := binance.NewWithConfig(&APIConfig{
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ClientType: "f",
	})
	Brackets := LoadBrackets(bnHttpWith)
	spot := &MockOrder{
		Direction: BUY,
		Type:      NewOrder_SPOT,
		Lever:     1,
		Quantity:  230,
		FeeRate:   0.075,
		CpUsd:     1,
	}
	spot.Buy(230)
	spot.Sell(231)
	log.Printf("现货测试%#v\r\n", spot.CalcCpUp())
	cm := &MockOrder{
		Direction: BUY,
		Type:      NewOrder_CM,
		Lever:     5,
		Quantity:  50,
		FeeRate:   0.05,
		CpUsd:     10,
	}
	cm.Buy(250.6)
	cm.Sell(251.6)
	log.Printf("币本位测试%#v\r\n", cm.CalcCpUp())
	um := &MockOrder{
		Direction: BUY,
		Type:      NewOrder_UM,
		Lever:     5,
		Quantity:  10000,
		FeeRate:   0.04,
	}
	um.Buy(230)
	um.CalcLiquidation(300, Brackets[goex.BNB_USDT.ToSymbol("")])
	um.Sell(231)
	log.Fatalf("U本位测试%#v\r\n", um.CalcCpUp())
}
