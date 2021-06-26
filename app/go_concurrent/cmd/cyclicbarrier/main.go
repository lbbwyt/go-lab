package main

import (
	"context"
	"fmt"
	"github.com/marusama/cyclicbarrier"
	"log"
	"sync"
	"sync/atomic"
)

func main() {

	var cnt int32 = 1
	b := cyclicbarrier.NewWithAction(4, func() error {
		fmt.Println(fmt.Sprintf("cnt %d", cnt))
		atomic.AddInt32(&cnt, 1)
		return nil
	})

	wg := sync.WaitGroup{}
	wg.Add(4)

	//只有barrier的最后一个才会去执行action

	for i := 0; i < 4; i++ {
		i := i
		go func() {
			defer wg.Done()
			log.Printf("goroutine %d waits", i)
			err := b.Await(context.TODO())
			log.Printf("goroutine %d is OK", i)
			if err != nil {
				panic(err)
			}

		}()
	}

	wg.Wait()

	fmt.Println(cnt)
}
