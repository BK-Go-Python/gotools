package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Filetomd5 把文件内容转成 MD5 十六进制字符串。
//
// Deprecated: 请使用 FileMD5，它会正确关闭文件并返回更清晰的错误。
func Filetomd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("hash.Filetomd5: open %q: %w", filename, err)
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("hash.Filetomd5: read %q: %w", filename, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// StringMd5 把字符串转成 MD5 十六进制字符串。
//
// Deprecated: 请使用 StringMD5。
func StringMd5(s string) string {
	return StringMD5(s)
}
