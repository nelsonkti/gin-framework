package cron

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/robfig/cron/v3"
	"go-framework/util/xlog"
	"time"
)

const (
	DefaultMutexPrefix = "media-matrix/cron"
	DefaultMutexFactor = 0.05
)

// DistributedLocker describes the behavior required for distributed locking.
type DistributedLocker interface {
	Lock(key string, expiry time.Duration) (bool, error)
	Unlock(key string) error
}

// Handler defines a cron job handler.
type Handler struct {
	cron     string
	handle   func()
	name     string
	schedule cron.Schedule
}

// NewHandler creates a new Handler instance.
func NewHandler(cronExpr, name string, f func()) (*Handler, error) {
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return nil, err
	}
	return &Handler{
		cron:     cronExpr,
		handle:   f,
		name:     name,
		schedule: schedule,
	}, nil
}

// Config contains the configuration for the distributed mutex.
type Config struct {
	Redis  *redis.Client
	logger *xlog.Log
	Prefix string
	Factor float64
}

// Cron represents a cron job scheduler with distributed locking.
type Cron struct {
	cronClient *cron.Cron
	Config     *Config
	sync       *redsync.Redsync
}

// NewCron creates a new Cron instance.
func NewCron(config *Config) *Cron {
	pool := goredis.NewPool(config.Redis)
	c := &Cron{
		Config:     config,
		sync:       redsync.New(pool),
		cronClient: cron.New(),
	}
	return c
}

// Register adds a new task to the cron scheduler.
func (c *Cron) Register(task Task) {
	handler, err := NewHandler(task.Rule(), task.Name(), task.Run)
	if err != nil {
	}

	_, _ = c.cronClient.AddFunc(task.Rule(), c.handle(handler))
}

// Run starts the cron scheduler.
func (c *Cron) Run() {
	c.cronClient.Run()
}

func (c *Cron) lock(h *Handler) (bool, error) {
	now := time.Now()
	d := h.schedule.Next(now).Sub(now)
	d = d - time.Duration(float64(d)*c.Config.Factor)

	key := fmt.Sprintf("%s/%s", c.Config.Prefix, h.name)
	mutex := c.sync.NewMutex(key, redsync.WithExpiry(d), redsync.WithTries(1))
	if err := mutex.Lock(); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Cron) handle(h *Handler) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				c.Config.logger.Errorf("task panic:%s %s %s\n", h.cron, h.name, err)
			}
		}()
		s, err := c.lock(h)
		if err != nil {
			c.Config.logger.Errorf("can't run task:%s %s %s\n", h.cron, h.name, err.Error())
			return
		}
		if !s {
			c.Config.logger.Errorf("task skipped by another instance:%s %s\n", h.cron, h.name)
			return
		}
		h.handle()
	}
}
