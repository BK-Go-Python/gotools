// Package base64 提供 Python base64 模块风格的编解码函数。
package base64

import stdbase64 "encoding/base64"

// B64encode 等价于 Python 的 base64.b64encode(data)，返回带 padding 的标准 Base64。
func B64encode(data []byte) string {
	return stdbase64.StdEncoding.EncodeToString(data)
}

// B64decode 等价于 Python 的 base64.b64decode(s)。
func B64decode(s string) ([]byte, error) {
	return stdbase64.StdEncoding.DecodeString(s)
}

// Urlsafe_b64encode 等价于 Python 的 base64.urlsafe_b64encode(data)。
func Urlsafe_b64encode(data []byte) string {
	return stdbase64.URLEncoding.EncodeToString(data)
}

// Urlsafe_b64decode 等价于 Python 的 base64.urlsafe_b64decode(s)。
func Urlsafe_b64decode(s string) ([]byte, error) {
	return stdbase64.URLEncoding.DecodeString(s)
}
