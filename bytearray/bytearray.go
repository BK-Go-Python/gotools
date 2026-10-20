// Package bytearray 提供 Python bytearray 风格的具名可变字节序列。
//
// ByteArray 底层是 []byte，方法在接收者上原地修改并返回 *ByteArray，
// 与 Python bytearray 的可变语义一致；需要只读副本时使用 Bytes 或 Copy。
package bytearray

import (
	"bytes"
	"fmt"

	"github.com/BK-Go-Python/gotools/conv"
	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// ByteArray 是 Python bytearray 的 Go 版本。
type ByteArray []byte

// New 创建一个包含 items 的 ByteArray。
//
// Python: bytearray(items)
func New(items ...byte) ByteArray {
	out := make(ByteArray, len(items))
	copy(out, items)
	return out
}

// From 等价于 Python 的 bytearray(v)：支持 string、[]byte、整数 n（n 个零字节）、
// 以及元素为 0..255 的切片/数组；具体规则复用 conv.Bytes。
func From(v any) (ByteArray, error) {
	data, err := conv.Bytes(v)
	if err != nil {
		return nil, err
	}
	return ByteArray(data), nil
}

// Bytes 返回底层字节的副本（对应 Python 的 bytes(b)）。
func (b ByteArray) Bytes() []byte {
	return append([]byte(nil), b...)
}

// Len 等价于 Python 的 len(b)。
func (b ByteArray) Len() int {
	return len(b)
}

// Append 等价于 Python 的 bytearray.append/extend，原地追加一个或多个字节。
func (b *ByteArray) Append(items ...byte) *ByteArray {
	*b = append(*b, items...)
	return b
}

// Extend 等价于 Python 的 bytearray.extend(iterable)，支持与 From 相同的输入。
func (b *ByteArray) Extend(v any) error {
	data, err := conv.Bytes(v)
	if err != nil {
		return err
	}
	*b = append(*b, data...)
	return nil
}

// Insert 等价于 Python 的 bytearray.insert(index, value)；负下标从末尾算，
// 越界会截断到 [0, len]。
func (b *ByteArray) Insert(index int, value byte) *ByteArray {
	n := len(*b)
	if index < 0 {
		index += n
	}
	if index < 0 {
		index = 0
	}
	if index > n {
		index = n
	}
	*b = append(*b, 0)
	copy((*b)[index+1:], (*b)[index:])
	(*b)[index] = value
	return b
}

// Pop 等价于 Python 的 bytearray.pop()：删除并返回最后一个字节。
func (b *ByteArray) Pop() (byte, bool) {
	if len(*b) == 0 {
		return 0, false
	}
	value := (*b)[len(*b)-1]
	*b = (*b)[:len(*b)-1]
	return value, true
}

// PopAt 等价于 Python 的 bytearray.pop(index)。
func (b *ByteArray) PopAt(index int) (byte, error) {
	n := len(*b)
	if n == 0 {
		return 0, fmt.Errorf("bytearray: pop from empty bytearray")
	}
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		return 0, fmt.Errorf("bytearray: pop index out of range: %d", index)
	}
	value := (*b)[index]
	*b = append((*b)[:index], (*b)[index+1:]...)
	return value, nil
}

// Remove 等价于 Python 的 bytearray.remove(value)，删除第一个匹配字节。
func (b *ByteArray) Remove(value byte) error {
	idx := bytes.IndexByte(*b, value)
	if idx < 0 {
		return fmt.Errorf("bytearray: byte %d not found", value)
	}
	*b = append((*b)[:idx], (*b)[idx+1:]...)
	return nil
}

// Count 等价于 Python 的 bytearray.count(value)。
func (b ByteArray) Count(value byte) int {
	return bytes.Count(b, []byte{value})
}

// Index 等价于 Python 的 bytearray.index(value)，返回第一个匹配下标。
func (b ByteArray) Index(value byte) (int, error) {
	idx := bytes.IndexByte(b, value)
	if idx < 0 {
		return -1, fmt.Errorf("bytearray: byte %d not found", value)
	}
	return idx, nil
}

