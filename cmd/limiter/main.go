package main

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

var limiter = ratelimit.NewBucketWithQuantum(time.Second, 100, 10)

func TokenRateLimiter() gin.HandlerFunc {

	return func(context *gin.Context) {
		if limiter.TakeAvailable(1) == 0 {
			context.AbortWithStatusJSON(http.StatusTooManyRequests, "too many request")
		} else {
			context.Next()
		}
	}
}

func main() {
	e := gin.Default()
	e.GET("/test", TokenRateLimiter(), func(context *gin.Context) {
		context.JSON(200, true)
	})
	e.Run(":7777")
}
