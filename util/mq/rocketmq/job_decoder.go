package rocketmq

import (
	"errors"
	"fmt"
	"go-framework/util/helper"
	"go-framework/util/mq/queue"
	"strings"
)

type JobDecoder struct {
	client *Client
}

func NewJobDecoder(client *Client) *JobDecoder {
	return &JobDecoder{client: client}
}

func (jd *JobDecoder) Marshal(job queue.Job, msg interface{}) (string, string, error) {
	var msgData MsgData
	jobName := job.Name()
	queueJob := jd.client.Jobs[jobName]
	if queueJob.Queue == nil {
		return "", "", errors.New("queue is nil")
	}
	queueInfo := jd.client.Jobs[jobName].Queue

	topic := queueInfo.Topic()
	groupId := queueInfo.GroupId()
	msgData.Topic = topic
	msgData.GroupId = groupId
	msgData.JobName = jobName
	msgData.Data = msg

	marshal, err := helper.Marshal(msgData)
	if err != nil {
		return "", "", err
	}

	return topic, fmt.Sprintf("%s%s%s", QueueMark, Separate, string(marshal)), nil
}

func (jd *JobDecoder) Check(msg string) bool {
	return strings.Contains(msg, QueueMark+Separate)
}

func (jd *JobDecoder) UnMarshal(msg string) (queue.Job, []byte, error) {
	if msg == "" {
		return nil, nil, fmt.Errorf("msg is empty")
	}
	subIndex := strings.Index(msg, Separate)
	body := msg[subIndex+1:]
	var msgData MsgData
	err := helper.UmMarshal([]byte(body), &msgData)
	if err != nil {
		return nil, nil, err
	}
	var res []byte
	if msgData.Data != nil {
		res, err = helper.Marshal(msgData.Data)
		if err != nil {
			return nil, nil, err
		}
	}

	queueJob := jd.client.Jobs[msgData.JobName]
	if queueJob.Queue == nil {
		return nil, nil, errors.New("queue is nil")
	}

	job := jd.client.Jobs[msgData.JobName].Job

	return job, res, nil
}
