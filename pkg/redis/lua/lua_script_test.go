package lua

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestRedis_DecrToMin(t *testing.T) {
	counter := New("127.0.0.1:6379", 0)
	count, err := counter.IncrBy(context.Background(), "counter", 50, 60*time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println(count)
	for {
		count, err = counter.DecrToMin(context.Background(), "counter", 10, 0, 60*time.Second)
		if err != nil {
			panic(err)
		}
		fmt.Println(count)
		if count == 0 {
			break
		}
	}
}

func TestLoadScript(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	defer c.Close()
	go producer(c, "apple")
	go producer(c, "pineapple")
	go consumer(c, "lbb", "apple")
	go consumer(c, "lbb", "apple")
	go consumer(c, "lbb", "apple")
	go consumer(c, "lbb", "apple")
	go consumer(c, "lbb", "pineapple")

	time.Sleep(time.Second * 100)

}

func producer(c *redis.Client, product string) error {
	for {
		_, err := c.Incr(context.Background(), product).Result()
		if err != nil {
			return err
		}
		//fmt.Println(product + "+1")
		time.Sleep(time.Second)
	}
}

func consumer(c *redis.Client, user string, product string) error {
	script, err := ioutil.ReadFile("./test_lua.lua")
	if err != nil {
		fmt.Println("Script read error", err)
		return err
	}
	lua := redis.NewScript(string(script))
	for {
		res, err := lua.Run(context.Background(), c, []string{user, product}).Result()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if r, _ := res.(int64); r == 1 {
			fmt.Println(user + " buy " + product)
		} else if r == 0 {
			fmt.Println("Unluckly " + user + "," + product + " is sold out!")
		} else {
			fmt.Println(user + "用户冷却中")
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)*100))
	}
}
