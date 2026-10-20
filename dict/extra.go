package dict

// FromKeys 等价于 Python 的 dict.fromkeys(keys, value=None)。
//
// Python: dict.fromkeys(keys, value)
func FromKeys(keys []string, value ...interface{}) Dict {
	val := interface{}(nil)
	if len(value) > 0 {
		val = value[0]
	}
	d := make(Dict, len(keys))
	for _, key := range keys {
		d[key] = val
	}
	return d
}
