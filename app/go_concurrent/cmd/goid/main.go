package main

import (
	"fmt"
	"go-lab/pkg/utils/runtime_utils"
	"sync"
)

func main() {
	fmt.Println("main", runtime_utils.GoID())
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i, runtime_utils.GoID())
		}()
	}
	wg.Wait()
}
