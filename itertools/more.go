package itertools

// Pair2 是 Starmap 的输入二元组。
type Pair2[A, B any] struct {
	First  A
	Second B
}

// Filterfalse 等价于 Python 的 itertools.filterfalse(pred, items)：
// 保留 pred 返回 false 的元素。
func Filterfalse[T any](items []T, pred func(T) bool) []T {
	out := make([]T, 0, len(items))
	if pred == nil {
		return out
	}
	for _, item := range items {
		if !pred(item) {
			out = append(out, item)
		}
	}
	return out
}

// Compress 等价于 Python 的 itertools.compress(items, selectors)：
// 只保留 selectors[i] 为 true 的元素，长度取两者较短者。
func Compress[T any](items []T, selectors []bool) []T {
	n := len(items)
	if len(selectors) < n {
		n = len(selectors)
	}
	out := make([]T, 0, n)
	for i := 0; i < n; i++ {
		if selectors[i] {
			out = append(out, items[i])
		}
	}
	return out
}

// Starmap 等价于 Python 的 itertools.starmap(fn, pairs)：
// 把每个二元组展开成 fn(first, second)。
func Starmap[A, B, R any](pairs []Pair2[A, B], fn func(A, B) R) []R {
	out := make([]R, 0, len(pairs))
	if fn == nil {
		return out
	}
	for _, pair := range pairs {
		out = append(out, fn(pair.First, pair.Second))
	}
	return out
}

// Count 是 Python itertools.count 的有限版本：从 start 开始、步长 step，
// 生成 n 个整数；n <= 0 返回空切片。
func Count(start, step, n int) []int {
	if n <= 0 {
		return []int{}
	}
	out := make([]int, n)
	value := start
	for i := 0; i < n; i++ {
		out[i] = value
		value += step
	}
	return out
}
