package resty

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	resty2 "github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"github.com/song/metas"
	"github.com/song/utils/caller"
	"github.com/spf13/cast"
)

const (
	HeaderKeyDid     = "X-Song-Did" // 依赖服务ID
	HeaderKeyProduct = "X-Song-Np"  // Product Name
	HeaderKeyApp     = "X-Song-Na"  // App Name
	HeaderKeyNode    = "X-Song-Nd"  // App Node
	HeaderKeyTraceId = "X-Song-Tid"
	HeaderKeySpanId  = "X-Song-Sid"
	HeaderKeyTs      = "X-Song-Ts"
	HeaderKeySign    = "X-Song-Sign"
	HeaderKeyRs      = "X-Song-Rs"
	HeaderKeyFl      = "X-Song-Fl"
)

var clients = &sync.Map{}

type Client struct {
	// client
	*resty2.Client
	// config key
	key string
	// config config
	config *Config
}

// Key return the normal config key
func (c *Client) Key() string {
	return c.key
}

// R return the resty request after set headers
func (c *Client) R(ctx context.Context) *resty2.Request {
	if c.config.Type == Extranet {
		// 外网请求
		return c.Client.R()
	}

	cl := caller.New(3)
	data := map[string]string{
		HeaderKeyDid:     c.config.Did, // 依赖服务ID
		HeaderKeyProduct: metas.Product(),
		HeaderKeyApp:     metas.App(),
		HeaderKeyNode:    metas.Node(),
		HeaderKeyTs:      cast.ToString(time.Now().Unix()),
		HeaderKeyRs:      lo.RandomString(32, lo.AlphanumericCharset),
		HeaderKeyFl:      fmt.Sprintf("%s-%d", cl.Func(), cl.Line()),
		HeaderKeySpanId:  "",
		HeaderKeyTraceId: "",
	}

	c.Client.SetHeaders(data)
	c.Client.SetPreRequestHook(c.setRequestSign)

	return c.Client.R()
}

// preRequestHook 发送请求前Hook(请求签名)
func (c *Client) setRequestSign(client *resty2.Client, request *http.Request) (err error) {
	data := map[string]string{
		HeaderKeyDid:     "",
		HeaderKeyProduct: "",
		HeaderKeyApp:     "",
		HeaderKeyNode:    "",
		HeaderKeyTraceId: "",
		HeaderKeySpanId:  "",
		HeaderKeyTs:      "",
		HeaderKeyRs:      "",
		HeaderKeyFl:      "",
	}

	for k := range data {
		data[k] = client.Header.Get(k)
	}

	qp := client.QueryParam.Encode()
	if len(qp) > 0 {
		sort.Slice([]rune(qp), func(i, j int) bool { return qp[i] > qp[j] })
		data["qp"] = qp
	}

	fd := client.FormData.Encode()
	if len(fd) > 0 {
		sort.Slice([]rune(fd), func(i, j int) bool { return fd[i] > fd[j] })
		data["fd"] = fd
	}

	client.SetHeader(HeaderKeySign, c.CreateSign(data))

	return
}
