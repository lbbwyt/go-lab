package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

//操作日志后置中间件
func ActionLogMiddleWare() gin.HandlerFunc {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	return func(ctx *gin.Context) {
		//写操作日志
	}
}
