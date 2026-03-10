package https

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/caller"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"go.uber.org/zap"
)

const (
	maxBodySize = 1024 * 1024 // 1MB 大小限制
)

type RecoveredData struct {
	// The error which triggered the panic
	Err error `json:"err"`
	// Caller is the caller of the function that triggered the panic
	Caller *caller.Caller `json:"caller"`
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	// If the connection is dead, we can't write a status to it.
	BrokenPipe bool `json:"broken_pipe"`
}

func (s *Server) setupRecoverAndTraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newCtx := erlogs.StartTrace(ctx.Request.Context(), "doHttpRequest")
		ctx.Request = ctx.Request.WithContext(newCtx)

		requestBody := s.captureRequestBody(ctx)
		responseBody := ""

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		defer func() {
			fields := []zap.Field{
				zap.String("status", strconv.Itoa(ctx.Writer.Status())),
				zap.String("method", ctx.Request.Method),
				zap.String("host", ctx.Request.Host),
				zap.String("path", ctx.Request.URL.Path),
				zap.String("query", ctx.Request.URL.RawQuery),
				zap.String("request_body", requestBody),
				zap.String("response_body", responseBody),
			}

			var recoveredData *RecoveredData
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := errors.AsType[*os.SyscallError](ne); ok {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				recoveredData = &RecoveredData{
					Err:        fmt.Errorf("%v", err),
					Caller:     caller.New(3),
					BrokenPipe: brokenPipe,
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					if s.OnRecovered != nil {
						s.OnRecovered(ctx, recoveredData)
					}
				} else {
					if s.OnRecovered != nil {
						s.OnRecovered(ctx, recoveredData)
					}
				}
			}

			err := erlogs.New("")
			if recoveredData != nil {
				err = erlogs.Convert(recoveredData.Err).Wrap("[Recovery] http request panic").Erorr()
				fields = append(fields, zap.Bool("broken_pipe", recoveredData.BrokenPipe))
				fields = append(fields, zap.String("caller", recoveredData.Caller.String()))
			}

			erlogs.EndTrace(newCtx, err.AppendFields(fields...))

			if recoveredData != nil {
				ResponseError(ctx, err)
				ctx.Abort()
			}
		}()

		ctx.Next()

		responseBody = s.truncateBody(blw.body.String())
	}
}

func (s *Server) captureRequestBody(ctx *gin.Context) string {
	if ctx.Request.Body == nil {
		return ""
	}

	contentType := ctx.Request.Header.Get("Content-Type")
	if isBinaryContent(contentType) {
		return "[binary content]"
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(ctx.Request.Body, maxBodySize+1))
	if err != nil {
		return "[read error]"
	}

	if len(bodyBytes) > maxBodySize {
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return "[too large]"
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if !utf8.Valid(bodyBytes) {
		return "[invalid utf8]"
	}

	body := string(bodyBytes)
	return body
}

func (s *Server) truncateBody(body string) string {
	if len(body) > maxBodySize {
		return body[:maxBodySize] + "...[truncated]"
	}
	return body
}

func isBinaryContent(contentType string) bool {
	contentType = strings.ToLower(contentType)
	binaryTypes := []string{
		"image/", "audio/", "video/", "application/pdf",
		"application/zip", "application/gzip", "application/octet-stream",
	}
	for _, t := range binaryTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}
	return false
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	if w.body.Len() < maxBodySize {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}
