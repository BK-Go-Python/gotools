package str

import (
	"fmt"
	"strings"
)

// Maketrans 等价于 Python 的 str.maketrans(from, to)：按 rune 逐位对应，
// 返回可直接传给 Translate 的映射表。
//
// from 与 to 的 rune 长度必须相同，否则返回 error 和 nil 表。from 中重复的
// rune 以最后一次映射为准（与 Python dict 的行为一致）。两者都为空时返回
// 空表。
//
// Python: str.maketrans(from, to)
// DIFF: Python 的 maketrans 还支持单个 dict 参数等重载形式，这里只实现
// from/to 两个字符串的用法。
func Maketrans(from, to string) (map[rune]string, error) {
	fromRunes := []rune(from)
	toRunes := []rune(to)
	if len(fromRunes) != len(toRunes) {
		return nil, fmt.Errorf("str: maketrans arguments have different lengths: %d != %d",
			len(fromRunes), len(toRunes))
	}
	table := make(map[rune]string, len(fromRunes))
	for i, r := range fromRunes {
		table[r] = string(toRunes[i])
	}
	return table, nil
}

// Translate 等价于 Python 的 s.translate(table)：逐个 rune 查表，命中时用
// 表中的值替换，未命中时保留原字符。值可以是多个 rune。
//
// table 为 nil 或空表时原样返回 s。
//
// Python: s.translate(table)
// DIFF: Python 用 None 表示删除字符，而 Go 的 map[rune]string 不能存 None，
// 因此这里约定空字符串 "" 表示删除该字符。
func Translate(s string, table map[rune]string) string {
	if len(table) == 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		replacement, ok := table[r]
		if !ok {
			b.WriteRune(r)
			continue
		}
		b.WriteString(replacement) // "" 表示删除，见上面的 DIFF。
	}
	return b.String()
}
