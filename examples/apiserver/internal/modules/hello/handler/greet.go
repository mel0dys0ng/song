package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/core/erlogs"
	"github.com/mel0dys0ng/song/core/https"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/modules/hello/service"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/status"
)

func (i *Instance) SayHello(ctx *gin.Context) {
	request := &service.SayHelloRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		err = status.InvalidArguments.Info(ctx, erlogs.ContentError(err))
		https.ResponseError(ctx, err)
		return
	}

	response, err := i.service.SayHello(ctx, request)
	if err != nil {
		https.ResponseError(ctx, err)
		return
	}

	https.ResponseSuccess(ctx, response)
}

func (i *Instance) SayHi(ctx *gin.Context) {
	request := &service.SayHiRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		err = status.InvalidArguments.Info(ctx, erlogs.ContentError(err))
		https.ResponseError(ctx, err)
		return
	}

	response, err := i.service.SayHi(ctx, request)
	if err != nil {
		https.ResponseError(ctx, err)
		return
	}

	https.ResponseSuccess(ctx, response)
}
