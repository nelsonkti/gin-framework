package mysql

import (
	"fmt"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/opentelemetry/tracing"
)

// 数据库连接的标准格式字符串
const DBConnectionFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"

// DB 管理MySQL数据库连接，支持读写分离
type DB struct {
	config          config.DBConfig
	primaryDialect  gorm.Dialector
	sourceDialects  []gorm.Dialector
	replicaDialects []gorm.Dialector
}

// NewDB 创建一个新的NewDB实例
func NewDB() *DB {
	return &DB{}
}

// Conn 初始化数据库连接
func (db *DB) Conn(config config.DBConfig) (*gorm.DB, error) {
	db.config = config

	db.setupDialects()

	gormDB, err := gorm.Open(db.primaryDialect, &gorm.Config{
		Logger: log.NewLogger(logger.Default.LogMode(logger.Info)),
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = gormDB.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	if err != nil {
		return nil, err
	}

	db.configureResolver(gormDB)

	return gormDB, nil
}

// setupDialects 配置主数据库和从数据库的方言
func (db *DB) setupDialects() {
	if db.config.Host != "" {
		db.primaryDialect = mysql.Open(db.generateDSN(db.config.Host))
	}

	// 多主
	if len(db.config.Sources) > 0 {
		if db.primaryDialect == nil {
			db.primaryDialect = mysql.Open(db.generateDSN(db.config.Sources[0]))
			db.sourceDialects = db.dialectsFromHosts(db.config.Sources[1:])
		} else {
			db.sourceDialects = db.dialectsFromHosts(db.config.Sources)
		}
	}

	// 多从
	if len(db.config.Replicas) > 0 {
		db.replicaDialects = db.dialectsFromHosts(db.config.Replicas)
	}
}

// generateDSN 生成数据库的DSN字符串
func (db *DB) generateDSN(host string) string {
	return fmt.Sprintf(DBConnectionFormat, db.config.Username, db.config.Password, host, db.config.Port, db.config.Database)
}

// dialectsFromHosts 从主机地址列表生成方言列表
func (db *DB) dialectsFromHosts(hosts []string) []gorm.Dialector {
	var dialects []gorm.Dialector
	for _, host := range hosts {
		dialects = append(dialects, mysql.Open(db.generateDSN(host)))
	}
	return dialects
}

// configureResolver 配置数据库解析器，支持读写分离
func (db *DB) configureResolver(gormDB *gorm.DB) {
	if len(db.sourceDialects) > 0 && len(db.replicaDialects) > 0 {
		return
	}
	resolverConfig := dbresolver.Config{
		Sources:  db.sourceDialects, // 假设至少有一个主数据库
		Replicas: db.replicaDialects,
		Policy:   dbresolver.RandomPolicy{},
	}
	gormDB.Use(dbresolver.Register(resolverConfig))
}
