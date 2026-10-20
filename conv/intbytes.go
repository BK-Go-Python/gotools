package conv

import (
	"fmt"
	"math"
	"strings"
)

// maxIntToBytesLen 限制 IntToBytes 能产生的字节数，避免超大 length 直接
// 耗尽内存或触发不可恢复的运行时错误。Python 允许任意长度。
const maxIntToBytesLen = 1 << 20

// IntToBytes 等价于 Python 的 int.to_bytes(length, byteorder, signed=signed)，
// 把 x 按大端或小端写成 length 个字节。
//
// order 只接受 "big"/"little"（大小写不敏感）。signed 省略时默认 false：
//   - 无符号：x 必须 >= 0，且 x < 256^length；
//   - 有符号：x 必须落在 [-2^(8*length-1), 2^(8*length-1)-1]，按二进制补码写入。
//
// Python: x.to_bytes(length, order, signed=signed)
// DIFF: length 必须 > 0（Python 允许 length=0 得到 b""）；length 上限为
// maxIntToBytesLen，超出返回 error；signed 传入多个值时只取第一个。
func IntToBytes(x int64, length int, order string, signed ...bool) ([]byte, error) {
	isSigned := false
	if len(signed) > 0 {
		isSigned = signed[0]
	}

	if length <= 0 {
		return nil, fmt.Errorf("conv.IntToBytes: length must be positive, got %d", length)
	}
	if length > maxIntToBytesLen {
		return nil, fmt.Errorf("conv.IntToBytes: length %d exceeds limit %d", length, maxIntToBytesLen)
	}
	little, err := parseIntBytesOrder("conv.IntToBytes", order)
	if err != nil {
		return nil, err
	}

	padding := byte(0)
	if isSigned {
		if !fitsIntBytesSigned(x, length) {
			return nil, fmt.Errorf("conv.IntToBytes: value %d overflows %d-byte signed integer", x, length)
		}
		if x < 0 {
			padding = 0xff
		}
	} else {
		if x < 0 {
			return nil, fmt.Errorf("conv.IntToBytes: negative value %d cannot be converted to unsigned bytes", x)
		}
		if !fitsIntBytesUnsigned(uint64(x), length) {
			return nil, fmt.Errorf("conv.IntToBytes: value %d overflows %d-byte unsigned integer", x, length)
		}
	}

	out := make([]byte, length)
	for i := range out {
		out[i] = padding
	}

	// 至多低 8 个字节携带有效数字，其余是零或符号扩展。
	n := length
	if n > 8 {
		n = 8
	}
	u := uint64(x)
	if little {
		for i := 0; i < n; i++ {
			out[i] = byte(u >> (8 * i))
		}
	} else {
		for i := 0; i < n; i++ {
			out[length-1-i] = byte(u >> (8 * i))
		}
	}
	return out, nil
}

// IntFromBytes 等价于 Python 的 int.from_bytes(data, order, signed=signed)，
// 把 1~8 个字节按大端或小端解释为 int64。
//
// order 只接受 "big"/"little"（大小写不敏感）。signed 省略时默认 false：
//   - 无符号：结果超过 math.MaxInt64 返回 error，因为无法用 int64 表示；
//   - 有符号：按二进制补码解释，结果一定是 int64。
//
// Python: int.from_bytes(data, order, signed=signed)
// DIFF: data 长度必须为 1~8（Python 接受更长的 bytes 和空 bytes）；
// signed 传入多个值时只取第一个。
func IntFromBytes(data []byte, order string, signed ...bool) (int64, error) {
	isSigned := false
	if len(signed) > 0 {
		isSigned = signed[0]
	}

	if len(data) < 1 || len(data) > 8 {
		return 0, fmt.Errorf("conv.IntFromBytes: data must be 1 to 8 bytes, got %d", len(data))
	}
	little, err := parseIntBytesOrder("conv.IntFromBytes", order)
	if err != nil {
		return 0, err
	}

	var u uint64
	if little {
		for i := len(data) - 1; i >= 0; i-- {
			u = u<<8 | uint64(data[i])
		}
	} else {
		for _, b := range data {
			u = u<<8 | uint64(b)
		}
	}

	if isSigned {
		if bits := uint(8 * len(data)); bits < 64 && u&(uint64(1)<<(bits-1)) != 0 {
			u |= ^uint64(0) << bits
		}
		return int64(u), nil
	}
	if u > math.MaxInt64 {
		return 0, fmt.Errorf("conv.IntFromBytes: value %d overflows int64", u)
	}
	return int64(u), nil
}

// parseIntBytesOrder 解析大小写不敏感的字节序，返回是否为小端。
func parseIntBytesOrder(caller, order string) (little bool, err error) {
	switch {
	case strings.EqualFold(order, "big"):
		return false, nil
	case strings.EqualFold(order, "little"):
		return true, nil
	default:
		return false, fmt.Errorf("%s: invalid byte order %q (want \"big\" or \"little\")", caller, order)
	}
}

// fitsIntBytesSigned 判断 x 是否能放进 length 字节的二进制补码。
func fitsIntBytesSigned(x int64, length int) bool {
	if length >= 8 {
		return true
	}
	bits := uint(8 * length)
	min := -int64(1) << (bits - 1)
	max := int64(1)<<(bits-1) - 1
	return x >= min && x <= max
}

// fitsIntBytesUnsigned 判断非负的 x 是否能放进 length 字节的无符号表示。
func fitsIntBytesUnsigned(x uint64, length int) bool {
	if length >= 8 {
		return true
	}
	return x < uint64(1)<<(8*length)
}
