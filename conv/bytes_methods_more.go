package conv

import (
	"bytes"
	"fmt"
)

// Bytes_strip_chars 等价于 Python 的 bytes.strip(chars)：去掉首尾出现在
// cutset 中的字节。cutset 为空时返回原内容的副本。
func Bytes_strip_chars(data, cutset []byte) []byte {
	start, end := trimBounds(data, cutset)
	return bytes.Clone(data[start:end])
}

// Bytes_lstrip 等价于 Python 的 bytes.lstrip()，去掉开头的 ASCII 空白。
func Bytes_lstrip(data []byte) []byte {
	start := 0
	for start < len(data) && isASCIIWhitespace(data[start]) {
		start++
	}
	return bytes.Clone(data[start:])
}

// Bytes_rstrip 等价于 Python 的 bytes.rstrip()，去掉末尾的 ASCII 空白。
func Bytes_rstrip(data []byte) []byte {
	end := len(data)
	for end > 0 && isASCIIWhitespace(data[end-1]) {
		end--
	}
	return bytes.Clone(data[:end])
}

// Bytes_lstrip_chars 等价于 Python 的 bytes.lstrip(chars)。
func Bytes_lstrip_chars(data, cutset []byte) []byte {
	start := 0
	for start < len(data) && inByteSet(cutset, data[start]) {
		start++
	}
	return bytes.Clone(data[start:])
}

// Bytes_rstrip_chars 等价于 Python 的 bytes.rstrip(chars)。
func Bytes_rstrip_chars(data, cutset []byte) []byte {
	end := len(data)
	for end > 0 && inByteSet(cutset, data[end-1]) {
		end--
	}
	return bytes.Clone(data[:end])
}

func trimBounds(data, cutset []byte) (start, end int) {
	start, end = 0, len(data)
	for start < end && inByteSet(cutset, data[start]) {
		start++
	}
	for end > start && inByteSet(cutset, data[end-1]) {
		end--
	}
	return start, end
}

func inByteSet(set []byte, b byte) bool {
	for _, c := range set {
		if c == b {
			return true
		}
	}
	return false
}

// Bytes_removeprefix 等价于 Python 3.9+ 的 bytes.removeprefix(prefix)。
func Bytes_removeprefix(data, prefix []byte) []byte {
	if bytes.HasPrefix(data, prefix) {
		return bytes.Clone(data[len(prefix):])
	}
	return bytes.Clone(data)
}

// Bytes_removesuffix 等价于 Python 3.9+ 的 bytes.removesuffix(suffix)。
func Bytes_removesuffix(data, suffix []byte) []byte {
	if bytes.HasSuffix(data, suffix) {
		return bytes.Clone(data[:len(data)-len(suffix)])
	}
	return bytes.Clone(data)
}

// Bytes_center 等价于 Python 的 bytes.center(width[, fill])。
// width 按字节计；fill 省略时为空格，多传只取第一个。
func Bytes_center(data []byte, width int, fill ...byte) []byte {
	if width <= len(data) {
		return bytes.Clone(data)
	}
	f := byte(' ')
	if len(fill) > 0 {
		f = fill[0]
	}
	padding := width - len(data)
	left := padding / 2
	if padding%2 == 1 && width%2 == 1 {
		left++
	}
	out := make([]byte, 0, width)
	out = append(out, bytes.Repeat([]byte{f}, left)...)
	out = append(out, data...)
	out = append(out, bytes.Repeat([]byte{f}, padding-left)...)
	return out
}

// Bytes_ljust 等价于 Python 的 bytes.ljust(width[, fill])。
func Bytes_ljust(data []byte, width int, fill ...byte) []byte {
	if width <= len(data) {
		return bytes.Clone(data)
	}
	f := byte(' ')
	if len(fill) > 0 {
		f = fill[0]
	}
	out := append(bytes.Clone(data), bytes.Repeat([]byte{f}, width-len(data))...)
	return out
}

// Bytes_rjust 等价于 Python 的 bytes.rjust(width[, fill])。
func Bytes_rjust(data []byte, width int, fill ...byte) []byte {
	if width <= len(data) {
		return bytes.Clone(data)
	}
	f := byte(' ')
	if len(fill) > 0 {
		f = fill[0]
	}
	out := make([]byte, 0, width)
	out = append(out, bytes.Repeat([]byte{f}, width-len(data))...)
	out = append(out, data...)
	return out
}

// Bytes_expandtabs 等价于 Python 的 bytes.expandtabs(tabsize=8)。
// 制表符按当前列补空格；遇到 \n 或 \r 列号归零；tabsize<=0 时删除制表符。
func Bytes_expandtabs(data []byte, tabsize ...int) []byte {
	size := 8
	if len(tabsize) > 0 {
		size = tabsize[0]
	}
	out := make([]byte, 0, len(data))
	column := 0
	for _, b := range data {
		switch b {
		case '\t':
			spaces := 0
			if size > 0 {
				spaces = size - column%size
			}
			out = append(out, bytes.Repeat([]byte{' '}, spaces)...)
			column += spaces
		case '\n', '\r':
			out = append(out, b)
			column = 0
		default:
			out = append(out, b)
			column++
		}
	}
	return out
}

// Bytes_split_n 等价于 Python 的 bytes.split(sep, maxsplit)。
// maxsplit<0 表示不限制；sep 为空返回 error。
func Bytes_split_n(data, sep []byte, maxsplit int) ([][]byte, error) {
	if len(sep) == 0 {
		return nil, fmt.Errorf("conv.Bytes_split_n: empty separator")
	}
	if maxsplit < 0 {
		return copyParts(bytes.Split(data, sep)), nil
	}
	return copyParts(bytes.SplitN(data, sep, maxsplit+1)), nil
}

// Bytes_rsplit 等价于 Python 的 bytes.rsplit(sep, maxsplit)。
// maxsplit<0 表示不限制；sep 为空返回 error。
func Bytes_rsplit(data, sep []byte, maxsplit int) ([][]byte, error) {
	if len(sep) == 0 {
		return nil, fmt.Errorf("conv.Bytes_rsplit: empty separator")
	}
	if maxsplit < 0 {
		return copyParts(bytes.Split(data, sep)), nil
	}
	if maxsplit == 0 {
		return [][]byte{bytes.Clone(data)}, nil
	}
	parts := make([][]byte, 0, maxsplit+1)
	end := len(data)
	for i := 0; i < maxsplit; i++ {
		idx := bytes.LastIndex(data[:end], sep)
		if idx < 0 {
			break
		}
		parts = append(parts, bytes.Clone(data[idx+len(sep):end]))
		end = idx
	}
	parts = append(parts, bytes.Clone(data[:end]))
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return parts, nil
}

func copyParts(parts [][]byte) [][]byte {
	out := make([][]byte, len(parts))
	for i, part := range parts {
		out[i] = bytes.Clone(part)
	}
	return out
}
