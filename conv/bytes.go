package conv

import (
	"fmt"
	"reflect"
)

// Bytes 等价于 Python 的 bytes(v)。
//
// 支持 string（UTF-8 编码）、[]byte、整数（n 个零字节）、
// 以及元素为 0..255 的切片/数组。
//
// Python: bytes(v)
// DIFF: 不支持 encoding 参数；Go 的 []byte 可变，不可变场景请使用 string。
func Bytes(v any) ([]byte, error) {
	return bytesFrom(v)
}

// ByteArray 等价于 Python 的 bytearray(v)，返回可变的字节切片。
//
// Python: bytearray(v)
// DIFF: 不提供 append/extend 等方法，直接使用 Go 的 append 与切片操作。
func ByteArray(v any) ([]byte, error) {
	return bytesFrom(v)
}

func bytesFrom(v any) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("conv.Bytes: cannot convert NoneType to bytes")
	}
	switch x := v.(type) {
	case string:
		return []byte(x), nil
	case []byte:
		return append([]byte(nil), x...), nil
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil, fmt.Errorf("conv.Bytes: cannot convert nil to bytes")
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		out := make([]byte, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			n, err := byteValue(rv.Index(i))
			if err != nil {
				return nil, fmt.Errorf("conv.Bytes: element %d: %w", i, err)
			}
			out[i] = n
		}
		return out, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := rv.Int()
		if n < 0 {
			return nil, fmt.Errorf("conv.Bytes: negative length %d", n)
		}
		return make([]byte, n), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n := rv.Uint()
		if n > uint64(maxInt()) {
			return nil, fmt.Errorf("conv.Bytes: length %d overflows int", n)
		}
		return make([]byte, int(n)), nil
	}
	return nil, fmt.Errorf("conv.Bytes: cannot convert %s to bytes", reflect.TypeOf(v))
}

func byteValue(rv reflect.Value) (byte, error) {
	for rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return 0, fmt.Errorf("cannot convert NoneType to byte")
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := rv.Int()
		if n < 0 || n > 255 {
			return 0, fmt.Errorf("byte must be in range(0, 256), got %d", n)
		}
		return byte(n), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n := rv.Uint()
		if n > 255 {
			return 0, fmt.Errorf("byte must be in range(0, 256), got %d", n)
		}
		return byte(n), nil
	}
	return 0, fmt.Errorf("expected an integer, got %s", rv.Type())
}

func maxInt() int {
	return int(^uint(0) >> 1)
}
