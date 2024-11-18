package internal

import (
	"go-framework/internal/container/service"
	"go-framework/internal/server"
)

type AppContent struct {
	Svc     *server.SvcContext
	Service *service.Container
}

func Register(svc *server.SvcContext) *AppContent {
	return &AppContent{
		Svc:     svc,
		Service: service.Register(svc),
	}
}
