package model

import "go-framework/util/xsql/databese"

// BaseModel 基础模型
type BaseModel struct {
	db        *databese.Engine
	database  string
	tableName string
}

func NewBaseModel(db *databese.Engine, database, tableName string) *BaseModel {
	return &BaseModel{db, database, tableName}
}

func (m *BaseModel) DB() *databese.Engine {
	return m.db
}

func (m *BaseModel) Connection() string {
	return m.database
}

func (m *BaseModel) Table() string {
	return m.tableName
}
