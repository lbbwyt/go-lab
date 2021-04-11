package global_web_var

import (
	"github.com/jmoiron/sqlx"
	"go-lab/app/go_web/go_gin/conf"
	"go-lab/pkg/database/clickhouse"
	database "go-lab/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebVars struct {
	MongoClient   *mongo.Client   // mongodb client
	MongoDb       *mongo.Database // mongodb的database
	StatisticCKDb *sqlx.DB        //statistic 数据库客户端
}

var GWebVars WebVars

// 初始化全局变量
func Init() error {
	// mongodb
	MongoDbConfig := conf.GConfig.MongoDbConfigs.DB
	mongodbClient, err := database.GetMongodbClient(&MongoDbConfig)
	if err != nil {
		return err
	}
	GWebVars.MongoClient = mongodbClient
	GWebVars.MongoDb = mongodbClient.Database(MongoDbConfig.Db)

	click := conf.GConfig.ClickHouse
	db, err := clickhouse.NewClickHouseClient(click.Addrs, click.User, click.Password, click.StatisticDb)
	if err != nil {
		return err
	}
	GWebVars.StatisticCKDb = db

	return nil
}
