package goex

import (
	"github.com/shopspring/decimal"
	"math"
	"trade-robot-bd/app/grid-strategy-svc/util/utils"
	"trade-robot-bd/libs/helper"
)

//"github.com/d4l3k/talib"

import (
	"github.com/markcheno/go-talib"
)

type PriceType int

const (
	InClose PriceType = iota + 1
	InHigh
	InLow
	InOpen
	InVolume
)

func Crossunder(source1, source2 []float64) bool {
	if len(source1) < 2 || len(source2) < 2 {
		return false
	}
	_source1R := source1[len(source1)-1]
	_source1L := source1[len(source1)-2]
	_source2R := source2[len(source2)-1]
	_source2L := source2[len(source2)-2]
	return _source1R < _source2R && _source1L >= _source2L
}
func Crossover(source1, source2 []float64) bool {
	if len(source1) < 2 || len(source2) < 2 {
		return false
	}
	_source1R := source1[len(source1)-1]
	_source1L := source1[len(source1)-2]
	_source2R := source2[len(source2)-1]
	_source2L := source2[len(source2)-2]
	return _source1R > _source2R && _source1L <= _source2L
}
func KlineToRenKo(k []*Kline, atrLen int, _type PriceType, places int32) (Renko []*Kline) {
	lr := CalcAtr(k, atrLen)
	avg := helper.AvgInArray(lr[len(lr)-helper.IfThen(len(lr) < atrLen, len(lr), atrLen):])
	//min := helper.AvgInArray(lr)
	Surge, _ := decimal.NewFromFloat(helper.Round(avg, helper.PlacesFloat(int(places)))).Round(places).Float64()
	for i := 0; i < len(k); i++ {
		item := k[i]
		if i == 0 {
			rk := &Kline{
				Open:      item.Open,
				Close:     item.Close,
				High:      item.Open,
				Low:       item.Close,
				Timestamp: item.Timestamp,
			}
			Renko = append(Renko, rk)
			continue
		}
		var priceDiff float64
		fGetA_B := func() (float64, float64) {
			var soureA float64
			var soureB float64
			last := Renko[len(Renko)-1]
			switch _type {
			case InOpen:
				soureA = item.Open
				soureB = last.Close
				break
			case InClose:
				soureA = item.Close
				soureB = last.Close
				break
			case InLow:
				soureA = item.Low
				soureB = last.Close
				break
			case InHigh:
				soureA = item.High
				soureB = last.Close
				break
			}
			return soureA, soureB
		}
		var soureA, soureB = fGetA_B()
		priceDiff, _ = decimal.NewFromFloat(soureA - soureB).Round(places).Float64()
		blocks := int(math.Abs(priceDiff) / Surge) //计算 ATR值分割后的Renko图数量
		change := helper.IfThen(priceDiff > 0, Surge, -Surge)
		addRenko := func() {
			var _, _soureB = fGetA_B()
			_open := _soureB
			_close, _ := decimal.NewFromFloat(_soureB + change).Round(places).Float64()
			Renko = append(Renko, &Kline{
				Open:      _open,
				Close:     _close,
				High:      _open,
				Low:       _close,
				Timestamp: item.Timestamp,
				Vol:       priceDiff,
			})
		}
		if blocks >= 1 {
			for j := 0; j < blocks; j++ {
				addRenko()
			}
		} else if math.Abs(priceDiff) >= Surge {
			addRenko()
		}
	}
	return
}

func CalcRsi(kLineData []*Kline, length int) (_rsi []float64) {
	prices := utils.SplitToFlotList(RealData(kLineData, InClose), "Reverse")
	return talib.Rsi(prices, (length))
}

