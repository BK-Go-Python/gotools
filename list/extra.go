package list

import (
	"fmt"
	"sort"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// IndexRange 等价于 Python 的 list.index(value, start, stop)。
//
// start/stop 支持负下标，stop 不包含在搜索范围内。
//
// Python: list.index(value, start, stop)
func (l List) IndexRange(value interface{}, start, stop int) (int, error) {
	n := len(l)
	if start < 0 {
		start += n
	}
	if stop < 0 {
		stop += n
	}
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}
	if stop < 0 {
		stop = 0
	}
	if stop > n {
		stop = n
	}
	for i := start; i < stop; i++ {
		if pyval.Equal(l[i], value) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("list: %s is not in list", pyval.Repr(value))
}

// SortKey 等价于 Python 的 list.sort(key=..., reverse=...)。
//
// 每个元素的 key 只计算一次（与 Python 一致），排序是稳定的。
func (l *List) SortKey(key func(interface{}) interface{}, reverse ...bool) *List {
	type entry struct {
		value interface{}
		key   interface{}
	}
	entries := make([]entry, len(*l))
	for i, value := range *l {
		entries[i] = entry{value: value, key: key(value)}
	}

	desc := len(reverse) > 0 && reverse[0]
	sort.SliceStable(entries, func(i, j int) bool {
		c, err := pyval.Compare(entries[i].key, entries[j].key)
		if err != nil {
			left, right := pyval.Repr(entries[i].key), pyval.Repr(entries[j].key)
			if desc {
				return left > right
			}
			return left < right
		}
		if desc {
			return c > 0
		}
		return c < 0
	})

	for i, entry := range entries {
		(*l)[i] = entry.value
	}
	return l
}
