package crypto

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/tjfoc/gmsm/sm3"
)

func MD5(data any) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	h := md5.New()
	h.Write(bytes)
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
