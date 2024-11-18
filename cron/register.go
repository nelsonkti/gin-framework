package cron

import (
	"github.com/go-redis/redis/v8"
	"go-framework/cron/task"
	"go-framework/util/cron"
	"go-framework/util/xlog"
)

func Register(redis *redis.Client, appName string, logger *xlog.Log) {
	c := cron.StartCronTab(redis, appName, logger)

	c.Register(&task.AutoGenerateMigrateTask{})
	c.Register(&task.DemoTask{})

	go c.Run()
}
