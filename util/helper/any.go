package helper

import (
	"reflect"
	"strconv"
)

// AnyToFloat64 any类型转换为float64（小数点后两位）
func AnyToFloat64(a any) (result float64) {
	if reflect.TypeOf(a).Name() == "float64" {
		result = a.(float64)
	} else {
		result, _ = strconv.ParseFloat(a.(string), 32)
		// 浮点精度转换
		result = RoundDecimal(result, 2)
	}
	return
}

// AnyToInt any类型转换为int
func AnyToInt[T any](input T) int {
	switch v := any(input).(type) {
	case float32:
		return int(v)
	case float64:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	default:
		return 0
	}
}
