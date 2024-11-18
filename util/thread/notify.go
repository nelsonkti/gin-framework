package thread

import "go-framework/util/dingtalk"

type Notifier struct {
	AppName       string
	Env           string
	ServerNumber  int
	DingtalkRobot *dingtalk.Robot
}

var notificationService *Notifier

func SetNotifier(appName, env string, serverNumber int, n *dingtalk.Robot) {
	notificationService = &Notifier{
		AppName:       appName,
		Env:           env,
		ServerNumber:  serverNumber,
		DingtalkRobot: n,
	}
}
