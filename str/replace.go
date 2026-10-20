package str

import "strings"

// Replace 等价于 Python 的 s.replace(old, new[, count])。
//
// count 省略或小于 0 时替换全部；count == 0 时不替换。old 为空字符串时，
// 在字符（rune）边界插入 new，与 Python 一致。count 统计的是不重叠的出现
// 次数，而不是 rune 数。多传的 count 参数会被忽略。
//
// Python: s.replace(old, new, count=-1)
func Replace(s, old, new string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	return strings.Replace(s, old, new, n)
}
