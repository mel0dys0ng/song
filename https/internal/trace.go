package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/song/erlogs"
)

func (s *Server) trace(ctx *gin.Context) {
	newCtx := erlogs.StartTrace(ctx.Request.Context(), "doHttpRequest")
	ctx.Request = ctx.Request.WithContext(newCtx)
	defer erlogs.EndTrace(newCtx, nil)
	ctx.Next()
}
