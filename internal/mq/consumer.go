package mq

import "go-framework/util/mq/rocketmq"

func ConsumerHandler(client *rocketmq.Client) {
	//rocketmq.ConsumerMessage(client, &queues.OrderQueue{}, rocketmq.WithConcurrency(10), rocketmq.WithRetryTimes(3)) // 订单队列消费者
	//rocketmq.ConsumerMessage(client, &queues.OrderQueue{}, rocketmq.WithConcurrency(10), rocketmq.WithRetryTimes(0)) // 订单队列消费者
	//rocketmq.ConsumerMessage(client, &queues.ShopQueue{}, rocketmq.WithConcurrency(3))                               // 商家队列消费者
}
