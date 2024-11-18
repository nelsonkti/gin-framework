package cron

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-framework/util/xlog"
)

var defaultCron *Cron

// StartCronTab initializes and starts the default Cron with a specific locker implementation and configuration.
// This function is meant to simplify the creation and usage of a Cron instance with default settings.
func StartCronTab(redis *redis.Client, appName string, logger *xlog.Log) *Cron {
	return InitDefaultCron(&Config{
		Redis:  redis,
		logger: logger,
		Prefix: fmt.Sprintf("%s/%s", appName, "cron"),
		Factor: 0.01,
	})
}

// InitDefaultCron initializes the default Cron instance with the given configuration.
// It panics if called more than once to prevent unintentional reinitialization.
func InitDefaultCron(config *Config) *Cron {
	if defaultCron != nil {
		config.logger.Panic("defaultCron init twice.")
	}
	defaultCron = NewCron(config)
	return defaultCron
}
