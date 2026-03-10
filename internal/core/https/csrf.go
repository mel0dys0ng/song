package https

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/samber/lo"
)

const (
	CSRFTokenKey = "X-Song-Csrf-Token"

	CSRFDefaultEnable         = false
	CSRFDefaultLookupType     = LookupTypeHeader
	CSRFDefaultLookupName     = CSRFTokenKey
	CSRFDefaultCookieName     = CSRFTokenKey
	CSRFDefaultCookieDomain   = ""
	CSRFDefaultCookiePath     = "/"
	CSRFDefaultCookieMaxAge   = 3600 // 1 hour
	CSRFDefaultCookieSecure   = false
	CSRFDefaultCookieHttpOnly = true

	CSRFTokenPrefix          = CSRFTokenKey
	CSRFTokenContextValueKey = CSRFTokenKey

	LookupTypeHeader = "header"
	LookupTypeForm   = "form"
	LookupTypeQuery  = "query"

	LookupNameHeader = CSRFTokenKey
	LookupNameForm   = "x-song-csrf-token"
	LookupNameQuery  = "x-song-csrf-token"
)

// setupCSRFMiddleware 设置CSRF中间件
func (s *Server) setupCSRFMiddleware() gin.HandlerFunc {
	if s.Csrf == nil || !s.Csrf.Enable {
		return func(c *gin.Context) {
			// 如果没有启用CSRF验证，直接返回一个不执行任何操作的中间件
			c.Next()
		}
	}

	if s.Csrf.LookupType == "" {
		s.Csrf.LookupType = CSRFDefaultLookupType
	}

	if s.Csrf.LookupName == "" {
		s.Csrf.LookupName = CSRFDefaultLookupName
	}

	if s.Csrf.CookieName == "" {
		s.Csrf.CookieName = CSRFDefaultCookieName
	}

	if s.Csrf.CookieDomain == "" {
		s.Csrf.CookieDomain = CSRFDefaultCookieDomain
	}

	if s.Csrf.CookieMaxAge == 0 {
		s.Csrf.CookieMaxAge = CSRFDefaultCookieMaxAge
	}

	// 使用CSRF中间件
	return newCSRFMIddleware(*s.Csrf)
}

// newCSRFMIddleware CSRF中间件实现
func newCSRFMIddleware(config CSRF) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查请求方法，如果是安全方法（GET、HEAD、OPTIONS、TRACE），跳过验证
		if isNoCheckTokenMethods(ctx) {
			// 为安全方法生成CSRF令牌
			token := buildCSRFToken()

			// 设置CSRF令牌到cookie
			ctx.SetCookie(
				config.CookieName,
				token,
				config.CookieMaxAge,
				config.CookiePath,
				config.CookieDomain,
				config.CookieSecure,
				config.CookieHttpOnly,
			)

			// 将token添加到上下文中供后续使用
			ctx.Set(CSRFTokenContextValueKey, token)
			ctx.Next()

			return
		}

		// 验证CSRF令牌
		if !isValidCSRFToken(ctx, config) {
			// 令牌无效，返回错误
			err := erlogs.InvalidCSRFToken.Warn().Log(ctx)
			ResponseError(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// 不需要校验csrf-token的请求方法
func isNoCheckTokenMethods(ctx *gin.Context) bool {
	switch ctx.Request.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	}
	return false
}

// isValidCSRFToken 检查CSRF令牌是否有效
func isValidCSRFToken(c *gin.Context, config CSRF) bool {
	switch config.LookupType {
	case LookupTypeHeader:
		token := c.GetHeader(config.LookupName)
		if token != "" {
			cookieToken, err := c.Cookie(config.CookieName)
			return err == nil && token == cookieToken
		}
	case LookupTypeForm:
		token := c.PostForm(config.LookupName)
		if token != "" {
			cookieToken, err := c.Cookie(config.CookieName)
			return err == nil && token == cookieToken
		}
	case LookupTypeQuery:
		token := c.Query(config.LookupName)
		if token != "" {
			cookieToken, err := c.Cookie(config.CookieName)
			return err == nil && token == cookieToken
		}
	}
	return false
}

// buildCSRFToken 生成CSRF令牌
func buildCSRFToken() string {
	// 使用时间戳和随机数生成安全的CSRF令牌
	return fmt.Sprintf("%s.%s.%d", CSRFTokenPrefix, lo.RandomString(32, lo.AllCharset), time.Now().UnixNano())
}
