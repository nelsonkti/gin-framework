package mq

import (
	"go-framework/internal/mq/queues"
	"go-framework/util/mq/rocketmq"
)

func RegisterQueue(client *rocketmq.Client) {
	client.AddQueue(&queues.OrderQueue{}, &queues.ShopQueue{})
}
