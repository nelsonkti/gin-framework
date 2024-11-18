package middleware

import (
	"github.com/gin-gonic/gin"
	"go-framework/pkg/aegis/ratelimit"
	"go-framework/pkg/aegis/ratelimit/bbr"
	"go-framework/util/xhttp"
	"net/http"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allow, err := bbr.NewLimiter().Allow()

		if err != nil {
			c.JSON(http.StatusOK, xhttp.Error(err))
			c.Abort()
			return
		}

		defer func() {
			if err == nil {
				allow(ratelimit.DoneInfo{})
			}
		}()

		c.Next()
	}
}
