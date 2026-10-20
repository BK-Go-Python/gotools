package py

import "math/bits"

// BitCount 等价于 Python 的 int.bit_count()，返回 x 的绝对值在二进制中 1 的个数。
//
// 先转换为无符号数再对负数取负，因此 math.MinInt64 等最小值也不会溢出。
//
// Python: int.bit_count()
// DIFF: Go 整数是定宽的；无符号整数按自身位模式统计。
func BitCount[T Integer](x T) int {
	u := uint64(x)
	if x < 0 {
		u = -u
	}
	return bits.OnesCount64(u)
}
