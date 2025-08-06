package utils

import (
	"github.com/mel0dys0ng/song/core/utils/crypto"
	"github.com/mel0dys0ng/song/core/utils/singleton"
	"github.com/mel0dys0ng/song/core/vipers"
)

func NewAESGCMCrypter() *crypto.AESGCMCrypter {
	return singleton.GetInstance("crypto.aesgcm", func() *crypto.AESGCMCrypter {
		return crypto.NewAESGCMCrypter(
			[]byte(vipers.GetString("crypto.aesKey", "")),
			[]byte(vipers.GetString("crypto.aesSalt", "")),
		)
	})
}
