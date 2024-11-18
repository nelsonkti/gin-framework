package zap_writer

import (
	"go.uber.org/zap"
	"io"
	"strings"
)

// zapWriter 是一个 io.Writer 的适配器，用于将日志写入 zap logger。
type zapWriter struct {
	logger *zap.SugaredLogger
}

// Write 将日志数据写入 zap logger。
func (w *zapWriter) Write(p []byte) (n int, err error) {
	logMessage := string(p)
	logMessage = strings.TrimSpace(logMessage)
	// 根据日志级别选择合适的 zap 方法
	if strings.Contains(logMessage, "error") {
		w.logger.Error(logMessage)
	} else if strings.Contains(logMessage, "warn") {
		w.logger.Warn(logMessage)
	} else if strings.Contains(logMessage, "info") {
		w.logger.Info(logMessage)
	} else {
		w.logger.Debug(logMessage)
	}
	return len(p), nil
}

// NewWriter 创建一个新的 zapWriter。
func NewWriter(logger *zap.SugaredLogger) io.Writer {
	return &zapWriter{logger: logger}
}
