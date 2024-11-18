package xoss

import (
	"go-framework/util/helper"
	"go-framework/util/xredis"
)

type Aliyun struct {
	conf        *config
	RedisClient *xredis.RedisClient
	appName     string
}

type config struct {
	Endpoint        string `json:"endpoint"`          // 地址
	AccessKey       string `json:"access_key"`        // key
	AccessSecret    string `json:"access_secret"`     // 秘钥
	Bucket          string `json:"bucket"`            // （桶）存储空间
	RegionId        string `json:"region_id"`         // 地域
	RoleSessionName string `json:"role_session_name"` // 角色名称
	RoleArn         string `json:"role_arn"`          // 角色arn
	Url             string `json:"url"`               // oss外网访问地址
}

type AliyunImpl interface {
	Upload(filePath string, opts ...Options) (*UploadResult, error)
	UploadBinaryData(data []byte, opts ...Options) (*UploadResult, error)
	GenerateOSSToken(opts ...TokenOptions) (*OSSTokenResult, error)
}

// NewAliyun 创建NewAliyun的新实例
func NewAliyun(conf interface{}, redisClient *xredis.RedisClient) AliyunImpl {
	var configs *config
	err := helper.UnMarshalWithInterface(conf, &configs)
	if err != nil {
		panic(err)
	}
	return &Aliyun{conf: configs, appName: "mashang", RedisClient: redisClient}
}
