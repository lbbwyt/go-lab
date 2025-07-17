package main

import (
	"fmt"
	"strconv"
	"sync"
)

var (
	mutex     = sync.Mutex{}
	sum   int = 0
)

func main() {
	// 示例字符串
	str := "12345"

	// 转换（10 进制, 64 位）
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		// 处理错误（如非数字、溢出等）
		panic(err)
	}
	// num 现在是 int64 类型
	fmt.Printf("%d, %T", num, num) // 输出: 12345, int64
}
