// Package html 提供 Python html 模块风格的转义/反转义。
package html

import (
	stdhtml "html"
	"strings"
)

var (
	escapeWithQuotes = strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#x27;",
	)
	escapeWithoutQuotes = strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
)

// Escape 等价于 Python 的 html.escape(s, quote=True)。
//
// quote 省略时默认为 true；false 时只转义 &、<、>。
func Escape(s string, quote ...bool) string {
	escapeQuotes := true
	if len(quote) > 0 {
		escapeQuotes = quote[0]
	}
	if escapeQuotes {
		return escapeWithQuotes.Replace(s)
	}
	return escapeWithoutQuotes.Replace(s)
}

// Unescape 等价于 Python 的 html.unescape(s)。
func Unescape(s string) string {
	return stdhtml.UnescapeString(s)
}
