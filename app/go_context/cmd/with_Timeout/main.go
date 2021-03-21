package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel() //超时之前可以调用cancel手动取消
	go task(cancelCtx)
	time.Sleep(time.Second * 4)
}

func task(ctx context.Context) {
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Gracefully exit")
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Println(i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
