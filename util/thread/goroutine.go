package thread

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Go(fn func()) {
	go runSafe(fn)
}

func runSafe(fn func()) {
	defer recoverLog()
	fn()
}

func recoverLog() {
	if p := recover(); p != nil {
		message := fmt.Sprintf("【错误信息】: %v；\n【错误内容】: %s\n", p, getStackTrace(3))
		log.Println(message)
		if notificationService != nil {
			env := notificationService.Env
			if notificationService.ServerNumber > 0 {
				env = fmt.Sprintf("%s%d", env, notificationService.ServerNumber)
			}
			message = fmt.Sprintf("【应用名称】: %s\n【环境】：%s\n%s；", notificationService.AppName, env, message)
			if err := notificationService.DingtalkRobot.SendText(message); err != nil {
				log.Printf("failed to send notification: %v", err)
			}
		}
	}
}

// getDetailedStackTrace 获取当前协程的详细堆栈跟踪信息
func getStackTrace(skip int) string {
	const depth = 32
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip, pcs)
	pcs = pcs[:n]

	var sb strings.Builder
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}
