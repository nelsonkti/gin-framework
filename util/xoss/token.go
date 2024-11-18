package xoss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/segmentio/ksuid"
	"path/filepath"
	"strings"
	"time"
)

// PolicyCondition 结构体定义上传策略的过期时间和条件
type PolicyCondition struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

// OSSTokenResult 结构体定义返回的OSS Token信息
type OSSTokenResult struct {
	AccessId  string `json:"access_id"`
	Host      string `json:"host"`
	Policy    string `json:"policy"`
	Signature string `json:"signature"`
	Expire    int64  `json:"expire"`
	Dir       string `json:"dir"`
	Token     string `json:"token"`
	Date      string `json:"date"`
	RegionId  string `json:"region_id"`
	FileName  string `json:"file_name"`
	Url       string `json:"url"`
}

// 获取当前时间的ISO8601格式
func gmtISO8601(expireEnd int64) string {
	t := time.Unix(expireEnd, 0).UTC()
	return t.Format("2006-01-02T15:04:05Z")
}

// 生成HMAC-SHA1签名
func hmacSHA1(key []byte, data string) []byte {
	h := hmac.New(sha1.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

type TokenOption struct {
	maxSize  int64 // 文件大小
	filetype int64 // 文件类型
}

type TokenOptions func(*TokenOption)

func WithMaxSize(maxSize int64) TokenOptions {
	return func(c *TokenOption) {
		c.maxSize = maxSize
	}
}

func WithFileType(filetype int64) TokenOptions {
	return func(c *TokenOption) {
		c.filetype = filetype
	}
}

// GenerateOSSToken 生成OSS上传凭证
func (a *Aliyun) GenerateOSSToken(opts ...TokenOptions) (*OSSTokenResult, error) {
	o := a.tokenOption(opts)

	stsResponse, err := a.generateSTSTokenWithCache()
	if err != nil {
		return nil, err
	}

	host := fmt.Sprintf("https://%s.%s", a.conf.Bucket, strings.TrimPrefix(a.conf.Endpoint, "https://"))

	// 设置策略过期时间
	now := time.Now().Unix()
	expire := a.getExpire(o.filetype) // 设置policy超时时间为1小时
	expireEnd := now + expire
	expiration := gmtISO8601(expireEnd)

	fileName := ksuid.New().String()
	dir := a.uploadFilePathPrefix(o.filetype, a.appName)
	date := carbon.Now().Format("Ymd")

	dir = filepath.Join(dir, date, fileName)
	dir = strings.ReplaceAll(dir, "\\", "/")
	// 设置上传策略
	conditions := []interface{}{
		[]interface{}{"content-length-range", 0, o.maxSize},
		[]interface{}{"starts-with", "$key", dir},
		[]interface{}{"in", "$content-type", []string{
			"image/jpg", "image/jpeg", "image/png", "image/gif", "image/bmp",
			"text/plain",
			"application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/pdf",
			"application/json", "text/csv",
			"application/octet-stream", "application/zip", "application/x-tar", "application/x-gzip",
			"audio/mpeg", "video/x-theora+ogg", "audio/mp4", "audio/wav", "audio/x-wav", "audio/x-m4a",
			"video/3gpp", "video/mp4", "video/mpeg", "video/quicktime", "video/x-flv", "video/x-ms-wmv", "video/avi", "video/x-msvideo",
			"application/vnd.android.package-archive",
		}}, // 添加文件类型限制
	}

	// 如果指定了文件名前缀，添加到条件中
	policyCondition := PolicyCondition{
		Expiration: expiration,
		Conditions: conditions,
	}

	// 将上传策略转为JSON格式
	policyJSON, err := json.Marshal(policyCondition)
	if err != nil {
		return nil, fmt.Errorf("oss failed to marshal policy: %v", err)
	}

	// 将策略进行Base64编码
	base64Policy := base64.StdEncoding.EncodeToString(policyJSON)

	// 生成签名
	signature := base64.StdEncoding.EncodeToString(hmacSHA1([]byte(stsResponse.Credentials.AccessKeySecret), base64Policy))

	// 返回OSS上传凭证
	return &OSSTokenResult{
		AccessId:  stsResponse.Credentials.AccessKeyId,
		Host:      host,
		Policy:    base64Policy,
		Signature: signature,
		Expire:    expireEnd,
		Dir:       dir,
		Token:     stsResponse.Credentials.SecurityToken,
		Date:      time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		RegionId:  a.conf.RegionId,
		FileName:  fileName,
		Url:       a.conf.Url,
	}, nil
}

func (a *Aliyun) tokenOption(opts []TokenOptions) TokenOption {
	var o TokenOption
	for _, opt := range opts {
		opt(&o)
	}
	if o.maxSize == 0 {
		o.maxSize = 5242880
	}
	return o
}

func (a *Aliyun) getExpire(fileType int64) int64 {
	if fileType == FileTypeAudio || fileType == FileTypeVideo {
		return 1800
	}
	return 300
}
