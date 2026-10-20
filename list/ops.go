package list

// Concat 等价于 Python 的 l1 + l2 + ...，把接收者和 others 依次拼接后返回新列表。
//
// 不修改任何入参；nil 接收者和 nil 参数都按空列表处理。
//
// Python: l1 + l2
func (l List) Concat(others ...List) List {
	size := len(l)
	for _, other := range others {
		size += len(other)
	}
	out := make(List, 0, size)
	out = append(out, l...)
	for _, other := range others {
		out = append(out, other...)
	}
	return out
}

// Repeat 等价于 Python 的 l * n，返回把 l 重复 n 次后的新列表。
//
// n <= 0 时返回空列表；空列表重复任意次仍为空列表；不修改入参。
//
// Python: l * n
func (l List) Repeat(n int) List {
	if n <= 0 || len(l) == 0 {
		return List{}
	}
	out := make(List, 0, len(l)*n)
	for i := 0; i < n; i++ {
		out = append(out, l...)
	}
	return out
}
