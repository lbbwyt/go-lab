package database

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-lab/app/go_web/go_gin/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

const connStringFmtWithoutPwd = "mongo_client://%s/?replicaSet=%s"    // 没有密码时用的连接字符串
const connStringFmtWithPwd = "mongo_client://%s:%s@%s/?replicaSet=%s" // 有密码时用的连接字符串

var once sync.Once // 保证全局只有1个mongodb client
var mongodbClient *mongo.Client
var mongodbClientInitErr error

// 闭包
func encapsulateBuildClientFunc(config *conf.MongoDb) func() {
	// db connection string using replica set
	var connectionString string
	if config.Password != "" {
		connectionString = fmt.Sprintf(connStringFmtWithPwd, config.User, config.Password, config.Addr, config.Rs)
	} else {
		connectionString = fmt.Sprintf(connStringFmtWithoutPwd, config.Addr, config.Rs)
	}

	return func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(connectionString)

		// Connect to MongoDB
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.WithFields(log.Fields{"user": config.User, "addr": config.Addr, "rs": config.Rs}).
				WithError(err).
				Error("mongo_client client init error")
			mongodbClientInitErr = err
			return
		}

		// Test connection
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.WithFields(log.Fields{"user": config.User, "addr": config.Addr, "rs": config.Rs}).
				WithError(err).
				Error("mongo_client connection ping error")
			mongodbClientInitErr = err
			return
		}

		log.WithFields(log.Fields{"user": config.User, "addr": config.Addr, "rs": config.Rs}).
			Info("mongo_client client init succeed")
		mongodbClient = client
	}
}

// 获取一个mongodb client(单例)
func GetMongodbClient(config *conf.MongoDb) (*mongo.Client, error) {
	if mongodbClient != nil {
		return mongodbClient, nil
	}
	once.Do(encapsulateBuildClientFunc(config))
	if mongodbClientInitErr != nil {
		return nil, mongodbClientInitErr
	} else {
		return mongodbClient, nil
	}
}
