package controller

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// saveSession 保存用户登录会话
func saveSession(ctx *gin.Context, key, value string) {
	session := sessions.Default(ctx)
	session.Set(key, value)
	session.Save()
}
