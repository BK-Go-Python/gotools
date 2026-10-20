package set

import "fmt"

// Discard 删除元素，不存在的元素会被忽略；nil Set 安全。返回集合便于链式调用。
//
// 等价于已有 Set.Remove 的 Python discard 语义，命名与 Python 对齐。
//
// Python: set.discard(item)
func (s Set[T]) Discard(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// RemoveOne 删除单个元素；元素存在时删除并返回 nil，不存在时返回错误。
//
// DIFF: 现有 Set.Remove(items ...T) 是 Python set.discard 语义，删除不存在的元素
// 不报错并返回集合；RemoveOne 对应 Python set.remove，元素缺失时返回错误
// （Python 抛 KeyError）。
func (s Set[T]) RemoveOne(item T) error {
	if !s.Has(item) {
		return fmt.Errorf("set: remove %v: element not found", item)
	}
	delete(s, item)
	return nil
}

// Issuperset 等价于 Python 的 set.issuperset(other)：
// other 的所有元素都在 s 中，空集是任何集合的超集。
func (s Set[T]) Issuperset(other Set[T]) bool {
	if len(s) < len(other) {
		return false
	}
	for item := range other {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// Issuperset 等价于 Python 的 frozenset.issuperset(other)。
func (f Frozen[T]) Issuperset(other Frozen[T]) bool {
	if len(f.m) < len(other.m) {
		return false
	}
	for item := range other.m {
		if !f.Has(item) {
			return false
		}
	}
	return true
}
