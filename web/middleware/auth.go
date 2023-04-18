package middleware

import (
	"bj38web/web/controller"
	"bj38web/web/dao/redis"
	"bj38web/web/utils"
	"strings"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		// 数据格式 Authorization:bearer xxxx.xxx.xx
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// 请求头没发token 就是未登录
			zap.L().Error("请求头没发token")
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		// token格式错误
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			zap.L().Error("token格式错误")
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		_, err := utils.ParseToken(parts[1])
		if err != nil {
			zap.L().Error("token解析错误")
			c.Abort()
			return
		}
		c.Next() // 后续的处理函数可以用过c.Get(CtxUserIDKey)来获取当前请求的用户信息
	}
}

// SingleLoginMiddleware 单点登录校验中间件
func SingleLoginMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 获取客户端的token
		authHeader := c.Request.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)

		token := parts[1]

		mc, _ := utils.ParseToken(parts[1])
		// 与redis的token做校验
		res, err := redis.HgetUsernameToken(mc.Mobile)
		if err != nil {
			zap.L().Error("token不存在")
			controller.ResponseError(c, utils.RECODE_INVALIDTOKENERR)
			c.Abort()
			return
		}
		if token != res {
			// 使用的token不一样 则为两地登录
			zap.L().Error("使用的token不一样")
			controller.ResponseError(c, utils.RECODE_INVALIDTOKENERR)
			c.Abort()
			return
		}
		c.Next() // 后续的处理函数可以用过c.Get(CtxUserIDKey)来获取当前请求的用户信息
	}
}
