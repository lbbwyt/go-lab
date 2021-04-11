package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"go-lab/app/go_web/go_gin/conf"
	"go-lab/app/go_web/go_gin/global_web_var"
	"go-lab/app/go_web/go_gin/web_router"
	"net/http"
	"os"
	"time"
)

var (
	cfgPath string
)

func init() {
	flag.StringVar(&cfgPath, "c", "/Users/mac/libaobao/github/private/go-lab/app/go_web/go_gin/etc/config.yaml", "")
	flag.Parse()
}

// @name gin_swagger
// @title gin_swagger API
// @version 2.0
// @description gin_swagger 接口文档.
// @in header
// @produces application/json
// @schemes http https
// @termsOfService test-gin-swagger
// @host test-localhost.cn
// @BasePath /
// @contact.name lbbwyt
// @contact.email lbbwyt@126.com
func main() {
	log.Info("[main] starting project")

	// 解析配置文件
	err := conf.Init(cfgPath)
	if err != nil {
		log.WithFields(log.Fields{"cfg_path": cfgPath}).WithError(err).Error("[main] config init error")
		return
	}

	// 初始化log框架
	InitLog(conf.GConfig.Logger.Debug, conf.GConfig.Logger.ReportCaller)

	//初始化依赖组件
	global_web_var.Init()

	// start web server
	router, err := web_router.InitRouter()
	if err != nil {
		log.WithError(err).Error("[main] init router error")
		return
	}
	server := &http.Server{
		Addr:           ":" + conf.GConfig.WebServer.Port,
		Handler:        router,
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.WithError(err).Error("[main] start web server error")
		return
	}
}

// 配置log框架参数
func InitLog(debug, reportCaller bool) {
	log.SetFormatter(&log.JSONFormatter{}) // 以json格式输出日志
	// 输出到控制台, 从控制台收集日志
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	// 打印行号
	log.SetReportCaller(reportCaller)
}
