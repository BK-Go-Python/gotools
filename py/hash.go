package py

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"strconv"
	stdtime "time"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// Hash 返回值的确定性哈希。它不是 CPython 的 hash()，而是跨进程稳定的
// FNV-1a 64 替代，用于分桶、缓存键等场景。
//
// Python: hash(v)
// DIFF: CPython 使用进程随机化种子，且不同版本的整数/字符串算法不同；
// 这里不保证与 Python 数值一致，也不保证等于对象身份。
func Hash(v any) (uint64, error) {
	if v != nil {
		rv := reflect.ValueOf(v)
		for rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				break
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Slice, reflect.Map, reflect.Func, reflect.Chan,
			reflect.Ptr, reflect.UnsafePointer:
			return 0, fmt.Errorf("py.Hash: unhashable type %s: %w", pyval.TypeName(v), ErrType)
		}
		if rv.IsValid() && !rv.Type().Comparable() {
			return 0, fmt.Errorf("py.Hash: unhashable type %s: %w", pyval.TypeName(v), ErrType)
		}
	}

	h := fnv.New64a()
	_, _ = h.Write([]byte(hashKey(v)))
	return h.Sum64(), nil
}

// hashKey 为数值生成统一的键，保证 Python 的相等对象具有相同哈希：
// hash(1) == hash(1.0) == hash(True)。
func hashKey(v any) string {
	switch x := v.(type) {
	case bool:
		if x {
			return "num:1"
		}
		return "num:0"
	case int:
		return numericKey(float64(x))
	case int8:
		return numericKey(float64(x))
	case int16:
		return numericKey(float64(x))
	case int32:
		return numericKey(float64(x))
	case int64:
		return numericKey(float64(x))
	case uint:
		return numericKey(float64(x))
	case uint8:
		return numericKey(float64(x))
	case uint16:
		return numericKey(float64(x))
	case uint32:
		return numericKey(float64(x))
	case uint64:
		return numericKey(float64(x))
	case uintptr:
		return numericKey(float64(x))
	case float32:
		return numericKey(float64(x))
	case float64:
		return numericKey(x)
	case stdtime.Time:
		return "time:" + x.Round(0).Format(stdtime.RFC3339Nano)
	case stdtime.Duration:
		return numericKey(x.Seconds())
	}
	return pyval.Repr(v)
}

func numericKey(f float64) string {
	if f == 0 {
		f = 0 // 把 -0.0 归一化为 0.0，Python hash(-0.0) == hash(0.0)
	}
	return "num:" + strconv.FormatFloat(f, 'g', -1, 64)
}
