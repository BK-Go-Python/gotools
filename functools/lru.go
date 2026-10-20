package functools

import (
	"container/list"
	"sync"
)

type lruEntry[K comparable, V any] struct {
	key   K
	value V
}

// LruCache 是固定容量的 LRU 缓存，对应 Python 的 functools.lru_cache(maxsize=n)。
//
// 零值不可用，请使用 NewLruCache；capacity <= 0 时缓存不存储任何元素。
type LruCache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	order    *list.List
	items    map[K]*list.Element
}

// NewLruCache 创建容量为 capacity 的 LRU 缓存。
func NewLruCache[K comparable, V any](capacity int) *LruCache[K, V] {
	if capacity < 0 {
		capacity = 0
	}
	return &LruCache[K, V]{
		capacity: capacity,
		order:    list.New(),
		items:    make(map[K]*list.Element),
	}
}

// Get 返回 key 对应的值，并把该 key 标记为最近使用。
func (c *LruCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if element, ok := c.items[key]; ok {
		c.order.MoveToFront(element)
		entry := element.Value.(*lruEntry[K, V])
		return entry.value, true
	}
	var zero V
	return zero, false
}

// Put 写入 key/value；容量为 0 时不存储。
func (c *LruCache[K, V]) Put(key K, value V) {
	if c.capacity <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if element, ok := c.items[key]; ok {
		entry := element.Value.(*lruEntry[K, V])
		entry.value = value
		c.order.MoveToFront(element)
		return
	}
	element := c.order.PushFront(&lruEntry[K, V]{key: key, value: value})
	c.items[key] = element
	for c.order.Len() > c.capacity {
		back := c.order.Back()
		if back == nil {
			break
		}
		entry := back.Value.(*lruEntry[K, V])
		delete(c.items, entry.key)
		c.order.Remove(back)
	}
}

// Remove 删除 key，返回是否存在。
func (c *LruCache[K, V]) Remove(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.items[key]
	if !ok {
		return false
	}
	delete(c.items, key)
	c.order.Remove(element)
	return true
}

// Len 返回当前缓存元素数量。
func (c *LruCache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

// Cap 返回缓存容量。
func (c *LruCache[K, V]) Cap() int {
	return c.capacity
}

// Keys 按最近使用到最久未使用返回键。
func (c *LruCache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]K, 0, c.order.Len())
	for element := c.order.Front(); element != nil; element = element.Next() {
		keys = append(keys, element.Value.(*lruEntry[K, V]).key)
	}
	return keys
}

// Clear 清空缓存。
func (c *LruCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.order.Init()
	c.items = make(map[K]*list.Element)
}
