package resty

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mel0dys0ng/song/pkg/utils/crypto"
)

// CreateSign 签名
func (c *Client) CreateSign(data map[string]string) (sign string) {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	ps := []string{fmt.Sprintf("%s@%d#", c.config.SignSerect, c.config.SignTTL)}
	for _, v := range keys {
		ps = append(ps, fmt.Sprintf("%s=%s", v, data[v]))
	}

	// 签名
	sign = crypto.MD5(strings.Join(ps, "&"))

	return
}

// VerifySign 验签
func (c *Client) VerifySign(sign string, data map[string]string) (res bool) {
	return c.CreateSign(data) == sign
}
