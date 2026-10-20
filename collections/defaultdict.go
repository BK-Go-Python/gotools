package collections

import (
	"fmt"
	"reflect"
	"sort"
)

// DefaultItem 是 DefaultDict.Items 返回的键值对。
type DefaultItem[K comparable, V any] struct {
	Key   K
	Value V
}

// DefaultDict 是 Python collections.defaultdict 的泛型实现：
// 访问不存在的键时调用 factory 生成默认值并写回。
type DefaultDict[K comparable, V any] struct {
	m       map[K]V
	factory func() V
}

// NewDefaultDict 创建默认值工厂为 factory 的 DefaultDict。
// factory 可以为 nil，此时缺失键返回 V 的零值。
func NewDefaultDict[K comparable, V any](factory func() V) *DefaultDict[K, V] {
	return &DefaultDict[K, V]{m: make(map[K]V), factory: factory}
}

// Get 等价于 Python 的 d[key]：键缺失时调用工厂、写回并返回。
func (d *DefaultDict[K, V]) Get(key K) V {
	if d == nil {
		var zero V
		return zero
	}
	if value, ok := d.m[key]; ok {
		return value
	}
	var value V
	if d.factory != nil {
		value = d.factory()
	}
	if d.m == nil {
		d.m = make(map[K]V)
	}
	d.m[key] = value
	return value
}

// Set 等价于 Python 的 d[key] = value。
func (d *DefaultDict[K, V]) Set(key K, value V) *DefaultDict[K, V] {
	if d == nil {
		return nil
	}
	if d.m == nil {
		d.m = make(map[K]V)
	}
	d.m[key] = value
	return d
}

// Has 等价于 Python 的 key in d（不会触发工厂）。
func (d *DefaultDict[K, V]) Has(key K) bool {
	if d == nil {
		return false
	}
	_, ok := d.m[key]
	return ok
}

// Delete 等价于 Python 的 del d[key]。
func (d *DefaultDict[K, V]) Delete(key K) bool {
	if d == nil {
		return false
	}
	if _, ok := d.m[key]; !ok {
		return false
	}
	delete(d.m, key)
	return true
}

// Keys 返回所有键，按 fmt.Sprint 排序保证确定性。
func (d *DefaultDict[K, V]) Keys() []K {
	if d == nil {
		return []K{}
	}
	keys := make([]K, 0, len(d.m))
	for key := range d.m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j]) })
	return keys
}

// Values 按 Keys 的顺序返回值。
func (d *DefaultDict[K, V]) Values() []V {
	keys := d.Keys()
	values := make([]V, len(keys))
	for i, key := range keys {
		values[i] = d.m[key]
	}
	return values
}

// Items 按 Keys 的顺序返回键值对。
func (d *DefaultDict[K, V]) Items() []DefaultItem[K, V] {
	keys := d.Keys()
	items := make([]DefaultItem[K, V], len(keys))
	for i, key := range keys {
		items[i] = DefaultItem[K, V]{Key: key, Value: d.m[key]}
	}
	return items
}

// Len 返回键数量（不会触发工厂）。
func (d *DefaultDict[K, V]) Len() int {
	if d == nil {
		return 0
	}
	return len(d.m)
}

// Clear 清空内容，保留工厂。
func (d *DefaultDict[K, V]) Clear() *DefaultDict[K, V] {
	if d != nil {
		d.m = make(map[K]V)
	}
	return d
}

// Copy 返回浅拷贝，工厂共享。
func (d *DefaultDict[K, V]) Copy() *DefaultDict[K, V] {
	if d == nil {
		return NewDefaultDict[K, V](nil)
	}
	out := NewDefaultDict[K, V](d.factory)
	for key, value := range d.m {
		out.m[key] = value
	}
	return out
}

// Equal 比较键、值和工厂是否存在；值使用 reflect.DeepEqual。
func (d *DefaultDict[K, V]) Equal(other *DefaultDict[K, V]) bool {
	if d == nil || other == nil {
		return d == nil && other == nil
	}
	if len(d.m) != len(other.m) || (d.factory == nil) != (other.factory == nil) {
		return false
	}
	for key, value := range d.m {
		otherValue, ok := other.m[key]
		if !ok || !reflect.DeepEqual(value, otherValue) {
			return false
		}
	}
	return true
}
