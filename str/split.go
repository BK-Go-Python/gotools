package str

import (
	"strings"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// Split 等价于 Python 的 s.split(sep)。
//
// sep 为空字符串时按任意空白分割（等价于 s.split()），并忽略首尾空白。
func Split(s, sep string) []string {
	if sep == "" {
		return splitWhitespaceMax(s, -1)
	}
	return strings.Split(s, sep)
}

// SplitN 等价于 Python 的 s.split(sep, maxsplit)。
//
// 注意这里沿用 Python 的 maxsplit 语义：最多分割 maxsplit 次，返回最多
// maxsplit+1 个元素；maxsplit 为负数表示不限制。这与 strings.SplitN 的
// n 参数含义不同。
func SplitN(s, sep string, maxsplit int) []string {
	if maxsplit < 0 {
		return Split(s, sep)
	}
	if sep == "" {
		return splitWhitespaceMax(s, maxsplit)
	}
	return strings.SplitN(s, sep, maxsplit+1)
}

// Rsplit 等价于 Python 的 s.rsplit(sep, maxsplit)：从右往左分割。
func Rsplit(s, sep string, maxsplit int) []string {
	if maxsplit < 0 {
		return Split(s, sep)
	}
	if sep == "" {
		return rsplitWhitespaceMax(s, maxsplit)
	}
	if maxsplit == 0 {
		return []string{s}
	}
	parts := make([]string, 0, maxsplit+1)
	end := len(s)
	for count := 0; count < maxsplit; count++ {
		idx := strings.LastIndex(s[:end], sep)
		if idx < 0 {
			break
		}
		parts = append(parts, s[idx+len(sep):end])
		end = idx
	}
	parts = append(parts, s[:end])
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return parts
}

// SplitLines 等价于 Python 的 s.splitlines()：按 \n、\r、\r\n 等换行符
// 切分，不保留行尾，也不会因为末尾换行而多出一个空字符串。
func SplitLines(s string) []string {
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
		lines = append(lines, string(r[start:i]))
		if r[i] == '\r' && i+1 < len(r) && r[i+1] == '\n' {
			i++
		}
		start = i + 1
	}
	if start < len(r) {
		lines = append(lines, string(r[start:]))
	}
	return lines
}

func isLineBreak(r rune) bool {
	switch r {
	case '\n', '\r', '\v', '\f', '\x1c', '\x1d', '\x1e', '\u0085', '\u2028', '\u2029':
		return true
	}
	return false
}

// Fields 等价于 Python 的 s.split()：按任意空白分割并忽略空字段。
func Fields(s string) []string {
	return splitWhitespaceMax(s, -1)
}

// Join 等价于 Python 的 sep.join(list)，参数以 []string 传入。
func Join(sep string, elems []string) string {
	return strings.Join(elems, sep)
}

// JoinAny 等价于 Python 的 sep.join(iterable)，元素会先按 py.Str 转换。
func JoinAny(sep string, elems []interface{}) string {
	parts := make([]string, len(elems))
	for i, elem := range elems {
		parts[i] = pyval.Str(elem)
	}
	return strings.Join(parts, sep)
}

// Partition 等价于 Python 的 s.partition(sep)，返回 (head, sep, tail)。
// 未找到时返回 (s, "", "")；sep 为空时同样返回 (s, "", "")。
func Partition(s, sep string) (string, string, string) {
	if sep == "" {
		return s, "", ""
	}
	idx := strings.Index(s, sep)
	if idx < 0 {
		return s, "", ""
	}
	return s[:idx], sep, s[idx+len(sep):]
}

// Rpartition 等价于 Python 的 s.rpartition(sep)，从右往左查找。
func Rpartition(s, sep string) (string, string, string) {
	if sep == "" {
		return "", "", s
	}
	idx := strings.LastIndex(s, sep)
	if idx < 0 {
		return "", "", s
	}
	return s[:idx], sep, s[idx+len(sep):]
}

func splitWhitespaceMax(s string, maxsplit int) []string {
	r := []rune(s)
	parts := make([]string, 0, maxsplit+1)
	i := 0
	for {
		for i < len(r) && isPythonSpace(r[i]) {
			i++
		}
		if i >= len(r) {
			break
		}
		if len(parts) == maxsplit {
			parts = append(parts, strings.TrimLeftFunc(string(r[i:]), isPythonSpace))
			return parts
		}
		start := i
		for i < len(r) && !isPythonSpace(r[i]) {
			i++
		}
		parts = append(parts, string(r[start:i]))
	}
	return parts
}

func rsplitWhitespaceMax(s string, maxsplit int) []string {
	r := []rune(s)
	end := len(r)
	if maxsplit == 0 {
		trimmed := strings.TrimRightFunc(s, isPythonSpace)
		if trimmed == "" {
			return []string{}
		}
		return []string{trimmed}
	}
	for end > 0 && isPythonSpace(r[end-1]) {
		end--
	}
	if end == 0 {
		return []string{}
	}

	parts := make([]string, 0, maxsplit+1)
	count := 0
	i := end
	for i > 0 {
		j := i
		for j > 0 && !isPythonSpace(r[j-1]) {
			j--
		}
		if count == maxsplit {
			parts = append(parts, strings.TrimRightFunc(string(r[:i]), isPythonSpace))
			break
		}
		parts = append(parts, string(r[j:i]))
		count++
		i = j
		for i > 0 && isPythonSpace(r[i-1]) {
			i--
		}
	}
	for a, b := 0, len(parts)-1; a < b; a, b = a+1, b-1 {
		parts[a], parts[b] = parts[b], parts[a]
	}
	return parts
}
