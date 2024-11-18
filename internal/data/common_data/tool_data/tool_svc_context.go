package tool_data

import (
	"go-framework/config"
	"go-framework/internal/container/repository"
	"go-framework/util/mq/rocketmq"
	"go-framework/util/xlog"
	"go-framework/util/xredis"
)

type SvcContext struct {
	Conf        config.Conf
	RedisClient *xredis.RedisClient
	Logger      *xlog.Log
	Repo        *repository.Container
	MQClient    *rocketmq.Client
}
