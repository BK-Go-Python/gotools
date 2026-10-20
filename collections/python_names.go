package collections

// 本文件提供 Python 拼写的别名。Go 要求导出标识符大写开头，因此采用
// “首字母大写、其余保持 Python 拼写”的规则（例如 most_common -> Most_common）。

// Most_common 是 MostCommon 的 Python 拼写别名。
func (c Counter[T]) Most_common(n ...int) []CountEntry[T] {
	return c.MostCommon(n...)
}

// Appendleft 是 AppendLeft 的 Python 拼写别名。
func (d *Deque[T]) Appendleft(items ...T) {
	d.AppendLeft(items...)
}

// Extendleft 是 ExtendLeft 的 Python 拼写别名。
func (d *Deque[T]) Extendleft(items ...T) {
	d.ExtendLeft(items...)
}

// Popleft 是 PopLeft 的 Python 拼写别名。
func (d *Deque[T]) Popleft() (T, bool) {
	return d.PopLeft()
}

// Peekleft 是 PeekLeft 的别名；Python deque 没有 peek，这只是本项目扩展。
func (d *Deque[T]) Peekleft() (T, bool) {
	return d.PeekLeft()
}

// Popitem 是 PopItem 的 Python 拼写别名。
func (d *OrderedDict[K, V]) Popitem(last ...bool) (OrderedItem[K, V], bool) {
	return d.PopItem(last...)
}

// Move_to_end 是 MoveToEnd 的 Python 拼写别名。
func (d *OrderedDict[K, V]) Move_to_end(key K) bool {
	return d.MoveToEnd(key)
}
