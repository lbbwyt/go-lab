package cache

import (
	"fmt"
	"testing"
)

func TestRedisCache_Put(t *testing.T) {
	c := GetRedisCache()
	err := c.Put("a", "ab", 0)
	if err != nil {
		panic(err)
	}

	value, _, err := c.Get("a")
	if err != nil {
		panic(err)
	}
	fmt.Println(value.(string))
}

func TestLocalCache_Put(t *testing.T) {
	c := GetLocalCache()
	err := c.Put("a", "ab", 0)
	if err != nil {
		panic(err)
	}

	value, _, err := c.Get("a")
	if err != nil {
		panic(err)
	}
	fmt.Println(value.(string))
}
