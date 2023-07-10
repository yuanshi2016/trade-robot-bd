package helper

import (
	"github.com/shopspring/decimal"
	"math"
)

// PlacesFloat 传入指定小数位 返回0.?1
func PlacesFloat(decimalPlaces int) float64 {
	pow := math.Pow10(decimalPlaces)
	num := decimal.NewFromFloat((pow + 1) / math.Pow10(decimalPlaces))
	a, _ := num.Sub(decimal.NewFromInt(1)).Float64()
	return a
}
