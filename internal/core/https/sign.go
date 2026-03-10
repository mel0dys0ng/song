package https

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/crypto"
	"github.com/mel0dys0ng/song/pkg/erlogs"
)

const (
	// HeaderKeySign 签名头部字段
	HeaderKeySign = "X-Song-Sign"
	// DefaultSignSecret 默认签名密钥
	DefaultSignSecret = "default_sign_secret"
	// DefaultSignTTL 默认签名有效期，单位秒
	DefaultSignTTL = 300 // 5分钟
	// TimestampKey 时间戳字段名
	TimestampKey = "timestamp"
)

// setupSignMiddleware 设置签名验证中间件
func (s *Server) setupSignMiddleware() gin.HandlerFunc {
	if s.Sign == nil || !s.Sign.Enable {
		// 如果没有启用签名验证，直接返回一个不执行任何操作的中间件
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 使用签名验证中间件
	return NewSignMiddleware(*s.Sign)
}

// NewSignMiddleware 创建签名验证中间件
func NewSignMiddleware(config Sign) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isNoCheckSignMethods(ctx) {
			// GET、HEAD、OPTIONS、TRACE 方法不需要验证签名
			ctx.Next()
			return
		}

		// 验证签名
		if !isValidSign(ctx, config) {
			// 签名无效，返回错误
			err := erlogs.InvalidSign.Warn().Log(ctx)
			ResponseError(ctx, err)
			ctx.Abort()
			return
		}

		// 签名有效，继续处理
		ctx.Next()
	}
}

// 不需要校验签名的请求方法
func isNoCheckSignMethods(ctx *gin.Context) bool {
	switch ctx.Request.Method {
	case http.MethodGet, http.MethodHead, http.MethodTrace, http.MethodOptions:
		return true
	}
	return false
}

// isValidSign 验证请求签名是否有效
func isValidSign(ctx *gin.Context, config Sign) bool {
	// 获取请求中的签名
	reqSign := ctx.GetHeader(HeaderKeySign)
	if reqSign == "" {
		// 尝试从查询参数获取
		reqSign = ctx.Query(HeaderKeySign)
		if reqSign == "" {
			// 尝试从表单获取
			reqSign = ctx.PostForm(HeaderKeySign)
		}
	}

	if reqSign == "" {
		return false
	}

	// 生成签名数据
	signData := make(map[string]string)

	// 添加查询参数
	if config.Query {
		for key, values := range ctx.Request.URL.Query() {
			if key != HeaderKeySign { // 排除签名本身
				if len(values) > 0 {
					signData[key] = values[0] // 只取第一个值
				}
			}
		}
	}

	// 添加表单数据
	if config.FormData {
		_ = ctx.Request.ParseForm() // 确保表单数据已解析
		for key, values := range ctx.Request.PostForm {
			if key != HeaderKeySign { // 排除签名本身
				if len(values) > 0 {
					signData[key] = values[0] // 只取第一个值
				}
			}
		}
	}

	// 添加请求头
	if config.Header {
		for key, values := range ctx.Request.Header {
			lowerKey := strings.ToLower(key)
			if lowerKey != strings.ToLower(HeaderKeySign) { // 排除签名头
				if len(values) > 0 {
					signData[lowerKey] = values[0] // 只取第一个值
				}
			}
		}
	}

	// 验证时间戳（如果存在）
	if timestampStr, exists := signData[TimestampKey]; exists {
		expectedTTL := config.TTL
		if expectedTTL == 0 {
			expectedTTL = DefaultSignTTL
		}

		if !validateTimestamp(timestampStr, expectedTTL) {
			return false
		}
	}

	// 生成期望的签名
	expectedSign := createSign(signData, config.Secret, config.TTL)

	// 比较签名
	return reqSign == expectedSign
}

// validateTimestamp 验证时间戳是否在有效期内
func validateTimestamp(timestampStr string, ttl int) bool {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		// 如果无法解析时间戳，则视为无效
		return false
	}

	// 检查时间戳是否在有效期内（允许前后误差1秒）
	now := time.Now().Unix()
	minTime := now - int64(ttl) - 1
	maxTime := now + int64(ttl) + 1

	return timestamp >= minTime && timestamp <= maxTime
}

// createSign 生成签名
func createSign(data map[string]string, secret string, ttl int) string {
	if secret == "" {
		secret = DefaultSignSecret
	}

	if ttl == 0 {
		ttl = DefaultSignTTL
	}

	var keys []string
	for k := range data {
		if k != HeaderKeySign { // 排除签名字段本身
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	// 格式: secret@ttl#key1=value1&key2=value2...
	parts := []string{fmt.Sprintf("%s@%d#", secret, ttl)}
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, data[k]))
	}

	signStr := strings.Join(parts, "&")
	return crypto.MD5(signStr)
}

// GenerateSign 为给定的数据生成签名，可用于客户端生成签名
func GenerateSign(data map[string]string, secret string, ttl int) string {
	return createSign(data, secret, ttl)
}
