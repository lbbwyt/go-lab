package web_router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-lab/app/go_web/go_gin/handler/constant_handler"
	"go-lab/app/go_web/go_gin/handler/login_handler"
	_ "go-lab/docs"
)

//初始化router時必須,引入doc 文件夾
//_ "go-lab/docs"

func InitRouter() (*gin.Engine, error) {
	r := gin.Default()
	gin.SetMode("debug")

	pprof.Register(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	login := r.Group("/login")
	_ = RegisterLoginRouter(login)
	api := r.Group("/api")
	// 获取所有常量
	_ = RegisterConstantRouter(api)
	pprof.RouteRegister(api, "pprof")

	return r, nil
}

// Login api
func RegisterLoginRouter(rg *gin.RouterGroup) error {
	loginHandler := login_handler.NewLoginHandler()

	// routes
	{
		rg.GET("/check_login", loginHandler.UserLogin)
	}

	return nil
}

func RegisterConstantRouter(rg *gin.RouterGroup) error {

	// routes
	{
		rg.GET("/all_constants", constant_handler.GetAllConstants) // 获取所有常量
	}

	return nil
}
