package signer

import (
	"crypto/md5"
	"encoding/hex"
	"go-framework/util/helper"
	"strconv"
	"time"
)

// GetAuthorization 生成授权认证字符串 (临时密钥，后面会废弃)
// 参数 userId 用户ID，用于构建授权认证字符串的一部分
// 返回值 Authorization 授权认证字符串，包含时间戳、加密的认证密钥、随机数和用户ID
func GetAuthorization(authKey, userId string) string {
	// 获取当前时间的Unix时间戳
	date1e10 := strconv.FormatInt(time.Now().Unix(), 10)
	// 生成8位随机数
	random1e8 := helper.Rand(8)
	// 初始化MD5加密算法
	hash := md5.New()
	// 将时间戳、认证密钥、随机数和用户ID拼接成字符串，并写入MD5加密算法中
	hash.Write([]byte(date1e10 + authKey + random1e8 + userId))
	// 将MD5加密结果转换为十六进制字符串
	authKeyMd5 := hex.EncodeToString(hash.Sum(nil))
	// 构建授权认证字符串，包含时间戳、加密的认证密钥、随机数和用户ID
	Authorization := date1e10 + authKeyMd5 + random1e8 + userId
	// 返回授权认证字符串
	return Authorization
}
