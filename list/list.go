// Package list 提供 Python list 风格的容器与操作方法。
//
// List 是 []interface{} 的命名类型，因此可以直接 range、len、copy，
// 也可以使用本包提供的 append/extend/insert/pop/sort 等 Python 同名方法。
//
// 与 Go 切片的习惯不同，修改方法（Append/Insert/Remove/Sort/...）会直接
// 修改接收者，并返回 *List 以便链式调用，这与 Python list 的原地修改一致。
package list

import (
	"fmt"
	"sort"
	"strings"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// List 是一个可以存放任意类型元素的列表。
type List []interface{}

// New 创建一个 List，参数会按顺序放入列表。
func New(items ...interface{}) List {
	l := make(List, len(items))
	copy(l, items)
	return l
}

// From 把可迭代对象（切片、数组、字符串、map 的键）转换为 List。
func From(v interface{}) (List, error) {
	items, err := pyval.Iter(v)
	if err != nil {
		return nil, err
	}
	return List(items), nil
}

// MustFrom 与 From 相同，但转换失败时 panic，适合确定类型的场景。
func MustFrom(v interface{}) List {
	l, err := From(v)
	if err != nil {
		panic(err)
	}
	return l
}

// Append 等价于 Python 的 list.append/extend，追加一个或多个元素。
func (l *List) Append(items ...interface{}) *List {
	*l = append(*l, items...)
	return l
}

// Extend 等价于 Python 的 list.extend(iterable)。
func (l *List) Extend(v interface{}) error {
	items, err := pyval.Iter(v)
	if err != nil {
		return err
	}
	*l = append(*l, items...)
	return nil
}

// Insert 等价于 Python 的 list.insert(index, value)。
// 下标为负时从末尾计算，越界会自动截断到 [0, len]。
func (l *List) Insert(index int, value interface{}) *List {
	n := len(*l)
	if index < 0 {
		index += n
	}
	if index < 0 {
		index = 0
	}
	if index > n {
		index = n
	}
	*l = append(*l, nil)
	copy((*l)[index+1:], (*l)[index:])
	(*l)[index] = value
	return l
}

// Remove 等价于 Python 的 list.remove(value)，删除第一个相等元素。
func (l *List) Remove(value interface{}) error {
	for i, item := range *l {
		if pyval.Equal(item, value) {
			*l = append((*l)[:i], (*l)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("list: %s not in list", pyval.Repr(value))
}

// Pop 等价于 Python 的 list.pop()：删除并返回最后一个元素。
func (l *List) Pop() (interface{}, error) {
	return l.PopAt(-1)
}

// PopAt 等价于 Python 的 list.pop(index)。
func (l *List) PopAt(index int) (interface{}, error) {
	n := len(*l)
	if n == 0 {
		return nil, fmt.Errorf("list: pop from empty list")
	}
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		return nil, fmt.Errorf("list: pop index out of range: %d", index)
	}
	value := (*l)[index]
	*l = append((*l)[:index], (*l)[index+1:]...)
	return value, nil
}

// Clear 等价于 Python 的 list.clear()。
func (l *List) Clear() *List {
	*l = (*l)[:0]
	return l
}

// Copy 等价于 Python 的 list.copy()，返回浅拷贝。
func (l *List) Copy() List {
	c := make(List, len(*l))
	copy(c, *l)
	return c
}

// Get 等价于 Python 的 list[index]，支持负数下标。
func (l List) Get(index int) (interface{}, error) {
	n := len(l)
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		return nil, fmt.Errorf("list: index out of range: %d", index)
	}
	return l[index], nil
}

// Set 等价于 Python 的 list[index] = value，支持负数下标。
func (l List) Set(index int, value interface{}) error {
	n := len(l)
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		return fmt.Errorf("list: index out of range: %d", index)
	}
	l[index] = value
	return nil
}

// Index 等价于 Python 的 list.index(value)，返回第一个相等元素的下标。
func (l List) Index(value interface{}) (int, error) {
	for i, item := range l {
		if pyval.Equal(item, value) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("list: %s is not in list", pyval.Repr(value))
}

// Count 等价于 Python 的 list.count(value)。
func (l List) Count(value interface{}) int {
	n := 0
	for _, item := range l {
		if pyval.Equal(item, value) {
			n++
		}
	}
	return n
}

// Contains 等价于 Python 的 value in list。
func (l List) Contains(value interface{}) bool {
	return l.Count(value) > 0
}

// Reverse 等价于 Python 的 list.reverse()，原地反转。
func (l *List) Reverse() *List {
	for i, j := 0, len(*l)-1; i < j; i, j = i+1, j-1 {
		(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
	}
	return l
}

// Sort 等价于 Python 的 list.sort()，原地稳定排序。
//
// 不传比较函数时使用 Python 风格的比较（数字、字符串、列表可比较，
// 不可比较的类型退化为 repr 比较，不会 panic）。传入 less 时使用自定义比较。
func (l *List) Sort(less ...func(a, b interface{}) bool) *List {
	if len(less) > 0 && less[0] != nil {
		sort.SliceStable(*l, func(i, j int) bool {
			return less[0]((*l)[i], (*l)[j])
		})
		return l
	}
	pyval.SortValues(*l)
	return l
}

// SortBy 等价于 Python 的 list.sort(key=...)（可选 reverse）。
func (l *List) SortBy(key func(interface{}) interface{}, reverse ...bool) *List {
	desc := len(reverse) > 0 && reverse[0]
	sort.SliceStable(*l, func(i, j int) bool {
		a, b := key((*l)[i]), key((*l)[j])
		c, err := pyval.Compare(a, b)
		if err != nil {
			sa, sb := pyval.Repr(a), pyval.Repr(b)
			if desc {
				return sa > sb
			}
			return sa < sb
		}
		if desc {
			return c > 0
		}
		return c < 0
	})
	return l
}

// Slice 等价于 Python 的 l[start:end]，支持负数下标，越界自动截断。
func (l List) Slice(start, end int) List {
	n := len(l)
	if start < 0 {
		start += n
	}
	if end < 0 {
		end += n
	}
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}
	if end < 0 {
		end = 0
	}
	if end > n {
		end = n
	}
	if start >= end {
		return List{}
	}
	out := make(List, end-start)
	copy(out, l[start:end])
	return out
}

// ToStrings 把每个元素按 Python str() 转成字符串。
func (l List) ToStrings() []string {
	out := make([]string, len(l))
	for i, item := range l {
		out[i] = pyval.Str(item)
	}
	return out
}

// Join 等价于 Python 的 sep.join(list)。
func (l List) Join(sep string) string {
	return strings.Join(l.ToStrings(), sep)
}

// Map 等价于 Python 的 [f(x) for x in list]。
func (l List) Map(fn func(interface{}) interface{}) List {
	out := make(List, len(l))
	for i, item := range l {
		out[i] = fn(item)
	}
	return out
}

// Filter 等价于 Python 的 [x for x in list if f(x)]。
func (l List) Filter(fn func(interface{}) bool) List {
	out := make(List, 0, len(l))
	for _, item := range l {
		if fn(item) {
			out = append(out, item)
		}
	}
	return out
}

// Each 遍历列表，等价于 Python 的 for i, x in enumerate(list)。
func (l List) Each(fn func(index int, value interface{})) {
	for i, item := range l {
		fn(i, item)
	}
}

// Equal 深比较两个列表，数字按 Python 语义跨类型比较（1 == 1.0）。
func (l List) Equal(other List) bool {
	return pyval.Equal([]interface{}(l), []interface{}(other))
}

// Len 等价于 Python 的 len(list)。
func (l List) Len() int {
	return len(l)
}

// String 返回 Python 风格的列表表示，例如 [1, 'a', None]。
func (l List) String() string {
	return pyval.Repr([]interface{}(l))
}
