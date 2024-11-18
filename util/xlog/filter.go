package xlog

import (
	"go.uber.org/zap/zapcore"
)

// FilterOption is filter option.
type FilterOption func(*Filter)

const fuzzyStr = "***"

type Filter struct {
	Key   map[interface{}]struct{}
	Value map[interface{}]struct{}
}

// FilterKey with filter key.
func FilterKey(key ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range key {
			o.Key[v] = struct{}{}
		}
	}
}

// FilterValue with filter value.
func FilterValue(value ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range value {
			o.Value[v] = struct{}{}
		}
	}
}

// Check 判断给定的字段列表是否满足过滤条件
func (f Filter) Check(fields []zapcore.Field) bool {
	for _, field := range fields {
		if _, ok := f.Key[field.Key]; ok {
			field.Key = fuzzyStr
		}

		if _, ok := f.Value[field.Key]; ok {
			field.Key = fuzzyStr
		}
	}
	return false
}

// Apply 如果日志条目包含指定的键，则将其值替换为星号。
func (f *Filter) Apply(fields []zapcore.Field) []zapcore.Field {
	for _, field := range fields {
		if _, ok := f.Key[field.Key]; ok {
			field.Key = fuzzyStr
		}

		if _, ok := f.Value[field.Key]; ok {
			field.Key = fuzzyStr
		}
	}
	return fields
}

func (f *Filter) FilterFiled(fields ...any) []interface{} {
	if len(f.Key) > 0 || len(f.Value) > 0 {
		for i := 0; i < len(fields); i += 2 {
			v := i + 1
			if v >= len(fields) {
				continue
			}
			if _, ok := f.Key[fields[i]]; ok {
				fields[v] = fuzzyStr
			}
			if _, ok := f.Value[fields[v]]; ok {
				fields[v] = fuzzyStr
			}
		}
	}
	return fields
}
