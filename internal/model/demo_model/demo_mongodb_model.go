package demo_model

import (
	"go-framework/internal/model"
	"go-framework/util/xsql/databese"
)

// DemoMongoDBModel 是一个示例模型
type DemoMongoDBModel struct {
	model.MongoDBModel
}

func NewDemoMongoDBModel(db *databese.Engine) *DemoMongoDBModel {
	return &DemoMongoDBModel{*model.NewMongoDBModel(db, "default", "test")}
}
