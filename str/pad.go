package str

import (
	"strings"
	"unicode/utf8"
)

// Center 等价于 Python 的 s.center(width[, fillchar])：把 s 居中到指定宽度。
//
// 与 Python 一致，多出来的填充字符放在左侧。
func Center(s string, width int, fill ...rune) string {
	n := utf8.RuneCountInString(s)
	if width <= n {
		return s
	}
	f := ' '
	if len(fill) > 0 {
		f = fill[0]
	}
	padding := width - n
	left := padding / 2
	if padding%2 == 1 && width%2 == 1 {
		left++
	}
	return strings.Repeat(string(f), left) + s + strings.Repeat(string(f), padding-left)
}

// Ljust 等价于 Python 的 s.ljust(width[, fillchar])：左对齐，右侧填充。
func Ljust(s string, width int, fill ...rune) string {
	n := utf8.RuneCountInString(s)
	if width <= n {
		return s
	}
	f := ' '
	if len(fill) > 0 {
		f = fill[0]
	}
	return s + strings.Repeat(string(f), width-n)
}

// Rjust 等价于 Python 的 s.rjust(width[, fillchar])：右对齐，左侧填充。
func Rjust(s string, width int, fill ...rune) string {
	n := utf8.RuneCountInString(s)
	if width <= n {
		return s
	}
	f := ' '
	if len(fill) > 0 {
		f = fill[0]
	}
	return strings.Repeat(string(f), width-n) + s
}

// Zfill 等价于 Python 的 s.zfill(width)：左侧补 0，并保留正负号位置。
//
//	Zfill("-42", 5) == "-0042"
func Zfill(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if width <= n {
		return s
	}
	zeros := strings.Repeat("0", width-n)
	if s != "" && (s[0] == '+' || s[0] == '-') {
		return s[:1] + zeros + s[1:]
	}
	return zeros + s
}
