#!/bin/bash

  swag init -d ../../app/go_web/go_gin/cmd/gin_swagger -o ../../docs --parseDependency




#  安装swag
#$ go get -u github.com/swaggo/swag/cmd/swag@v1.6.5
#
#验证是否安装成功
#检查 $GOBIN 下是否有 swag 文件，如下：
#$ swag -v
#swag version v1.6.5
#
#安装 gin-swagger
#$ go get -u github.com/swaggo/gin-swagger@v1.2.0
#$ go get -u github.com/swaggo/files
#$ go get -u github.com/alecthomas/template

#sh ./swag.sh

#启动服务，访问：http://127.0.0.1:8080/swagger/index.html#/