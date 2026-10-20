package set

// SymmetricDifference 等价于 Python 的 set ^ other。
func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	out := make(Set[T], len(s)+len(other))
	for item := range s {
		if !other.Has(item) {
			out[item] = struct{}{}
		}
	}
	for item := range other {
		if !s.Has(item) {
			out[item] = struct{}{}
		}
	}
	return out
}

// IsDisjoint 等价于 Python 的 set.isdisjoint(other)。
func (s Set[T]) IsDisjoint(other Set[T]) bool {
	small, large := s, other
	if len(small) > len(large) {
		small, large = large, small
	}
	for item := range small {
		if large.Has(item) {
			return false
		}
	}
	return true
}

// Update 等价于 Python 的 set.update(*others)，原地并入并返回集合。
func (s Set[T]) Update(others ...Set[T]) Set[T] {
	if s == nil {
		s = make(Set[T])
	}
	for _, other := range others {
		for item := range other {
			s[item] = struct{}{}
		}
	}
	return s
}

// IntersectionUpdate 等价于 Python 的 set.intersection_update(other)。
func (s Set[T]) IntersectionUpdate(other Set[T]) Set[T] {
	for item := range s {
		if !other.Has(item) {
			delete(s, item)
		}
	}
	return s
}

// DifferenceUpdate 等价于 Python 的 set.difference_update(other)。
func (s Set[T]) DifferenceUpdate(other Set[T]) Set[T] {
	for item := range other {
		delete(s, item)
	}
	return s
}

// SymmetricDifferenceUpdate 等价于 Python 的 set.symmetric_difference_update(other)。
func (s Set[T]) SymmetricDifferenceUpdate(other Set[T]) Set[T] {
	if s == nil {
		s = make(Set[T])
	}
	for item := range other {
		if s.Has(item) {
			delete(s, item)
		} else {
			s[item] = struct{}{}
		}
	}
	return s
}

// Pop 等价于 Python 的 set.pop()，移除并返回一个任意元素。
// 空集合返回零值和 false。
func (s Set[T]) Pop() (T, bool) {
	for item := range s {
		delete(s, item)
		return item, true
	}
	var zero T
	return zero, false
}

// Clear 等价于 Python 的 set.clear()。
func (s Set[T]) Clear() Set[T] {
	for item := range s {
		delete(s, item)
	}
	return s
}

// Copy 等价于 Python 的 set.copy()。
func (s Set[T]) Copy() Set[T] {
	return s.Clone()
}

// SymmetricDifference 等价于 Python 的 frozenset ^ other。
func (f Frozen[T]) SymmetricDifference(other Frozen[T]) Frozen[T] {
	out := make(map[T]struct{}, len(f.m)+len(other.m))
	for item := range f.m {
		if !other.Has(item) {
			out[item] = struct{}{}
		}
	}
	for item := range other.m {
		if !f.Has(item) {
			out[item] = struct{}{}
		}
	}
	return Frozen[T]{m: out}
}

// IsDisjoint 等价于 Python 的 frozenset.isdisjoint(other)。
func (f Frozen[T]) IsDisjoint(other Frozen[T]) bool {
	small, large := f, other
	if len(small.m) > len(large.m) {
		small, large = large, small
	}
	for item := range small.m {
		if large.Has(item) {
			return false
		}
	}
	return true
}

// Copy 等价于 Python 的 frozenset.copy()；Frozen 不可变，直接返回自身。
func (f Frozen[T]) Copy() Frozen[T] {
	return f
}
