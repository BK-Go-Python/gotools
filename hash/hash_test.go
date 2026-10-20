package hash

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestStringKnownVectors 覆盖 MD5/SHA-1/SHA-256/SHA-512 的标准测试向量。
// 期望值来自 Python 3.13 hashlib。
func TestStringKnownVectors(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		md5    string
		sha1   string
		sha256 string
		sha512 string
	}{
		{
			name:   "empty",
			input:  "",
			md5:    "d41d8cd98f00b204e9800998ecf8427e",
			sha1:   "da39a3ee5e6b4b0d3255bfef95601890afd80709",
			sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			sha512: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			name:   "hello",
			input:  "hello",
			md5:    "5d41402abc4b2a76b9719d911017c592",
			sha1:   "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
			sha256: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			sha512: "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
		},
		{
			name:   "abc",
			input:  "abc",
			md5:    "900150983cd24fb0d6963f7d28e17f72",
			sha1:   "a9993e364706816aba3e25717850c26c9cd0d89d",
			sha256: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
			sha512: "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f",
		},
		{
			name:   "utf8",
			input:  "你好",
			md5:    "7eca689f0d3389d9dea66ae112e5cfd7",
			sha1:   "440ee0853ad1e99f962b63e459ef992d7c211722",
			sha256: "670d9743542cae3ea7ebe36af56bd53648b0a1126162e78d81a32934a711302e",
			sha512: "5232181bc0d9888f5c9746e410b4740eb461706ba5dacfbc93587cecfc8d068bac7737e92870d6745b11a25e9cd78b55f4ffc706f73cfcae5345f1b53fb8f6b5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checks := []struct {
				name string
				got  string
				want string
			}{
				{"StringMD5", StringMD5(tt.input), tt.md5},
				{"StringSHA1", StringSHA1(tt.input), tt.sha1},
				{"StringSHA256", StringSHA256(tt.input), tt.sha256},
				{"StringSHA512", StringSHA512(tt.input), tt.sha512},
			}
			for _, c := range checks {
				if c.got != c.want {
					t.Errorf("%s(%q) = %q, want %q", c.name, tt.input, c.got, c.want)
				}
			}
		})
	}
}

// TestGenericHashAPI 覆盖 Algorithm/Sum/Hex/String 通用接口及大小写算法名。
func TestGenericHashAPI(t *testing.T) {
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	algorithmTests := []struct {
		name string
		alg  Algorithm
	}{
		{"lowercase", SHA256},
		{"uppercase", Algorithm("SHA256")},
		{"mixed case", Algorithm("ShA256")},
	}

	for _, tt := range algorithmTests {
		t.Run(tt.name, func(t *testing.T) {
			sum, err := Sum(tt.alg, []byte("hello"))
			if err != nil {
				t.Fatalf("Sum(%q) error: %v", tt.alg, err)
			}
			if got := hex.EncodeToString(sum); got != want {
				t.Errorf("Sum(%q) = %q, want %q", tt.alg, got, want)
			}

			gotHex, err := Hex(tt.alg, []byte("hello"))
			if err != nil {
				t.Fatalf("Hex(%q) error: %v", tt.alg, err)
			}
			if gotHex != want {
				t.Errorf("Hex(%q) = %q, want %q", tt.alg, gotHex, want)
			}

			gotString, err := String(tt.alg, "hello")
			if err != nil {
				t.Fatalf("String(%q) error: %v", tt.alg, err)
			}
			if gotString != want {
				t.Errorf("String(%q) = %q, want %q", tt.alg, gotString, want)
			}
		})
	}
}

func TestUnsupportedAlgorithm(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		call func() error
	}{
		{"Sum", func() error { _, err := Sum(Algorithm("sha384"), []byte("hello")); return err }},
		{"Hex", func() error { _, err := Hex(Algorithm("sha384"), []byte("hello")); return err }},
		{"String", func() error { _, err := String(Algorithm("sha384"), "hello"); return err }},
		{"File", func() error { _, err := File(Algorithm("sha384"), path); return err }},
		{"empty algorithm", func() error { _, err := Hex(Algorithm(""), []byte("hello")); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected unsupported algorithm error, got nil")
			}
			if !strings.Contains(err.Error(), "unsupported algorithm") {
				t.Errorf("error = %v, want it to mention unsupported algorithm", err)
			}
		})
	}
}

