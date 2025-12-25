package middlewares

import (
	"golang_per_day_30/internal/utils/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/spf13/viper"
)

func RateLimitMiddleware() gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(time.Second, viper.GetInt64(`ratelimit.cap`), viper.GetInt64(`ratelimit.quantum`))
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Error(`rate limit ...`, http.StatusForbidden, c)
			c.Abort()
			return
		}
		c.Next()
	}
}
