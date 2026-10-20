package str

import (
	"strings"
	"unicode"
)

// Upper 等价于 Python 的 s.upper()。
//
// DIFF: Go 的简单大小写映射不做多字符展开；这里额外处理了常见的 ß→SS，
// 其余完整 Unicode 特殊映射（连字等）尚未覆盖。
func Upper(s string) string {
	return upperSpecial(strings.ToUpper(s))
}

// Lower 等价于 Python 的 s.lower()。
//
// DIFF: 土耳其大写 İ 按 Python 规则展开为 i + U+0307；其他特殊映射未完全覆盖。
func Lower(s string) string {
	if !strings.ContainsRune(s, 'İ') {
		return strings.ToLower(s)
	}
	var b strings.Builder
	b.Grow(len(s) + 2)
	for _, r := range s {
		if r == 'İ' {
			b.WriteString("i\u0307")
			continue
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// Capitalize 等价于 Python 的 s.capitalize()：首字符大写，其余全部小写。
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	if expanded, ok := titleLigature(r[0]); ok {
		return expanded + strings.ToLower(string(r[1:]))
	}
	if r[0] == 'ß' {
		return "Ss" + strings.ToLower(string(r[1:]))
	}
	r[0] = unicode.ToTitle(r[0])
	for i := 1; i < len(r); i++ {
		r[i] = unicode.ToLower(r[i])
	}
	return string(r)
}

// Title 等价于 Python 的 s.title()：每个单词的首字母大写，其余字母小写。
//
// 单词边界是非字母字符，因此 "they're".Title() 得到 "They'Re"，与 Python 一致。
func Title(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			if prevLetter {
				b.WriteRune(unicode.ToLower(r))
			} else if expanded, ok := titleLigature(r); ok {
				b.WriteString(expanded)
			} else if r == 'ß' {
				b.WriteString("Ss")
			} else {
				b.WriteRune(unicode.ToTitle(r))
			}
			prevLetter = true
			continue
		}
		b.WriteRune(r)
		prevLetter = false
	}
	return b.String()
}

// Swapcase 等价于 Python 的 s.swapcase()：大写变小写，小写变大写。
func Swapcase(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if expanded, ok := swapcaseSpecial(r); ok {
			b.WriteString(expanded)
			continue
		}
		switch {
		case unicode.IsUpper(r):
			b.WriteRune(unicode.ToLower(r))
		case unicode.IsLower(r):
			b.WriteRune(unicode.ToUpper(r))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
