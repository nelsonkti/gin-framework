package grpc

import (
	"context"
	"go-framework/config"
)

type Container struct {
}

func Register(c config.Conf, ctx context.Context) *Container {
	return &Container{}
}
