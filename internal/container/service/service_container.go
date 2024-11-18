package service

import (
	"go-framework/internal/server"
	"go-framework/internal/service/demo_service"
)

type Container struct {
	DemoService demo_service.DemoServiceImpl
}

func Register(svc *server.SvcContext) *Container {
	return &Container{
		DemoService: demo_service.NewDemoService(svc),
	}
}
