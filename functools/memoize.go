// Package functools 提供 Python functools 风格的缓存工具。
package functools

import "sync"

// Memoize 返回带缓存的 fn，等价于 Python 的 functools.cache（无参数版本）。
//
// 并发安全；同一 key 的并发首次调用可能重复执行 fn，但最终缓存一个结果。
func Memoize[K comparable, V any](fn func(K) V) func(K) V {
	var mu sync.Mutex
	cache := make(map[K]V)
	return func(key K) V {
		mu.Lock()
		if value, ok := cache[key]; ok {
			mu.Unlock()
			return value
		}
		mu.Unlock()

		value := fn(key)

		mu.Lock()
		cache[key] = value
		mu.Unlock()
		return value
	}
}
