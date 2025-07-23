package internal

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/erlogs"
	"github.com/mel0dys0ng/song/metas"
	"github.com/mel0dys0ng/song/utils/aob"
	"go.uber.org/zap"
)

const (
	EventName = "httpAccessLog"
)

type (
	RequestResponseData struct {
		Metadata  metas.MetadataInterface
		StartTime time.Time
		EndTime   time.Time
		Cost      int64
		ClientIp  string
		RemoteIp  string
		Status    int
		Method    string
		Proto     string
		Host      string
		Path      string
		Form      url.Values
		Headers   map[string]string
		BodySize  int
		Body      *Response
		TraceId   string
	}
)

func (s *Server) getRequestResponseData(ctx *gin.Context) (res *RequestResponseData) {

	start := time.Now()
	if _elgSys == nil { // 日志未初始化
		ctx.Next()
		return
	}

	path := ctx.Request.URL.Path
	rawQuery := ctx.Request.URL.RawQuery

	headers := make(map[string]string)
	const uaKey, ccKey = "User-Agent", "Cache-Control"
	s.LoggerHeaderKeys = append(s.LoggerHeaderKeys, uaKey, ccKey)
	for _, v := range s.LoggerHeaderKeys {
		if len(v) > 0 {
			headers[v] = ctx.Request.Header.Get(v)
		}
	}

	ctx.Next()

	end := time.Now()
	rsp := ResponseFromContext(ctx)

	return &RequestResponseData{
		Metadata:  metas.Mt(),
		StartTime: start,
		EndTime:   end,
		Cost:      end.Sub(start).Milliseconds(),
		ClientIp:  ctx.ClientIP(),
		RemoteIp:  ctx.RemoteIP(),
		Method:    ctx.Request.Method,
		Proto:     ctx.Request.Proto,
		Host:      ctx.Request.Host,
		Path:      aob.Aorb(len(rawQuery) > 0, path+"?"+rawQuery, path),
		Form:      ctx.Request.Form,
		Status:    ctx.Writer.Status(),
		BodySize:  ctx.Writer.Size(),
		Body:      rsp,
		Headers:   headers,
		TraceId:   rsp.GetTraceId(),
	}
}

// 记录请求日志并执行responded回调
func (s *Server) responded(ctx *gin.Context) {
	data := s.getRequestResponseData(ctx)

	fields := []zap.Field{
		erlogs.Event(EventName),
		zap.String("startTs", data.StartTime.String()),
		zap.String("endTs", data.EndTime.String()),
		zap.Int64("cost", data.Cost),
		zap.String("clientIp", data.ClientIp),
		zap.String("remoteIp", data.RemoteIp),
		zap.Int("status", data.Status),
		zap.String("proto", data.Proto),
		zap.String("host", data.Host),
		zap.String("method", data.Method),
		zap.String("path", data.Path),
		zap.String("traceId", data.TraceId),
		zap.Any("form", data.Form),
		zap.Any("headers", data.Headers),
		zap.Int("bodySize", data.BodySize),
	}

	msgv := "failure"
	if ctx.Writer.Status() == http.StatusOK {
		msgv = "success"
	}

	if data.Body != nil {
		if data.Body.GetCode() != ResponseSuccessCode {
			msgv = "failure"
		}
		fields = append(fields,
			zap.Int64("bodyCode", data.Body.Code),
			zap.String("bodyMsg", data.Body.Msg),
		)
	}

	// record log
	_elgSys.InfoL(ctx.Request.Context(), erlogs.Msgv(msgv), erlogs.Fields(fields...))

	// 自定义请求响应后
	if s.OnResponded != nil {
		s.OnResponded(ctx.Request.Context(), data)
	}

	ctx.Next()
}
