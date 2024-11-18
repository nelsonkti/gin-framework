package rocketmq

import (
	"errors"
	"go-framework/util/mq/queue"
)

type QueueJob struct {
	Queue queue.Queue
	Job   queue.Job
}

// AddQueue 添加队列
func (c *Client) AddQueue(queues ...queue.Queue) {
	for _, q := range queues {
		if c.queues[q.Topic()] != nil {
			c.Logger.Panicf("%s queue already", q.Topic())
		}
		c.queues[q.Topic()] = q
	}
}

// RegisterJob 注册队列
func (c *Client) RegisterJob() error {
	queueJobs := make(map[string]*QueueJob)
	for _, queueInfo := range c.queues {
		for _, jobInfo := range queueInfo.Enqueue() {
			key := jobInfo.Name()
			if _, exists := queueJobs[key]; exists {
				return errors.New("failed to register: duplicate job found for key " + key)
			}
			queueJobs[key] = &QueueJob{
				Queue: queueInfo,
				Job:   jobInfo,
			}
		}
	}
	c.Jobs = queueJobs
	return nil
}
