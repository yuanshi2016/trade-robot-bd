package helper

import (
	"fmt"
	"github.com/shopspring/decimal"
	"sort"
	"strconv"
	"strings"
)

type MapValue struct {
	Key   string
	Value []int
}

func SplitNumberToNumber(s []int) {
	strSlice := []string{"1", "2", "3", "4", "5"}
	numSlice := make([]int, len(strSlice))
	for i, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("转换错误：", err)
			return
		}
		numSlice[i] = num
	}
	fmt.Println(numSlice)
}

// SplitNumber 将一个数字转为切片
func SplitNumber(n int) []int {
	s := strconv.Itoa(n)
	res := make([]int, len(s))
	for i := range s {
		res[i], _ = strconv.Atoi(string(s[i]))
	}
	return res
}
func NumberIncOrDecTolist(number int64, length, types int) []int64 {
	list := make([]int64, 0)
	center := number
	for i := 0; i < length; i++ {
		if types < 0 {
			center -= 1
			if center <= 0 {
				break
			}
		} else {
			center += 1
		}
		list = append(list, center)
	}
	return list
}
func SplitDecimalSum[T Number](i32List []T) (total T) {
	for _, t := range i32List {
		total = total + t
	}
	return total
}
func SplitNumberJoin[T Number](arr []T, sep string) (i32List string) {
	if len(arr) == 0 {
		return ""
	}

	for _, val := range arr {
		i32List += fmt.Sprintf("%v%v", val, sep)
	}
	return i32List[0 : len(i32List)-1]
}

// SplitToDecimalList 传入字符串以及分隔符 返回Int64切片
func SplitToDecimalList(str string, sep string) (i32List []decimal.Decimal) {
	if str == "" {
		return
	}
	strList := strings.Split(str, sep)
	if len(strList) == 0 {
		return
	}
	for _, item := range strList {
		if item == "" {
			continue
		}
		val, err := strconv.ParseFloat(item, 10)
		if err != nil {
			// logs.CtxError(ctx, "ParseInt fail, err=%v, str=%v, sep=%v", err, str, sep) // 此处打印出log报错信息
			continue
		}
		i32List = append(i32List, decimal.NewFromFloat(val))
	}
	return i32List
}

// SplitToFloatList 传入字符串以及分隔符 返回Int64切片
func SplitToFloatList(str string, sep string) (i32List []float64) {
	if str == "" {
		return
	}
	strList := strings.Split(str, sep)
	if len(strList) == 0 {
		return
	}
	for _, item := range strList {
		if item == "" {
			continue
		}
		val, err := strconv.ParseFloat(item, 10)
		if err != nil {
			// logs.CtxError(ctx, "ParseInt fail, err=%v, str=%v, sep=%v", err, str, sep) // 此处打印出log报错信息
			continue
		}
		i32List = append(i32List, val)
	}
	return i32List
}

// SplitToIntList 传入字符串以及分隔符 返回Int64切片
func SplitToIntList(str string, sep string, actions string) (i32List sort.IntSlice) {
	if str == "" {
		return
	}
	strList := strings.Split(str, sep)
	if len(strList) == 0 {
		return
	}
	for _, item := range strList {
		if item == "" {
			continue
		}
		val, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			// logs.CtxError(ctx, "ParseInt fail, err=%v, str=%v, sep=%v", err, str, sep) // 此处打印出log报错信息
			continue
		}
		i32List = append(i32List, int(val))
	}
	switch actions {
	case "asc":
		sort.Slice(i32List, func(i, j int) bool {
			return i32List[i] < i32List[j]
		})
		break
	case "desc":
		sort.Slice(i32List, func(i, j int) bool {
			return i32List[i] > i32List[j]
		})
		break
	case "Reverse":
		for i, j := 0, len(i32List)-1; i < j; i, j = i+1, j-1 {
			i32List[i], i32List[j] = i32List[j], i32List[i]
		}
		break
	default:
		break
	}
	return i32List
}

// SplitToStringList 传入字符串以及分隔符 返回Int64切片
func SplitToStringList(str string, sep string, actions string) (i32List sort.StringSlice) {
	if str == "" {
		return
	}
	strList := strings.Split(str, sep)
	if len(strList) == 0 {
		return
	}
	for _, item := range strList {
		if item == "" {
			continue
		}
		i32List = append(i32List, item)
	}
	switch actions {
	case "asc":
		sort.Slice(i32List, func(i, j int) bool {
			return i32List[i] < i32List[j]
		})
		break
	case "desc":
		sort.Slice(i32List, func(i, j int) bool {
			return i32List[i] > i32List[j]
		})
		break
	case "Reverse":
		for i, j := 0, len(i32List)-1; i < j; i, j = i+1, j-1 {
			i32List[i], i32List[j] = i32List[j], i32List[i]
		}
		break
	default:
		break
	}
	return i32List
}

// SplitToFlotList 传入字符串以及分隔符 返回Int64切片
func SplitToFlotList[T Number](str []T, actions string) (i32List []T) {
	if len(str) == 0 {
		return
	}
	for _, item := range str {
		i32List = append(i32List, item)
	}
	switch actions {
	case "asc":
		sort.Slice(i32List, func(i, j int) bool {
			return i32List[i] < i32List[j]
		})
		break
	case "desc":
		sort.Slice(i32List, func(i, j int) bool {
			return i32List[i] > i32List[j]
		})
		break
	case "Reverse":
		for i, j := 0, len(i32List)-1; i < j; i, j = i+1, j-1 {
			i32List[i], i32List[j] = i32List[j], i32List[i]
		}
		break
	default:
		break
	}
	return i32List
}
func MaxInArray[T Number](arr []T) T {
	if len(arr) == 0 {
		return 0
	}
	max := arr[0] // 假设第一个元素是最大值
	for i := 1; i < len(arr); i++ {
		if arr[i] > max { // 如果当前元素比假设的最大值还大，则更新最大值
			max = arr[i]
		}
	}
	return max
}
func MinInArray[T Number](arr []T) T {
	if len(arr) == 0 {
		return 0
	}
	max := arr[0] // 假设第一个元素是最大值
	for i := 1; i < len(arr); i++ {
		if arr[i] < max { // 如果当前元素比假设的最大值还大，则更新最大值
			max = arr[i]
		}
	}
	return max
}
func SumInArray[T Number](arr []T) (r T) {
	for i := 0; i < len(arr); i++ {
		r += arr[i]
	}
	return r
}
func AvgInArray[T Number](arr []T) (r T) {
	return SumInArray(arr) / T(len(arr))
}
func MulInArray[T Number](arr []T) (r T) {
	r = 1
	for i := 0; i < len(arr); i++ {
		r *= IfThen(arr[i] > 0, arr[i], 1)
	}
	return r
}
