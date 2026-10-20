package conv

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Bytes_rindex 等价于 Python 的 bytes.rindex(sub)：从右往左查找，
// 未找到返回 error。
func Bytes_rindex(data, sub []byte) (int, error) {
	idx := Bytes_rfind(data, sub)
	if idx < 0 {
		return -1, fmt.Errorf("conv.Bytes_rindex: subsection not found")
	}
	return idx, nil
}

// Bytes_splitlines_keepends 等价于 Python 的 bytes.splitlines(keepends=True)：
// 按 \n、\r、\r\n 切分并保留行尾换行符。
func Bytes_splitlines_keepends(data []byte) [][]byte {
	if len(data) == 0 {
		return [][]byte{}
	}
	out := make([][]byte, 0, 4)
	for i := 0; i < len(data); {
		start := i
		for i < len(data) && data[i] != '\n' && data[i] != '\r' {
			i++
		}
		end := i
		if i < len(data) {
			if data[i] == '\r' && i+1 < len(data) && data[i+1] == '\n' {
				i += 2
			} else {
				i++
			}
			end = i
		}
		out = append(out, append([]byte(nil), data[start:end]...))
	}
	return out
}

// Bytes_hex_sep 等价于 Python 的 bytes.hex(sep, bytes_per_sep=1)。
//
// sep 为空时等价于 Bytes_hex；bytes_per_sep 为负数时按其绝对值分组
// （与 CPython 3.13 的实测行为一致）；为 0 返回 error。
func Bytes_hex_sep(data []byte, sep string, bytesPerSep ...int) (string, error) {
	n := 1
	if len(bytesPerSep) > 0 {
		n = bytesPerSep[0]
	}
	if n == 0 {
		return "", fmt.Errorf("conv.Bytes_hex_sep: bytes_per_sep must not be zero")
	}
	if n < 0 {
		n = -n
	}
	if sep == "" {
		return hex.EncodeToString(data), nil
	}

	var b strings.Builder
	b.Grow(len(data)*2 + (len(data)/n)*len(sep))
	for i, by := range data {
		if i > 0 && i%n == 0 {
			b.WriteString(sep)
		}
		const digits = "0123456789abcdef"
		b.WriteByte(digits[by>>4])
		b.WriteByte(digits[by&0x0f])
	}
	return b.String(), nil
}
