package resty

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mel0dys0ng/song/pkg/crypto"
)

const (
	// TimestampKey 时间戳字段名
	TimestampKey = "timestamp"
)

// CreateSign 创建请求签名
// 参数 data 包含需要参与签名计算的所有键值对
// 返回计算出的签名字符串
func (c *Client) CreateSign(data map[string]string) (sign string) {
	return createSign(data, c.config.SignSecret, c.config.SignTTL)
}

// VerifySign 验证请求签名
// 参数 sign 是要验证的签名字符串
// 参数 data 是参与签名计算的原始数据
// 返回验证结果，true表示签名正确
func (c *Client) VerifySign(sign string, data map[string]string) (res bool) {
	expectedSign := createSign(data, c.config.SignSecret, c.config.SignTTL)
	return sign == expectedSign
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
