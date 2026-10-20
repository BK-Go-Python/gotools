package dict

import "maps"

// Of 返回 m 的浅拷贝。
//
// Python: dict(m)
// DIFF: 仅适用于 Go map；需要 string 键的动态字典请使用 dict.From。
func Of[K comparable, V any](m map[K]V) map[K]V {
	return maps.Clone(m)
}
