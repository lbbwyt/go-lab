package runtime_utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGoWithRecover(t *testing.T) {
	handler := func() {
		panic("1")
	}
	recoverHandler := func(r interface{}) {
		panic("2")
	}
	GoWithRecover(handler, recoverHandler)

	time.Sleep(5 * time.Second)
	fmt.Println("123")
}
