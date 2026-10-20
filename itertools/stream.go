package itertools

import "fmt"

// Chain 依次拼接多个切片，等价于 Python 的 itertools.chain(*seqs)。
//
// nil 或空切片视为空切片；无参数时返回非 nil 空切片。
// 返回新切片，不修改任何入参。
//
// Python: list(chain(*seqs))
func Chain[T any](seqs ...[]T) []T {
	total := 0
	for _, seq := range seqs {
		total += len(seq)
	}
	out := make([]T, 0, total)
	for _, seq := range seqs {
		out = append(out, seq...)
	}
	return out
}

// Repeat 将 value 重复 n 次，等价于 Python 的 itertools.repeat(value, n)。
//
// n <= 0 时返回非 nil 空切片。
//
// Python: list(repeat(value, n))
func Repeat[T any](value T, n int) []T {
	if n <= 0 {
		return make([]T, 0)
	}
	out := make([]T, n)
	for i := range out {
		out[i] = value
	}
	return out
}

// Cycle 把 items 循环取 n 个元素，
// 等价于 Python 的 list(islice(itertools.cycle(items), n))。
//
// n <= 0 或 items 为空时返回非 nil 空切片；不修改入参。
//
// Python: list(islice(cycle(items), n))
func Cycle[T any](items []T, n int) []T {
	if n <= 0 || len(items) == 0 {
		return make([]T, 0)
	}
	out := make([]T, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, items[i%len(items)])
	}
	return out
}

// Islice 等价于 Python 的 itertools.islice(items, start, stop, step)。
//
// start、stop 必须非负，step 必须为正，否则返回非 nil 空切片和 error。
// stop < start 或 start 越过末尾时返回空切片；stop 超过长度时截断。
//
// Python: list(islice(items, start, stop, step))
func Islice[T any](items []T, start, stop, step int) ([]T, error) {
	if start < 0 {
		return make([]T, 0), fmt.Errorf("itertools.Islice: start must be non-negative, got %d", start)
	}
	if stop < 0 {
		return make([]T, 0), fmt.Errorf("itertools.Islice: stop must be non-negative, got %d", stop)
	}
	if step <= 0 {
		return make([]T, 0), fmt.Errorf("itertools.Islice: step must be positive, got %d", step)
	}

	limit := stop
	if limit > len(items) {
		limit = len(items)
	}
	if start >= limit {
		return make([]T, 0), nil
	}

	count := 1 + (limit-1-start)/step
	out := make([]T, 0, count)
	for i := start; i < limit; {
		out = append(out, items[i])
		if step > limit-1-i {
			break
		}
		i += step
	}
	return out, nil
}

// Takewhile 返回满足 pred 的最长前缀，等价于 Python 的 itertools.takewhile。
//
// pred 为 nil 时视为没有元素满足条件，返回非 nil 空切片。
//
// Python: list(takewhile(pred, items))
//
// DIFF: pred 为 nil 时返回空切片，Python 会抛 TypeError。
func Takewhile[T any](items []T, pred func(T) bool) []T {
	out := make([]T, 0, len(items))
	if pred == nil {
		return out
	}
	for _, item := range items {
		if !pred(item) {
			break
		}
		out = append(out, item)
	}
	return out
}

// Dropwhile 跳过满足 pred 的最长前缀并返回剩余元素，
// 等价于 Python 的 itertools.dropwhile。
//
// pred 为 nil 时视为没有元素满足条件，返回全部元素的副本。
//
// Python: list(dropwhile(pred, items))
//
// DIFF: pred 为 nil 时返回全部元素，Python 会抛 TypeError。
func Dropwhile[T any](items []T, pred func(T) bool) []T {
	i := 0
	if pred != nil {
		for i < len(items) && pred(items[i]) {
			i++
		}
	}
	out := make([]T, len(items)-i)
	copy(out, items[i:])
	return out
}

// Accumulate 依次用 fn 合并元素并返回每步的累积结果，
// 等价于 Python 的 itertools.accumulate(items, fn)。
//
// 空输入返回非 nil 空切片；首元素原样保留，之后为 fn(acc, x)。
// fn 为 nil 时不做合并，返回 items 的副本。
//
// Python: list(accumulate(items, fn))
//
// DIFF: 泛型 T 无法提供默认加法，fn 为 nil 时返回原值副本，
// Python 此时默认使用 operator.add。
func Accumulate[T any](items []T, fn func(acc, x T) T) []T {
	if len(items) == 0 {
		return make([]T, 0)
	}
	out := make([]T, len(items))
	out[0] = items[0]
	if fn == nil {
		copy(out[1:], items[1:])
		return out
	}
	acc := items[0]
	for i := 1; i < len(items); i++ {
		acc = fn(acc, items[i])
		out[i] = acc
	}
	return out
}

// Pairwise 返回相邻元素对，等价于 Python 的 itertools.pairwise。
//
// 长度小于 2 时返回非 nil 空切片；不修改入参。
//
// Python: list(pairwise(items))
func Pairwise[T any](items []T) [][2]T {
	if len(items) < 2 {
		return make([][2]T, 0)
	}
	out := make([][2]T, 0, len(items)-1)
	for i := 1; i < len(items); i++ {
		out = append(out, [2]T{items[i-1], items[i]})
	}
	return out
}

// Zip_longest 按最长序列拉链合并，缺失位置填 fill，
// 等价于 Python 的 itertools.zip_longest(*seqs, fillvalue=fill)。
//
// 无参数或所有序列为空时返回非 nil 空切片；
// 否则每行长度等于序列个数，外层长度等于最长序列长度。
//
// Python: list(zip_longest(*seqs, fillvalue=fill))
func Zip_longest[T any](fill T, seqs ...[]T) [][]T {
	if len(seqs) == 0 {
		return make([][]T, 0)
	}
	longest := 0
	for _, seq := range seqs {
		if len(seq) > longest {
			longest = len(seq)
		}
	}
	if longest == 0 {
		return make([][]T, 0)
	}

	out := make([][]T, longest)
	for i := 0; i < longest; i++ {
		row := make([]T, len(seqs))
		for j, seq := range seqs {
			if i < len(seq) {
				row[j] = seq[i]
			} else {
				row[j] = fill
			}
		}
		out[i] = row
	}
	return out
}
