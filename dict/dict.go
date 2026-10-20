// Package dict 提供 Python dict 风格的容器与操作方法。
//
// Dict 是 map[string]interface{} 的命名类型，可以直接 range、delete，
// 也可以使用 get/setdefault/keys/values/items/update 等 Python 同名方法。
//
// 由于 Go 的 map 不保留插入顺序，Keys/Values/Items 一律按键排序返回，
// 保证结果稳定、可测试。
package dict

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// Dict 是一个键为 string、值为任意类型的字典。
type Dict map[string]interface{}

// Item 是 Items() 返回的键值对。
type Item struct {
	Key   string
	Value interface{}
}

// New 创建一个空字典。
func New() Dict {
	return Dict{}
}

// From 等价于 Python 的 dict(other)，对 map 做浅拷贝。
func From(m map[string]interface{}) Dict {
	d := make(Dict, len(m))
	for k, v := range m {
		d[k] = v
	}
	return d
}

// Get 等价于 Python 的 d.get(key, default)，默认值为 nil。
func (d Dict) Get(key string, def ...interface{}) interface{} {
	if v, ok := d[key]; ok {
		return v
	}
	return defaultValue(def)
}

// Set 等价于 Python 的 d[key] = value，返回 d 以支持链式调用。
func (d Dict) Set(key string, value interface{}) Dict {
	d[key] = value
	return d
}

// SetDefault 等价于 Python 的 d.setdefault(key, default)：键不存在时写入默认值。
func (d Dict) SetDefault(key string, def interface{}) interface{} {
	if v, ok := d[key]; ok {
		return v
	}
	d[key] = def
	return def
}

// Has 等价于 Python 的 key in d。
func (d Dict) Has(key string) bool {
	_, ok := d[key]
	return ok
}

// Delete 等价于 Python 的 del d[key]，返回该键原先是否存在。
func (d Dict) Delete(key string) bool {
	if _, ok := d[key]; !ok {
		return false
	}
	delete(d, key)
	return true
}

// Pop 等价于 Python 的 d.pop(key[, default])。键不存在时返回默认值或 nil。
func (d Dict) Pop(key string, def ...interface{}) interface{} {
	if v, ok := d[key]; ok {
		delete(d, key)
		return v
	}
	return defaultValue(def)
}

// PopItem 等价于 Python 3.7+ 的 d.popitem()，删除并返回一个键值对。
// 由于 Go map 无序，这里按键排序后取第一个，保证行为确定。
func (d Dict) PopItem() (Item, error) {
	if len(d) == 0 {
		return Item{}, fmt.Errorf("dict: popitem(): dictionary is empty")
	}
	keys := d.Keys()
	key := keys[0]
	value := d[key]
	delete(d, key)
	return Item{Key: key, Value: value}, nil
}

// Clear 等价于 Python 的 d.clear()。
func (d Dict) Clear() Dict {
	for k := range d {
		delete(d, k)
	}
	return d
}

// Keys 等价于 Python 的 d.keys()，按键排序返回。
func (d Dict) Keys() []string {
	keys := make([]string, 0, len(d))
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Values 等价于 Python 的 d.values()，顺序与 Keys 一致。
func (d Dict) Values() []interface{} {
	keys := d.Keys()
	values := make([]interface{}, len(keys))
	for i, k := range keys {
		values[i] = d[k]
	}
	return values
}

// Items 等价于 Python 的 d.items()，顺序与 Keys 一致。
func (d Dict) Items() []Item {
	keys := d.Keys()
	items := make([]Item, len(keys))
	for i, k := range keys {
		items[i] = Item{Key: k, Value: d[k]}
	}
	return items
}

// Update 等价于 Python 的 d.update(other)，支持 Dict、map[string]interface{}
// 以及其他键为 string 的 map。
func (d Dict) Update(other interface{}) error {
	if other == nil {
		return nil
	}
	switch m := other.(type) {
	case Dict:
		for k, v := range m {
			d[k] = v
		}
		return nil
	case map[string]interface{}:
		for k, v := range m {
			d[k] = v
		}
		return nil
	}

	rv := reflect.ValueOf(other)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("dict: update() expected a mapping, got %s", pyval.TypeName(other))
	}
	for _, key := range rv.MapKeys() {
		d[key.String()] = rv.MapIndex(key).Interface()
	}
	return nil
}

// Merge 返回一个新的 Dict，包含 d 与 others 的所有键（后者覆盖前者）。
func (d Dict) Merge(others ...Dict) Dict {
	out := d.Copy()
	for _, other := range others {
		for k, v := range other {
			out[k] = v
		}
	}
	return out
}

// Copy 等价于 Python 的 d.copy()，返回浅拷贝。
func (d Dict) Copy() Dict {
	out := make(Dict, len(d))
	for k, v := range d {
		out[k] = v
	}
	return out
}

// Len 等价于 Python 的 len(d)。
func (d Dict) Len() int {
	return len(d)
}

// GetPath 按点号路径读取嵌套值，例如 "user.address.city"；
// 路径中的数字会被当作切片下标。不存在时返回默认值。
func (d Dict) GetPath(path string, def ...interface{}) interface{} {
	var current interface{} = map[string]interface{}(d)
	if path == "" {
		return defaultValue(def)
	}
	for _, part := range strings.Split(path, ".") {
		switch node := current.(type) {
		case Dict:
			value, ok := node[part]
			if !ok {
				return defaultValue(def)
			}
			current = value
		case map[string]interface{}:
			value, ok := node[part]
			if !ok {
				return defaultValue(def)
			}
			current = value
		case []interface{}:
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 0 || idx >= len(node) {
				return defaultValue(def)
			}
			current = node[idx]
		default:
			// 支持命名切片（如 list.List）和任意 string 键 map。
			rv := reflect.ValueOf(current)
			for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
				if rv.IsNil() {
					return defaultValue(def)
				}
				rv = rv.Elem()
			}
			switch rv.Kind() {
			case reflect.Map:
				if rv.Type().Key().Kind() != reflect.String {
					return defaultValue(def)
				}
				value := rv.MapIndex(reflect.ValueOf(part).Convert(rv.Type().Key()))
				if !value.IsValid() {
					return defaultValue(def)
				}
				current = value.Interface()
			case reflect.Slice, reflect.Array:
				idx, err := strconv.Atoi(part)
				if err != nil || idx < 0 || idx >= rv.Len() {
					return defaultValue(def)
				}
				current = rv.Index(idx).Interface()
			default:
				return defaultValue(def)
			}
		}
	}
	return current
}

// ToJSON 等价于 Python 的 json.dumps(d)，不转义 HTML 字符。
func (d Dict) ToJSON() (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(map[string]interface{}(d)); err != nil {
		return "", err
	}
	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// Equal 深比较两个字典。
func (d Dict) Equal(other Dict) bool {
	return pyval.Equal(map[string]interface{}(d), map[string]interface{}(other))
}

// String 返回 Python 风格的字典表示，例如 {'a': 1, 'b': 'x'}。
func (d Dict) String() string {
	return pyval.Repr(map[string]interface{}(d))
}

func defaultValue(def []interface{}) interface{} {
	if len(def) > 0 {
		return def[0]
	}
	return nil
}
