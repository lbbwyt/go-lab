package main

import "fmt"

//通过带缓冲的channel 控制最大并发数

var limit = make(chan struct{}, 100)

func requests() chan int {
	// 将从外界中接受到的请求放入到channel里
	return make(chan int)
}

func main() {
	for r := range requests() {
		go func() {
			limit <- struct{}{}
			fmt.Println(fmt.Sprintf("处理请求%d", r))
			<-limit
		}()
	}
}
