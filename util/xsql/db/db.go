package db

import (
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/databese"
	"go-framework/util/xsql/mongodb"
	"go-framework/util/xsql/xgorm"
)

type DB struct {
	C           map[string]config.DBConfig
	config      map[string]map[string]config.DBConfig
	gormConfig  map[string]config.DBConfig
	MongoConfig map[string]config.DBConfig
	Client      map[string]databese.DatabaseClient
	engine      *databese.Engine
}

func NewDB(c map[string]config.DBConfig) *DB {
	return &DB{
		engine: &databese.Engine{},
		C:      c,
		Client: make(map[string]databese.DatabaseClient),
	}
}

func (db *DB) Register() {
	db.register(xgorm.NewGorm())
	db.register(mongodb.NewMongoDB())
}

func (db *DB) register(database databese.DatabaseClient) {
	db.Client[database.Name()] = database
}

func (db *DB) InitDatabases() (*databese.Engine, error) {
	db.Register()
	configs := db.loadConfig()

	for name, client := range db.Client {
		if configs[name] == nil {
			continue
		}
		client.Connect(configs[name])
		client.Result(db.engine)
	}

	return db.engine, nil
}

func (db *DB) loadConfig() map[string]map[string]config.DBConfig {
	configs := make(map[string]map[string]config.DBConfig)
	for name, client := range db.Client {
		driver := make(map[string]config.DBConfig)
		for key, dbConfig := range db.C {
			if client.ConnType(dbConfig.Driver) {
				driver[key] = dbConfig
				configs[name] = driver
			}
		}
	}
	return configs
}
