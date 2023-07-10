package helper

import (
	"math"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type Float interface {
	~float32 | ~float64
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Integer interface {
	Signed | Unsigned
}
type Number interface {
	~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 |
		~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32
}
type Ordered interface {
	Integer | Float
}
type Slice[T Ordered | Integer | Unsigned] []T

func GetTimeNowOfUinx() int64 {
	return time.Now().UTC().Unix()
}

func GetTimeNow() time.Time {
	return time.Now().UTC()
}

func StringToFloat64(s string) float64 {
	d, _ := decimal.NewFromString(s)
	f, _ := d.Float64()
	return f
}

func Float64ToString(f float64) string {
	return decimal.NewFromFloat(f).String()
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Float32ToString(f float32) string {
	return decimal.NewFromFloat32(f).String()
}
func IfThen[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
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

func MakeMathQuantity[T Ordered](_length int, start, inc T) (r []T) {
	r = append(r, start)
	for {
		if len(r) < int(_length) {
			r = append(r, r[len(r)-1]+inc)
		} else {
			break
		}
	}
	return r
}

// DeleteSliceElms 删除切片指定元素（不许改原切片）
func DeleteSliceElms[T Ordered](sl []T, elms ...T) []T {
	if len(sl) == 0 || len(elms) == 0 {
		return sl
	}
	// 先将元素转为 set
	m := make(map[T]struct{})
	for _, v := range elms {
		m[v] = struct{}{}
	}
	// 过滤掉指定元素
	res := make([]T, 0, len(sl))
	for _, v := range sl {
		if _, ok := m[v]; !ok {
			res = append(res, v)
		}
	}
	return res
}
func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
