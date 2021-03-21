package main

import (
	"sync/atomic"
	"time"
)

func loadConfig() map[string]string {
	return make(map[string]string)
}

func requests() chan int {
	// 将从外界中接受到的请求放入到channel里
	return make(chan int)
}

func main() {
	// config变量用来存放该服务的配置信息, 比如说配置服务
	var config atomic.Value
	config.Store(loadConfig())

	// 每10秒钟定时的拉取最新的配置信息，并且更新到config变量里
	go func() {
		for {
			time.Sleep(10 * time.Second)
			config.Store(loadConfig())
		}
	}()
	// 创建工作线程，每个工作线程都会根据它所读取到的最新的配置信息来处理请求
	for i := 0; i < 10; i++ {
		go func() {
			for r := range requests() {
				// 对应于取值操作 c := config
				// 由于Load()返回的是一个interface{}类型，所以我们要先强制转换一下
				c := config.Load().(map[string]string)
				// 这里是根据配置信息处理请求的逻辑...
				_, _ = r, c
			}
		}()
	}

}
