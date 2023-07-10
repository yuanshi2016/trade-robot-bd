package goex

import "strings"

type Currency struct {
	Symbol string
	Desc   string
}

func (c Currency) String() string {
	return c.Symbol
}

func (c Currency) Eq(c2 Currency) bool {
	return c.Symbol == c2.Symbol
}

// A->B(A兑换为B)
type CurrencyPair struct {
	CurrencyA Currency
	CurrencyB Currency
}
type AssetTransferType struct {
	Type string
	Desc string
}

func (c AssetTransferType) String() string {
	return c.Type
}

var (
	UNKNOWN  = Currency{"UNKNOWN", ""}
	CNY      = Currency{"CNY", ""}
	USD      = Currency{"USD", ""}
	USDPERP  = Currency{"USD_PERP", ""}
	USDTPERP = Currency{"USDT_PERP", ""}
	USDT     = Currency{"USDT", ""}
	PAX      = Currency{"PAX", "https://www.paxos.com/"}
	USDC     = Currency{"USDC", "https://www.centre.io/"}
	EUR      = Currency{"EUR", ""}
	KRW      = Currency{"KRW", ""}
	JPY      = Currency{"JPY", ""}
	BTC      = Currency{"BTC", "https://bitcoin.org/"}
	XBT      = Currency{"XBT", ""}
	BCC      = Currency{"BCC", ""}
	BCH      = Currency{"BCH", ""}
	FIL      = Currency{"FIL", ""}
	PEPE     = Currency{"1000PEPE", ""}
	SOL      = Currency{"SOL", ""}
	BCX      = Currency{"BCX", ""}
	LTC      = Currency{"LTC", ""}
	ETH      = Currency{"ETH", ""}
	CTSI     = Currency{"CTSI", ""}
	ONT      = Currency{"ONT", ""}
	ETC      = Currency{"ETC", ""}
	EOS      = Currency{"EOS", ""}
	BTS      = Currency{"BTS", ""}
	QTUM     = Currency{"QTUM", ""}
	SC       = Currency{"SC", ""}
	ANS      = Currency{"ANS", ""}
	ZEC      = Currency{"ZEC", ""}
	DCR      = Currency{"DCR", ""}
	XRP      = Currency{"XRP", ""}
	BTG      = Currency{"BTG", ""}
	BCD      = Currency{"BCD", ""}
	NEO      = Currency{"NEO", ""}
	HSR      = Currency{"HSR", ""}
	BSV      = Currency{"BSV", ""}
	OKB      = Currency{"OKB", "OKB is a global utility token issued by OK Blockchain Foundation"}
	HT       = Currency{"HT", "HuoBi Token"}
	BNB      = Currency{"BNB", "BNB, or Binance Coin, is a cryptocurrency created by Binance."}
	TRX      = Currency{"TRX", ""}
	MTL      = Currency{"MTL", ""}

	//currency pair

	BTC_CNY  = CurrencyPair{BTC, CNY}
	LTC_CNY  = CurrencyPair{LTC, CNY}
	BCC_CNY  = CurrencyPair{BCC, CNY}
	ETH_CNY  = CurrencyPair{ETH, CNY}
	ETC_CNY  = CurrencyPair{ETC, CNY}
	EOS_CNY  = CurrencyPair{EOS, CNY}
	BTS_CNY  = CurrencyPair{BTS, CNY}
	QTUM_CNY = CurrencyPair{QTUM, CNY}
	SC_CNY   = CurrencyPair{SC, CNY}
	ANS_CNY  = CurrencyPair{ANS, CNY}
	ZEC_CNY  = CurrencyPair{ZEC, CNY}

	BTC_KRW = CurrencyPair{BTC, KRW}
	ETH_KRW = CurrencyPair{ETH, KRW}
	ETC_KRW = CurrencyPair{ETC, KRW}
	LTC_KRW = CurrencyPair{LTC, KRW}
	BCH_KRW = CurrencyPair{BCH, KRW}

	BTC_USD       = CurrencyPair{BTC, USD}
	BNB_USD       = CurrencyPair{BNB, USD}
	BNB_USD_PERP  = CurrencyPair{BNB, USDPERP}
	BNB_USDT_PERP = CurrencyPair{BNB, USDTPERP}
	LTC_USD       = CurrencyPair{LTC, USD}
	ETH_USD       = CurrencyPair{ETH, USD}
	ETC_USD       = CurrencyPair{ETC, USD}
	BCH_USD       = CurrencyPair{BCH, USD}
	BCC_USD       = CurrencyPair{BCC, USD}
	XRP_USD       = CurrencyPair{XRP, USD}
	BCD_USD       = CurrencyPair{BCD, USD}
	EOS_USD       = CurrencyPair{EOS, USD}
	BTG_USD       = CurrencyPair{BTG, USD}
	BSV_USD       = CurrencyPair{BSV, USD}

	MTL_USDT  = CurrencyPair{MTL, USDT}
	BTC_USDT  = CurrencyPair{BTC, USDT}
	LTC_USDT  = CurrencyPair{LTC, USDT}
	BCH_USDT  = CurrencyPair{BCH, USDT}
	FIL_USDT  = CurrencyPair{FIL, USDT}
	PEPE_USDT = CurrencyPair{PEPE, USDT}
	SOL_USDT  = CurrencyPair{SOL, USDT}
	BCC_USDT  = CurrencyPair{BCC, USDT}
	ETC_USDT  = CurrencyPair{ETC, USDT}
	ETH_USDT  = CurrencyPair{ETH, USDT}
	CTSI_USDT = CurrencyPair{CTSI, USDT}
	ONT_USDT  = CurrencyPair{ONT, USDT}
	BCD_USDT  = CurrencyPair{BCD, USDT}
	NEO_USDT  = CurrencyPair{NEO, USDT}
	EOS_USDT  = CurrencyPair{EOS, USDT}
	XRP_USDT  = CurrencyPair{XRP, USDT}
	HSR_USDT  = CurrencyPair{HSR, USDT}
	BSV_USDT  = CurrencyPair{BSV, USDT}
	OKB_USDT  = CurrencyPair{OKB, USDT}
	HT_USDT   = CurrencyPair{HT, USDT}
	BNB_USDT  = CurrencyPair{BNB, USDT}
	PAX_USDT  = CurrencyPair{PAX, USDT}
	TRX_USDT  = CurrencyPair{TRX, USDT}

	XRP_EUR = CurrencyPair{XRP, EUR}

	BTC_JPY = CurrencyPair{BTC, JPY}
	LTC_JPY = CurrencyPair{LTC, JPY}
	ETH_JPY = CurrencyPair{ETH, JPY}
	ETC_JPY = CurrencyPair{ETC, JPY}
	BCH_JPY = CurrencyPair{BCH, JPY}

	LTC_BTC = CurrencyPair{LTC, BTC}
	ETH_BTC = CurrencyPair{ETH, BTC}
	ETC_BTC = CurrencyPair{ETC, BTC}
	BCC_BTC = CurrencyPair{BCC, BTC}
	BCH_BTC = CurrencyPair{BCH, BTC}
	DCR_BTC = CurrencyPair{DCR, BTC}
	XRP_BTC = CurrencyPair{XRP, BTC}
	BTG_BTC = CurrencyPair{BTG, BTC}
	BCD_BTC = CurrencyPair{BCD, BTC}
	NEO_BTC = CurrencyPair{NEO, BTC}
	EOS_BTC = CurrencyPair{EOS, BTC}
	HSR_BTC = CurrencyPair{HSR, BTC}
	BSV_BTC = CurrencyPair{BSV, BTC}
	OKB_BTC = CurrencyPair{OKB, BTC}
	HT_BTC  = CurrencyPair{HT, BTC}
	BNB_BTC = CurrencyPair{BNB, BTC}
	TRX_BTC = CurrencyPair{TRX, BTC}

	ETC_ETH = CurrencyPair{ETC, ETH}
	EOS_ETH = CurrencyPair{EOS, ETH}
	ZEC_ETH = CurrencyPair{ZEC, ETH}
	NEO_ETH = CurrencyPair{NEO, ETH}
	HSR_ETH = CurrencyPair{HSR, ETH}
	LTC_ETH = CurrencyPair{LTC, ETH}

	UNKNOWN_PAIR                           = CurrencyPair{UNKNOWN, UNKNOWN}
	Transfer_MAIN_UMFUTURE                 = AssetTransferType{"MAIN_UMFUTURE", "现货钱包转向U本位合约钱包"}
	Transfer_MAIN_CMFUTURE                 = AssetTransferType{"MAIN_CMFUTURE", "现货钱包转向币本位合约钱包"}
	Transfer_MAIN_MARGIN                   = AssetTransferType{"MAIN_MARGIN", "现货钱包转向杠杆全仓钱包"}
	Transfer_UMFUTURE_MAIN                 = AssetTransferType{"UMFUTURE_MAIN", "U本位合约钱包转向现货钱包"}
	Transfer_UMFUTURE_MARGIN               = AssetTransferType{"UMFUTURE_MARGIN", "U本位合约钱包转向杠杆全仓钱包"}
	Transfer_CMFUTURE_MAIN                 = AssetTransferType{"CMFUTURE_MAIN", "币本位合约钱包转向现货钱包"}
	Transfer_MARGIN_MAIN                   = AssetTransferType{"MARGIN_MAIN", "杠杆全仓钱包转向现货钱包"}
	Transfer_MARGIN_UMFUTURE               = AssetTransferType{"MARGIN_UMFUTURE", "杠杆全仓钱包转向U本位合约钱包"}
	Transfer_MARGIN_CMFUTURE               = AssetTransferType{"MARGIN_CMFUTURE", "杠杆全仓钱包转向币本位合约钱包"}
	Transfer_CMFUTURE_MARGIN               = AssetTransferType{"CMFUTURE_MARGIN", "币本位合约钱包转向杠杆全仓钱包"}
	Transfer_ISOLATEDMARGIN_MARGIN         = AssetTransferType{"ISOLATEDMARGIN_MARGIN", "杠杆逐仓钱包转向杠杆全仓钱包"}
	Transfer_MARGIN_ISOLATEDMARGIN         = AssetTransferType{"MARGIN_ISOLATEDMARGIN", "杠杆全仓钱包转向杠杆逐仓钱包"}
	Transfer_ISOLATEDMARGIN_ISOLATEDMARGIN = AssetTransferType{"ISOLATEDMARGIN_ISOLATEDMARGIN", "杠杆逐仓钱包转向杠杆逐仓钱包"}
	Transfer_MAIN_FUNDING                  = AssetTransferType{"MAIN_FUNDING", "现货钱包转向资金钱包"}
	Transfer_FUNDING_MAIN                  = AssetTransferType{"FUNDING_MAIN", "资金钱包转向现货钱包"}
	Transfer_FUNDING_UMFUTURE              = AssetTransferType{"FUNDING_UMFUTURE", "资金钱包转向U本位合约钱包"}
	Transfer_UMFUTURE_FUNDING              = AssetTransferType{"UMFUTURE_FUNDING", "U本位合约钱包转向资金钱包"}
	Transfer_MARGIN_FUNDING                = AssetTransferType{"MARGIN_FUNDING", "杠杆全仓钱包转向资金钱包"}
	Transfer_FUNDING_MARGIN                = AssetTransferType{"FUNDING_MARGIN", "资金钱包转向杠杆全仓钱包"}
	Transfer_FUNDING_CMFUTURE              = AssetTransferType{"FUNDING_CMFUTURE", "资金钱包转向币本位合约钱包"}
	Transfer_CMFUTURE_FUNDING              = AssetTransferType{"CMFUTURE_FUNDING", "币本位合约钱包转向资金钱包"}
	Transfer_MAIN_OPTION                   = AssetTransferType{"MAIN_OPTION", "现货钱包转向期权钱包"}
	Transfer_OPTION_MAIN                   = AssetTransferType{"OPTION_MAIN", "期权钱包转向现货钱包"}
	Transfer_UMFUTURE_OPTION               = AssetTransferType{"UMFUTURE_OPTION", "U本位合约钱包转向期权钱包"}
	Transfer_OPTION_UMFUTURE               = AssetTransferType{"OPTION_UMFUTURE", "期权钱包转向U本位合约钱包"}
	Transfer_MARGIN_OPTION                 = AssetTransferType{"MARGIN_OPTION", "杠杆全仓钱包转向期权钱包"}
	Transfer_OPTION_MARGIN                 = AssetTransferType{"OPTION_MARGIN", "期权全仓钱包转向杠杆钱包"}
	Transfer_FUNDING_OPTION                = AssetTransferType{"FUNDING_OPTION", "资金钱包转向期权钱包"}
	Transfer_OPTION_FUNDING                = AssetTransferType{"OPTION_FUNDING", "期权钱包转向资金钱包"}
	Transfer_MAIN_PORTFOLIO_MARGIN         = AssetTransferType{"MAIN_PORTFOLIO_MARGIN", "现货钱包转向统一账户钱包"}
	Transfer_PORTFOLIO_MARGIN_MAIN         = AssetTransferType{"PORTFOLIO_MARGIN_MAIN", "统一账户钱包转向现货钱包"}
)

