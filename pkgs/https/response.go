package https

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkgs/erlogs"
	"github.com/mel0dys0ng/song/pkgs/https/internal"
)

type ContextResponse struct {
	ctx *gin.Context
}

func WithContext(ctx *gin.Context) *ContextResponse {
	return &ContextResponse{ctx: ctx}
}

// Response 请求响应。
// @Param data any 响应数据
// @param err error 请求错误。成功时，nil或者级别低于warning的错误；失败时，不为nil且级别高于warning的错误
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据err和data设定的选项
func (r *ContextResponse) Response(data any, err error, opts ...internal.ResponseOption) {
	r.ResponseWithStatus(http.StatusOK, append([]internal.ResponseOption{
		{
			Apply: func(rsp *internal.Response) {
				rsp.Data = data
				rsp.Code = internal.ResponseSuccessCode
				if err == nil {
					return
				}

				var v erlogs.ErLogInterface
				if !errors.As(err, &v) {
					v = erlogs.Unknown.WithOptions(erlogs.ContentError(err))
				}

				rsp.Msg = v.Msg()
				rsp.Code = v.Code()

				v.RecordLog(r.ctx.Request.Context())
			},
		},
	}, opts...)...)
}

// ResponseSuccess 请求成功时的响应
// @param data any 响应数据
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据data设定的选项
func (r *ContextResponse) ResponseSuccess(data any, opts ...internal.ResponseOption) {
	r.ResponseWithStatus(http.StatusOK, append([]internal.ResponseOption{
		{
			Apply: func(rsp *internal.Response) {
				rsp.Code = internal.ResponseSuccessCode
				switch v := data.(type) {
				case string:
					rsp.Msg = v
				default:
					rsp.Data = data
				}
			},
		},
	}, opts...)...)
}

// ResponseError 请求失败（发生错误或异常）时的响应
// @param err error 请求错误。成功时，nil或者级别低于warning的错误；失败时，不为nil且级别高于warning的错误
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据data设定的选项
func (r *ContextResponse) ResponseError(err error, opts ...internal.ResponseOption) {
	r.ResponseWithStatus(http.StatusOK, append([]internal.ResponseOption{
		{
			Apply: func(rsp *internal.Response) {
				var v erlogs.ErLogInterface
				if err == nil || !errors.As(err, &v) {
					v = erlogs.Unknown.WithOptions(erlogs.ContentError(err))
				}

				rsp.Msg = v.Msg()
				rsp.Code = v.Code()
				if rsp.Code == internal.ResponseSuccessCode {
					rsp.Code = erlogs.Unknown.Code()
				}

				v.RecordLog(r.ctx.Request.Context())
			},
		},
	}, opts...)...)
}

// ResponseWithStatus 自定义http status响应，默认JSON格式
func (r *ContextResponse) ResponseWithStatus(status int, opts ...internal.ResponseOption) {
	rsp := internal.NewResponse(opts...)

	rsp.Ts = time.Now().String()
	rsp.TraceId = erlogs.TraceSpanFromContext(r.ctx.Request.Context()).GetTraceID()
	r.ctx.Header(internal.ResponseTraceIdHeaderKey, rsp.TraceId)
	r.ctx.Set(internal.ResponseCtxValueKey, rsp)

	switch rsp.Type {
	case internal.ResponseTypeJSON:
		r.ctx.JSON(status, rsp)
	case internal.ResponseTypeSTREAM:
		r.ctx.Stream(func(w io.Writer) bool {
			_, err := w.Write([]byte(rsp.String()))
			if err != nil {
				return false
			}
			return rsp.Code == internal.ResponseSuccessCode
		})
	case internal.ResponseTypeAsciiJSON:
		r.ctx.AsciiJSON(status, rsp)
	case internal.ResponseTypeJSONP:
		r.ctx.JSONP(status, rsp)
	case internal.ResponseTypeHTML:
		r.ctx.HTML(status, rsp.Msg, rsp)
	default:
		r.ctx.JSON(status, rsp)
	}

	r.ctx.Abort()
}

// Response 请求响应。
// @Param data any 响应数据
// @param err error 请求错误。成功时，nil或者级别低于warning的错误；失败时，不为nil且级别高于warning的错误
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据err和data设定的选项
func Response(ctx *gin.Context, data any, err error, opts ...internal.ResponseOption) {
	WithContext(ctx).Response(data, err, opts...)
}

// ResponseSuccess 请求成功时的响应
// @param data any 响应数据
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据data设定的选项
func ResponseSuccess(ctx *gin.Context, data any, opts ...internal.ResponseOption) {
	WithContext(ctx).ResponseSuccess(data, opts...)
}

// ResponseError 请求失败（发生错误或异常）时的响应
// @param err error 请求错误。成功时，nil或者级别低于warning的错误；失败时，不为nil且级别高于warning的错误
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据data设定的选项
func ResponseError(ctx *gin.Context, err error, opts ...internal.ResponseOption) {
	WithContext(ctx).ResponseError(err, opts...)
}

// ResponseWithStatus 自定义http status响应，默认JSON格式
func ResponseWithStatus(ctx *gin.Context, status int, opts ...internal.ResponseOption) {
	WithContext(ctx).ResponseWithStatus(status, opts...)
}

func ResponseOptionCode(code int64) internal.ResponseOption {
	return internal.ResponseOption{
		Apply: func(rsp *internal.Response) {
			rsp.Code = code
		},
	}
}

func ResponseOptionMsg(msg string) internal.ResponseOption {
	return internal.ResponseOption{
		Apply: func(rsp *internal.Response) {
			rsp.Msg = msg
		},
	}
}

func ResponseOptionData(data any) internal.ResponseOption {
	return internal.ResponseOption{
		Apply: func(rsp *internal.Response) {
			rsp.Data = data
		},
	}
}

func ResponseOptionTypeJSON() internal.ResponseOption {
	return internal.ResponseOption{
		Apply: func(rsp *internal.Response) {
			rsp.Type = internal.ResponseTypeJSON
		},
	}
}

func ResponseOptionTypeJSONP() internal.ResponseOption {
	return internal.ResponseOption{
		Apply: func(rsp *internal.Response) {
			rsp.Type = internal.ResponseTypeJSONP
		},
	}
}

func ResponseOptionTypeAsciiJSON() internal.ResponseOption {
	return internal.ResponseOption{
		Apply: func(rsp *internal.Response) {
			rsp.Type = internal.ResponseTypeAsciiJSON
		},
	}
}

func ResponseOptionTypeHTML() internal.ResponseOption {
	return internal.ResponseOption{
		Apply: func(rsp *internal.Response) {
			rsp.Type = internal.ResponseTypeHTML
		},
	}
}

func ResponseOptionTypeStream() internal.ResponseOption {
	return internal.ResponseOption{
		Apply: func(rsp *internal.Response) {
			rsp.Type = internal.ResponseTypeSTREAM
		},
	}
}
