package test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestGin(t *testing.T) {
	r := gin.Default()
	// 设置session的存储对象 可以配置用redis mongodb 等存储session

	redisStore, err := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		fmt.Println("err", err)
	}
	store := redisStore
	// 使用options配置session对应的cookie 和gin中设置cookie是一样的
	store.Options(sessions.Options{
		MaxAge: 60 * 60,
	})
	// session根据cookie生成 session的name是cookie的name
	// 在gin中注册session中间件 指定session中间件名字（对应cookie的名字） 和 具体的session存储对象
	r.Use(sessions.Sessions("mysession", store))
	// 使用了session中间件的路由，回复会自带和session对应的cookie
	r.GET("/incr", func(c *gin.Context) {
		// gin中获取自己注册的session中间件
		session := sessions.Default(c)

		var count int
		// 获取session中的key
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}
		session.Set("count", count)
		// 保存session
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	r.Run(":8000")
}
