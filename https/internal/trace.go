package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/erlogs"
)

func (s *Server) trace(ctx *gin.Context) {
	newCtx := erlogs.StartTrace(ctx.Request.Context(), "doHttpRequest")
	ctx.Request = ctx.Request.WithContext(newCtx)
	defer erlogs.EndTrace(newCtx, nil)
	ctx.Next()
}
