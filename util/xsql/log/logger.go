package log

import (
	"context"
	"errors"
	"fmt"
	"go-framework/internal/common/tool/dingtalk_tool"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"runtime"
	"time"
)

var DingtalkTool *dingtalk_tool.Dingtalk

type Logger struct {
	gormLogger logger.Interface
}

func NewLogger(logger logger.Interface) *Logger {
	return &Logger{gormLogger: logger}
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.gormLogger = l.gormLogger.LogMode(level)
	return &newLogger
}

func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.gormLogger.Info(ctx, msg, data...)
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.gormLogger.Warn(ctx, msg, data...)
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.gormLogger.Error(ctx, msg, data...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	_, file, line, _ := runtime.Caller(4)
	fmt.Printf("%s %s:%d", time.Now().Format(time.DateTime), file, line)
	l.gormLogger.Trace(ctx, begin, fc, err)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		buf := make([]byte, 2048)
		n := runtime.Stack(buf, false)
		pc, file2, line2, _ := runtime.Caller(2)
		fn := runtime.FuncForPC(pc)
		msg := fmt.Sprintf("error message: %+v; \nline: %s:%d; function: %s; \nstackTrace: %s", err, file2, line2, fn.Name(), buf[:n])

		if DingtalkTool != nil {
			_ = DingtalkTool.SendAlarm(ctx, msg)
		} else {
			log.Printf(msg)
		}
	}
}
