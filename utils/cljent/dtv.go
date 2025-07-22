package cljent

import (
	"encoding/json"
	"errors"

	"github.com/mel0dys0ng/song/utils/crypto"
)

type DTV struct {
	Did           string `json:"did"`
	ClientType    uint8  `json:"client_type"`
	ClientVersion string `json:"client_version"`
	Salt          string `json:"salt"`
}

func ParseClientDTV(dtv string, crypter *crypto.AESGCMCrypter) (res *DTV, err error) {
	if len(dtv) == 0 {
		return
	}

	if crypter == nil {
		err = errors.New("crypter is nil")
		return
	}

	plaintext, err := crypter.Decrypt(dtv, nil)
	if err != nil {
		return
	}

	if len(plaintext) == 0 {
		err = errors.New("dtv empty")
		return
	}

	res = &DTV{}
	err = json.Unmarshal(plaintext, res)
	return
}

func BuildClientDTV(dtv *DTV, crypter *crypto.AESGCMCrypter) (res string, err error) {
	if dtv == nil {
		err = errors.New("dtv is nil")
		return
	}
	return crypto.AESGCMEncrypt(crypter, dtv, nil)
}

func CompareCientDTV(a, b *DTV) bool {
	return a != nil && b != nil && a.Did == b.Did && a.ClientVersion == b.ClientVersion &&
		a.ClientType == b.ClientType
}
