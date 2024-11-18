package databese

import (
	"context"
	"fmt"
	"go-framework/util/types"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/transaction"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Engine struct {
	Gorm  map[string]*gorm.DB
	Mongo map[string]*mongo.Database
}

type DatabaseClient interface {
	Name() string
	Connect(c map[string]config.DBConfig)
	ConnType(database string) bool
	Result(c *Engine)
}

func (e *Engine) Close() {
	e.gormClose()
	e.mongodbClose()
}

func (e *Engine) gormClose() {
	for _, g := range e.Gorm {
		db, err := g.DB()
		if err != nil || db.Ping() != nil {
			fmt.Println(err)
			continue
		}
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (e *Engine) mongodbClose() {
	for _, m := range e.Mongo {
		err := m.Client().Disconnect(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// NewTransaction 创建一个新的事务上下文的 DBRepository 实例
func (r *Engine) NewTransaction(dbNames ...types.DB) (*transaction.Transaction, error) {
	txs := map[string]*gorm.DB{}
	for _, dbName := range dbNames {
		dbNameStr := string(dbName)
		txs[dbNameStr] = r.Gorm[dbNameStr].Begin()
	}

	return &transaction.Transaction{Tx: txs}, nil
}
