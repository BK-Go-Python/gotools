// Package urllib 提供 Python urllib.parse 风格的 URL 编解码与解析工具。
package urllib

import (
	"fmt"
	"net/url"
	"strings"
)

// Quote 等价于 Python 的 urllib.parse.quote(s, safe="/")：
// 保留字母、数字、`_.-~` 和 safe 中的字符，其余按 UTF-8 字节百分号编码。
// safe 省略时默认为 "/"。
func Quote(s string, safe ...string) string {
	allowed := "/"
	if len(safe) > 0 {
		allowed = strings.Join(safe, "")
	}
	var out strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if isUnreserved(b) || strings.IndexByte(allowed, b) >= 0 {
			out.WriteByte(b)
			continue
		}
		fmt.Fprintf(&out, "%%%02X", b)
	}
	return out.String()
}

// Unquote 等价于 Python 的 urllib.parse.unquote(s)。
//
// 非法百分号转义返回 error；Python 会保留原样。
func Unquote(s string) (string, error) {
	return url.PathUnescape(s)
}

// QuotePlus 等价于 Python 的 urllib.parse.quote_plus(s, safe="")：
// 空格编码为 '+'，默认不保留 '/'。
func QuotePlus(s string, safe ...string) string {
	allowed := ""
	if len(safe) > 0 {
		allowed = strings.Join(safe, "")
	}
	return strings.ReplaceAll(Quote(s, allowed), "%20", "+")
}

// UnquotePlus 等价于 Python 的 urllib.parse.unquote_plus(s)：'+' 还原为空格。
func UnquotePlus(s string) (string, error) {
	return url.QueryUnescape(s)
}

// Encode 等价于 Python 的 urllib.parse.urlencode(query)：按键排序并转义。
func Encode(values map[string][]string) string {
	return url.Values(values).Encode()
}

// ParseQuery 等价于 Python 的 urllib.parse.parse_qs(query)。
func ParseQuery(query string) (map[string][]string, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	return map[string][]string(values), nil
}

// Parse 等价于 Python 的 urllib.parse.urlparse(url)，返回 Go *url.URL。
func Parse(raw string) (*url.URL, error) {
	return url.Parse(raw)
}

// Join 等价于 Python 的 urllib.parse.urljoin(base, *parts)。
func Join(base string, parts ...string) (string, error) {
	return url.JoinPath(base, parts...)
}

func isUnreserved(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') ||
		b == '-' || b == '_' || b == '.' || b == '~'
}
