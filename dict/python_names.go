package dict

// 本文件提供 Python 拼写别名。

// Fromkeys 是 FromKeys 的 Python 拼写别名。
func Fromkeys(keys []string, value ...interface{}) Dict {
	return FromKeys(keys, value...)
}

// Setdefault 是 SetDefault 的 Python 拼写别名；默认值为 nil。
func (d Dict) Setdefault(key string, def ...interface{}) interface{} {
	var value interface{}
	if len(def) > 0 {
		value = def[0]
	}
	return d.SetDefault(key, value)
}

// Popitem 是 PopItem 的 Python 拼写别名。
func (d Dict) Popitem() (Item, error) {
	return d.PopItem()
}
