package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	defer func() {
		fmt.Println("done")
	}()
	//令牌桶大小为 100, 以每秒 10 个 Token 的速率向桶中放置 Token。
	limit := rate.NewLimiter(10, 100)

	for i := 0; i < 10; i++ {
		err := limit.WaitN(context.Background(), 20)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println(fmt.Sprintf("index [%d] cost 20 tokens  at [%s]", i, time.Now().String()))
		}
	}
	//截止到某一时刻，目前桶中数目是否至少为 n 个，满足则返回 true，同时从桶中消费 n 个 token。反之不消费桶中的Token，返回false。
	if limit.AllowN(time.Now(), 20) {
		fmt.Println("event allowed")
	} else {
		fmt.Println("event not allowed")
	}

	for i := 10; i < 20; i++ {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		err := limit.WaitN(ctx, 10)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println(fmt.Sprintf("index [%d] cost 20 tokens  at [%s]", i, time.Now().String()))

		}
	}

}
