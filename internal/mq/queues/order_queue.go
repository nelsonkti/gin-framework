package queues

import (
	job2 "go-framework/internal/mq/job"
	"go-framework/util/mq/queue"
)

var _ queue.Queue = (*OrderQueue)(nil)

type OrderQueue struct {
}

func (o *OrderQueue) Topic() string {
	return "Order"
}

func (o *OrderQueue) GroupId() string {
	return "GID_Order"
}

func (o *OrderQueue) Enqueue() []queue.Job {
	var jobs []queue.Job
	jobs = append(jobs, &job2.OrderJob{})

	return jobs
}
