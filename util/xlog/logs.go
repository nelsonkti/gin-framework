package xlog

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

var _ Logger = (*Log)(nil)

type Log struct {
	Logger *zap.Logger
	filter Filter
	sync.Mutex
	lastRotate time.Time
	logDir     string
	logName    string
}

func NewLog(log *zap.Logger, filter Filter) *Log {
	return &Log{Logger: log, filter: filter}
}

func (l *Log) With(a ...zap.Field) *Log {
	l.Logger = l.Logger.With(a...)
	return l
}

func (l *Log) Log(level Level, a ...any) {
	l.Logger.Sugar().Log(zapcore.Level(level), l.translateFields(a...))
}

func (l *Log) Debug(a ...any) {
	l.Logger.Sugar().Debug(l.translateFields(a...))
}

func (l *Log) Debugf(format string, a ...any) {
	l.Logger.Sugar().Debugf(l.translateFields(fmt.Sprintf(format, a...)))
}

func (l *Log) Info(a ...any) {
	l.Logger.Sugar().Info(l.translateFields(a...))
}

func (l *Log) Infof(format string, a ...any) {
	l.Logger.Sugar().Info(l.translateFields(fmt.Sprintf(format, a...)))
}

func (l *Log) Warn(a ...any) {
	l.Logger.Sugar().Warn(l.translateFields(a...))
}

func (l *Log) Warnf(format string, a ...any) {
	l.Logger.Sugar().Warn(l.translateFields(fmt.Sprintf(format, a...)))
}

func (l *Log) Panic(a ...any) {
	l.Logger.Sugar().Panic(a...)
}

func (l *Log) Panicf(format string, a ...any) {
	l.Logger.Sugar().Panic(l.translateFields(fmt.Sprintf(format, a...)))
}

func (l *Log) Error(a ...any) {
	l.Logger.Sugar().Error(l.translateFields(a...))
}

func (l *Log) Errorf(format string, a ...any) {
	l.Logger.Sugar().Error(l.translateFields(fmt.Sprintf(format, a...)))
}

func (l *Log) Fatal(a ...any) {
	l.Logger.Sugar().Fatal(l.translateFields(a...))
}

func (l *Log) Fatalf(format string, a ...any) {
	l.Logger.Sugar().Fatal(l.translateFields(fmt.Sprintf(format, a...)))
}

func (l *Log) translateFields(a ...any) string {
	fields := l.filter.FilterFiled(a...)
	fieldLen := len(fields)
	if fieldLen == 0 {
		return ""
	}
	if fieldLen == 1 {
		return fmt.Sprintf("%v", fields[0])
	}
	if (fieldLen & 1) == 1 {
		fields = append(fields, "UNPAIRED")
	}

	var spacer string
	buf := new(bytes.Buffer)
	for i := 0; i < len(fields); i += 2 {
		key := i + 1
		_, _ = fmt.Fprintf(buf, "%s%s=%v", spacer, fields[i], fields[key])
		spacer = " "
	}

	return buf.String()
}