func CalcMfi(kLineData []*Kline, length int) (_mfi []float64) {
	high := utils.SplitToFlotList(RealData(kLineData, InHigh), "Reverse")
	low := utils.SplitToFlotList(RealData(kLineData, InLow), "Reverse")
	_close := utils.SplitToFlotList(RealData(kLineData, InClose), "Reverse")
	volume := utils.SplitToFlotList(RealData(kLineData, InVolume), "Reverse")
	_mfi = talib.Mfi(high, low, _close, volume, (length))
	return _mfi
}
func CalcAtr(kLineData []*Kline, length int) (_mfi []float64) {
	high := utils.SplitToFlotList(RealData(kLineData, InHigh), "Reverse")
	low := utils.SplitToFlotList(RealData(kLineData, InLow), "Reverse")
	_close := utils.SplitToFlotList(RealData(kLineData, InClose), "Reverse")
	_mfi = talib.Atr(high, low, _close, (length))
	return _mfi
}
func CalcRvgi(kLineData []*Kline, length int) (r []float64) {
	open := utils.SplitToFlotList(RealData(kLineData, InOpen), "Reverse")
	high := utils.SplitToFlotList(RealData(kLineData, InHigh), "Reverse")
	low := utils.SplitToFlotList(RealData(kLineData, InLow), "Reverse")
	_close := utils.SplitToFlotList(RealData(kLineData, InClose), "Reverse")
	var rvi1, rvi2 []float64
	var swmi int
	for i := 0; i <= len(low); i++ {
		if i-swmi < 0 {
			swmi = 0
		} else {
			swmi = i - swmi
		}
		rvi1 = append(rvi1, swma(_close[i-swmi:i], open[i-swmi:i]))
		rvi2 = append(rvi2, swma(high[i-swmi:i], low[i-swmi:i]))
		s1 := sum(rvi1, length)
		s2 := sum(rvi2, length)
		if s1 == 0 || s2 == 0 {
			r = append(r, 0)
			continue
		}
		rvi, _ := decimal.NewFromFloat(s1).Div(decimal.NewFromFloat(s2)).Round(4).Float64()
		r = append(r, rvi)
	}
	return r
}
func sum[T helper.Number](in []T, lenth int) (r T) {
	var inIndex = 0
	for i := len(in) - 1; i > 0; i-- {
		if inIndex >= lenth {
			return r
		}
		_r, _ := decimal.NewFromFloat(float64(in[i])).Add(decimal.NewFromFloat(float64(r))).Float64()
		r = T(_r)
		inIndex++
	}
	return r
}
func swma[T helper.Number](a, b []T) T {
	if len(a) == 0 {
		return 0
	}
	var a0, a1, a2, a3 T = a[len(a)-1], 0, 0, 0
	var b0, b1, b2, b3 T = b[len(b)-1], 0, 0, 0

	if len(a) >= 2 {
		a1 = a[len(a)-2]
		b1 = b[len(b)-2]
	}
	if len(a) >= 3 {
		a2 = a[len(a)-3]
		b2 = b[len(b)-3]
	}
	if len(a) >= 4 {
		a3 = a[len(a)-4]
		b3 = b[len(b)-4]
	}
	//log.Printf("A:%v %v %v %v B:%v %v %v %v", a0, a1, a2, a3, b0, b1, b2, b3)
	f, _ := decimal.NewFromFloat(float64(((a0 - b0) + (2 * (a1 - b1)) + (2 * (a2 - b2)) + (a3 - b3)) / 6)).Round(4).Float64()
	return T(f)
}
func RealData(list []*Kline, priceTy PriceType) []float64 {
	var inReal []float64
	for i := len(list) - 1; i >= 0; i-- {
		k := list[i]
		switch priceTy {
		case InClose:
			inReal = append(inReal, k.Close)
			break
		case InHigh:
			inReal = append(inReal, k.High)
			break
		case InLow:
			inReal = append(inReal, k.Low)
			break
		case InOpen:
			inReal = append(inReal, k.Open)
			break
		case InVolume:
			inReal = append(inReal, k.Vol)
			break
		default:
			panic("please set ema type")
		}
	}
	return inReal
}