// TestFileHashes 覆盖四个文件哈希入口，并与旧版 Filetomd5、String 哈希交叉验证。
// 临时文件统一放在 t.TempDir 下。
func TestFileHashes(t *testing.T) {
	files := []struct {
		name   string
		data   []byte
		md5    string
		sha1   string
		sha256 string
		sha512 string
	}{
		{
			name:   "hello",
			data:   []byte("hello"),
			md5:    "5d41402abc4b2a76b9719d911017c592",
			sha1:   "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
			sha256: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			sha512: "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
		},
		{
			name:   "empty",
			data:   nil,
			md5:    "d41d8cd98f00b204e9800998ecf8427e",
			sha1:   "da39a3ee5e6b4b0d3255bfef95601890afd80709",
			sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			sha512: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
	}

	for _, f := range files {
		f := f
		t.Run(f.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "data.bin")
			if err := os.WriteFile(path, f.data, 0o600); err != nil {
				t.Fatal(err)
			}

			checks := []struct {
				name string
				fn   func(string) (string, error)
				want string
			}{
				{"FileMD5", FileMD5, f.md5},
				{"FileSHA1", FileSHA1, f.sha1},
				{"FileSHA256", FileSHA256, f.sha256},
				{"FileSHA512", FileSHA512, f.sha512},
				{"File(MD5)", func(p string) (string, error) { return File(MD5, p) }, f.md5},
			}
			for _, c := range checks {
				got, err := c.fn(path)
				if err != nil {
					t.Errorf("%s(%q) error: %v", c.name, path, err)
					continue
				}
				if got != c.want {
					t.Errorf("%s(%q) = %q, want %q", c.name, path, got, c.want)
				}
			}

			legacy, err := Filetomd5(path)
			if err != nil {
				t.Fatalf("Filetomd5(%q) error: %v", path, err)
			}
			if legacy != f.md5 {
				t.Errorf("Filetomd5(%q) = %q, want %q", path, legacy, f.md5)
			}

			// 同内容的内存哈希必须与文件哈希一致。
			if got := StringMD5(string(f.data)); got != f.md5 {
				t.Errorf("StringMD5(file content) = %q, want %q", got, f.md5)
			}
		})
	}
}

func TestFileHashErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.txt")
	tests := []struct {
		name string
		fn   func(string) (string, error)
	}{
		{"FileMD5", FileMD5},
		{"FileSHA1", FileSHA1},
		{"FileSHA256", FileSHA256},
		{"FileSHA512", FileSHA512},
		{"Filetomd5", Filetomd5},
		{"File(unsupported)", func(p string) (string, error) { return File(Algorithm("sha384"), p) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.fn(missing); err == nil {
				t.Fatalf("%s(%q) = nil error, want open error", tt.name, missing)
			}
		})
	}

	t.Run("directory", func(t *testing.T) {
		if _, err := FileSHA256(t.TempDir()); err == nil {
			t.Fatal("FileSHA256(directory) = nil error, want read error")
		}
	})
}

func TestDeprecatedStringMd5(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"hello", "hello"},
		{"utf8", "你好"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, want := StringMd5(tt.input), StringMD5(tt.input); got != want {
				t.Errorf("StringMd5(%q) = %q, want StringMD5 = %q", tt.input, got, want)
			}
		})
	}
}

