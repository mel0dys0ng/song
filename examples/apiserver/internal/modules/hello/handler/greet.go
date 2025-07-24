package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (i *Instance) SayHello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hello world")
}

func (i *Instance) SayHi(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hi world")
}
