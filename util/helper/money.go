package helper

import (
	"fmt"
	"math"
)

// 金额转换中文描述
func Money2Str(money interface{}) (moneyStr string) {
	if money == 0 || money == "" {
		moneyStr = "--"
	}
	moneyStr = fmt.Sprintf("￥%v元", money)
	return
}

// Bccomp
// -1 if a < b, 0 if a == b, and 1 if a > b.
func Bccomp(a, b float64, precision int) int {
	scale := math.Pow(10, float64(precision))
	aScaled := int64(a * scale)
	bScaled := int64(b * scale)

	if aScaled < bScaled {
		return -1
	} else if aScaled > bScaled {
		return 1
	}
	return 0
}
