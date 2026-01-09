package examples

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	ctx := context.Background()
	res, err := DoWithData(ctx,
		func() (int, error) {
			return callRemoteWithData() // 可能失败
		},
		WithAttempts(5),
		WithDelay(time.Second*1),
		WithBackoff(time.Second*10),
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)
}

func callRemote() error {
	fmt.Println("do error")
	return errors.New("call Remote error")
	//return nil

}

func callRemoteWithData() (int, error) {
	fmt.Println("do error")

	return 0, errors.New("call Remote error")
}
