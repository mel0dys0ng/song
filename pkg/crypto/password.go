package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// PasswordHasher 密码哈希接口
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) error
}

// BcryptHasher bcrypt 密码哈希器
type BcryptHasher struct {
	Cost int
}

// NewBcryptHasher 创建 bcrypt 密码哈希器
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	return &BcryptHasher{Cost: cost}
}

// Hash 使用 bcrypt 哈希密码
func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify 验证 bcrypt 哈希密码
func (h *BcryptHasher) Verify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Argon2IDParams argon2id 参数
type Argon2IDParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgon2IDParams 默认 argon2id 参数
var DefaultArgon2IDParams = &Argon2IDParams{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,         // 3 轮
	Parallelism: 2,         // 2 线程
	SaltLength:  16,        // 16 字节盐
	KeyLength:   32,        // 32 字节密钥
}

// Argon2IDHasher argon2id 密码哈希器
type Argon2IDHasher struct {
	Params *Argon2IDParams
}

// NewArgon2IDHasher 创建 argon2id 密码哈希器
func NewArgon2IDHasher(params *Argon2IDParams) *Argon2IDHasher {
	if params == nil {
		params = DefaultArgon2IDParams
	}
	return &Argon2IDHasher{Params: params}
}

// Hash 使用 argon2id 哈希密码
func (h *Argon2IDHasher) Hash(password string) (string, error) {
	params := h.Params

	// 生成随机盐
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 使用 argon2id 哈希
	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	// 格式化结果: $argon2id$v=19$m=65536,t=3,p=2$c2FsdA$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// Verify 验证 argon2id 哈希密码
func (h *Argon2IDHasher) Verify(hashedPassword, password string) error {
	// 解析哈希字符串
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return fmt.Errorf("invalid hash format")
	}

	// 检查算法
	if parts[1] != "argon2id" {
		return fmt.Errorf("invalid algorithm")
	}

	// 检查版本
	if parts[2] != "v=19" {
		return fmt.Errorf("invalid version")
	}

	// 解析参数
	var memory, iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return err
	}

	// 解码盐和哈希
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return err
	}

	// 计算哈希
	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expectedHash)))

	// 比较哈希
	if subtle.ConstantTimeCompare(computedHash, expectedHash) != 1 {
		return fmt.Errorf("password mismatch")
	}

	return nil
}

// ScryptParams scrypt 参数
type ScryptParams struct {
	N       int
	R       int
	P       int
	SaltLen int
	KeyLen  int
}

// DefaultScryptParams 默认 scrypt 参数
var DefaultScryptParams = &ScryptParams{
	N:       32768, // 2^15
	R:       8,     // 8
	P:       1,     // 1
	SaltLen: 16,    // 16 字节盐
	KeyLen:  32,    // 32 字节密钥
}

// ScryptHasher scrypt 密码哈希器
type ScryptHasher struct {
	Params *ScryptParams
}

// NewScryptHasher 创建 scrypt 密码哈希器
func NewScryptHasher(params *ScryptParams) *ScryptHasher {
	if params == nil {
		params = DefaultScryptParams
	}
	return &ScryptHasher{Params: params}
}

// Hash 使用 scrypt 哈希密码
func (h *ScryptHasher) Hash(password string) (string, error) {
	params := h.Params

	// 生成随机盐
	salt := make([]byte, params.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 使用 scrypt 哈希
	hash, err := scrypt.Key([]byte(password), salt, params.N, params.R, params.P, params.KeyLen)
	if err != nil {
		return "", err
	}

	// 格式化结果: $scrypt$n=32768,r=8,p=1$c2FsdA$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$scrypt$n=%d,r=%d,p=%d$%s$%s",
		params.N, params.R, params.P, b64Salt, b64Hash)

	return encoded, nil
}

// Verify 验证 scrypt 哈希密码
func (h *ScryptHasher) Verify(hashedPassword, password string) error {
	// 解析哈希字符串
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 5 {
		return fmt.Errorf("invalid hash format")
	}

	// 检查算法
	if parts[1] != "scrypt" {
		return fmt.Errorf("invalid algorithm")
	}

	// 解析参数
	var n, r, p int
	_, err := fmt.Sscanf(parts[2], "n=%d,r=%d,p=%d", &n, &r, &p)
	if err != nil {
		return err
	}

	// 解码盐和哈希
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}

	// 计算哈希
	computedHash, err := scrypt.Key([]byte(password), salt, n, r, p, len(expectedHash))
	if err != nil {
		return err
	}

	// 比较哈希
	if subtle.ConstantTimeCompare(computedHash, expectedHash) != 1 {
		return fmt.Errorf("password mismatch")
	}

	return nil
}

// HashPassword 使用默认算法哈希密码
func HashPassword(password string) (string, error) {
	hasher := NewArgon2IDHasher(nil)
	return hasher.Hash(password)
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) error {
	// 根据哈希前缀选择验证方法
	if strings.HasPrefix(hashedPassword, "$2a$") || strings.HasPrefix(hashedPassword, "$2b$") || strings.HasPrefix(hashedPassword, "$2y$") {
		hasher := NewBcryptHasher(0)
		return hasher.Verify(hashedPassword, password)
	} else if strings.HasPrefix(hashedPassword, "$argon2id$") {
		hasher := NewArgon2IDHasher(nil)
		return hasher.Verify(hashedPassword, password)
	} else if strings.HasPrefix(hashedPassword, "$scrypt$") {
		hasher := NewScryptHasher(nil)
		return hasher.Verify(hashedPassword, password)
	}
	return fmt.Errorf("unknown hash format")
}
