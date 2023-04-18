package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	config := cors.Config{
		// 允许所有域名访问
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		// 允许方法
		AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE", "PATCH"},
		// 允许的客户端发送的header类型
		AllowHeaders: []string{"Authorization", "Content-Length", "X-CSRF-Token", "Token,session", "X_Requested_With,Accept", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Pragma", "Origin"},
		// 允许客户端浏览器解析的请求头 跨域关键设置 让浏览器JS可以解析的请求头
		ExposeHeaders: []string{"Access-Control-Expose-Headers", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma", "FooBar", "Content-Length"},
		// 允许cookie和认证信息
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(config)
}
