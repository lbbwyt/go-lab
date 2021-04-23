package seata_golang

//1:mysql 建库建表，脚本见源码(包括seata脚本和examples的脚本)

//2： 启动seata-golang , 开启事务协调器服务（tc）
//
//编译：
//seata-golang 项目下：
//cd ./cmd/tc
//go build
//修改mysql配置：
//
//启动脚本：
//tc.exe start -config "C:\Users\30LHV53\Desktop\kcmdb\本地测试环境\seata-golang\config.yml"

//3：修改订单，产品，聚合服务配置， mysql和tc的服务地址并分别启动各个微服务，

//4：通过聚合服务暴露的api接口测试error时候分支事务回滚的情况
