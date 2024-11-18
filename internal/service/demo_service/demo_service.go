package demo_service

import (
	"context"
	"go-framework/internal/server"
)

type DemoServiceImpl interface {
	Demo(c context.Context) (interface{}, error)
}

type DemoService struct {
	svc *server.SvcContext
}

func NewDemoService(svc *server.SvcContext) *DemoService {
	return &DemoService{svc: svc}
}

func (s *DemoService) Demo(ctx context.Context) (interface{}, error) {

	return make(map[string]interface{}), nil
}
