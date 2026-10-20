package py

import "math/bits"

// BitLength 等价于 Python 的 int.bit_length()，返回绝对值所需的二进制位数。
//
// Python: n.bit_length()
func BitLength[T Integer](x T) int {
	switch v := any(x).(type) {
	case int:
		return signedBitLength(int64(v))
	case int8:
		return signedBitLength(int64(v))
	case int16:
		return signedBitLength(int64(v))
	case int32:
		return signedBitLength(int64(v))
	case int64:
		return signedBitLength(v)
	case uint:
		return bits.Len64(uint64(v))
	case uint8:
		return bits.Len64(uint64(v))
	case uint16:
		return bits.Len64(uint64(v))
	case uint32:
		return bits.Len64(uint64(v))
	case uint64:
		return bits.Len64(v)
	case uintptr:
		return bits.Len64(uint64(v))
	}
	return 0
}

func signedBitLength(n int64) int {
	if n >= 0 {
		return bits.Len64(uint64(n))
	}
	// 对最小整数取负会溢出，先转成无符号补码再加一得到绝对值。
	return bits.Len64(uint64(-(n + 1)) + 1)
}
