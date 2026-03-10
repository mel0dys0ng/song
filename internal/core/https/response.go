package https

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/erlogs"
)

const (
	ResponseSuccessCode      = 0
	ResponseUnknownCode      = -1
	ResponseTraceIdHeaderKey = "X-Song-Trace-Id"
	ResponseCtxValueKey      = "response"
	ResponseTypeJSON         = "JSON"
	ResponseTypeJSONP        = "JSONP"
	ResponseTypeAsciiJSON    = "ASCIIJSON"
	ResponseTypeSTREAM       = "STREAM"
	ResponseTypeHTML         = "HTML"
)

type (
	ResponseData struct {
		Type    string `json:"-"` // data stringfy type
		Code    int64  `json:"code"`
		Msg     string `json:"msg"`
		Data    any    `json:"data"`
		Biz     string `json:"biz,omitempty"`
		TraceId string `json:"request_id"`
		Ts      string `json:"ts"`
	}

	ResponseOption func(rsp *ResponseData)
)

func NewResponseData(opts ...ResponseOption) *ResponseData {
	rsp := &ResponseData{Code: ResponseSuccessCode, Type: ResponseTypeJSON}
	for _, opt := range opts {
		if opt != nil {
			opt(rsp)
		}
	}
	return rsp
}

func (r *ResponseData) GetBiz() string {
	if r != nil {
		return r.Biz
	}
	return ""
}

func (r *ResponseData) GetTraceId() string {
	if r != nil {
		return r.TraceId
	}
	return ""
}

func (r *ResponseData) GetCode() int64 {
	if r != nil {
		return r.Code
	}
	return 0
}

func (r *ResponseData) GetMsg() string {
	if r != nil {
		return r.Msg
	}
	return ""
}

func (r *ResponseData) GetData() any {
	if r != nil {
		return r.Data
	}
	return nil
}

func (r *ResponseData) GetTs() string {
	if r != nil {
		return r.Ts
	}
	return ""
}

func (r *ResponseData) GetType() string {
	if r != nil {
		return r.Type
	}
	return ""
}

func (r *ResponseData) String() string {
	if r != nil {
		bytes, err := json.Marshal(r)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
	return ""
}

// Response 响应
func Response(ctx *gin.Context, data any, err error, opts ...ResponseOption) {
	if err != nil {
		ResponseError(ctx, err, opts...)
	} else {
		ResponseSuccess(ctx, data, opts...)
	}
}

// ResponseSuccess 响应成功
func ResponseSuccess(ctx *gin.Context, data any, opts ...ResponseOption) {
	ResponseWithStatus(ctx, http.StatusOK, append([]ResponseOption{
		func(rsp *ResponseData) {
			rsp.Code = ResponseSuccessCode
			switch v := data.(type) {
			case string:
				rsp.Msg = v
			default:
				rsp.Data = data
			}
		},
	}, opts...)...)
}

// ResponseError 响应错误
func ResponseError(ctx *gin.Context, err error, opts ...ResponseOption) {
	el := erlogs.Convert(err)

	var status int
	switch el.GetCode() {
	case ResponseSuccessCode:
		status = http.StatusOK
	case ResponseUnknownCode:
		status = http.StatusInternalServerError
	default:
		status = http.StatusBadRequest
	}

	ResponseWithStatus(ctx, status, append([]ResponseOption{
		func(rsp *ResponseData) {
			rsp.Msg = el.GetMsg()
			rsp.Code = el.GetCode()
			if rsp.Code == ResponseSuccessCode {
				rsp.Code = erlogs.ServerError.GetCode()
			}
			el.RecordLog(ctx.Request.Context())
		},
	}, opts...)...)
}

// ResponseWithStatus 自定义http status响应，默认JSON格式
func ResponseWithStatus(ctx *gin.Context, status int, opts ...ResponseOption) {
	rsp := NewResponseData(opts...)

	rsp.Ts = time.Now().String()
	rsp.TraceId = erlogs.TraceSpanFromContext(ctx.Request.Context()).GetTraceID()
	ctx.Header(ResponseTraceIdHeaderKey, rsp.TraceId)
	ctx.Set(ResponseCtxValueKey, rsp)

	switch rsp.Type {
	case ResponseTypeJSON:
		ctx.JSON(status, rsp)
	case ResponseTypeSTREAM:
		ctx.Stream(func(w io.Writer) bool {
			_, err := w.Write([]byte(rsp.String()))
			if err != nil {
				return false
			}
			return rsp.Code == ResponseSuccessCode
		})
	case ResponseTypeAsciiJSON:
		ctx.AsciiJSON(status, rsp)
	case ResponseTypeJSONP:
		ctx.JSONP(status, rsp)
	case ResponseTypeHTML:
		ctx.HTML(status, rsp.Msg, rsp)
	default:
		ctx.JSON(status, rsp)
	}

	ctx.Abort()
}

func ResponseFromContext(ctx *gin.Context) *ResponseData {
	response, _ := ctx.Get(ResponseCtxValueKey)
	rsp, _ := response.(*ResponseData)
	return rsp
}
