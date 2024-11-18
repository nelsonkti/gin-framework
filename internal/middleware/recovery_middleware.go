package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-framework/internal/server"
	"go-framework/util/xhttp"
	"net/http"
	"runtime"
)

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware(svc *server.SvcContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				pc, file, line, _ := runtime.Caller(2)
				fn := runtime.FuncForPC(pc)
				errMsg := fmt.Sprintf("message: %+v; \nline: %s:%d; function: %s; \nstackTrace: %s", err, file, line, fn.Name(), buf[:n])
				svc.Logger.Errorf(errMsg)

				err2 := svc.Tool.DingtalkTool.SendAlarm(c.Request.Context(), errMsg)
				if err2 != nil {
					svc.Logger.Errorf("Error sending alarm %+v", err2)
				}

				if svc.Conf.App.Env == "local" {
					fmt.Println(errMsg)
				}
				c.AbortWithStatusJSON(http.StatusOK, xhttp.ErrMsg("系统异常", 500).SetCode(500))
			}
		}()

		c.Next()
	}
}
