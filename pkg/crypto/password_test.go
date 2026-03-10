package crypto

import (
	"strings"
	"testing"
)

func TestBcryptHasher(t *testing.T) {
	hasher := NewBcryptHasher(10)
	password := "test_password123"

	// 测试哈希
	hashed, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if hashed == "" {
		t.Fatal("Hash returned empty string")
	}

	// 测试验证正确密码
	err = hasher.Verify(hashed, password)
	if err != nil {
		t.Fatalf("Verify failed for correct password: %v", err)
	}

	// 测试验证错误密码
	err = hasher.Verify(hashed, "wrong_password")
	if err == nil {
		t.Fatal("Verify should have failed for wrong password")
	}
}

func TestArgon2IDHasher(t *testing.T) {
	hasher := NewArgon2IDHasher(nil)
	password := "test_password123"

	// 测试哈希
	hashed, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if hashed == "" {
		t.Fatal("Hash returned empty string")
	}

	if !strings.HasPrefix(hashed, "$argon2id$") {
		t.Fatal("Hash should start with $argon2id$")
	}

	// 测试验证正确密码
	err = hasher.Verify(hashed, password)
	if err != nil {
		t.Fatalf("Verify failed for correct password: %v", err)
	}

	// 测试验证错误密码
	err = hasher.Verify(hashed, "wrong_password")
	if err == nil {
		t.Fatal("Verify should have failed for wrong password")
	}
}

func TestScryptHasher(t *testing.T) {
	hasher := NewScryptHasher(nil)
	password := "test_password123"

	// 测试哈希
	hashed, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if hashed == "" {
		t.Fatal("Hash returned empty string")
	}

	if !strings.HasPrefix(hashed, "$scrypt$") {
		t.Fatal("Hash should start with $scrypt$")
	}

	// 测试验证正确密码
	err = hasher.Verify(hashed, password)
	if err != nil {
		t.Fatalf("Verify failed for correct password: %v", err)
	}

	// 测试验证错误密码
	err = hasher.Verify(hashed, "wrong_password")
	if err == nil {
		t.Fatal("Verify should have failed for wrong password")
	}
}

func TestHashPassword(t *testing.T) {
	password := "test_password123"

	// 测试默认哈希
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashed == "" {
		t.Fatal("HashPassword returned empty string")
	}

	if !strings.HasPrefix(hashed, "$argon2id$") {
		t.Fatal("HashPassword should return argon2id hash")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "test_password123"

	// 测试 bcrypt
	bcryptHasher := NewBcryptHasher(10)
	bcryptHashed, err := bcryptHasher.Hash(password)
	if err != nil {
		t.Fatalf("Bcrypt hash failed: %v", err)
	}

	err = VerifyPassword(bcryptHashed, password)
	if err != nil {
		t.Fatalf("VerifyPassword failed for bcrypt: %v", err)
	}

	// 测试 argon2id
	argon2idHasher := NewArgon2IDHasher(nil)
	argon2idHashed, err := argon2idHasher.Hash(password)
	if err != nil {
		t.Fatalf("Argon2id hash failed: %v", err)
	}

	err = VerifyPassword(argon2idHashed, password)
	if err != nil {
		t.Fatalf("VerifyPassword failed for argon2id: %v", err)
	}

	// 测试 scrypt
	scryptHasher := NewScryptHasher(nil)
	scryptHashed, err := scryptHasher.Hash(password)
	if err != nil {
		t.Fatalf("Scrypt hash failed: %v", err)
	}

	err = VerifyPassword(scryptHashed, password)
	if err != nil {
		t.Fatalf("VerifyPassword failed for scrypt: %v", err)
	}

	// 测试错误密码
	err = VerifyPassword(argon2idHashed, "wrong_password")
	if err == nil {
		t.Fatal("VerifyPassword should have failed for wrong password")
	}
}
