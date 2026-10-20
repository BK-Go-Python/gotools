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

// Min 等价于 Python 的 min(seq)，返回最小的元素；空序列返回 ErrEmpty。
//
// 支持任意 cmp.Ordered 的切片，包括命名切片类型：
//
//	py.Min([]int{3, 1, 2}) // (1, nil)
//
// Python: min(seq)
// DIFF: 仅支持同构有序元素；混合类型或需要 key 时使用 MinAny。
func Min[S ~[]E, E cmp.Ordered](seq S) (E, error) {
	var zero E
	if len(seq) == 0 {
		return zero, fmt.Errorf("py.Min: %w", ErrEmpty)
	}
	return slices.Min(seq), nil
}

// Max 等价于 Python 的 max(seq)，返回最大的元素；空序列返回 ErrEmpty。
//
// Python: max(seq)
// DIFF: 仅支持同构有序元素；混合类型或需要 key 时使用 MaxAny。
func Max[S ~[]E, E cmp.Ordered](seq S) (E, error) {
	var zero E
	if len(seq) == 0 {
		return zero, fmt.Errorf("py.Max: %w", ErrEmpty)
	}
	return slices.Max(seq), nil
}

// Sorted 等价于 Python 的 sorted(seq)，返回新的有序切片。
//
// Python: sorted(seq)
// DIFF: 仅支持同构有序元素；混合类型、key/reverse 请使用 SortedBy 或 SortedAny。
func Sorted[S ~[]E, E cmp.Ordered](seq S) []E {
	out := slices.Clone(seq)
	if out == nil {
		out = make(S, 0)
	}
	slices.SortStableFunc(out, cmp.Compare)
	return []E(out)
}

// MinOf 是 Min 的旧名字，保留用于兼容。
//
// Deprecated: 使用 Min。
func MinOf[T cmp.Ordered](seq []T) (T, error) {
	return Min(seq)
}

// MaxOf 是 Max 的旧名字，保留用于兼容。
//
// Deprecated: 使用 Max。
func MaxOf[T cmp.Ordered](seq []T) (T, error) {
	return Max(seq)
}

// SortedOf 是 Sorted 的旧名字，保留用于兼容。
//
// Deprecated: 使用 Sorted。
func SortedOf[T cmp.Ordered](seq []T) []T {
	return Sorted(seq)
}
