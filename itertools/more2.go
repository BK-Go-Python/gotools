package itertools

// Tee 是 Python itertools.tee 的有限版本：返回 n 份独立的切片副本。
// n <= 0 返回空切片。
func Tee[T any](items []T, n int) [][]T {
	if n <= 0 {
		return [][]T{}
	}
	out := make([][]T, n)
	for i := range out {
		out[i] = append([]T(nil), items...)
	}
	return out
}

// Combinations_with_replacement 等价于 Python 的
// itertools.combinations_with_replacement(items, r)：允许重复取元素的组合。
//
// r<0 或（r>0 且 items 为空）返回空；r==0 返回一个空组合。
func Combinations_with_replacement[T any](items []T, r int) [][]T {
	n := len(items)
	if r < 0 || (r > 0 && n == 0) {
		return [][]T{}
	}
	if r == 0 {
		return [][]T{{}}
	}

	result := make([][]T, 0)
	indices := make([]int, r)
	for {
		row := make([]T, r)
		for i, index := range indices {
			row[i] = items[index]
		}
		result = append(result, row)

		i := r - 1
		for i >= 0 && indices[i] == n-1 {
			i--
		}
		if i < 0 {
			break
		}
		next := indices[i] + 1
		for j := i; j < r; j++ {
			indices[j] = next
		}
	}
	return result
}
