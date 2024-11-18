package locker

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type Mutex struct {
	sync  *redsync.Redsync
	mutex *redsync.Mutex
}

func NewMutex(redis *redis.Client) *Mutex {
	pool := goredis.NewPool(redis)
	rs := redsync.New(pool)
	return &Mutex{sync: rs}
}

// Lock acquires a distributed lock for the given key and duration.
func (m *Mutex) Lock(key string, opts ...redsync.Option) error {
	m.mutex = m.sync.NewMutex(key, opts...)
	return m.mutex.Lock()
}

// Unlock 解锁
func (m *Mutex) Unlock() (bool, error) {
	return m.mutex.Unlock()
}
