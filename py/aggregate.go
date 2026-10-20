package py

import (
	"cmp"
	"fmt"
	"slices"
)

// Sum 等价于 Python 的 sum(seq)，返回与元素相同的数值类型。
//
// 支持任意数值切片，包括命名切片类型（例如 type Ints []int）：
//
//	py.Sum([]int{1, 2, 3})       // 6
//	py.Sum([]float64{1.5, 2.5})  // 4
//	py.Sum(Ints{1, 2, 3})        // 6
//
// Python: sum(seq)
// DIFF: Go 整数溢出会回绕；Python int 是任意精度。异构 []any 请使用 SumAny。
func Sum[S ~[]E, E Number](seq S) E {
	var total E
	for _, value := range seq {
		total += value
	}
	return total
}

// SumOf 是 Sum 的旧名字，保留用于兼容。
//
// Deprecated: 使用 Sum。
func SumOf[S ~[]E, E Number](seq S) E {
	return Sum(seq)
}

// MinOf 等价于 Python 的 min(seq)，空序列返回 ErrEmpty。
//
// Python: min(seq)
// DIFF: 仅支持 cmp.Ordered 同构元素；混合类型请使用动态版 py.Min。
func MinOf[T cmp.Ordered](seq []T) (T, error) {
	var zero T
	if len(seq) == 0 {
		return zero, fmt.Errorf("py.MinOf: %w", ErrEmpty)
	}
	return slices.Min(seq), nil
}

// MaxOf 等价于 Python 的 max(seq)，空序列返回 ErrEmpty。
//
// Python: max(seq)
// DIFF: 仅支持 cmp.Ordered 同构元素；混合类型请使用动态版 py.Max。
func MaxOf[T cmp.Ordered](seq []T) (T, error) {
	var zero T
	if len(seq) == 0 {
		return zero, fmt.Errorf("py.MaxOf: %w", ErrEmpty)
	}
	return slices.Max(seq), nil
}

// SortedOf 等价于 Python 的 sorted(seq)，返回新的有序切片。
//
// Python: sorted(seq)
// DIFF: 仅支持 cmp.Ordered 同构元素；动态版 py.Sorted 支持 key/reverse。
func SortedOf[T cmp.Ordered](seq []T) []T {
	out := slices.Clone(seq)
	if out == nil {
		out = []T{}
	}
	slices.SortStableFunc(out, cmp.Compare)
	return out
}
