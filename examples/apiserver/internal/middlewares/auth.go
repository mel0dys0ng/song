package middlewares

import (
	"github.com/gin-gonic/gin"
)

// MustAuth 强登录：必须登录，否则返回请登录错误信息，前端收到后跳转至登录页
func MustAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// auth logic
		ctx.Next()
	}
}

// WeaAuth 弱登录：可以不登录，也可以登录
func WeaAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// auth logic
		ctx.Next()
	}
}
