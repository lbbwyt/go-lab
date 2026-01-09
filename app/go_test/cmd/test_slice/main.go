package main

import "sync"

func main() {
	mp := sync.Map{}
	keys := []int64{122, 211, 122, 122, 122, 211, 5}

	for index, v := range keys {
		if value, okk := mp.Load(v); okk {
			if arr, ok := value.([]int); ok {
				arr = append(arr, index)
				mp.Store(v, arr)
				if len(arr) == 4 {
					println(v)
				}

			}
		} else {
			arr := make([]int, 0)
			arr = append(arr, index)
			mp.Store(v, arr)
		}
	}
}
