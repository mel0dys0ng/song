package internal

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
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
	Response struct {
		Type    string `json:"-"` // data stringfy type
		Code    int64  `json:"code"`
		Msg     string `json:"msg"`
		Data    any    `json:"data"`
		Biz     string `json:"biz,omitempty"`
		TraceId string `json:"request_id"`
		Ts      string `json:"ts"`
	}

	ResponseOption struct {
		Apply func(rsp *Response)
	}
)

func NewResponse(opts ...ResponseOption) *Response {
	rsp := &Response{Code: ResponseSuccessCode, Type: ResponseTypeJSON}
	for _, opt := range opts {
		if opt.Apply != nil {
			opt.Apply(rsp)
		}
	}
	return rsp
}

func (r *Response) GetBiz() string {
	if r != nil {
		return r.Biz
	}
	return ""
}

func (r *Response) GetTraceId() string {
	if r != nil {
		return r.TraceId
	}
	return ""
}

func (r *Response) GetCode() int64 {
	if r != nil {
		return r.Code
	}
	return 0
}

func (r *Response) GetMsg() string {
	if r != nil {
		return r.Msg
	}
	return ""
}

func (r *Response) GetData() any {
	if r != nil {
		return r.Data
	}
	return nil
}

func (r *Response) GetTs() string {
	if r != nil {
		return r.Ts
	}
	return ""
}

func (r *Response) GetType() string {
	if r != nil {
		return r.Type
	}
	return ""
}

func (r *Response) String() string {
	if r != nil {
		bytes, err := json.Marshal(r)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
	return ""
}

func ResponseFromContext(ctx *gin.Context) *Response {
	response, _ := ctx.Get(ResponseCtxValueKey)
	rsp, _ := response.(*Response)
	return rsp
}
