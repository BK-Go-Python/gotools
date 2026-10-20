// Package string 提供 Python string 模块的高频工具：capwords、Template 与常量。
package string

import (
	"fmt"
	"strings"
	"unicode"
)

// 常用字符集合，对应 Python string 模块常量。
const (
	Ascii_lowercase = "abcdefghijklmnopqrstuvwxyz"
	Ascii_uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Ascii_letters   = Ascii_lowercase + Ascii_uppercase
	Digits          = "0123456789"
	Hexdigits       = "0123456789abcdefABCDEF"
	Octdigits       = "01234567"
	Punctuation     = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Whitespace      = " \t\n\r\v\f"
	Printable       = Digits + Ascii_letters + Punctuation + Whitespace
)

// Capwords 等价于 Python 的 string.capwords(s)：按空白切词，每个词首字母
// 大写、其余小写，再用单个空格连接。
func Capwords(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = capitalize(word)
	}
	return strings.Join(words, " ")
}

func capitalize(word string) string {
	runes := []rune(strings.ToLower(word))
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Template 等价于 Python 的 string.Template。
type Template struct {
	text string
}

// NewTemplate 创建模板。
func NewTemplate(text string) *Template {
	return &Template{text: text}
}

// Template 返回原始模板文本。
func (t *Template) Template() string { return t.text }

// Substitute 等价于 Python 的 Template.substitute(mapping)：
// 缺失变量或非法占位符返回 error。
func (t *Template) Substitute(mapping map[string]string) (string, error) {
	return t.render(mapping, false)
}

// Safe_substitute 等价于 Python 的 Template.safe_substitute(mapping)：
// 缺失变量或非法占位符保持原样。
func (t *Template) Safe_substitute(mapping map[string]string) string {
	out, _ := t.render(mapping, true)
	return out
}

func (t *Template) render(mapping map[string]string, safe bool) (string, error) {
	if t == nil {
		return "", nil
	}
	var out strings.Builder
	for i := 0; i < len(t.text); {
		if t.text[i] != '$' {
			out.WriteByte(t.text[i])
			i++
			continue
		}
		start := i
		i++
		if i >= len(t.text) {
			if safe {
				out.WriteString(t.text[start:])
				break
			}
			return "", fmt.Errorf("string.Template: trailing '$'")
		}
		if t.text[i] == '$' {
			out.WriteByte('$')
			i++
			continue
		}
		if t.text[i] == '{' {
			end := strings.IndexByte(t.text[i+1:], '}')
			if end < 0 {
				if safe {
					out.WriteString(t.text[start:])
					break
				}
				return "", fmt.Errorf("string.Template: missing '}'")
			}
			name := t.text[i+1 : i+1+end]
			if !validIdentifier(name) {
				if safe {
					out.WriteString(t.text[start : i+1+end+1])
					i = i + 1 + end + 1
					continue
				}
				return "", fmt.Errorf("string.Template: invalid placeholder %q", name)
			}
			if value, ok := mapping[name]; ok {
				out.WriteString(value)
			} else if safe {
				out.WriteString(t.text[start : i+1+end+1])
			} else {
				return "", fmt.Errorf("string.Template: missing key %q", name)
			}
			i = i + 1 + end + 1
			continue
		}
		if !isIdentifierStart(t.text[i]) {
			if safe {
				out.WriteByte('$')
				i = start + 1
				continue
			}
			return "", fmt.Errorf("string.Template: invalid '$'")
		}
		j := i + 1
		for j < len(t.text) && isIdentifierPart(t.text[j]) {
			j++
		}
		name := t.text[i:j]
		if value, ok := mapping[name]; ok {
			out.WriteString(value)
		} else if safe {
			out.WriteString(t.text[start:j])
		} else {
			return "", fmt.Errorf("string.Template: missing key %q", name)
		}
		i = j
	}
	return out.String(), nil
}

func validIdentifier(name string) bool {
	if name == "" || !isIdentifierStart(name[0]) {
		return false
	}
	for i := 1; i < len(name); i++ {
		if !isIdentifierPart(name[i]) {
			return false
		}
	}
	return true
}

func isIdentifierStart(b byte) bool {
	return b == '_' || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isIdentifierPart(b byte) bool {
	return isIdentifierStart(b) || (b >= '0' && b <= '9')
}
