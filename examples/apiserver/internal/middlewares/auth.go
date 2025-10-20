package middlewares

import (
	"github.com/gin-gonic/gin"
)

// Auth 强登录：必须登录，否则返回请登录错误信息，前端收到后跳转至登录页
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// auth logic
		ctx.Next()
	}
}

// WeakAuth 弱登录：可以不登录，也可以登录
func WeakAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// auth logic
		ctx.Next()
	}
}
