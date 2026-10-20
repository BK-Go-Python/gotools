package conv

import "bytes"

// Bytes_startswith 等价于 Python 的 bytes.startswith(prefix)。
func Bytes_startswith(data, prefix []byte) bool {
	return bytes.HasPrefix(data, prefix)
}

// Bytes_endswith 等价于 Python 的 bytes.endswith(suffix)。
func Bytes_endswith(data, suffix []byte) bool {
	return bytes.HasSuffix(data, suffix)
}

// Bytes_find 等价于 Python 的 bytes.find(sub)，未找到返回 -1。
func Bytes_find(data, sub []byte) int {
	return bytes.Index(data, sub)
}

// Bytes_rfind 等价于 Python 的 bytes.rfind(sub)，未找到返回 -1。
func Bytes_rfind(data, sub []byte) int {
	return bytes.LastIndex(data, sub)
}

// Bytes_partition 等价于 Python 的 bytes.partition(sep)。
//
// 未找到时返回 (data, nil, nil)；sep 为空时同样返回 (data, nil, nil)
// （Python 会抛 ValueError，写 DIFF）。
func Bytes_partition(data, sep []byte) (before, middle, after []byte) {
	if len(sep) == 0 {
		return append([]byte(nil), data...), nil, nil
	}
	idx := bytes.Index(data, sep)
	if idx < 0 {
		return append([]byte(nil), data...), nil, nil
	}
	return append([]byte(nil), data[:idx]...),
		append([]byte(nil), sep...),
		append([]byte(nil), data[idx+len(sep):]...)
}

// Bytes_rpartition 等价于 Python 的 bytes.rpartition(sep)。
//
// 未找到时返回 (nil, nil, data)；sep 为空时同样返回 (nil, nil, data)。
func Bytes_rpartition(data, sep []byte) (before, middle, after []byte) {
	if len(sep) == 0 {
		return nil, nil, append([]byte(nil), data...)
	}
	idx := bytes.LastIndex(data, sep)
	if idx < 0 {
		return nil, nil, append([]byte(nil), data...)
	}
	return append([]byte(nil), data[:idx]...),
		append([]byte(nil), sep...),
		append([]byte(nil), data[idx+len(sep):]...)
}

// Bytes_splitlines 等价于 Python 的 bytes.splitlines()：按 \n、\r、\r\n
// 切分，不保留行尾；空输入返回空切片。
func Bytes_splitlines(data []byte) [][]byte {
	if len(data) == 0 {
		return [][]byte{}
	}
	out := make([][]byte, 0, bytes.Count(data, []byte("\n"))+1)
	for i := 0; i < len(data); {
		start := i
		for i < len(data) && data[i] != '\n' && data[i] != '\r' {
			i++
		}
		out = append(out, append([]byte(nil), data[start:i]...))
		if i < len(data) {
			if data[i] == '\r' && i+1 < len(data) && data[i+1] == '\n' {
				i += 2
			} else {
				i++
			}
		}
	}
	return out
}

// Bytes_capitalize 等价于 Python 的 bytes.capitalize()：首字节 ASCII 小写
// 字母转大写，其余 ASCII 字母转小写；0x80 以上字节保持不变。
func Bytes_capitalize(data []byte) []byte {
	out := append([]byte(nil), data...)
	for i, b := range out {
		if i == 0 {
			out[i] = asciiUpper(b)
		} else {
			out[i] = asciiLower(b)
		}
	}
	return out
}

// Bytes_title 等价于 Python 的 bytes.title()：每个单词的首个 ASCII 字母
// 大写，其余 ASCII 字母小写；非 ASCII 字节不参与大小写转换。
func Bytes_title(data []byte) []byte {
	out := append([]byte(nil), data...)
	prevLetter := false
	for i, b := range out {
		if isASCIILetter(b) {
			if prevLetter {
				out[i] = asciiLower(b)
			} else {
				out[i] = asciiUpper(b)
			}
			prevLetter = true
			continue
		}
		prevLetter = false
	}
	return out
}

// Bytes_swapcase 等价于 Python 的 bytes.swapcase()：仅交换 ASCII 字母大小写。
func Bytes_swapcase(data []byte) []byte {
	out := append([]byte(nil), data...)
	for i, b := range out {
		switch {
		case b >= 'a' && b <= 'z':
			out[i] = b - 'a' + 'A'
		case b >= 'A' && b <= 'Z':
			out[i] = b - 'A' + 'a'
		}
	}
	return out
}

// Bytes_zfill 等价于 Python 的 bytes.zfill(width)：左侧补 '0'，并保留
// 开头的 '+'/'-' 符号位置。
func Bytes_zfill(data []byte, width int) []byte {
	if width <= len(data) {
		return append([]byte(nil), data...)
	}
	out := make([]byte, 0, width)
	zeros := bytes.Repeat([]byte{'0'}, width-len(data))
	if len(data) > 0 && (data[0] == '+' || data[0] == '-') {
		out = append(out, data[0])
		out = append(out, zeros...)
		out = append(out, data[1:]...)
		return out
	}
	out = append(out, zeros...)
	out = append(out, data...)
	return out
}

func isASCIILetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func asciiUpper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	}
	return b
}

func asciiLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b - 'A' + 'a'
	}
	return b
}
