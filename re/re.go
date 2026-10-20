// Package re 提供 Python re 模块风格的高频正则 API。
//
// 底层使用 Go regexp（RE2）：不支持回溯引用、lookaround、条件表达式等
// Python 扩展语法；替换模板使用 Go 的 $1/${name} 形式而不是 \1。
package re

import (
	"fmt"
	stdre "regexp"
)

// Regexp 是编译后的正则表达式。
type Regexp struct {
	expr *stdre.Regexp
}

// Compile 等价于 Python 的 re.compile(pattern)。
func Compile(pattern string) (*Regexp, error) {
	expr, err := stdre.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("re.Compile: %w", err)
	}
	return &Regexp{expr: expr}, nil
}

// SearchString 等价于 Python 的 pattern.search(s) 的布尔结果。
func (r *Regexp) SearchString(s string) bool {
	return r != nil && r.expr.MatchString(s)
}

// FindString 返回第一个匹配；未匹配返回 ok=false。
func (r *Regexp) FindString(s string) (string, bool) {
	if r == nil {
		return "", false
	}
	match := r.expr.FindString(s)
	return match, match != ""
}

// FindAllString 返回所有匹配；n<0 表示全部。
func (r *Regexp) FindAllString(s string, n int) []string {
	if r == nil {
		return []string{}
	}
	return r.expr.FindAllString(s, n)
}

// Sub 等价于 Python 的 re.sub(pattern, repl, s, count=0)。
//
// 替换模板使用 Go 语法（$1、${name}），不是 Python 的 \1。
func (r *Regexp) Sub(repl, s string, count ...int) string {
	if r == nil {
		return s
	}
	limit := -1
	if len(count) > 0 && count[0] > 0 {
		limit = count[0]
	}
	if limit < 0 {
		return r.expr.ReplaceAllString(s, repl)
	}
	matches := r.expr.FindAllStringSubmatchIndex(s, limit)
	if len(matches) == 0 {
		return s
	}
	var out []byte
	last := 0
	for _, match := range matches {
		out = append(out, s[last:match[0]]...)
		out = r.expr.ExpandString(out, repl, s, match)
		last = match[1]
	}
	out = append(out, s[last:]...)
	return string(out)
}

// Split 等价于 Python 的 re.split(pattern, s, maxsplit=0)。
func (r *Regexp) Split(s string, maxsplit ...int) []string {
	if r == nil {
		return []string{s}
	}
	n := -1
	if len(maxsplit) > 0 {
		switch {
		case maxsplit[0] == 0:
			n = 1
		case maxsplit[0] > 0:
			n = maxsplit[0] + 1
		}
	}
	return r.expr.Split(s, n)
}

// Match 等价于 Python 的 re.match(pattern, s)：只匹配字符串开头。
func Match(pattern, s string) (bool, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return false, err
	}
	index := compiled.expr.FindStringIndex(s)
	return index != nil && index[0] == 0, nil
}

// FullMatch 等价于 Python 的 re.fullmatch(pattern, s)：整串匹配。
func FullMatch(pattern, s string) (bool, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return false, err
	}
	index := compiled.expr.FindStringIndex(s)
	return index != nil && index[0] == 0 && index[1] == len(s), nil
}

// Search 等价于 Python 的 re.search(pattern, s) 的布尔结果。
func Search(pattern, s string) (bool, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return false, err
	}
	return compiled.expr.MatchString(s), nil
}

// FindAll 等价于 Python 的 re.findall(pattern, s)；只返回整体匹配，不返回分组。
func FindAll(pattern, s string) ([]string, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return nil, err
	}
	return compiled.expr.FindAllString(s, -1), nil
}

// FindString 等价于 Python 的 re.search(...).group(0)。
func FindString(pattern, s string) (string, bool, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return "", false, err
	}
	match := compiled.expr.FindString(s)
	return match, match != "", nil
}

// Sub 等价于 Python 的 re.sub(pattern, repl, s, count=0)。
func Sub(pattern, repl, s string, count ...int) (string, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return "", err
	}
	return compiled.Sub(repl, s, count...), nil
}

// Split 等价于 Python 的 re.split(pattern, s, maxsplit=0)。
func Split(pattern, s string, maxsplit ...int) ([]string, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return nil, err
	}
	return compiled.Split(s, maxsplit...), nil
}

// QuoteMeta 等价于 Python 的 re.escape(s)。
func QuoteMeta(s string) string {
	return stdre.QuoteMeta(s)
}

// FindStringSubmatch 返回第一个匹配及其分组；未匹配返回 nil。
//
// 下标 0 是整体匹配，1..n 是捕获组，语义与 Python 的 match.group(0..n) 一致。
func (r *Regexp) FindStringSubmatch(s string) []string {
	if r == nil {
		return nil
	}
	return r.expr.FindStringSubmatch(s)
}

// FindAllStringSubmatch 返回所有匹配及分组；n<0 表示全部。
func (r *Regexp) FindAllStringSubmatch(s string, n int) [][]string {
	if r == nil {
		return nil
	}
	return r.expr.FindAllStringSubmatch(s, n)
}

// FindStringSubmatch 等价于 Python 的 re.search(pattern, s).groups() 加整体匹配。
func FindStringSubmatch(pattern, s string) ([]string, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return nil, err
	}
	return compiled.FindStringSubmatch(s), nil
}

// FindAllStringSubmatch 等价于 Python 的 re.findall(pattern, s) 的分组形式。
func FindAllStringSubmatch(pattern, s string) ([][]string, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return nil, err
	}
	return compiled.FindAllStringSubmatch(s, -1), nil
}

// SubexpNames 返回捕获组名称；未命名组为空字符串，下标与 FindStringSubmatch 对齐。
func (r *Regexp) SubexpNames() []string {
	if r == nil {
		return nil
	}
	return r.expr.SubexpNames()
}

// FindStringSubmatchMap 返回命名分组的 map；未匹配的命名分组值为空字符串。
func (r *Regexp) FindStringSubmatchMap(s string) map[string]string {
	groups := r.FindStringSubmatch(s)
	if groups == nil {
		return nil
	}
	names := r.SubexpNames()
	result := make(map[string]string)
	for i, name := range names {
		if i == 0 || name == "" {
			continue
		}
		if i < len(groups) {
			result[name] = groups[i]
		} else {
			result[name] = ""
		}
	}
	return result
}

// FindAllStringSubmatchMap 返回所有匹配的命名分组 map；n<0 表示全部。
func (r *Regexp) FindAllStringSubmatchMap(s string, n int) []map[string]string {
	all := r.FindAllStringSubmatch(s, n)
	if all == nil {
		return nil
	}
	names := r.SubexpNames()
	out := make([]map[string]string, 0, len(all))
	for _, groups := range all {
		result := make(map[string]string)
		for i, name := range names {
			if i == 0 || name == "" {
				continue
			}
			if i < len(groups) {
				result[name] = groups[i]
			} else {
				result[name] = ""
			}
		}
		out = append(out, result)
	}
	return out
}

// FindStringSubmatchMap 等价于 Python 的 match.groupdict()。
func FindStringSubmatchMap(pattern, s string) (map[string]string, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return nil, err
	}
	return compiled.FindStringSubmatchMap(s), nil
}

// FindAllStringSubmatchMap 返回所有匹配的命名分组 map。
func FindAllStringSubmatchMap(pattern, s string) ([]map[string]string, error) {
	compiled, err := Compile(pattern)
	if err != nil {
		return nil, err
	}
	return compiled.FindAllStringSubmatchMap(s, -1), nil
}
