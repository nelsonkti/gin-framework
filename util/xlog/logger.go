package xlog

import "go.uber.org/zap"

// A Logger represents a logger.
type Logger interface {
	With(a ...zap.Field) *Log
	Log(Level, ...any)
	Debug(...any)
	Debugf(string, ...any)
	Info(...any)
	Infof(string, ...any)
	Warn(...any)
	Warnf(string, ...any)
	Panic(...any)
	Panicf(string, ...any)
	Error(...any)
	Errorf(string, ...any)
	Fatal(...any)
	Fatalf(string, ...any)
}
