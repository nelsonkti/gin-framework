package demo_model

import (
	"go-framework/internal/model"
	"go-framework/util/xsql/databese"
)

// DemoModel 是一个示例模型
type DemoModel struct {
	model.DBModel
}

func NewDemoModel(db *databese.Engine) *DemoModel {
	return &DemoModel{*model.NewDBModel(db, "default", "test")}
}
