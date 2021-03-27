package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-lab/pkg/utils/limiter"
	"time"
)

const defaultMaxConn = 1

//全局限流
func LimitMiddleware(limiter *limiter.ConnLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		if !limiter.GetConn() {
			fmt.Println(fmt.Sprintf("当前连接数：%d", len(limiter.Bucket)))
			c.JSON(400, gin.H{"限流策略": len(limiter.Bucket)})
			c.Abort()
			return
		}
		c.Next()
		latency := time.Since(t).Milliseconds()
		fmt.Println(fmt.Sprintf("耗时%d毫秒", latency))
		limiter.ReleaseConn()
	}
}
