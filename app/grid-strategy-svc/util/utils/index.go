package utils

import (
	"github.com/shopspring/decimal"
	"strconv"
)

var Div100 = decimal.NewFromInt(100)

func IfThen[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}
func IfFunction[T any](expr bool, a, b func()) {
	if expr {
		a()
	}
	b()
}

// FindFloat64 获取一个切片并在其中查找元素。如果找到它，它将返回它的密钥，否则它将返回-1和一个错误的bool。
func FindFloat64(slice []float64, val string) bool {
	for _, item := range slice {
		s := strconv.FormatFloat(item, 'g', -1, 64)
		if s == val {
			return true
		}
	}
	return false
}

// DeleteSliceElms 删除切片指定元素（不许改原切片）
func DeleteSliceElms(sl []float64, elms ...float64) []float64 {
	if len(sl) == 0 || len(elms) == 0 {
		return sl
	}
	// 先将元素转为 set
	m := make(map[float64]struct{})
	for _, v := range elms {
		m[v] = struct{}{}
	}
	// 过滤掉指定元素
	res := make([]float64, 0, len(sl))
	for _, v := range sl {
		if _, ok := m[v]; !ok {
			res = append(res, v)
		}
	}
	return res
}
