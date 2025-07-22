package crypto

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/tjfoc/gmsm/sm3"
)

// Md5 生成md5字符串
// @args[0] 是否加入32位随机字符串
// @args[1] 是否加入当前纳秒时间戳
func Md5(data string, args ...bool) string {
	h := md5.New()
	datas := []string{data}

	argsLen := len(args)
	if argsLen >= 1 && args[0] {
		datas = append(datas, lo.RandomString(32, lo.AllCharset))
	}

	if argsLen >= 2 && args[1] {
		datas = append(datas, cast.ToString(time.Now().UnixNano()))
	}

	h.Write([]byte(strings.Join(datas, ":")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SM3 返回64位国密3加密后的字符串
func SM3(data string) string {
	h := sm3.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Base64URLEncode(data []byte) string {
	if len(data) > 0 {
		return base64.URLEncoding.EncodeToString(data)
	}
	return ""
}

func Base64URLDecode(data string) (res []byte, err error) {
	if len(data) > 0 {
		return base64.URLEncoding.DecodeString(data)
	}
	return
}

func MD5(data any) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	h := md5.New()
	h.Write(bytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}
