package demo_repository

import (
	"go-framework/internal/model/demo_model"
	"go-framework/internal/repository"
	"go-framework/util/xlog"
)

type DemoMongoDBRepository struct {
	*repository.MongoDBRepository
}

func NewDemoMongoDBRepository(model *demo_model.DemoMongoDBModel, log *xlog.Log) *DemoMongoDBRepository {
	return &DemoMongoDBRepository{repository.NewMongoDBRepository(model, log)}
}
