package main

import (
	"fmt"
	"sync"
)

var (
	mutex     = sync.Mutex{}
	sum   int = 0
)

func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			mutex.Lock()
			defer mutex.Unlock()
			sum = sum + i
			fmt.Println(fmt.Sprintf("%d, %d", i, sum))
		}(i)
	}
	select {}
}
