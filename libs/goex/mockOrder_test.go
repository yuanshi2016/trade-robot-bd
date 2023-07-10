/**
 * @Notes:
 * @class mockOrder_test
 * @package
 * @author: 原始
 * @Time: 2023/6/13   02:20
 */
package goex

import "log"

func mains() {
	spot := &MockOrder{
		Direction: NewOrder_Buy,
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
		Direction: NewOrder_Buy,
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
		Direction: NewOrder_Buy,
		Type:      NewOrder_UM,
		Lever:     5,
		Quantity:  10000,
		FeeRate:   0.04,
	}
	um.Buy(230)
	um.Sell(231)
	log.Fatalf("U本位测试%#v\r\n", um.CalcCpUp())
}
