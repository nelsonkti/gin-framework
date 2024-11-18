package verify

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/util"
	"net/url"
	"sort"
	"strings"
)

// GenerateHash 生成参数的哈希值
func GenerateHash(params url.Values) string {
	// 将参数键排序以保证哈希的一致性
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 创建要哈希的字符串
	var paramString string
	for _, key := range keys {
		paramString += key + "=" + params.Get(key) + "&"
	}
	paramString = strings.TrimRight(paramString, "&")

	return util.Md5(util.Md5(paramString))
}

// Md5 md5加密
func Md5(params map[string]interface{}) string {
	return util.Md5(util.Md5(createEncryptStr(params)))
}

// createEncryptStr 创建加密字符串
func createEncryptStr(params map[string]interface{}) string {
	var key []string
	var str = ""
	for k := range params {
		if k != "sn" && k != "debug" {
			key = append(key, k)
		}
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params[key[i]])
		} else {
			str = str + fmt.Sprintf("&%v=%v", key[i], params[key[i]])
		}
	}
	return str
}
