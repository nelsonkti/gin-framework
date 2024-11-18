package rocketmq

const (
	Separate  = "@"
	QueueMark = "app_name"
)

type MsgData struct {
	Topic   string
	GroupId string
	JobName string
	Data    interface{}
}