// TestHMACKnownVectors 覆盖 RFC 4231 / RFC 2202 向量以及空 key/message 边界。
func TestHMACKnownVectors(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "HMACSHA256 RFC4231 case2",
			got:  HMACSHA256([]byte("Jefe"), []byte("what do ya want for nothing?")),
			want: "5bdcc146bf60754e6a042426089575c75a003f089d2739839dec58b964ec3843",
		},
		{
			name: "HMACSHA1 RFC2202 case2",
			got:  HMACSHA1([]byte("Jefe"), []byte("what do ya want for nothing?")),
			want: "effcdf6ae5eb2fa2d27416d5f184df9c259a7c79",
		},
		{
			name: "HMACSHA256 empty key and message",
			got:  HMACSHA256(nil, nil),
			want: "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad",
		},
		{
			name: "HMACSHA1 empty key and message",
			got:  HMACSHA1(nil, nil),
			want: "fbdb1d1b18aa6c08324b7d64b71fb76370690e1d",
		},
		{
			name: "HMACSHA256 RFC4231 case1",
			got:  HMACSHA256([]byte(strings.Repeat("\x0b", 20)), []byte("Hi There")),
			want: "b0344c61d8db38535ca8afceaf0bf12b881dc200c9833da726e9376c2e32cff7",
		},
		{
			name: "HMACSHA1 RFC2202 case1",
			got:  HMACSHA1([]byte(strings.Repeat("\x0b", 20)), []byte("Hi There")),
			want: "b617318655057264e28bc0b6fb378c8ef146be00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("HMAC = %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestRandomBytes(t *testing.T) {
	for _, n := range []int{1, 16, 32, 64} {
		n := n
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			got, err := RandomBytes(n)
			if err != nil {
				t.Fatalf("RandomBytes(%d) error: %v", n, err)
			}
			if len(got) != n {
				t.Errorf("len(RandomBytes(%d)) = %d, want %d", n, len(got), n)
			}

			other, err := RandomBytes(n)
			if err != nil {
				t.Fatalf("second RandomBytes(%d) error: %v", n, err)
			}
			if hex.EncodeToString(got) == hex.EncodeToString(other) {
				t.Errorf("two RandomBytes(%d) calls returned identical bytes", n)
			}
		})
	}

	for _, n := range []int{0, -1, -16} {
		n := n
		t.Run(fmt.Sprintf("invalid n=%d", n), func(t *testing.T) {
			if _, err := RandomBytes(n); err == nil {
				t.Fatalf("RandomBytes(%d) = nil error, want error", n)
			}
		})
	}
}

func TestRandomToken(t *testing.T) {
	for _, n := range []int{1, 16, 32} {
		n := n
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			got, err := RandomToken(n)
			if err != nil {
				t.Fatalf("RandomToken(%d) error: %v", n, err)
			}
			if len(got) != 2*n {
				t.Errorf("len(RandomToken(%d)) = %d, want %d", n, len(got), 2*n)
			}
			decoded, err := hex.DecodeString(got)
			if err != nil {
				t.Errorf("RandomToken(%d) = %q is not valid hex: %v", n, got, err)
			}
			if len(decoded) != n {
				t.Errorf("decoded len = %d, want %d", len(decoded), n)
			}

			other, err := RandomToken(n)
			if err != nil {
				t.Fatalf("second RandomToken(%d) error: %v", n, err)
			}
			if got == other {
				t.Errorf("two RandomToken(%d) calls returned identical tokens", n)
			}
		})
	}

	for _, n := range []int{0, -1, -8} {
		n := n
		t.Run(fmt.Sprintf("invalid n=%d", n), func(t *testing.T) {
			if _, err := RandomToken(n); err == nil {
				t.Fatalf("RandomToken(%d) = nil error, want error", n)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	seen := make(map[string]bool, 32)
	for i := 0; i < 32; i++ {
		got, err := UUID()
		if err != nil {
			t.Fatalf("UUID() error: %v", err)
		}
		if len(got) != 36 {
			t.Fatalf("len(UUID()) = %d, want 36: %q", len(got), got)
		}
		for _, pos := range []int{8, 13, 18, 23} {
			if got[pos] != '-' {
				t.Errorf("UUID() = %q, want '-' at index %d", got, pos)
			}
		}
		compact := strings.ReplaceAll(got, "-", "")
		if len(compact) != 32 {
			t.Errorf("UUID() = %q, want 32 hex digits", got)
		}
		if _, err := hex.DecodeString(compact); err != nil {
			t.Errorf("UUID() = %q, hex decode error: %v", got, err)
		}
		if got[14] != '4' {
			t.Errorf("UUID() = %q, version nibble = %q, want '4'", got, got[14])
		}
		if !strings.ContainsRune("89ab", rune(got[19])) {
			t.Errorf("UUID() = %q, variant nibble = %q, want one of 8/9/a/b", got, got[19])
		}
		if strings.ToLower(got) != got {
			t.Errorf("UUID() = %q, want lowercase hex", got)
		}
		if seen[got] {
			t.Errorf("UUID() returned duplicate %q", got)
		}
		seen[got] = true
	}
}
