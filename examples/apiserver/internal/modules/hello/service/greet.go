package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mel0dys0ng/song/core/erlogs"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/status"
)

type (
	SayHelloRequest struct {
		Name string `json:"name" form:"name" validate:"required,max=10" msg:"请输入名称,最多10个字符"`
	}

	SayHelloResponse struct {
		Msg string `json:"msg"`
	}
)

func (i *Instance) SayHello(ctx *gin.Context, request *SayHelloRequest) (response *SayHelloResponse, err error) {
	reqCtx := ctx.Request.Context()
	if err = validator.New().Struct(request); err != nil {
		err = status.InvalidArguments.Info(reqCtx, erlogs.ValidateError(request, err))
		return
	}

	// other logic ......

	return
}

type (
	SayHiRequest struct {
		Name string `json:"name" form:"name" validate:"required,max=10" msg:"请输入名称,最多10个字符"`
	}

	SayHiResponse struct {
		Msg string `json:"msg"`
	}
)

func (i *Instance) SayHi(ctx *gin.Context, request *SayHiRequest) (response *SayHiResponse, err error) {
	reqCtx := ctx.Request.Context()
	if err = validator.New().Struct(request); err != nil {
		err = status.InvalidArguments.Info(reqCtx, erlogs.ValidateError(request, err))
		return
	}

	// other logic ......

	return
}
