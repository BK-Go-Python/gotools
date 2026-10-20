package conv

import (
	"bytes"
	"fmt"
)

// Bytes_contains 等价于 Python 的 sub in data。
//
// Python: sub in data
func Bytes_contains(data, sub []byte) bool {
	return bytes.Contains(data, sub)
}

// Bytes_index 等价于 Python 的 data.find(sub)：返回 sub 第一次出现的起始
// 字节下标，未找到返回 -1。sub 为空时返回 0。
//
// Python: data.find(sub)
// DIFF: 未找到返回 -1（与 find 一致），不抛 ValueError。
func Bytes_index(data, sub []byte) int {
	return bytes.Index(data, sub)
}

// Bytes_count 等价于 Python 的 data.count(sub)：统计 sub 在 data 中不重叠
// 出现的次数。sub 为空时返回 len(data)+1，与 Python 的字节语义一致。
//
// Python: data.count(sub)
func Bytes_count(data, sub []byte) int {
	if len(sub) == 0 {
		return len(data) + 1
	}
	return bytes.Count(data, sub)
}

// Bytes_replace 等价于 Python 的 data.replace(old, new, count)。
//
// count 省略或为负数时替换全部；count 为 0 时不替换，直接返回 data 的
// 副本。old 为空时在 data 的字节边界（包括首尾）插入 new。
//
// Python: data.replace(old, new, count)
// DIFF: count 传入多个值时只取第一个。old 为空时按字节边界插入，这也是
// Python bytes 的语义；Go 的 bytes.Replace 对非 ASCII 数据按 UTF-8 字符
// 边界插入，二者不同。
func Bytes_replace(data, old, new []byte, count ...int) []byte {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	if n == 0 {
		return bytes.Clone(data)
	}
	if len(old) == 0 {
		return replaceEmptyOld(data, new, n)
	}
	return bytes.Replace(data, old, new, n)
}

// replaceEmptyOld 在 data 的每个字节边界插入 new，最多插入 count 次；
// count 为负数表示插入全部 len(data)+1 个位置。结果不共享 data。
func replaceEmptyOld(data, new []byte, count int) []byte {
	replacements := len(data) + 1
	if count >= 0 && count < replacements {
		replacements = count
	}

	out := make([]byte, 0, len(data)+replacements*len(new))
	for i := 0; i < replacements; i++ {
		out = append(out, new...)
		if i < len(data) {
			out = append(out, data[i])
		}
	}
	if replacements < len(data) {
		out = append(out, data[replacements:]...)
	}
	return out
}

// Bytes_split 等价于 Python 的 data.split(sep)，按 sep 切分 data。
//
// sep 为空时返回 error（对应 Python 的 ValueError: empty separator）。
// 返回的每个片段都是独立副本，不与 data 或其他片段共享底层数组。
//
// Python: data.split(sep)
// DIFF: ValueError 用 error 表示；不提供 maxsplit 参数。
func Bytes_split(data, sep []byte) ([][]byte, error) {
	if len(sep) == 0 {
		return nil, fmt.Errorf("conv.Bytes_split: empty separator")
	}

	parts := make([][]byte, 0, bytes.Count(data, sep)+1)
	for {
		i := bytes.Index(data, sep)
		if i < 0 {
			parts = append(parts, append([]byte{}, data...))
			return parts, nil
		}
		parts = append(parts, append([]byte{}, data[:i]...))
		data = data[i+len(sep):]
	}
}

// Bytes_join 等价于 Python 的 sep.join(parts)，用 sep 连接各片段。
//
// parts 为空时返回空切片。不修改 sep 或 parts，结果不与入参共享底层数组。
// 空输入返回非 nil 的空切片。
//
// Python: sep.join(parts)
func Bytes_join(sep []byte, parts [][]byte) []byte {
	if len(parts) == 0 {
		return []byte{}
	}

	total := len(sep) * (len(parts) - 1)
	for _, part := range parts {
		total += len(part)
	}

	out := make([]byte, 0, total)
	for i, part := range parts {
		if i > 0 {
			out = append(out, sep...)
		}
		out = append(out, part...)
	}
	return out
}

// Bytes_strip 等价于 Python 的 data.strip()：去掉首尾的 ASCII 空白
// " \t\n\r\v\f"。与 bytes.TrimSpace 不同，它不会把 Unicode 空白或无效
// UTF-8 字节当作空白处理。
//
// Python: data.strip()
// DIFF: 不提供 chars 参数；返回新切片。
func Bytes_strip(data []byte) []byte {
	start, end := 0, len(data)
	for start < end && isASCIIWhitespace(data[start]) {
		start++
	}
	for end > start && isASCIIWhitespace(data[end-1]) {
		end--
	}
	return bytes.Clone(data[start:end])
}

// isASCIIWhitespace 判断 b 是否属于 Python bytes.strip() 去掉的 ASCII 空白。
func isASCIIWhitespace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	default:
		return false
	}
}

// Bytes_upper 等价于 Python 的 data.upper()：只把 ASCII 小写字母 a-z 转成
// 大写，其他字节（包括 UTF-8 多字节序列）原样保留。
//
// Python: data.upper()
// DIFF: 返回新切片；不解析 UTF-8，也不做 Unicode 大小写映射。
func Bytes_upper(data []byte) []byte {
	out := bytes.Clone(data)
	for i, b := range out {
		if b >= 'a' && b <= 'z' {
			out[i] = b - ('a' - 'A')
		}
	}
	return out
}

// Bytes_lower 等价于 Python 的 data.lower()：只把 ASCII 大写字母 A-Z 转成
// 小写，其他字节（包括 UTF-8 多字节序列）原样保留。
//
// Python: data.lower()
// DIFF: 返回新切片；不解析 UTF-8，也不做 Unicode 大小写映射。
func Bytes_lower(data []byte) []byte {
	out := bytes.Clone(data)
	for i, b := range out {
		if b >= 'A' && b <= 'Z' {
			out[i] = b + ('a' - 'A')
		}
	}
	return out
}
