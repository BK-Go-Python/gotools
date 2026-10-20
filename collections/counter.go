package collections

import (
	"fmt"
	"sort"
)

// CountEntry 是 MostCommon 返回的元素及其计数，对应 Python 的 (元素, 计数) 元组。
type CountEntry[T comparable] struct {
	Value T
	Count int
}

// Counter 是 Python collections.Counter 的泛型实现：把可比较的值映射到计数。
//
// 零值可用，nil map 等价于空计数器。Add、Update、Subtract 在 nil 接收者上
// 会分配新 map 并返回，因此必须写回：
//
//	var c collections.Counter[string]
//	c = c.Add("a").Add("b")
//
// 更推荐使用 NewCounter。
type Counter[T comparable] map[T]int

// NewCounter 统计 items 中每个元素出现的次数。
//
// Python: collections.Counter(items)
func NewCounter[T comparable](items ...T) Counter[T] {
	c := make(Counter[T], len(items))
	for _, item := range items {
		c[item]++
	}
	return c
}

// Add 把 item 的计数增加 n，省略 n 时增加 1。
//
// n 可以为负数；只使用第一个 n，多余参数会被忽略。即使 n 为 0，也会创建
// 计数为 0 的键，这与 Python 的 c[item] += 0 一致。
//
// Python: c[item] += n
func (c Counter[T]) Add(item T, n ...int) Counter[T] {
	if c == nil {
		c = make(Counter[T], 1)
	}
	delta := 1
	if len(n) > 0 {
		delta = n[0]
	}
	c[item] += delta
	return c
}

// Count 返回 item 的计数，不存在时返回 0。
//
// Python: c[item]
func (c Counter[T]) Count(item T) int {
	return c[item]
}

// Total 返回所有计数之和，包含 0 和负数计数。
//
// Python: sum(c.values())
func (c Counter[T]) Total() int {
	total := 0
	for _, n := range c {
		total += n
	}
	return total
}

// MostCommon 按计数从高到低返回至多 n 个 CountEntry。
//
// n 省略或为负数时返回全部；n 为 0 时返回空切片；n 超过键数时返回全部。
// 只使用第一个 n，多余参数会被忽略。计数相同时按 fmt.Sprint(Value) 升序，
// 若该结果仍相同则按类型名升序，以保证结果确定。
//
// Python: c.most_common(n)
//
// DIFF: Python 在计数并列时保持字典插入顺序；Go map 没有插入顺序，因此这里
// 按 fmt.Sprint(Value) 升序排序，并列项的顺序可能与 Python 不同。
// DIFF: 返回 []CountEntry[T]，而不是 Python 的 (元素, 计数) 元组列表。
func (c Counter[T]) MostCommon(n ...int) []CountEntry[T] {
	limit := -1
	if len(n) > 0 {
		limit = n[0]
	}
	if limit == 0 {
		return []CountEntry[T]{}
	}

	entries := make([]CountEntry[T], 0, len(c))
	for value, count := range c {
		entries = append(entries, CountEntry[T]{Value: value, Count: count})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		left, right := fmt.Sprint(entries[i].Value), fmt.Sprint(entries[j].Value)
		if left != right {
			return left < right
		}
		// fmt.Sprint 相同（例如混合类型中的 1 与 "1"）时再按类型名排序，
		// 避免 map 遍历顺序影响结果。
		return fmt.Sprintf("%T", entries[i].Value) < fmt.Sprintf("%T", entries[j].Value)
	})

	if limit < 0 || limit > len(entries) {
		limit = len(entries)
	}
	return entries[:limit]
}

// Elements 按计数重复元素并返回切片，计数为 0 或负数的元素不会输出。
//
// 元素按 fmt.Sprint(Value) 升序分组排列；若该结果仍相同则按类型名升序，
// 以保证结果确定。同一元素的重复项相邻。
//
// Python: list(c.elements())
//
// DIFF: Python 的 elements() 按底层字典的迭代顺序产出（文档说明顺序任意），
// 不额外排序；Go map 没有顺序，因此这里按 fmt.Sprint(Value) 升序排序，
// 输出顺序可能与 Python 不同。
// DIFF: 返回 []T，而不是 Python 的迭代器。
func (c Counter[T]) Elements() []T {
	values := make([]T, 0, len(c))
	for value, count := range c {
		if count > 0 {
			values = append(values, value)
		}
	}
	sort.Slice(values, func(i, j int) bool {
		left, right := fmt.Sprint(values[i]), fmt.Sprint(values[j])
		if left != right {
			return left < right
		}
		return fmt.Sprintf("%T", values[i]) < fmt.Sprintf("%T", values[j])
	})

	// 不按计数预分配容量：计数之和可能溢出 int。重复 append 更安全。
	out := make([]T, 0)
	for _, value := range values {
		for n := 0; n < c[value]; n++ {
			out = append(out, value)
		}
	}
	return out
}

// Update 把 items 中每个元素的计数加 1，允许计数从负数回升。
//
// Python: c.update(items)
func (c Counter[T]) Update(items ...T) Counter[T] {
	if c == nil {
		c = make(Counter[T], len(items))
	}
	for _, item := range items {
		c[item]++
	}
	return c
}

// Subtract 把 items 中每个元素的计数减 1。计数可以为负数。
//
// Python: c.subtract(items)
func (c Counter[T]) Subtract(items ...T) Counter[T] {
	if c == nil {
		c = make(Counter[T], len(items))
	}
	for _, item := range items {
		c[item]--
	}
	return c
}

// Copy 返回计数器的浅拷贝。nil 接收者返回空计数器，而不是 nil。
//
// Python: c.copy()
func (c Counter[T]) Copy() Counter[T] {
	out := make(Counter[T], len(c))
	for value, count := range c {
		out[value] = count
	}
	return out
}

// Equal 比较两个计数器的键集合与计数是否完全相同。
//
// 两个 nil/空计数器相等；只有键和对应计数都相同才返回 true。
//
// Python: c1 == c2
//
// DIFF: Python 把缺失键视为计数 0，因此 Counter(a=0) == Counter() 为 True；
// 这里按本包约定严格比较键集合与计数，两者不相等。
func (c Counter[T]) Equal(other Counter[T]) bool {
	if len(c) != len(other) {
		return false
	}
	for value, count := range c {
		otherCount, ok := other[value]
		if !ok || count != otherCount {
			return false
		}
	}
	return true
}
