package chanx

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewUnboundChanSize(t *testing.T) {
	unboundChan := NewUnboundedChan(2)

	go func() {
		for v := range unboundChan.Out {
			fmt.Println(fmt.Sprintf("read:%d", v.(int64)))
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	go func(context context.Context) {
		for {
			unboundChan.In <- time.Now().Unix()
			time.Sleep(1 * time.Second)
			select {
			case <-context.Done():
				fmt.Println("done")
				return
			default:

			}
		}
	}(ctx)

	time.Sleep(10 * time.Second)
	cancel()
	fmt.Println(unboundChan.Len())
	fmt.Println(unboundChan.BufLen())

	for {

	}
}
