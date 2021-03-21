package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	rootCtx := context.Background()
	childCtx := context.WithValue(rootCtx, "request_Id", "123456")
	childOfChildCtx, cancelFunc := context.WithCancel(childCtx)
	go task(childOfChildCtx)
	time.Sleep(time.Second * 3)
	fmt.Println("number of goroutine: ", runtime.NumGoroutine()) // 协
	cancelFunc()
	time.Sleep(time.Second * 1)
	fmt.Println("number of goroutine: ", runtime.NumGoroutine()) // 协
}

func task(ctx context.Context) {
	fmt.Println(ctx.Value("request_Id").(string))
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Gracefully exit")
			fmt.Println(ctx.Err()) // 取消原因
			return
		default:
			time.Sleep(time.Second * 1)
		}
	}
}
