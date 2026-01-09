package main

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"time"
)

func main() {
	err := retry.Do(
		func() error {
			return callRemote() // 可能失败
		},
		retry.Attempts(5),                     // 最多 5 次
		retry.Delay(time.Second),              // 首次间隔 1s
		retry.DelayType(retry.BackOffDelay),   // 指数退避
		retry.MaxJitter(time.Millisecond*100), // 随机抖动
	)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func callRemote() error {
	fmt.Println("do error")
	return errors.New("call Remote error")
	//return nil

}
