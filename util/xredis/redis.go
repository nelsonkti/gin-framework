package xredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisClient struct {
	client map[string]*redis.Client
}

// Default 获取默认Redis客户端
func (c *RedisClient) Default() *redis.Client {
	if c.client["default"] == nil {
		log.Panic("RedisClient: default client not found")
	}
	return c.client["default"]
}

// Conn 获取指定别名的Redis客户端
func (c *RedisClient) Conn(name string) *redis.Client {
	if c.client[name] == nil {
		log.Panicf("RedisClient: client %s not found", name)
	}
	return c.client[name]
}

// NewClient 初始化多个Redis客户端
func NewClient(c interface{}) *RedisClient {
	cByte, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	var configs []Config
	err = json.Unmarshal(cByte, &configs)
	if err != nil {
		panic(err)
	}

	clients := make(map[string]*redis.Client)
	for _, v := range configs {
		add := fmt.Sprintf("%s:%d", v.Host, v.Port)
		options := &redis.Options{
			Addr:         add,
			Username:     v.UserName,
			Password:     v.Password,
			DB:           v.Database,
			PoolSize:     10,
			MinIdleConns: 3,
			WriteTimeout: 3 * time.Second,
		}
		clients[v.Alias], err = connect(options)
		if err != nil {
			panic(err)
		}
	}

	return &RedisClient{client: clients}
}

func connect(options *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(options)

	// 测试连接是否有效
	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		// 若某个实例连接失败，需根据需求决定是立即返回错误还是继续尝试连接其他实例
		return nil, err
	}
	if pong != "PONG" {
		return nil, errors.New("unexpected PONG response")
	}
	return client, nil
}

func (c *RedisClient) Close() {
	for _, client := range c.client {
		err := client.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
