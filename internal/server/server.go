package server

import (
	"context"
	"go-framework/config"
	"go-framework/internal/container/common/tool"
	"go-framework/internal/container/grpc"
	"go-framework/internal/container/repository"
	"go-framework/internal/data/common_data/tool_data"
	"go-framework/internal/mq"
	"go-framework/util/mq/rocketmq"
	"go-framework/util/thread"
	"go-framework/util/tracer"
	"go-framework/util/xlog"
	"go-framework/util/xredis"
	"go-framework/util/xsql"
	"go-framework/util/xsql/databese"
)

type SvcContext struct {
	Ctx         context.Context
	Conf        config.Conf
	DBEngine    *databese.Engine
	RedisClient *xredis.RedisClient
	Logger      *xlog.Log
	MQClient    *rocketmq.Client
	Repo        *repository.Container
	Tool        *tool.Container
	Grpc        *grpc.Container
}

func NewSvcContext(c config.Conf, logger *xlog.Log) *SvcContext {

	svc := &SvcContext{
		Conf:        c,
		Logger:      logger,
		Ctx:         context.Background(),
		DBEngine:    xsql.NewClient(c.DB),
		RedisClient: xredis.NewClient(c.Redis),
	}
	svc.MQClient = rocketmq.NewClient(c, logger, svc.RedisClient.Default(), mq.RegisterQueue)
	svc.Repo = repository.Register(svc.DBEngine, svc.Logger)

	svc.Tool = tool.Register(&tool_data.SvcContext{
		Conf:        c,
		Logger:      svc.Logger,
		RedisClient: svc.RedisClient,
		Repo:        svc.Repo,
		MQClient:    svc.MQClient,
	})

	xsql.SetNotifier(svc.Tool.DingtalkTool)

	svc.MQClient.SetNotifier(svc.Tool.DingtalkTool)
	svc.MQClient.ConsumerRun(mq.ConsumerHandler)
	// 客户端
	grpcClient := grpc.Register(c, svc.Ctx)

	// grpc客户端
	svc.Grpc = grpcClient

	tracer.NewOpentelemetry(c.App.Name, c.App.Env, c.Trace.Endpoint, c.Trace.UrlPath)

	thread.SetNotifier(c.App.Name, c.App.Env, c.App.ServerNumber, svc.Tool.DingtalkTool.AlarmRobot)

	return svc
}
