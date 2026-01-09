package main

import (
	"fmt"
	"github.com/r3labs/sse/v2"
)

func main() {
	client := sse.NewClient("http://localhost:8080/events")

	// 1) 如果本地有上次记录，塞到头部
	//    这里演示写死 5，实际从 DB/文件读
	client.Headers["Last-Event-ID"] = "5"

	// 2) 订阅并打印
	client.Subscribe("test", func(ev *sse.Event) {
		fmt.Printf("got ID=%s  Data=%s\n", ev.ID, ev.Data)
		// 3) 落库：把 ev.ID 存盘，下次重启再带过去
	})
}
