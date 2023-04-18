package test

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func TestRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 50,
	})
	result, err := client.Ping().Result()
	if err != nil {
		fmt.Println("client.Ping().Result():", err)
		return
	}
	defer client.Close()
	fmt.Println("redis 连接成功", result)

	res, err := client.Set("itcast", "itheima", 0).Result()
	if err != nil {
		fmt.Println("client.Set(\"itcast\", \"itheima\", 0).Result():", err)
		return
	}
	fmt.Println(res)
	s, err := client.Get("itcast").Result()
	if err != nil {
		fmt.Println("client.Set(\"itcast\", \"itheima\", 0).Result():", err)
		return
	}
	fmt.Println(s)
}
