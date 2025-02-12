package goex

import (
	"fmt"
)

type TradeSide int
type ActionType int

const (
	OPEN_BUY   ActionType = 1 + iota //开多
	OPEN_SELL                        //开空
	CLOSE_BUY                        //平多
	CLOSE_SELL                       //平空
)
const (
	SELL TradeSide = -1 - iota
	BUY  TradeSide = iota
	BUY_MARKET
	SELL_MARKET
)

func (ts TradeSide) String() string {
	switch ts {
	case BUY:
		return "BUY"
	case SELL:
		return "SELL"
	case BUY_MARKET:
		return "BUY_MARKET"
	case SELL_MARKET:
		return "SELL_MARKET"
	default:
		return "UNKNOWN"
	}
}

type TradeStatus int

func (ts TradeStatus) String() string {
	return tradeStatusSymbol[ts]
}

var tradeStatusSymbol = [...]string{"UNFINISH", "PART_FINISH", "FINISH", "CANCEL", "REJECT", "CANCEL_ING", "FAIL"}

const (
	ORDER_UNFINISH TradeStatus = iota
	ORDER_PART_FINISH
	ORDER_FINISH
	ORDER_CANCEL
	ORDER_REJECT
	ORDER_CANCEL_ING
	ORDER_FAIL
)

type KlinePeriod int

// k线周期
const (
	KLINE_PERIOD_1MIN KlinePeriod = 1 + iota
	KLINE_PERIOD_3MIN
	KLINE_PERIOD_5MIN
	KLINE_PERIOD_15MIN
	KLINE_PERIOD_30MIN
	KLINE_PERIOD_60MIN
	KLINE_PERIOD_1H
	KLINE_PERIOD_2H
	KLINE_PERIOD_3H
	KLINE_PERIOD_4H
	KLINE_PERIOD_6H
	KLINE_PERIOD_8H
	KLINE_PERIOD_12H
	KLINE_PERIOD_1DAY
	KLINE_PERIOD_3DAY
	KLINE_PERIOD_1WEEK
	KLINE_PERIOD_1MONTH
	KLINE_PERIOD_1YEAR
)

type OrderFeature int

const (
	ORDER_FEATURE_ORDINARY = 0 + iota
	ORDER_FEATURE_POST_ONLY
	ORDER_FEATURE_FOK
	ORDER_FEATURE_IOC
	ORDER_FEATURE_FAK
	ORDER_FEATURE_LIMIT
)

func (of OrderFeature) String() string {
	if of > 0 && int(of) < len(orderFeatureSymbol) {
		return orderFeatureSymbol[of]
	}
	return fmt.Sprintf("UNKNOWN_ORDER_TYPE(%d)", of)
}

var orderFeatureSymbol = [...]string{"ORDINARY", "POST_ONLY", "FOK", "IOC", "FAK", "LIMIT"}

type OrderType int

func (ot OrderType) String() string {
	if ot > 0 && int(ot) <= len(orderTypeSymbol) {
		return orderTypeSymbol[ot-1]
	}
	return fmt.Sprintf("UNKNOWN_ORDER_TYPE(%d)", ot)
}

var orderTypeSymbol = [...]string{"LIMIT", "MARKET"}

const (
	ORDER_TYPE_LIMIT = 1 + iota
	ORDER_TYPE_MARKET
)

var (
	THIS_WEEK_CONTRACT  = "this_week"  //周合约
	NEXT_WEEK_CONTRACT  = "next_week"  //次周合约
	QUARTER_CONTRACT    = "quarter"    //季度合约
	BI_QUARTER_CONTRACT = "bi_quarter" // NEXT QUARTER
	SWAP_CONTRACT       = "swap"       //永续合约
)

// exchanges const
const (
	KUCOIN       = "kucoin.com"
	OKCOIN_COM   = "okcoin.com"
	OKEX         = "okex.com"
	OKEX_V3      = "okex.com_v3"
	OKEX_FUTURE  = "okex.com_future"
	OKEX_SWAP    = "okex.com_swap"
	HUOBI        = "huobi.com"
	HUOBI_PRO    = "huobi.pro"
	BITSTAMP     = "bitstamp.net"
	KRAKEN       = "kraken.com"
	ZB           = "zb.com"
	BITFINEX     = "bitfinex.com"
	BINANCE      = "binance.com"
	BINANCE_SWAP = "binance.com_swap"
	POLONIEX     = "poloniex.com"
	COINEX       = "coinex.com"
	BITHUMB      = "bithumb.com"
	GATEIO       = "gate.io"
	BITTREX      = "bittrex.com"
	GDAX         = "gdax.com"
	BIGONE       = "big.one"
	FCOIN        = "fcoin.com"
	FCOIN_MARGIN = "fcoin.com_margin"
	FMEX         = "fmex.com"
	HITBTC       = "hitbtc.com"
	BITMEX       = "bitmex.com"
	BITMEX_TEST  = "testnet.bitmex.com"
	CRYPTOPIA    = "cryptopia.co.nz"
	HBDM         = "hbdm.com"
	COINBENE     = "coinbene.com"
	ATOP         = "a.top"
)
