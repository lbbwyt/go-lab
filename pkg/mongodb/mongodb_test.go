package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

const connStringFmtWithoutPwd = "mongodb://%s/?replicaSet=%s"    // 没有密码时用的连接字符串
const connStringFmtWithPwd = "mongodb://%s:%s@%s/?replicaSet=%s" // 有密码时用的连接字符串
func TestTest(t *testing.T) {

	type role struct {
		IsMaster bool `json:"ismaster"`
	}

	var r role
	var (
		connectionString = ""
		password         = ""
		user             = ""
		addr             = "localhost:27017,localhost:27027,localhost:27037"
		rs               = "rs0"
	)

	if password != "" {
		connectionString = fmt.Sprintf(connStringFmtWithPwd, user, password, addr, rs)
	} else {
		connectionString = fmt.Sprintf(connStringFmtWithoutPwd, addr, rs)
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString).
		SetReadPreference(readpref.SecondaryPreferred()) // 设置读操作走Slave

	// Connect to MongoDB
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
		return
	}

	// Test connection
	err = c.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
		return
	}

	res := c.Database("admin").RunCommand(context.Background(), bson.M{"isMaster": 1}, &options.RunCmdOptions{})
	if res.Err() != nil {
		panic(err)
		return
	}

	if err = res.Decode(&r); err != nil {
		panic(err)
		return
	}

	s, _ := json.Marshal(r)
	fmt.Println(string(s))
}
