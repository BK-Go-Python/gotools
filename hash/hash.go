// Package hash 提供常用哈希与随机串工具：MD5、SHA-1、SHA-256、SHA-512、
// HMAC、文件哈希、随机 token 和 UUID v4。
package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

// Algorithm 是哈希算法名称。
type Algorithm string

// 支持的哈希算法。
const (
	MD5    Algorithm = "md5"
	SHA1   Algorithm = "sha1"
	SHA256 Algorithm = "sha256"
	SHA512 Algorithm = "sha512"
)

func newHash(algorithm Algorithm) (hash.Hash, error) {
	switch Algorithm(strings.ToLower(string(algorithm))) {
	case MD5:
		return md5.New(), nil
	case SHA1:
		return sha1.New(), nil
	case SHA256:
		return sha256.New(), nil
	case SHA512:
		return sha512.New(), nil
	}
	return nil, fmt.Errorf("hash: unsupported algorithm %q (want md5, sha1, sha256 or sha512)", algorithm)
}

// Sum 返回数据的原始摘要字节。
func Sum(algorithm Algorithm, data []byte) ([]byte, error) {
	h, err := newHash(algorithm)
	if err != nil {
		return nil, err
	}
	h.Write(data)
	return h.Sum(nil), nil
}

// Hex 返回数据的十六进制摘要字符串。
func Hex(algorithm Algorithm, data []byte) (string, error) {
	sum, err := Sum(algorithm, data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sum), nil
}

// String 返回字符串的十六进制摘要。
func String(algorithm Algorithm, s string) (string, error) {
	return Hex(algorithm, []byte(s))
}

// File 返回文件的十六进制摘要。
func File(algorithm Algorithm, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("hash.File: open %q: %w", path, err)
	}
	defer file.Close()

	h, err := newHash(algorithm)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("hash.File: read %q: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// StringMD5 返回字符串的 MD5 十六进制摘要。
func StringMD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// StringSHA1 返回字符串的 SHA-1 十六进制摘要。
func StringSHA1(s string) string {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// StringSHA256 返回字符串的 SHA-256 十六进制摘要。
func StringSHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// StringSHA512 返回字符串的 SHA-512 十六进制摘要。
func StringSHA512(s string) string {
	sum := sha512.Sum512([]byte(s))
	return hex.EncodeToString(sum[:])
}

// FileMD5 返回文件的 MD5 十六进制摘要。
func FileMD5(path string) (string, error) {
	return File(MD5, path)
}

// FileSHA1 返回文件的 SHA-1 十六进制摘要。
func FileSHA1(path string) (string, error) {
	return File(SHA1, path)
}

// FileSHA256 返回文件的 SHA-256 十六进制摘要。
func FileSHA256(path string) (string, error) {
	return File(SHA256, path)
}

// FileSHA512 返回文件的 SHA-512 十六进制摘要。
func FileSHA512(path string) (string, error) {
	return File(SHA512, path)
}

// HMACSHA256 返回 HMAC-SHA256 的十六进制摘要。
func HMACSHA256(key, message []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}

// HMACSHA1 返回 HMAC-SHA1 的十六进制摘要。
func HMACSHA1(key, message []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}

// RandomBytes 返回 n 个加密安全的随机字节。
func RandomBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("hash.RandomBytes: size must be positive, got %d", n)
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, fmt.Errorf("hash.RandomBytes: %w", err)
	}
	return b, nil
}

// RandomToken 返回 n 字节随机数的十六进制字符串，长度为 2n。
func RandomToken(n int) (string, error) {
	b, err := RandomBytes(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// UUID 生成 RFC 4122 版本 4 的 UUID。
func UUID() (string, error) {
	var b [16]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		return "", fmt.Errorf("hash.UUID: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
