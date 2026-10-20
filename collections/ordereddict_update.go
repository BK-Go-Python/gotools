package collections

// Setdefault 等价于 Python 的 OrderedDict.setdefault(key, default)：
// 键存在时返回原值；不存在时追加到末尾并返回默认值。
func (d *OrderedDict[K, V]) Setdefault(key K, def V) V {
	if n, ok := d.lookup(key); ok {
		return n.value
	}
	if d != nil {
		d.Set(key, def)
	}
	return def
}

// Update 等价于 Python 的 OrderedDict.update(other)：已有键只更新值、不移动，
// 新键按 other 的顺序追加到末尾。
func (d *OrderedDict[K, V]) Update(other *OrderedDict[K, V]) *OrderedDict[K, V] {
	if d == nil || other == nil {
		return d
	}
	for _, key := range other.Keys() {
		value, _ := other.lookup(key)
		d.Set(key, value.value)
	}
	return d
}

// Merge 返回一个新的 OrderedDict，包含 d 与 other 的内容；other 覆盖 d，
// 已有键保持 d 的位置，新键按 other 的顺序追加。
func (d *OrderedDict[K, V]) Merge(other *OrderedDict[K, V]) *OrderedDict[K, V] {
	out := d.Copy()
	out.Update(other)
	return out
}
