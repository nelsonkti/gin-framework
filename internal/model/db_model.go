package model

import (
	"fmt"
	"go-framework/util/xsql/databese"
	"gorm.io/gorm"
)

// DBModelImpl 数据库模型接口【mysql、Clickhouse、Hologres】
type DBModelImpl interface {
	DB() *databese.Engine
	Connection() string
	Table() string
	Model() *gorm.DB
}

func NewDBModel(db *databese.Engine, database, tableName string) *DBModel {
	return &DBModel{*NewBaseModel(db, database, tableName)}
}

type DBModel struct {
	BaseModel
}

func (m *DBModel) Model() *gorm.DB {
	if m.DB().Gorm[m.Connection()] == nil {
		panic(fmt.Sprintf("db【%s】connection is not initialized", m.Connection()))
	}
	return m.DB().Gorm[m.Connection()].Table(m.Table())
}
