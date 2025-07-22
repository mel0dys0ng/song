package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"github.com/song/utils/strjngs"
	"golang.org/x/crypto/pbkdf2"
)

type (
	AESGCMCrypter struct {
		Key []byte `json:"key"` // 密钥
	}
)

func NewAESGCMCrypter(key, salt []byte) *AESGCMCrypter {
	return &AESGCMCrypter{
		Key: pbkdf2.Key(key, salt, 4096, 32, sha256.New),
	}
}

// AESGCMEncrypt 加密数据
func AESGCMEncrypt[T any](crypter *AESGCMCrypter, data T, salt []byte) (res string, err error) {
	if crypter == nil {
		err = errors.New("crypter is nil")
		return
	}

	str, err := strjngs.JSONMarshal(data)
	if err != nil {
		return
	}

	return crypter.Encrypt([]byte(str), salt)
}

// Encrypt 加密
func (c *AESGCMCrypter) Encrypt(data, salt []byte) (res string, err error) {
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, er := io.ReadFull(rand.Reader, nonce); er != nil {
		err = er
		return
	}

	res = Base64URLEncode(append(nonce, gcm.Seal(nil, nonce, data, salt)...))
	return
}

// Decrypt 解密
func (c *AESGCMCrypter) Decrypt(data string, salt []byte) (res []byte, err error) {
	bytes, err := Base64URLDecode(data)
	if err != nil || len(bytes) == 0 {
		return
	}

	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	// 解密
	nonceSize := gcm.NonceSize()
	return gcm.Open(nil, bytes[:nonceSize], bytes[nonceSize:], salt)
}
