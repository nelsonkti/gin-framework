package rocketmq

import (
	"fmt"
	"strings"
)

const (
	groupPrefix = "GID"
)

func (c *Client) GetGroupName(topic string) string {
	if c.queues[topic] == nil || c.queues[topic].Topic() == "" {
		if strings.Contains(topic, groupPrefix) {
			return c.GetGroupNameByGroupId(topic)
		}
		return c.GetGroupNameByGroupId(fmt.Sprintf("%s_%s", groupPrefix, topic))
	}

	return c.GetGroupNameByGroupId(c.queues[topic].GroupId())
}

func (c *Client) GetGroupNameByGroupId(groupId string) string {
	if c.conf.MQ.Env == "" || c.conf.App.Env == "prod" {
		return groupId
	}

	return fmt.Sprintf("%s_%s", groupId, c.conf.MQ.Env)
}
