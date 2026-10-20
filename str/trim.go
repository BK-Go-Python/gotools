package str

import "strings"

// Strip 等价于 Python 的 s.strip()：去掉首尾的 Unicode 空白字符。
func Strip(s string) string {
	return strings.TrimFunc(s, isPythonSpace)
}

// StripChars 等价于 Python 的 s.strip(chars)：去掉首尾出现在 cutset 中的字符。
func StripChars(s, cutset string) string {
	return strings.Trim(s, cutset)
}

// Lstrip 等价于 Python 的 s.lstrip()。
func Lstrip(s string) string {
	return strings.TrimLeftFunc(s, isPythonSpace)
}

// LstripChars 等价于 Python 的 s.lstrip(chars)。
func LstripChars(s, cutset string) string {
	return strings.TrimLeft(s, cutset)
}

// Rstrip 等价于 Python 的 s.rstrip()。
func Rstrip(s string) string {
	return strings.TrimRightFunc(s, isPythonSpace)
}

// RstripChars 等价于 Python 的 s.rstrip(chars)。
func RstripChars(s, cutset string) string {
	return strings.TrimRight(s, cutset)
}
