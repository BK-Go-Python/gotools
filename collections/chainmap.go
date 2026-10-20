package collections

import (
	"fmt"
	"sort"
)

// ChainItem 是 ChainMap.Items 返回的键值对。
type ChainItem[K comparable, V any] struct {
	Key   K
	Value V
}

// ChainMap 是 Python collections.ChainMap 的泛型实现：
// 按顺序在多个 map 中查找，写操作只作用于第一个 map。
type ChainMap[K comparable, V any] struct {
	maps []map[K]V
}

// NewChainMap 按优先级创建 ChainMap，maps[0] 优先级最高。
func NewChainMap[K comparable, V any](maps ...map[K]V) *ChainMap[K, V] {
	out := make([]map[K]V, len(maps))
	copy(out, maps)
	return &ChainMap[K, V]{maps: out}
}

// Get 等价于 Python 的 cm[key]，键不存在时返回 def 或零值。
func (c *ChainMap[K, V]) Get(key K, def ...V) V {
	if c != nil {
		for _, m := range c.maps {
			if value, ok := m[key]; ok {
				return value
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	var zero V
	return zero
}

// Has 等价于 Python 的 key in cm。
func (c *ChainMap[K, V]) Has(key K) bool {
	if c == nil {
		return false
	}
	for _, m := range c.maps {
		if _, ok := m[key]; ok {
			return true
		}
	}
	return false
}

// Set 等价于 Python 的 cm[key] = value：写入第一个 map。
func (c *ChainMap[K, V]) Set(key K, value V) *ChainMap[K, V] {
	if c == nil {
		return nil
	}
	if len(c.maps) == 0 {
		c.maps = append(c.maps, make(map[K]V))
	}
	if c.maps[0] == nil {
		c.maps[0] = make(map[K]V)
	}
	c.maps[0][key] = value
	return c
}

// Delete 等价于 Python 的 del cm[key]：只从第一个 map 删除。
func (c *ChainMap[K, V]) Delete(key K) bool {
	if c == nil || len(c.maps) == 0 {
		return false
	}
	if _, ok := c.maps[0][key]; !ok {
		return false
	}
	delete(c.maps[0], key)
	return true
}

// Keys 返回所有 map 中的唯一键；顺序按 fmt.Sprint 排序以保证确定性。
func (c *ChainMap[K, V]) Keys() []K {
	return chainKeys(c)
}

// Items 返回所有 map 中的唯一键值对，顺序与 Keys 一致。
func (c *ChainMap[K, V]) Items() []ChainItem[K, V] {
	keys := c.Keys()
	items := make([]ChainItem[K, V], len(keys))
	for i, key := range keys {
		items[i] = ChainItem[K, V]{Key: key, Value: c.Get(key)}
	}
	return items
}

// Len 返回唯一键数量。
func (c *ChainMap[K, V]) Len() int {
	return len(c.Keys())
}

// Maps 返回底层 map 切片副本（元素仍共享）。
func (c *ChainMap[K, V]) Maps() []map[K]V {
	if c == nil {
		return nil
	}
	out := make([]map[K]V, len(c.maps))
	copy(out, c.maps)
	return out
}

// NewChild 返回一个新 ChainMap，最前面是新的空 map，后面共享当前 maps。
func (c *ChainMap[K, V]) NewChild() *ChainMap[K, V] {
	child := make([]map[K]V, 0, len(c.maps)+1)
	child = append(child, make(map[K]V))
	if c != nil {
		child = append(child, c.maps...)
	}
	return &ChainMap[K, V]{maps: child}
}

func chainKeys[K comparable, V any](c *ChainMap[K, V]) []K {
	if c == nil {
		return []K{}
	}
	seen := make(map[K]struct{})
	keys := make([]K, 0)
	for _, m := range c.maps {
		for key := range m {
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
	})
	return keys
}
