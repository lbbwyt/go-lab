package main

import (
	"github.com/gin-gonic/gin"
	"go-lab/app/go_web/go_gin/middleware"
	limiter2 "go-lab/pkg/utils/limiter"
	"time"
)

func main() {
	// 1.创建路由
	// 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	limiter := limiter2.NewConnLimiter(1)
	// 注册中间件
	r.Use(middleware.LimitMiddleware(limiter))
	// {}为了代码规范
	{
		r.GET("/ce", func(c *gin.Context) {
			// 取值
			req := c.Query("req")
			time.Sleep(time.Second * 1000)
			// 页面接收
			c.JSON(200, gin.H{"request": req})
		})

	}
	r.Run()
}
