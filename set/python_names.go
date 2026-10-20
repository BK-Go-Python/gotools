package set

// 本文件提供 Python 拼写别名。

// Issubset 是 IsSubset 的 Python 拼写别名。
func (s Set[T]) Issubset(other Set[T]) bool { return s.IsSubset(other) }

// Isdisjoint 是 IsDisjoint 的 Python 拼写别名。
func (s Set[T]) Isdisjoint(other Set[T]) bool { return s.IsDisjoint(other) }

// Symmetric_difference 是 SymmetricDifference 的 Python 拼写别名。
func (s Set[T]) Symmetric_difference(other Set[T]) Set[T] {
	return s.SymmetricDifference(other)
}

// Intersection_update 是 IntersectionUpdate 的 Python 拼写别名。
func (s Set[T]) Intersection_update(other Set[T]) Set[T] {
	return s.IntersectionUpdate(other)
}

// Difference_update 是 DifferenceUpdate 的 Python 拼写别名。
func (s Set[T]) Difference_update(other Set[T]) Set[T] {
	return s.DifferenceUpdate(other)
}

// Symmetric_difference_update 是 SymmetricDifferenceUpdate 的 Python 拼写别名。
func (s Set[T]) Symmetric_difference_update(other Set[T]) Set[T] {
	return s.SymmetricDifferenceUpdate(other)
}

// Issubset 是 Frozen.IsSubset 的 Python 拼写别名。
func (f Frozen[T]) Issubset(other Frozen[T]) bool { return f.IsSubset(other) }

// Isdisjoint 是 Frozen.IsDisjoint 的 Python 拼写别名。
func (f Frozen[T]) Isdisjoint(other Frozen[T]) bool { return f.IsDisjoint(other) }

// Symmetric_difference 是 Frozen.SymmetricDifference 的 Python 拼写别名。
func (f Frozen[T]) Symmetric_difference(other Frozen[T]) Frozen[T] {
	return f.SymmetricDifference(other)
}
