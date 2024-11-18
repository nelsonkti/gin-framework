package middleware

import (
	"go-framework/internal/server"
	"google.golang.org/grpc"
)

type Middleware func(svc *server.SvcContext) grpc.UnaryServerInterceptor
