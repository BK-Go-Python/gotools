package str

import (
	"strings"
	"unicode"
)

// Casefold 等价于 Python 的 s.casefold()，用于大小写无关比较。
//
// Python: s.casefold()
// DIFF: Go 没有完整的 Unicode case folding 表；这里覆盖了常见特殊映射
// （ß/ẞ → ss、İ → i + U+0307），其余字符退化为 ToLower。
func Casefold(s string) string {
	s = strings.ReplaceAll(s, "İ", "i\u0307")
	return casefoldSpecial(strings.ToLower(s))
}

// Expandtabs 等价于 Python 的 s.expandtabs(tabsize=8)。
//
// 制表符按当前列对齐补空格；遇到 \n 或 \r 时列号归零。
// tabsize <= 0 时删除制表符。
func Expandtabs(s string, tabsize ...int) string {
	size := 8
	if len(tabsize) > 0 {
		size = tabsize[0]
	}
	var b strings.Builder
	b.Grow(len(s) + 8)
	column := 0
	for _, r := range s {
		switch r {
		case '\t':
			spaces := 0
			if size > 0 {
				spaces = size - column%size
			}
			b.WriteString(strings.Repeat(" ", spaces))
			column += spaces
		case '\n', '\r':
			b.WriteRune(r)
			column = 0
		default:
			b.WriteRune(r)
			column++
		}
	}
	return b.String()
}

// IsIdentifier 等价于 Python 的 s.isidentifier()。
//
// Python: s.isidentifier()
// DIFF: 使用 Go 的 Letter/Digit 近似 XID_Start/XID_Continue，
// 未覆盖组合记号、连接符等全部 Unicode 标识符字符。
func IsIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if r != '_' && !unicode.IsLetter(r) {
				return false
			}
			continue
		}
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) &&
			!unicode.Is(unicode.Nl, r) && !unicode.Is(unicode.Mn, r) &&
			!unicode.Is(unicode.Mc, r) && !unicode.Is(unicode.Pc, r) {
			return false
		}
	}
	return true
}

// IsPrintable 等价于 Python 的 s.isprintable()：空字符串返回 true，
// 所有字符都是可打印字符或空格时返回 true。
//
// Python: s.isprintable()
// DIFF: 使用 unicode.IsPrint 近似 Python 的 printable 定义。
func IsPrintable(s string) bool {
	for _, r := range s {
		if r != ' ' && !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// SplitLinesKeepends 等价于 Python 的 s.splitlines(keepends=True)。
//
// 与 SplitLines 不同，返回的每一行保留原来的换行符（\n、\r\n、\r 等）。
func SplitLinesKeepends(s string) []string {
	if s == "" {
		return []string{}
	}
	r := []rune(s)
	lines := make([]string, 0, strings.Count(s, "\n")+1)
	start := 0
	for i := 0; i < len(r); i++ {
		if !isLineBreak(r[i]) {
			continue
		}
		end := i + 1
		if r[i] == '\r' && i+1 < len(r) && r[i+1] == '\n' {
			i++
			end = i + 1
		}
		lines = append(lines, string(r[start:end]))
		start = end
	}
	if start < len(r) {
		lines = append(lines, string(r[start:]))
	}
	return lines
}
