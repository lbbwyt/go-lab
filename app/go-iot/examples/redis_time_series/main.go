package main

import (
	"fmt"
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	client := redistimeseries.NewClient("localhost:6379", "", nil)

	key := "test:ac001:pi4"
	//opts := redistimeseries.DefaultCreateOptions
	//opts.Labels = map[string]string{"deviceId": "pi4", "product": "ac001"}
	//opts.RetentionMSecs = 30 * 24 * 60 * time.Minute
	//
	//if err := client.CreateKeyWithOptions(key, opts); err != nil {
	//	log.Fatal(err)
	//}

	// 3. 写入数据（* = 服务器时间戳）
	go func() {
		var temp = 25.0
		for {
			if _, err := client.Add(key, time.Now().UnixMilli(), temp); err != nil {
				log.Fatal(err)
			}
			temp = temp + 0.1
			time.Sleep(time.Millisecond * 200)
		}
	}()

	time.Sleep(time.Second * 5)

	// 4. 实时查询最新点
	dp, _ := client.Get(key)
	fmt.Printf("最新值：%.2f @ %d\n", dp.Value, dp.Timestamp)

	// 5. 范围查询（最近 1 分钟）
	now := time.Now().UnixMilli()
	dps, _ := client.RangeWithOptions(key, now-60000, now,
		redistimeseries.DefaultRangeOptions)
	for _, v := range dps {
		fmt.Printf("t=%d  v=%.2f\n", v.Timestamp, v.Value)
	}
}
