package list

// Of 返回 items 的浅拷贝。
//
// Python: list(items)
// DIFF: Go 泛型要求元素同构；异构列表请使用 List/From。
func Of[T any](items ...T) []T {
	out := make([]T, len(items))
	copy(out, items)
	return out
}
