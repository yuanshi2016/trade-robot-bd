package cron

import "encoding/json"

type Kline struct {
	Open      float64 `json:"open,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Close     float64 `json:"close,string"`
	Vol       float64 `json:"vol,string"`
	CloseTime uint64  `json:"time,string"`
	QuoteVol  float64 `json:"amount,string"`
}
type Klines struct {
	Type   string `json:"type,string"`
	Symbol string `json:"symbol,string"`
	Fin    int    `json:"fin,int"`
	Data   Kline  `json:"data"`
}

func (tk Kline) MarshalBinary() ([]byte, error) {
	return json.Marshal(tk)
}

func (tk Klines) MarshalBinary() ([]byte, error) {
	return json.Marshal(tk)
}
