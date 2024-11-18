package dingtalk_tool

import (
	"context"
	"errors"
	"fmt"
	"go-framework/config"
	"go-framework/internal/data/message_data/robot_message_data"
	"go-framework/util/dingtalk"
	"go.opentelemetry.io/otel/trace"
)

type Dingtalk struct {
	env        string
	AlarmRobot *dingtalk.Robot
}

func NewDingtalkTool(c config.Conf) *Dingtalk {
	return &Dingtalk{
		env:        c.App.Env,
		AlarmRobot: dingtalk.NewRobot(RobotURL, c.Dingtalk.Robots.AlarmSecret),
	}
}

// SendAlarm 发送告警消息
func (d *Dingtalk) SendAlarm(ctx context.Context, msg string) error {
	if d.env == "local" {
		return nil
	}

	if d.AlarmRobot == nil {
		return errors.New("AlarmRobot is nil")
	}

	var domain string
	if host, ok := ctx.Value("request_host").(string); ok {
		domain = host
	}

	var traceId string
	span := trace.SpanFromContext(ctx)
	if span != nil {
		traceId = span.SpanContext().TraceID().String()
	}

	return d.AlarmRobot.SendText(fmt.Sprintf(robot_message_data.AlarmMessage, "【系统异常】", d.env, domain, traceId, msg))
}
