package set

// Set 是可变集合，等价于 Python 的 set。
//
// 零值可用，但 Add 的返回值必须重新赋值：
//
//	var s set.Set[int]
//	s = s.Add(1)
//
// 更推荐 set.New。
type Set[T comparable] map[T]struct{}

// New 创建一个包含 items 的集合。
//
// Python: set(items)
func New[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Add 添加元素并返回集合，便于链式调用。
func (s Set[T]) Add(items ...T) Set[T] {
	if s == nil {
		s = make(Set[T], len(items))
	}
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Remove 删除元素；不存在的元素会被忽略。返回集合便于链式调用。
//
// Python: set.discard(items)
func (s Set[T]) Remove(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has 等价于 Python 的 item in set。
func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

// Len 等价于 Python 的 len(set)。
func (s Set[T]) Len() int {
	return len(s)
}

// Items 返回元素切片。由于 Go map 无序，顺序不保证；需要排序请自行 slices.Sort。
func (s Set[T]) Items() []T {
	out := make([]T, 0, len(s))
	for item := range s {
		out = append(out, item)
	}
	return out
}

// Union 等价于 Python 的 set | other。
func (s Set[T]) Union(other Set[T]) Set[T] {
	out := make(Set[T], len(s)+len(other))
	for item := range s {
		out[item] = struct{}{}
	}
	for item := range other {
		out[item] = struct{}{}
	}
	return out
}

// Intersection 等价于 Python 的 set & other。
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	out := make(Set[T])
	if len(s) > len(other) {
		s, other = other, s
	}
	for item := range s {
		if other.Has(item) {
			out[item] = struct{}{}
		}
	}
	return out
}

// Difference 等价于 Python 的 set - other。
func (s Set[T]) Difference(other Set[T]) Set[T] {
	out := make(Set[T], len(s))
	for item := range s {
		if !other.Has(item) {
			out[item] = struct{}{}
		}
	}
	return out
}

// IsSubset 等价于 Python 的 set <= other。
func (s Set[T]) IsSubset(other Set[T]) bool {
	if len(s) > len(other) {
		return false
	}
	for item := range s {
		if !other.Has(item) {
			return false
		}
	}
	return true
}

// Equal 等价于 Python 的 set == other。
func (s Set[T]) Equal(other Set[T]) bool {
	return len(s) == len(other) && s.IsSubset(other)
}

// Clone 返回浅拷贝。
func (s Set[T]) Clone() Set[T] {
	out := make(Set[T], len(s))
	for item := range s {
		out[item] = struct{}{}
	}
	return out
}

// Frozen 是不可变集合，等价于 Python 的 frozenset。
type Frozen[T comparable] struct {
	m map[T]struct{}
}

// NewFrozen 创建一个不可变集合。
//
// Python: frozenset(items)
func NewFrozen[T comparable](items ...T) Frozen[T] {
	m := make(map[T]struct{}, len(items))
	for _, item := range items {
		m[item] = struct{}{}
	}
	return Frozen[T]{m: m}
}

// Has 等价于 Python 的 item in frozenset。
func (f Frozen[T]) Has(item T) bool {
	_, ok := f.m[item]
	return ok
}

// Len 等价于 Python 的 len(frozenset)。
func (f Frozen[T]) Len() int {
	return len(f.m)
}

// Items 返回元素切片，顺序不保证。
func (f Frozen[T]) Items() []T {
	out := make([]T, 0, len(f.m))
	for item := range f.m {
		out = append(out, item)
	}
	return out
}

// Union 返回两个不可变集合的并集。
func (f Frozen[T]) Union(other Frozen[T]) Frozen[T] {
	out := make(map[T]struct{}, len(f.m)+len(other.m))
	for item := range f.m {
		out[item] = struct{}{}
	}
	for item := range other.m {
		out[item] = struct{}{}
	}
	return Frozen[T]{m: out}
}

// Intersection 返回两个不可变集合的交集。
func (f Frozen[T]) Intersection(other Frozen[T]) Frozen[T] {
	small, large := f, other
	if len(small.m) > len(large.m) {
		small, large = large, small
	}
	out := make(map[T]struct{}, len(small.m))
	for item := range small.m {
		if large.Has(item) {
			out[item] = struct{}{}
		}
	}
	return Frozen[T]{m: out}
}

// Difference 返回差集。
func (f Frozen[T]) Difference(other Frozen[T]) Frozen[T] {
	out := make(map[T]struct{}, len(f.m))
	for item := range f.m {
		if !other.Has(item) {
			out[item] = struct{}{}
		}
	}
	return Frozen[T]{m: out}
}

// IsSubset 判断是否为子集。
func (f Frozen[T]) IsSubset(other Frozen[T]) bool {
	if len(f.m) > len(other.m) {
		return false
	}
	for item := range f.m {
		if !other.Has(item) {
			return false
		}
	}
	return true
}

// Equal 判断两个不可变集合是否相等。
func (f Frozen[T]) Equal(other Frozen[T]) bool {
	return len(f.m) == len(other.m) && f.IsSubset(other)
}

// Clone 返回自身；Frozen 不可变，无需复制底层 map。
func (f Frozen[T]) Clone() Frozen[T] {
	return f
}
