package postgresql

import (
	"fmt"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConnectionFormat 数据库连接的标准格式字符串
// const DBConnectionFormat = "host=localhost order_model=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
const DBConnectionFormat = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai"

type DB struct {
	dns    string
	config config.DBConfig
}

// NewDB 创建一个新的NewDB实例
func NewDB() *DB {
	return &DB{}
}

// Conn 初始化数据库连接
func (db *DB) Conn(config config.DBConfig) (*gorm.DB, error) {
	db.config = config
	gormDB, err := gorm.Open(postgres.Open(db.generateDSN(db.config.Host)), &gorm.Config{
		Logger: log.NewLogger(logger.Default.LogMode(logger.Info)),
	})
	if err != nil {
		return nil, err
	}

	return gormDB, err
}

// generateDSN 生成数据库的DSN字符串
func (db *DB) generateDSN(host string) string {
	dns := fmt.Sprintf(DBConnectionFormat, host, db.config.Username, db.config.Password, db.config.Database, db.config.Port)
	return dns
}
