package repository

import (
	"go-framework/internal/model/demo_model"
	"go-framework/internal/repository/demo_repository"
	"go-framework/util/xlog"
	"go-framework/util/xsql/databese"
)

type Container struct {
	DemoRepository        *demo_repository.DemoRepository
	DemoMongoDBRepository *demo_repository.DemoMongoDBRepository
}

func Register(db *databese.Engine, log *xlog.Log) *Container {
	return &Container{
		DemoRepository:        demo_repository.NewDemoRepository(demo_model.NewDemoModel(db), log),
		DemoMongoDBRepository: demo_repository.NewDemoMongoDBRepository(demo_model.NewDemoMongoDBModel(db), log),
	}
}
