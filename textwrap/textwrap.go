// Package textwrap 提供 Python textwrap 模块的常用文本处理函数。
package textwrap

import "strings"

// Wrap 等价于 Python 的 textwrap.wrap(text, width)：按空白切词并贪心换行。
//
// width <= 0 返回空切片；超长单词会按 width 强制断开。
func Wrap(text string, width int) []string {
	if width <= 0 {
		return []string{}
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}
	lines := make([]string, 0)
	current := ""
	for _, word := range words {
		runes := []rune(word)
		for len(runes) > width {
			if current != "" {
				lines = append(lines, current)
				current = ""
			}
			lines = append(lines, string(runes[:width]))
			runes = runes[width:]
		}
		word = string(runes)
		if word == "" {
			continue
		}
		if current == "" {
			current = word
			continue
		}
		if len([]rune(current))+1+len([]rune(word)) <= width {
			current += " " + word
			continue
		}
		lines = append(lines, current)
		current = word
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// Fill 等价于 Python 的 textwrap.fill(text, width)。
func Fill(text string, width int) string {
	return strings.Join(Wrap(text, width), "\n")
}

// Dedent 等价于 Python 的 textwrap.dedent(text)：去掉非空行的公共前导空白。
func Dedent(text string) string {
	lines := strings.Split(text, "\n")
	common := ""
	first := true
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := leadingWhitespace(line)
		if first {
			common = indent
			first = false
		} else {
			common = commonPrefix(common, indent)
		}
		if common == "" {
			return text
		}
	}
	if common == "" {
		return text
	}
	for i, line := range lines {
		if strings.HasPrefix(line, common) {
			lines[i] = line[len(common):]
		} else if strings.TrimSpace(line) == "" {
			lines[i] = ""
		}
	}
	return strings.Join(lines, "\n")
}

// Indent 等价于 Python 的 textwrap.indent(text, prefix)：给非空白行加前缀。
func Indent(text, prefix string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

// Shorten 等价于 Python 的 textwrap.shorten(text, width, placeholder="...")。
//
// placeholder 省略时使用 "..."；按 rune 计算宽度。
func Shorten(text string, width int, placeholder ...string) string {
	tail := "..."
	if len(placeholder) > 0 {
		tail = placeholder[0]
	}
	runes := []rune(text)
	if len(runes) <= width {
		return text
	}
	available := width - len([]rune(tail))
	if available < 0 {
		available = 0
	}
	trimmed := strings.TrimRight(string(runes[:available]), " \t\n")
	return trimmed + tail
}

func leadingWhitespace(line string) string {
	for i, r := range line {
		if r != ' ' && r != '\t' {
			return line[:i]
		}
	}
	return line
}

func commonPrefix(a, b string) string {
	if len(a) > len(b) {
		a, b = b, a
	}
	i := 0
	for i < len(a) && a[i] == b[i] {
		i++
	}
	return a[:i]
}
