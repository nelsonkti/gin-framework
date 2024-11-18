package demo_repository

import (
	"go-framework/internal/model/demo_model"
	"go-framework/internal/repository"
	"go-framework/util/xlog"
)

type DemoRepository struct {
	*repository.DBRepository
}

func NewDemoRepository(model *demo_model.DemoModel, log *xlog.Log) *DemoRepository {
	return &DemoRepository{repository.NewDBRepository(model, log)}
}
