package mongo_client

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestNewMongoDBClient(t *testing.T) {
	c, err := NewMongoDBClient("mongodb://localhost:27017", "multimodal_db", "physiological_signals_libaobao")
	if err != nil {
		panic(err)
	}
	filter := bson.M{"timestep": 1752135868350}
	m := make(map[string]interface{})
	err = c.FindOne(filter, &m)
	if err != nil {
		panic(err)
	}
	fmt.Println("查询到的文档:", m)
}

//windows 创建mongo服务
//mongod.exe --dbpath D:\install\install_package\mongodb\data --logpath D:\install\install_package\mongodb\log\mongod.log --serviceName --logappend "MongoDB" --install
//net start MongoDB
//net stop MongoDB
