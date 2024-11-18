package router

import (
	"github.com/gin-gonic/gin"
	"go-framework/internal"
	"go-framework/internal/controller/demo_controller"
	"go-framework/internal/middleware"
)

func Register(app *gin.Engine, appCxt *internal.AppContent) {
	app.Use(
		middleware.OTELMiddleware(appCxt.Svc),
		middleware.RecoveryMiddleware(appCxt.Svc),
		middleware.RateLimiterMiddleware(),
	)

	app.GET("/demo", demo_controller.Demo(appCxt.Service))
}
