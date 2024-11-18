package queue

type Queue interface {
	Topic() string
	GroupId() string
	Enqueue() []Job
}

type Job interface {
	Name() string
	Execute([]byte) error
}
