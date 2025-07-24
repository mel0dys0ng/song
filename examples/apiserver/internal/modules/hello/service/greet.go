package service

import "github.com/gin-gonic/gin"

type (
	SayHelloRequest struct {
		Name string `json:"name" form:"name" validate:"required,max=10" msg:"请输入名称,最多10个字符"`
	}

	SayHelloResponse struct {
		Msg string `json:"msg"`
	}
)

func (i *Instance) SayHello(ctx *gin.Context, request *SayHelloRequest) (res *SayHelloResponse, err error) {
	return
}

type (
	SayHiRequest struct {
		Name string `json:"name" form:"name"`
	}

	SayHiResponse struct {
		Msg string `json:"msg"`
	}
)

func (i *Instance) SayHi(ctx *gin.Context, request *SayHiRequest) (res *SayHiResponse, err error) {
	return
}