func (c CurrencyPair) String() string {
	return c.ToSymbol("_")
}

func (c CurrencyPair) Eq(c2 CurrencyPair) bool {
	return c.String() == c2.String()
}

func (c Currency) AdaptBchToBcc() Currency {
	if c.Symbol == "BCH" || c.Symbol == "bch" {
		return BCC
	}
	return c
}

func (c Currency) AdaptBccToBch() Currency {
	if c.Symbol == "BCC" || c.Symbol == "bcc" {
		return BCH
	}
	return c
}

func NewCurrency(symbol, desc string) Currency {
	switch symbol {
	case "cny", "CNY":
		return CNY
	case "usdt", "USDT":
		return USDT
	case "usd", "USD":
		return USD
	case "usdc", "USDC":
		return USDC
	case "pax", "PAX":
		return PAX
	case "jpy", "JPY":
		return JPY
	case "krw", "KRW":
		return KRW
	case "eur", "EUR":
		return EUR
	case "btc", "BTC":
		return BTC
	case "xbt", "XBT":
		return XBT
	case "bch", "BCH":
		return BCH
	case "bcc", "BCC":
		return BCC
	case "ltc", "LTC":
		return LTC
	case "sc", "SC":
		return SC
	case "ans", "ANS":
		return ANS
	case "neo", "NEO":
		return NEO
	case "okb", "OKB":
		return OKB
	case "ht", "HT":
		return HT
	case "bnb", "BNB":
		return BNB
	case "trx", "TRX":
		return TRX
	default:
		return Currency{strings.ToUpper(symbol), desc}
	}
}

