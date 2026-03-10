package https

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/tjme"
)

// setupCORSMiddleware 设置跨域中间件
func (s *Server) setupCORSMiddleware() gin.HandlerFunc {
	if s.Cors == nil || !s.Cors.Enable {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:              s.Cors.AllowOrigins,
		AllowMethods:              s.Cors.AllowMethods,
		AllowHeaders:              s.Cors.AllowHeaders,
		AllowCredentials:          s.Cors.AllowCredentials,
		AllowWildcard:             s.Cors.AllowWildcard,
		ExposeHeaders:             s.Cors.ExposeHeaders,
		MaxAge:                    tjme.ParseDuration(s.Cors.MaxAge, DefaultCorsMaxAge),
		OptionsResponseStatusCode: s.Cors.OptionsResponseStatusCode,
	})
}
