package str

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Len 等价于 Python 的 len(s)，返回字符（rune）数量。
func Len(s string) int {
	return utf8.RuneCountInString(s)
}

// LenBytes 返回字符串的字节长度，等价于 Go 的 len(s)。
func LenBytes(s string) int {
	return len(s)
}

// At 等价于 Python 的 s[index]，支持负数下标（-1 表示最后一个字符）。
func At(s string, index int) (rune, error) {
	r := []rune(s)
	if index < 0 {
		index += len(r)
	}
	if index < 0 || index >= len(r) {
		return 0, fmt.Errorf("str: string index out of range: %d", index)
	}
	return r[index], nil
}

// Slice 等价于 Python 的 s[start:end]，支持负数下标，越界会自动截断。
//
//	end 为 -1 之类的负值表示从末尾计算；若要取到末尾，请传 len(s) 或更大的值。
func Slice(s string, start, end int) string {
	r := []rune(s)
	n := len(r)
	if start < 0 {
		start += n
	}
	if end < 0 {
		end += n
	}
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}
	if end < 0 {
		end = 0
	}
	if end > n {
		end = n
	}
	if start >= end {
		return ""
	}
	return string(r[start:end])
}

// Reverse 等价于 Python 的 s[::-1]。
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Ord 等价于 Python 的 ord(s)：返回单个字符的 Unicode 码点。
func Ord(s string) (rune, error) {
	r := []rune(s)
	if len(r) != 1 {
		return 0, fmt.Errorf("str: ord() expected a character, but string of length %d found", len(r))
	}
	return r[0], nil
}

// Chr 等价于 Python 的 chr(i)：把 Unicode 码点转成字符串。
// 非法码点（负数、超过 U+10FFFF 或代理区）会返回 error。
func Chr(i int) (string, error) {
	if i < 0 || i > unicode.MaxRune || (i >= 0xD800 && i <= 0xDFFF) {
		return "", fmt.Errorf("str: chr() arg not in range: %d", i)
	}
	return string(rune(i)), nil
}

// Repeat 等价于 Python 的 s * n；n <= 0 时返回空字符串。
func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}

// RemovePrefix 等价于 Python 3.9+ 的 s.removeprefix(prefix)。
func RemovePrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// RemoveSuffix 等价于 Python 3.9+ 的 s.removesuffix(suffix)。
func RemoveSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}
