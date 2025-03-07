package xlog

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go-framework/util/xlog/zap_writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"path/filepath"
	"time"
)

const (
	DateLayout      = "2006-01-02"
	FileDateLayout  = "20060102"
	FileMonthLayout = "200601"
	TsDateLayout    = "2006-01-02 15:04:05.000"
)

// LogOptionFunc 核心选项
type LogOptionFunc func(*LogOption)

type LogOption struct {
	Filter       Filter
	GlobalFields []zap.Field
}

// NewLogger 创建并返回一个配置好的zap.Logger实例
func NewLogger(logDir, logName string, opts ...LogOptionFunc) *Log {

	logger := &Log{
		logDir:     logDir,
		logName:    logName,
		lastRotate: time.Now(),
	}

	// Apply log options
	logOption := &LogOption{}
	// 应用配置选项
	for _, opt := range opts {
		opt(logOption)
	}

	logger.RotateLogger(logOption) // 初始日志
	go logger.Maintain(logOption)  // 启动维护程序以处理轮换

	return logger
}

// RotateLogger handles the log rotation
func (l *Log) RotateLogger(logOption *LogOption) {
	l.Lock()
	defer l.Unlock()

	currentDate := time.Now().Format(FileDateLayout)
	currentMonthDir := time.Now().Format(FileMonthLayout)
	logPath := filepath.Join(l.logDir, currentMonthDir, fmt.Sprintf("%s_%s.log", l.logName, currentDate))
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename: logPath,
		MaxSize:  1024,
		MaxAge:   30, // days
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(TsDateLayout)
	core := &Core{
		Core: zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			w,
			zap.DebugLevel,
		),
		GlobalFields: []zap.Field{},
	}

	core.GlobalFields = logOption.GlobalFields

	l.Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l.filter = logOption.Filter
	l.lastRotate = time.Now()

	zapLogger := l.Logger.Sugar()
	log.SetFlags(0) // 取消内置 log 包的时间戳，避免重复
	log.SetOutput(zap_writer.NewWriter(zapLogger))
}

// Maintain 如果日期发生变化并触发轮换，请定期进行检查
func (l *Log) Maintain(logOption *LogOption) {
	ticker := time.NewTicker(time.Second) // Check every 10 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 获取当前时间
			now := time.Now()
			// 如果当前日期不是最后一次记录的日期，进行日志轮转
			if now.Add(time.Minute).Format(DateLayout) != l.lastRotate.Format(DateLayout) {
				l.RotateLogger(logOption)
			}
		}
	}
}

func With(l *Log) *Log {
	// 创建一个新的Log实例
	log := &Log{Logger: l.Logger}
	return log
}

// WithGlobalFields 添加全局字段到Core。
func WithGlobalFields(fields ...zap.Field) LogOptionFunc {
	return func(c *LogOption) {
		c.GlobalFields = append(c.GlobalFields, fields...)
	}
}

// WithFilters 添加字段过滤器到Core。
func WithFilters(opts ...FilterOption) LogOptionFunc {
	return func(c *LogOption) {
		options := Filter{
			Key:   make(map[interface{}]struct{}),
			Value: make(map[interface{}]struct{}),
		}
		for _, o := range opts {
			o(&options)
		}
		c.Filter = options
	}
}