// Get 等价于 Python 的 b[index]，支持负数下标。
func (b ByteArray) Get(index int) (byte, error) {
	n := len(b)
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		return 0, fmt.Errorf("bytearray: index out of range: %d", index)
	}
	return b[index], nil
}

// Set 等价于 Python 的 b[index] = value，支持负数下标。
func (b ByteArray) Set(index int, value byte) error {
	n := len(b)
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		return fmt.Errorf("bytearray: index out of range: %d", index)
	}
	b[index] = value
	return nil
}

// Slice 等价于 Python 的 b[start:end]，支持负数下标并自动截断。
func (b ByteArray) Slice(start, end int) ByteArray {
	n := len(b)
	if start < 0 {
		start += n
	}
	if end < 0 {
		end += n
	}
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}
	if end < 0 {
		end = 0
	}
	if end > n {
		end = n
	}
	if start >= end {
		return ByteArray{}
	}
	out := make(ByteArray, end-start)
	copy(out, b[start:end])
	return out
}

// Reverse 等价于 Python 的 bytearray.reverse()，原地反转。
func (b *ByteArray) Reverse() *ByteArray {
	for i, j := 0, len(*b)-1; i < j; i, j = i+1, j-1 {
		(*b)[i], (*b)[j] = (*b)[j], (*b)[i]
	}
	return b
}

// Clear 等价于 Python 的 bytearray.clear()。
func (b *ByteArray) Clear() *ByteArray {
	*b = (*b)[:0]
	return b
}

// Copy 等价于 Python 的 bytearray.copy()。
func (b ByteArray) Copy() ByteArray {
	return append(ByteArray(nil), b...)
}

// Replace 等价于 Python 的 bytearray.replace(old, new[, count])，原地替换。
func (b *ByteArray) Replace(old, new []byte, count ...int) *ByteArray {
	*b = ByteArray(conv.Bytes_replace([]byte(*b), old, new, count...))
	return b
}

// Upper 等价于 Python 的 bytearray.upper()，返回新的 ByteArray。
func (b ByteArray) Upper() ByteArray {
	return ByteArray(conv.Bytes_upper([]byte(b)))
}

// Lower 等价于 Python 的 bytearray.lower()，返回新的 ByteArray。
func (b ByteArray) Lower() ByteArray {
	return ByteArray(conv.Bytes_lower([]byte(b)))
}

// Decode 等价于 Python 的 bytearray.decode("utf-8")。
func (b ByteArray) Decode() (string, error) {
	return conv.Decode([]byte(b))
}

// Equal 比较两个 ByteArray 的字节内容。
func (b ByteArray) Equal(other ByteArray) bool {
	return bytes.Equal(b, other)
}

// Repr 返回 Python 风格的 bytearray(b'...') 表示。
func (b ByteArray) Repr() string {
	return fmt.Sprintf("bytearray(%s)", pyval.Repr([]byte(b)))
}

// FromHex 等价于 Python 的 bytearray.fromhex(s)，忽略输入中的空白。
func FromHex(s string) (ByteArray, error) {
	data, err := conv.FromHex(s)
	if err != nil {
		return nil, err
	}
	return ByteArray(data), nil
}

// Concat 等价于 Python 的 b1 + b2 + ...，返回新的 ByteArray。
func (b ByteArray) Concat(others ...ByteArray) ByteArray {
	size := len(b)
	for _, other := range others {
		size += len(other)
	}
	out := make(ByteArray, 0, size)
	out = append(out, b...)
	for _, other := range others {
		out = append(out, other...)
	}
	return out
}

// Repeat 等价于 Python 的 b * n；n <= 0 时返回空 ByteArray。
func (b ByteArray) Repeat(n int) ByteArray {
	if n <= 0 || len(b) == 0 {
		return ByteArray{}
	}
	out := make(ByteArray, 0, len(b)*n)
	for i := 0; i < n; i++ {
		out = append(out, b...)
	}
	return out
}
