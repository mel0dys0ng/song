package https

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/metas"
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
		Body      *ResponseData
		TraceId   string
	}
)

func (s *Server) setupRespondedMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		if s.mt == nil {
			ctx.Next()
			return
		}

		path := ctx.Request.URL.Path
		rawQuery := ctx.Request.URL.RawQuery

		headers := make(map[string]string)
		const uaKey, ccKey = "User-Agent", "Cache-Control"

		headerKeys := make([]string, 0, len(s.LoggerHeaderKeys)+2)
		headerKeys = append(headerKeys, s.LoggerHeaderKeys...)
		headerKeys = append(headerKeys, uaKey, ccKey)

		for _, v := range headerKeys {
			if len(v) > 0 {
				headers[v] = ctx.Request.Header.Get(v)
			}
		}

		ctx.Next()

		end := time.Now()
		rsp := ResponseFromContext(ctx)

		var fullPath string
		if len(rawQuery) > 0 {
			var builder strings.Builder
			builder.Grow(len(path) + len(rawQuery) + 1)
			builder.WriteString(path)
			builder.WriteByte('?')
			builder.WriteString(rawQuery)
			fullPath = builder.String()
		} else {
			fullPath = path
		}

		data := &RequestResponseData{
			Metadata:  s.mt,
			StartTime: start,
			EndTime:   end,
			Cost:      end.Sub(start).Milliseconds(),
			ClientIp:  ctx.ClientIP(),
			RemoteIp:  ctx.RemoteIP(),
			Method:    ctx.Request.Method,
			Proto:     ctx.Request.Proto,
			Host:      ctx.Request.Host,
			Path:      fullPath,
			Form:      ctx.Request.Form,
			Status:    ctx.Writer.Status(),
			BodySize:  ctx.Writer.Size(),
			Body:      rsp,
			Headers:   headers,
			TraceId:   rsp.GetTraceId(),
		}

		fields := []zap.Field{
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

		msg := "failure"
		if ctx.Writer.Status() == http.StatusOK {
			msg = "success"
		}

		if data.Body != nil {
			if data.Body.GetCode() != ResponseSuccessCode {
				msg = "failure"
			}
			fields = append(fields,
				zap.Int64("bodyCode", data.Body.Code),
				zap.String("bodyMsg", data.Body.Msg),
			)
		}

		// record log
		erlogs.New(msg).InfoLog(ctx.Request.Context(), erlogs.OptionFields(fields...))

		// 自定义请求响应后
		if s.OnResponded != nil {
			s.OnResponded(ctx.Request.Context(), data)
		}

		ctx.Next()
	}
}