func NewCurrencyPair(currencyA Currency, currencyB Currency) CurrencyPair {
	return CurrencyPair{currencyA, currencyB}
}

func NewCurrencyPair2(currencyPairSymbol string) CurrencyPair {
	return NewCurrencyPair3(currencyPairSymbol, "_")
}

func NewCurrencyPair3(currencyPairSymbol string, sep string) CurrencyPair {
	currencys := strings.Split(currencyPairSymbol, sep)
	if len(currencys) >= 2 {
		return CurrencyPair{NewCurrency(currencys[0], ""),
			NewCurrency(currencys[1], "")}
	}
	return UNKNOWN_PAIR
}

func (pair CurrencyPair) ToSymbol(joinChar string) string {
	return strings.Join([]string{pair.CurrencyA.Symbol, pair.CurrencyB.Symbol}, joinChar)
}

func (pair CurrencyPair) ToSymbol2(joinChar string) string {
	return strings.Join([]string{pair.CurrencyB.Symbol, pair.CurrencyA.Symbol}, joinChar)
}

func (pair CurrencyPair) AdaptUsdtToUsd() CurrencyPair {
	CurrencyB := pair.CurrencyB
	if pair.CurrencyB.Eq(USDT) {
		CurrencyB = USD
	}
	return CurrencyPair{pair.CurrencyA, CurrencyB}
}

