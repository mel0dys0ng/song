package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/modules/hello/service"
	"github.com/mel0dys0ng/song/pkgs/https"
)

func (i *Instance) SayHello(ctx *gin.Context) {
	request := &service.SayHelloRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	data, err := i.service.SayHello(ctx, request)
	if err != nil {
		https.ResponseError(ctx, err)
		return
	}

	https.ResponseSuccess(ctx, data)
}

func (i *Instance) SayHi(ctx *gin.Context) {
	request := &service.SayHiRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	data, err := i.service.SayHi(ctx, request)
	if err != nil {
		https.ResponseError(ctx, err)
		return
	}

	https.ResponseSuccess(ctx, data)
}
