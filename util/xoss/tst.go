package xoss

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"time"
)

const (
	STSTokenKey = "sts:token:%s:key"
)

func (a *Aliyun) generateSTSTokenWithCache() (*sts.AssumeRoleResponse, error) {
	key := fmt.Sprintf(STSTokenKey, a.conf.AccessKey)
	data := a.RedisClient.Default().Get(context.Background(), key).Val()
	if data == "" {
		res, err := a.generateSTSToken()
		if err != nil {
			return nil, err
		}
		marshal, err := helper.Marshal(res)
		if err != nil {
			return nil, err
		}
		a.RedisClient.Default().Set(context.Background(), key, string(marshal), time.Second*3500)
		return res, nil
	}

	var res *sts.AssumeRoleResponse
	err := helper.UmMarshal([]byte(data), &res)
	if err != nil {
		return nil, fmt.Errorf("redis数据解析失败: %v", err)
	}

	return res, err
}

// 生成STS临时凭证
func (a *Aliyun) generateSTSToken() (*sts.AssumeRoleResponse, error) {
	client := a.getSTSClient()

	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = a.conf.RoleArn
	request.RoleSessionName = a.conf.RoleSessionName
	request.DurationSeconds = "3600"

	response, err := client.AssumeRole(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// 获取STS客户端
func (a *Aliyun) getSTSClient() *sts.Client {
	client, err := sts.NewClientWithAccessKey(a.conf.RegionId, a.conf.AccessKey, a.conf.AccessSecret)
	if err != nil {
		panic(err)
	}
	return client
}
