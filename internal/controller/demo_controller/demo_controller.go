package demo_controller

import (
	"github.com/gin-gonic/gin"
	"go-framework/internal/container/service"
	"go-framework/util/xhttp"
	"net/http"
)

func Demo(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		//var param validation.UserRequest
		// 接收请求参数
		//if err := c.ShouldBindJSON(&param); err != nil {
		//	c.JSON(http.StatusOK, xhttp.Error(err))
		//	return
		//}
		res, err := svc.DemoService.Demo(c.Request.Context())
		//res, err := svc.CmdEventService.GetOne()
		if err != nil {
			c.JSON(http.StatusOK, xhttp.Error(err))
			return
		}
		c.JSON(http.StatusOK, xhttp.Data(res))
	}
}
