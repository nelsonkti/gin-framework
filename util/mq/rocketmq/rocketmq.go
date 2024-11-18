package rocketmq

import (
	"context"
	"fmt"
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/go-redis/redis/v8"
	"go-framework/internal/common/tool/dingtalk_tool"
	"go-framework/util/helper"
	"go-framework/util/mq/queue"
	"go-framework/util/xlog"
)

type clientHandler func(client *Client)

type mqConf struct {
	App App
	MQ  config
}

type App struct {
	Name         string `json:"name"`          // 应用名称
	Env          string `json:"env"`           // 环境
	ServerNumber int    `json:"server_number"` // 服务器编号
}

type config struct {
	Endpoint  []string `json:"endpoint"`
	AccessKey string   `json:"access_key"`
	SecretKey string   `json:"secret_key"`
	Namespace string   `json:"namespace"`
	Env       string   `json:"env"`
}

type Client struct {
	conf         *mqConf
	Logger       *xlog.Log
	Producer     *Producer
	redisClient  *redis.Client
	dingtalkTool *dingtalk_tool.Dingtalk
	queues       map[string]queue.Queue
	Jobs         map[string]*QueueJob
	Decoder      Decoder
}

func NewClient(c interface{}, logger *xlog.Log, redisClient *redis.Client, fs ...clientHandler) (client *Client) {
	var conf *mqConf
	err := helper.UnMarshalWithInterface(c, &conf)
	if err != nil {
		logger.Panicf("rocketmq config error: %v", err)
	}
	client = &Client{
		conf:        conf,
		Logger:      logger,
		redisClient: redisClient,
		queues:      make(map[string]queue.Queue),
		Jobs:        make(map[string]*QueueJob),
	}

	for _, f := range fs {
		f(client)
	}

	client.Producer = NewProducer(client)

	err = client.RegisterJob()
	if err != nil {
		client.Logger.Panicf("register job error: %v", err)
	}

	client.Decoder = NewJobDecoder(client)
	return
}

// ConsumerRun 启动消费者
func (c *Client) ConsumerRun(handler func(client *Client)) {
	handler(c)
}

func (c *Client) SetNotifier(dingtalkTool *dingtalk_tool.Dingtalk) {
	c.dingtalkTool = dingtalkTool
}

func (c *Client) ErrorNotify(ctx context.Context, message string) {
	c.Logger.Errorf(message)
	if c.dingtalkTool == nil {
		return
	}
	err := c.dingtalkTool.SendAlarm(ctx, fmt.Sprintf("告警信息：\n%s\n", message))
	if err != nil {
		c.Logger.Errorf("【队列】钉钉机器人发送失败: %+v", err)
	}
}

// Client 创建阿里云客户端
func (c *Client) Client() mq_http_sdk.MQClient {
	client := mq_http_sdk.NewAliyunMQClient(c.conf.MQ.Endpoint[0], c.conf.MQ.AccessKey, c.conf.MQ.SecretKey, "")
	return client
}
