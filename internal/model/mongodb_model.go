package model

import (
	"go-framework/util/xsql/databese"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBModelImpl MongoDB数据库模型接口 【MongoDB】
type MongoDBModelImpl interface {
	DB() *databese.Engine
	Connection() string
	Table() string
	Model() *mongo.Collection
}

type MongoDBModel struct {
	BaseModel
}

func NewMongoDBModel(db *databese.Engine, database, tableName string) *MongoDBModel {
	return &MongoDBModel{*NewBaseModel(db, database, tableName)}
}

func (m *MongoDBModel) Model() *mongo.Collection {
	return m.DB().Mongo[m.Connection()].Collection(m.Table())
}
