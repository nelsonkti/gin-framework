package rocketmq

import (
	"context"
	"fmt"
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"go-framework/util/helper"
	"go-framework/util/mq/queue"
	"time"
)

const (
	GroupId = ""
)

type Producer struct {
	client *Client
}

func NewProducer(client *Client) *Producer {
	return &Producer{
		client: client,
	}
}

// SendJobMessage 发送任务消息
func (p *Producer) SendJobMessage(ctx context.Context, job queue.Job, msg interface{}) error {
	topic, msgd, err := p.client.Decoder.Marshal(job, msg)
	if err != nil {
		p.client.Logger.Errorf("SendJobMessage marshal job error, %v", err)
		return err
	}

	return p.sendMessage(ctx, topic, GroupId, msgd)
}

// SendJobDelayMessage 发送延时任务消息
func (p *Producer) SendJobDelayMessage(ctx context.Context, job queue.Job, msg interface{}, duration time.Duration) error {
	topic, msgd, err := p.client.Decoder.Marshal(job, msg)
	if err != nil {
		p.client.Logger.Errorf("SendJobDelayMessage marshal Decoder job error, %v", err)
		return err
	}

	return p.sendDelayMessage(ctx, topic, GroupId, msgd, duration)
}

// SendMessage 发送消息
func (p *Producer) SendMessage(ctx context.Context, topic string, groupId string, msg interface{}) error {
	marshalMsg, err := helper.Marshal(msg)
	if err != nil {
		p.client.Logger.Errorf("发送普通消息异常：%+v, %+v, %+v", topic, msg, err)
		return err
	}

	return p.sendMessage(ctx, topic, groupId, string(marshalMsg))
}

func (p *Producer) sendMessage(ctx context.Context, topic string, groupId string, marshalMsg string) error {
	msgRequest, err := p.publishMessageRequest(topic, groupId, marshalMsg)
	if err != nil {
		return err
	}

	return p.publishMessage(ctx, msgRequest, topic)
}

// SendDelayMessage 发送延时消息
func (p *Producer) SendDelayMessage(ctx context.Context, topic string, groupId string, msg interface{}, duration time.Duration) error {
	marshalMsg, err := helper.Marshal(msg)
	if err != nil {
		p.client.Logger.Errorf("发送延时消息异常：%+v, %+v, %+v", topic, msg, err)
		return err
	}

	return p.sendDelayMessage(ctx, topic, groupId, string(marshalMsg), duration)
}

func (p *Producer) sendDelayMessage(ctx context.Context, topic string, groupId string, marshalMsg string, duration time.Duration) error {
	msgRequest, err := p.publishMessageRequest(topic, groupId, marshalMsg)
	if err != nil {
		return err
	}

	msgRequest.StartDeliverTime = time.Now().Add(duration).Unix() * 1000

	err = p.publishMessage(ctx, msgRequest, topic)
	return err
}

// publishMessage 发送消息
func (p *Producer) publishMessage(ctx context.Context, msgRequest mq_http_sdk.PublishMessageRequest, topic string) error {
	producer := p.client.Client().GetProducer(p.client.conf.MQ.Namespace, topic)
	res, err := producer.PublishMessage(msgRequest)

	p.log(ctx, msgRequest, res, err)
	return err
}

// publishMessageRequest 封装消息
func (p *Producer) publishMessageRequest(topic, groupId string, msg string) (mq_http_sdk.PublishMessageRequest, error) {
	msgRequest := mq_http_sdk.PublishMessageRequest{}

	msgRequest.Properties = make(map[string]string)
	msgRequest.MessageBody = msg //消息内容。
	msgRequest.MessageKey = topic
	if groupId == "" {
		groupId = topic
	}
	msgRequest.Properties["groupId"] = p.client.GetGroupName(groupId)
	return msgRequest, nil
}

// log 日志
func (p *Producer) log(ctx context.Context, msgRequest mq_http_sdk.PublishMessageRequest, msgResponse mq_http_sdk.PublishMessageResponse, err error) {
	if err != nil {
		p.client.ErrorNotify(ctx, fmt.Sprintf("【队列生产者】发送异常：\n 错误信息:\n %+v \n请求数据：\n %+v \n, 返回数据：\n %+v \n", err, msgRequest, msgResponse))
	} else {
		p.client.Logger.Infof("队列发送成功：请求数据：\n%+v \n, 返回数据：\n%+v \n", msgRequest, msgResponse)
	}
}
