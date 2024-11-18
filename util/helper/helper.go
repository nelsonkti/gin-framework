package helper

import (
	"fmt"
	"go-framework/util/xlog"
	"runtime"
	"strconv"
	"strings"
)

// RecoverPanic 恢复panic
func RecoverPanic(logger *xlog.Log) {
	err := recover()
	if err != nil {
		logger.Error(err)

		buf := make([]byte, 2048)
		n := runtime.Stack(buf, false)
		logger.Errorf("%s", buf[:n])
	}
}

// HttpBuildQuery 构建http请求参数
func HttpBuildQuery(data map[string]interface{}) string {
	var query string
	for k, v := range data {
		query += fmt.Sprintf("%s=%v&", k, v)
	}
	return query[:len(query)-1]
}

// VersionCompare 比较两个版本号字符串
// 如果 v1 > v2，返回 1
// 如果 v1 < v2，返回 -1
// 如果相等，返回 0
func VersionCompare(v1, v2 string) int {
	// Split version numbers by dot
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Compare version number parts
	for i := 0; i < len(parts1) || i < len(parts2); i++ {
		// Convert string to int for comparison, assuming missing parts as 0
		var num1, num2 int
		var err error
		if i < len(parts1) {
			num1, err = strconv.Atoi(parts1[i])
			if err != nil {
				panic(err) // Or handle error appropriately
			}
		}
		if i < len(parts2) {
			num2, err = strconv.Atoi(parts2[i])
			if err != nil {
				panic(err) // Or handle error appropriately
			}
		}

		// Compare the two numbers
		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
		// If equal, move to the next parts
	}

	// Versions are equal
	return 0
}

// IsNewVersion 是否为最新版本
// 版本号.相隔的数不能超过三位，如果超过需要加大相乘的倍数
func IsNewVersion(version, miniVersion string) bool {
	arr1 := strings.Split(version, ".")
	arr2 := strings.Split(miniVersion, ".")

	if len(arr1) != 3 || len(arr2) != 3 {
		return false
	}

	num1, _ := strconv.Atoi(arr1[0])
	num2, _ := strconv.Atoi(arr1[1])
	num3, _ := strconv.Atoi(arr1[2])
	num4, _ := strconv.Atoi(arr2[0])
	num5, _ := strconv.Atoi(arr2[1])
	num6, _ := strconv.Atoi(arr2[2])

	if num4*1000000+num5*1000+num6 >= num1*1000000+num2*1000+num3 {
		return true
	}
	return false
}