func (pair CurrencyPair) AdaptUsdToUsdt() CurrencyPair {
	CurrencyB := pair.CurrencyB
	if pair.CurrencyB.Eq(USD) {
		CurrencyB = USDT
	}
	return CurrencyPair{pair.CurrencyA, CurrencyB}
}

// It is currently applicable to binance and zb
func (pair CurrencyPair) AdaptBchToBcc() CurrencyPair {
	CurrencyA := pair.CurrencyA
	if pair.CurrencyA.Eq(BCH) {
		CurrencyA = BCC
	}
	return CurrencyPair{CurrencyA, pair.CurrencyB}
}

func (pair CurrencyPair) AdaptBccToBch() CurrencyPair {
	if pair.CurrencyA.Eq(BCC) {
		return CurrencyPair{BCH, pair.CurrencyB}
	}
	return pair
}

// for to symbol lower , Not practical '==' operation method
func (pair CurrencyPair) ToLower() CurrencyPair {
	return CurrencyPair{Currency{strings.ToLower(pair.CurrencyA.Symbol), pair.CurrencyA.Desc},
		Currency{strings.ToLower(pair.CurrencyB.Symbol), pair.CurrencyB.Desc}}
}

func (pair CurrencyPair) Reverse() CurrencyPair {
	return CurrencyPair{pair.CurrencyB, pair.CurrencyA}
}
