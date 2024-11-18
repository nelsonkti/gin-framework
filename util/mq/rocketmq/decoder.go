package rocketmq

import "go-framework/util/mq/queue"

type Decoder interface {
	Marshal(job queue.Job, msg interface{}) (string, string, error)
	UnMarshal(msg string) (queue.Job, []byte, error)
	Check(msg string) bool
}
