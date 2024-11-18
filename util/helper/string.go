package helper

import (
	"math/rand"
	"strconv"
	"strings"
)

func StrToFloat64(str string) (float64, error) {
	if strings.EqualFold(str, "") || strings.EqualFold(str, "0.00") {
		return 0.00, nil
	}

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {

		return 0.00, err
	}

	return f, nil
}
func StrToInt(str string) (int, error) {
	f, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func ReplaceAll(s, new string, olds ...string) string {
	for _, old := range olds {
		s = strings.ReplaceAll(s, old, new)
	}
	return s
}

// StrToInt64 字符串转int64
func StrToInt64(str string) (int64, error) {
	if strings.EqualFold(str, "") {
		return 0, nil
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// StrToUint64 字符串转uint64
func StrToUint64(str string) (uint64, error) {
	if strings.EqualFold(str, "") {
		return 0, nil
	}

	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// IsEmptyString 判断字符串是否为空
func IsEmptyString(str string) bool {
	return strings.EqualFold(str, "")
}

// Float64ToStr float64转换为string
func Float64ToStr(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

// FormatFloat64WithCommas float64转换为字符串并加上千分位符号
func FormatFloat64WithCommas(num float64, prec int) string {
	// 将float64转换为string
	numStr := strconv.FormatFloat(num, 'f', prec, 64)

	// 分离整数部分和小数部分
	parts := strings.Split(numStr, ".")
	intPart := parts[0]
	fracPart := ""
	if len(parts) > 1 {
		fracPart = "." + parts[1]
	}

	// 对整数部分添加千分位分隔符
	var b strings.Builder
	length := len(intPart)
	for i := 0; i < length; i++ {
		if i > 0 && (length-i)%3 == 0 {
			b.WriteString(",")
		}
		b.WriteByte(intPart[i])
	}
	// 结合整数和小数部分
	return b.String() + fracPart
}

// 截取字符串，支持多字节字符
// start：起始下标，负数从从尾部开始，最后一个为-1
// length：截取长度，负数表示截取到末尾
func SubStr(str string, start int, length int) (result string) {
	s := []rune(str)
	total := len(s)
	if total == 0 {
		return
	}
	// 允许从尾部开始计算
	if start < 0 {
		start = total + start
		if start < 0 {
			return
		}
	}
	if start > total {
		return
	}
	// 到末尾
	if length < 0 {
		length = total
	}

	end := start + length
	if end > total {
		result = string(s[start:])
	} else {
		result = string(s[start:end])
	}

	return
}

// TruncateString 截取字符串并在超出length后面添加省略号
func TruncateString(str string, length int) string {
	if len(str) > length {
		return str[:length-3] + "..."
	}
	return str
}

func Rand(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
