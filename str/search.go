package str

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Contains 等价于 Python 的 sub in s。
func Contains(s, sub string) bool {
	return strings.Contains(s, sub)
}

// Find 等价于 Python 的 s.find(sub)，返回字符（rune）下标，未找到返回 -1。
func Find(s, sub string) int {
	idx := strings.Index(s, sub)
	if idx < 0 {
		return -1
	}
	return utf8.RuneCountInString(s[:idx])
}

// Rfind 等价于 Python 的 s.rfind(sub)，从右往左查找，返回字符下标。
func Rfind(s, sub string) int {
	idx := strings.LastIndex(s, sub)
	if idx < 0 {
		return -1
	}
	return utf8.RuneCountInString(s[:idx])
}

// Index 等价于 Python 的 s.index(sub)：找到返回字符下标，找不到返回 error。
func Index(s, sub string) (int, error) {
	if idx := Find(s, sub); idx >= 0 {
		return idx, nil
	}
	return -1, fmt.Errorf("str: substring %q not found", sub)
}

// Rindex 等价于 Python 的 s.rindex(sub)。
func Rindex(s, sub string) (int, error) {
	if idx := Rfind(s, sub); idx >= 0 {
		return idx, nil
	}
	return -1, fmt.Errorf("str: substring %q not found", sub)
}

// Count 等价于 Python 的 s.count(sub)，统计不重叠出现次数。
func Count(s, sub string) int {
	return strings.Count(s, sub)
}

// Startswith 等价于 Python 的 s.startswith(prefix)：任一前缀匹配即为 true。
func Startswith(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// Endswith 等价于 Python 的 s.endswith(suffix)：任一后缀匹配即为 true。
func Endswith(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}
