package resty

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	resty2 "github.com/go-resty/resty/v2"
	"github.com/mel0dys0ng/song/pkg/caller"
	"github.com/mel0dys0ng/song/pkg/metas"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const (
	HeaderKeyDid     = "X-Song-Did"  // 依赖服务ID
	HeaderKeyKind    = "X-Song-Kd"   // App Kind
	HeaderKeyApp     = "X-Song-Na"   // App Name
	HeaderKeyNode    = "X-Song-Nd"   // App Node
	HeaderKeyTraceId = "X-Song-Tid"  // Trace ID
	HeaderKeySpanId  = "X-Song-Sid"  // Span ID
	HeaderKeyTs      = "X-Song-Ts"   // 时间戳
	HeaderKeySign    = "X-Song-Sign" // 签名
	HeaderKeyRs      = "X-Song-Rs"   // 随机字符串
	HeaderKeyFl      = "X-Song-Fl"   // 调用位置
)

var clients = &sync.Map{}

// Client Resty 客户端结构
type Client struct {
	// resty 客户端实例
	*resty2.Client
	// 配置键
	key string
	// 配置对象
	config *Config
	// metadata
	metadata metas.MetadataInterface
}

// Key 返回配置键
func (c *Client) Key() string {
	return c.key
}

// R 返回带有预设头信息的 Resty 请求对象
func (c *Client) R(ctx context.Context) *resty2.Request {
	if c.config.Type == Extranet {
		// 外网请求不需要设置额外头部和签名
		return c.Client.R()
	}

	// 获取调用者信息
	cl := caller.New(3)

	// 构建请求头数据
	data := map[string]string{
		HeaderKeyDid:     c.config.Did, // 依赖服务ID
		HeaderKeyKind:    c.metadata.Kind().String(),
		HeaderKeyApp:     c.metadata.App(),
		HeaderKeyNode:    c.metadata.Node(),
		HeaderKeyTs:      cast.ToString(time.Now().Unix()),
		HeaderKeyRs:      lo.RandomString(32, lo.AlphanumericCharset),
		HeaderKeyFl:      fmt.Sprintf("%s-%d", cl.Func(), cl.Line()),
		HeaderKeySpanId:  "", // 可以从上下文中获取trace信息
		HeaderKeyTraceId: "", // 可以从上下文中获取trace信息
	}

	// 设置请求头
	c.Client.SetHeaders(data)
	// 设置签名钩子
	c.Client.SetPreRequestHook(c.setRequestSign)

	return c.Client.R()
}

// setRequestSign 在发送请求前设置签名
func (c *Client) setRequestSign(client *resty2.Client, request *http.Request) (err error) {
	// 收集需要签名的头部信息
	data := map[string]string{
		HeaderKeyDid:     "",
		HeaderKeyKind:    "",
		HeaderKeyApp:     "",
		HeaderKeyNode:    "",
		HeaderKeyTraceId: "",
		HeaderKeySpanId:  "",
		HeaderKeyTs:      "",
		HeaderKeyRs:      "",
		HeaderKeyFl:      "",
	}

	// 从请求头中提取上述字段的值
	for k := range data {
		if headerVal := client.Header.Get(k); headerVal != "" {
			data[k] = headerVal
		}
	}

	// 如果配置允许对查询参数签名，则处理查询参数
	if c.config.SignConfig.Query {
		queryParams := request.URL.Query()
		for key, values := range queryParams {
			if key != HeaderKeySign { // 排除签名参数本身
				if len(values) > 0 {
					data[key] = values[0] // 只取第一个值
				}
			}
		}
	}

	// 如果配置允许对表单数据签名，则处理表单数据
	if c.config.SignConfig.FormData {
		if request.Method == "POST" || request.Method == "PUT" || request.Method == "PATCH" {
			// 对于表单数据，需要检查Content-Type是否为application/x-www-form-urlencoded
			contentType := request.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				formValues := request.Form
				for key, values := range formValues {
					if key != HeaderKeySign { // 排除签名参数本身
						if len(values) > 0 {
							data[key] = values[0] // 只取第一个值
						}
					}
				}
			}
		}
	}

	// 如果配置允许对请求头签名，则处理所有请求头
	if c.config.SignConfig.Header {
		for key, values := range request.Header {
			lowerKey := strings.ToLower(key)
			if lowerKey != strings.ToLower(HeaderKeySign) { // 排除签名头
				if len(values) > 0 {
					data[lowerKey] = values[0] // 只取第一个值
				}
			}
		}
	}

	// 设置时间戳（如果未设置的话）
	if _, exists := data[TimestampKey]; !exists {
		data[TimestampKey] = strconv.FormatInt(time.Now().Unix(), 10)
	}

	// 设置签名头
	client.SetHeader(HeaderKeySign, c.CreateSign(data))

	return
}
