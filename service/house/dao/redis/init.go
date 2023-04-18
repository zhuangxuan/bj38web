package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// 由封装的服务来调用client
var client *redis.Client

// Init 初始化连接
func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password:     viper.GetString("redis.password"),                               // no password set
		DB:           viper.GetInt("redis.db"),                                        // use default DB
		PoolSize:     viper.GetInt("redis.pool_size"),                                 // 最大连接数量
		MinIdleConns: viper.GetInt("redis.min_idle_conns"),                            // 最少空闲连接数量
		MaxConnAge:   time.Second * time.Duration(viper.GetInt("redis.max_conn_age")), // 最大连接生命周期
		IdleTimeout:  time.Second * time.Duration(viper.GetInt("redis.idle_timeout")), // 空闲连接超时关闭时间
	})

	_, err = client.Ping().Result()
	return err
}

func Close() {
	_ = client.Close()
}
