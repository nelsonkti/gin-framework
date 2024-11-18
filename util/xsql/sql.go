package xsql

import (
	"encoding/json"
	"go-framework/internal/common/tool/dingtalk_tool"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/databese"
	"go-framework/util/xsql/db"
	"go-framework/util/xsql/log"
)

func NewClient(c interface{}) *databese.Engine {
	cByte, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	databases := make(map[string]config.DBConfig)
	err = json.Unmarshal(cByte, &databases)
	if err != nil {
		panic(err)
	}

	engine, err := db.NewDB(databases).InitDatabases()
	if err != nil {
		panic(err)
	}

	return engine
}

// SetNotifier 设置钉钉通知
func SetNotifier(dingtalkTool *dingtalk_tool.Dingtalk) {
	log.DingtalkTool = dingtalkTool
}
