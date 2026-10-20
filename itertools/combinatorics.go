package itertools

// Group 表示 Groupby 产生的一个相邻分组。
type Group[T any, K comparable] struct {
	Key    K
	Values []T
}

// Product 返回给定集合的笛卡尔积，结果按 sets 的顺序嵌套生成。
//
// 与 Python 的 itertools.product 一致：不传参数时返回只含一个空切片的结果
// （对应 product() 产生的空元组）；任一集合为空时返回非 nil 空切片。
// Product 不修改入参，返回的每个切片都独立分配。
func Product[T any](sets ...[]T) [][]T {
	result := [][]T{{}}
	for _, set := range sets {
		if len(set) == 0 {
			return [][]T{}
		}
		capacity := len(result) * len(set)
		if capacity/len(set) != len(result) {
			capacity = 0
		}
		next := make([][]T, 0, capacity)
		for _, prefix := range result {
			for _, value := range set {
				row := make([]T, len(prefix)+1)
				copy(row, prefix)
				row[len(prefix)] = value
				next = append(next, row)
			}
		}
		result = next
	}
	return result
}

// Combinations 返回从 items 中按索引顺序选取 r 个元素的所有组合。
//
// 结果按索引字典序升序排列；r < 0 或 r > len(items) 时返回非 nil 空切片；
// r == 0 时返回只含一个空切片的结果。重复元素按位置区分，
// 因此 [1, 1] 会产生重复的组合。
func Combinations[T any](items []T, r int) [][]T {
	n := len(items)
	if r < 0 || r > n {
		return [][]T{}
	}
	if r == 0 {
		return [][]T{{}}
	}

	result := make([][]T, 0)
	indices := make([]int, r)
	for i := range indices {
		indices[i] = i
	}
	for {
		row := make([]T, r)
		for i, index := range indices {
			row[i] = items[index]
		}
		result = append(result, row)

		i := r - 1
		for i >= 0 && indices[i] == i+n-r {
			i--
		}
		if i < 0 {
			break
		}
		indices[i]++
		for j := i + 1; j < r; j++ {
			indices[j] = indices[j-1] + 1
		}
	}
	return result
}

// Permutations 返回从 items 中按索引顺序选取 r 个元素的所有排列。
//
// 结果按索引字典序升序排列；r < 0 或 r > len(items) 时返回非 nil 空切片；
// r == 0 时返回只含一个空切片的结果。重复元素按位置区分，
// 因此会像 Python 一样产生重复排列。
func Permutations[T any](items []T, r int) [][]T {
	n := len(items)
	if r < 0 || r > n {
		return [][]T{}
	}
	if r == 0 {
		return [][]T{{}}
	}

	result := make([][]T, 0)
	indices := make([]int, r)
	used := make([]bool, n)
	for i := 0; i < r; i++ {
		indices[i] = i
		used[i] = true
	}

	for {
		row := make([]T, r)
		for i, index := range indices {
			row[i] = items[index]
		}
		result = append(result, row)

		pos := r - 1
		for ; pos >= 0; pos-- {
			used[indices[pos]] = false
			next := -1
			for candidate := indices[pos] + 1; candidate < n; candidate++ {
				if !used[candidate] {
					next = candidate
					break
				}
			}
			if next >= 0 {
				indices[pos] = next
				used[next] = true
				break
			}
		}
		if pos < 0 {
			break
		}
		fill := pos + 1
		for candidate := 0; candidate < n && fill < r; candidate++ {
			if !used[candidate] {
				indices[fill] = candidate
				used[candidate] = true
				fill++
			}
		}
	}
	return result
}

// Groupby 按相邻元素的 key 分组，行为与 Python 的 itertools.groupby 一致，
// 即只合并连续出现且 key 相同的元素，而不是全局分组。
//
// 结果保持输入顺序；key 对每个元素恰好调用一次；空输入返回非 nil 空切片。
// key 为 nil 时所有元素使用 K 的零值作为 key。Groupby 不修改入参。
func Groupby[T any, K comparable](items []T, key func(T) K) []Group[T, K] {
	groups := make([]Group[T, K], 0)
	if len(items) == 0 {
		return groups
	}
	if key == nil {
		var zero K
		key = func(T) K { return zero }
	}

	currentKey := key(items[0])
	current := Group[T, K]{
		Key:    currentKey,
		Values: []T{items[0]},
	}
	for _, item := range items[1:] {
		itemKey := key(item)
		if itemKey == currentKey {
			current.Values = append(current.Values, item)
			continue
		}
		groups = append(groups, current)
		currentKey = itemKey
		current = Group[T, K]{
			Key:    itemKey,
			Values: []T{item},
		}
	}
	return append(groups, current)
}
