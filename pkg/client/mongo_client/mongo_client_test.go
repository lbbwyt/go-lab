package mongo_client

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestNewMongoDBClient(t *testing.T) {

	c, err := NewMongoDBClient("mongodb://localhost:27017", "multimodal_db", "physiological_signals_libaobao")
	if err != nil {
		panic(err)
	}
	////filter := bson.M{"timestep": 1752135868350}
	////tmp := 1752163978
	//
	//tmpstr := "20250710161246"
	//
	//// 解析为时间对象
	//loc, _ := time.LoadLocation("Asia/Shanghai") // 或使用系统本地时区：time.Local
	//tm, err := time.ParseInLocation("20060102150405", tmpstr, loc)
	//if err != nil {
	//	panic(err)
	//}

	// 转换（10 进制, 64 位）
	//tmp := 1752135166000
	tmp := 1752135868350

	filter := bson.M{"timestep": bson.M{"$lt": tmp + 1000, "$gte": tmp}}

	//db["physiological_signals_libaobao"].find({"timestep", {$gte:1752163978000, $lt:1752163979000}})

	update := bson.M{"$set": bson.M{"label": "labelslabelslabelslabelslabelslabels"}}
	//1752163978
	start := time.Now().UnixMilli()
	//m := make(map[string]interface{})
	count, err := c.UpdateMany(filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Printf("更新文档数:【%d】", count))

	end := time.Now().UnixMilli()
	fmt.Println(fmt.Printf("查询到的文档耗时:【%d】", end-start))
}

func TestMongoDBClient_Close(t *testing.T) {
	timeStr := "20250710161246"
	layout := "20060102150405"
	loc, _ := time.LoadLocation("Asia/Shanghai") // 或使用系统本地时区：time.Local
	tm, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		panic(err)
	}
	println(tm.UnixMilli())

}

//windows 创建mongo服务
//mongod.exe --dbpath D:\install\install_package\mongodb\data --logpath D:\install\install_package\mongodb\log\mongod.log --serviceName --logappend "MongoDB" --install
//net start MongoDB
//net stop MongoDB
