package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Core struct {
	zapcore.Core
	GlobalFields []zap.Field
}

func (c *Core) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	// 判断日志条目的级别是否足够高，能够通过这个Core的级别检查
	if c.Enabled(ent.Level) {
		// 如果日志条目的级别足够高，调用底层Core的Check方法
		return ce.AddCore(ent, c)
		//return c.Core.Check(ent, ce)
	}
	// 如果日志条目的级别不够高，或底层Core决定不记录这个条目，返回nil
	return nil
}

func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := append(fields, c.translateGlobalFields()...)

	return c.Core.Write(entry, allFields)
}

func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	newGlobalFields := make([]zap.Field, len(c.GlobalFields), len(c.GlobalFields)+len(fields))
	copy(newGlobalFields, c.GlobalFields)
	for _, f := range fields {
		newGlobalFields = append(newGlobalFields, zap.Field{Key: f.Key, Type: f.Type, Integer: f.Integer, String: f.String, Interface: f.Interface})
	}
	return &Core{
		Core:         c.Core,
		GlobalFields: newGlobalFields,
	}
}

func (c *Core) translateGlobalFields() []zapcore.Field {
	fields := make([]zapcore.Field, len(c.GlobalFields))
	for i, f := range c.GlobalFields {
		fields[i] = zapcore.Field{
			Key:       f.Key,
			Type:      f.Type,
			Integer:   f.Integer,
			String:    f.String,
			Interface: f.Interface,
		}
	}
	return fields
}
