package redis

import (
	"bj38web/web/conf"
	"fmt"

	"github.com/go-redis/redis"
)

// 由封装的服务来调用client
var client *redis.Client

// Init 初始化连接
func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Conf.RedisConfig.Host, conf.Conf.RedisConfig.Port),
		Password:     conf.Conf.RedisConfig.Password,     // no password set
		DB:           conf.Conf.RedisConfig.DB,           // use default DB
		PoolSize:     conf.Conf.RedisConfig.PoolSize,     // 最大连接数量
		MinIdleConns: conf.Conf.RedisConfig.MinIdleConns, // 最少空闲连接数量
		//MaxConnAge:   time.Second * time.Duration(viper.GetInt("redis.max_conn_age")), // 最大连接生命周期
		//IdleTimeout:  time.Second * time.Duration(viper.GetInt("redis.idle_timeout")), // 空闲连接超时关闭时间
	})

	_, err = client.Ping().Result()
	return err
}

func Close() {
	_ = client.Close()
}
