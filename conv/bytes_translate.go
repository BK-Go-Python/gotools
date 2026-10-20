package conv

import "fmt"

// Bytes_maketrans 等价于 Python 的 bytes.maketrans(from, to)：
// 生成 256 字节的转换表，把 from 中第 i 个字节映射到 to 中第 i 个字节。
//
// from 与 to 长度必须相同，否则返回 identity 表和 error。
func Bytes_maketrans(from, to []byte) ([256]byte, error) {
	var table [256]byte
	for i := 0; i < 256; i++ {
		table[i] = byte(i)
	}
	if len(from) != len(to) {
		return table, fmt.Errorf("conv.Bytes_maketrans: from and to must have equal length: %d != %d", len(from), len(to))
	}
	for i, b := range from {
		table[b] = to[i]
	}
	return table, nil
}

// Bytes_translate 等价于 Python 的 bytes.translate(table, delete=b"")：
// 对每个字节查 256 字节表，delete 中的字节先被删除。
//
// delete 可传多个集合；返回新切片。
func Bytes_translate(data []byte, table [256]byte, delete ...[]byte) []byte {
	var deleted [256]bool
	for _, set := range delete {
		for _, b := range set {
			deleted[b] = true
		}
	}
	out := make([]byte, 0, len(data))
	for _, b := range data {
		if deleted[b] {
			continue
		}
		out = append(out, table[b])
	}
	return out
}
