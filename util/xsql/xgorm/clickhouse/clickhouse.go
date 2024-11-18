package clickhouse

import (
	"fmt"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/log"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConnectionFormat 数据库连接的标准格式字符串
// const DBConnectionFormat = "clickhouse://gorm:gorm@localhost:9942/gorm?dial_timeout=10s&read_timeout=20s"
const DBConnectionFormat = "clickhouse://%s:%s@%s:%d/%s?dial_timeout=10s&read_timeout=20s"

type DB struct {
	config config.DBConfig
}

// NewDB 创建一个新的NewDB实例
func NewDB() *DB {
	return &DB{}
}

// Conn 初始化数据库连接
func (db *DB) Conn(config config.DBConfig) (*gorm.DB, error) {
	db.config = config
	dns := db.generateDSN(db.config.Host)

	gormDB, err := gorm.Open(clickhouse.Open(dns), &gorm.Config{
		Logger: log.NewLogger(logger.Default.LogMode(logger.Info)),
	})
	if err != nil {
		return nil, err
	}

	return gormDB, err
}

// generateDSN 生成数据库的DSN字符串
func (db *DB) generateDSN(host string) string {
	return fmt.Sprintf(DBConnectionFormat, db.config.Username, db.config.Password, host, db.config.Port, db.config.Database)
}
