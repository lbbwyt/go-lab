package prase

import (
	"fmt"
	"testing"
)

func TestParseArgs(t *testing.T) {
	args, err := ParseArgs("SET redis redis.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%v", args))
}
