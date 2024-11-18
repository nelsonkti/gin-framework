package helper

import (
	"github.com/shopspring/decimal"
	"math"
)

func RoundDecimal(number float64, decimalPlaces int) float64 {
	multiplier := math.Pow(10, float64(decimalPlaces))
	rounded := math.Round(number*multiplier) / multiplier
	return rounded
}

// 浮点数相除 a / b
func FloatDiv(a float64, b float64) (res float64) {
	f1 := decimal.NewFromFloat(a)
	f2 := decimal.NewFromFloat(b)
	res, _ = f1.Div(f2).Float64()
	return
}

// 浮点数相乘 a * b
func FloatMul(a float64, b float64) (res float64) {
	f1 := decimal.NewFromFloat(a)
	f2 := decimal.NewFromFloat(b)
	res, _ = f1.Mul(f2).Float64()
	return
}

// 浮点数相减 a - b
func FloatSub(a float64, b float64) (res float64) {
	f1 := decimal.NewFromFloat(a)
	f2 := decimal.NewFromFloat(b)
	res, _ = f1.Sub(f2).Float64()
	return
}

// 浮点数相加 a + b
func FloatAdd(a float64, b float64) (res float64) {
	f1 := decimal.NewFromFloat(a)
	f2 := decimal.NewFromFloat(b)
	res, _ = f1.Add(f2).Float64()
	return
}
