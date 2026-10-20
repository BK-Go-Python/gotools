package conv

import (
	"fmt"
	"unicode/utf8"
)

// Encode 等价于 Python 的 str.encode("utf-8")，返回 s 的 UTF-8 字节。
//
// Python: s.encode("utf-8")
// DIFF: Go 字符串可以携带非法 UTF-8 字节，此时 Encode 原样返回，不报错。
func Encode(s string) []byte {
	return []byte(s)
}

// Decode 等价于 Python 的 bytes.decode("utf-8")：data 必须是合法 UTF-8，
// 否则返回 error。
//
// Python: data.decode("utf-8")
func Decode(data []byte) (string, error) {
	if !utf8.Valid(data) {
		return "", fmt.Errorf("conv.Decode: invalid UTF-8 byte sequence")
	}
	return string(data), nil
}
