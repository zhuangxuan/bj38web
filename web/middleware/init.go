package middleware

import (
	"bj38web/web/controller"
	"bj38web/web/utils"
	"fmt"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/gin-gonic/gin"
)

// InitMiddleware 接受服务实例，并存到gin.Key中，载入微服务
func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(context *gin.Context) {
		// 将实例存在gin.Keys中
		context.Keys = make(map[string]interface{})
		context.Keys["GetCaptcha"] = service[0]
		context.Keys["User"] = service[1]
		context.Keys["GetArea"] = service[2]
		context.Keys["House"] = service[3]
		context.Keys["Order"] = service[4]
		context.Next()
	}
}

// 错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				context.JSON(200, gin.H{
					"code": 404,
					"msg":  fmt.Sprintf("%s", r),
				})
				context.Abort()
			}
		}()
		context.Next()
	}
}

// LoginFilter 用户状态验证中间件
func LoginFilter(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userName := session.Get("userName")

	if userName == nil {
		fmt.Println("用户未登录")
		controller.ResponseError(ctx, utils.RECODE_LOGINERR)
		// 提前结束后面所有操作
		ctx.Abort()
	}
	ctx.Next()
}
