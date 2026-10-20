// Package conv 提供 Python 风格的类型转换：int()、float()、bool()、str()、
// list()，以及 JSON 与结构体/字典之间的转换。
//
// 解析失败一律返回 error（而不是 panic），适合处理外部输入；如果你确定
// 类型正确，可以使用 XXXDefault 系列或 py 包。
package conv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// Int 等价于 Python 的 int(v)，结果受 Go 的 int 位宽限制。
func Int(v interface{}) (int, error) {
	n, err := Int64(v)
	if err != nil {
		return 0, err
	}
	if strconv.IntSize == 32 && (n < math.MinInt32 || n > math.MaxInt32) {
		return 0, fmt.Errorf("conv.Int: %d overflows int32", n)
	}
	return int(n), nil
}

// Int64 等价于 Python 的 int(v)：
//   - 字符串按十进制解析并忽略首尾空白；
//   - 浮点数向零截断；
//   - bool 转成 1/0。
func Int64(v interface{}) (int64, error) {
	if v == nil {
		return 0, fmt.Errorf("conv.Int64: cannot convert NoneType to int")
	}
	if data, ok := v.([]byte); ok {
		v = string(data)
	}
	if num, ok := v.(json.Number); ok {
		n, err := num.Int64()
		if err != nil {
			return 0, fmt.Errorf("conv.Int64: invalid json.Number %q: %w", num, err)
		}
		return n, nil
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return 0, fmt.Errorf("conv.Int64: cannot convert NoneType to int")
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u := rv.Uint()
		if u > math.MaxInt64 {
			return 0, fmt.Errorf("conv.Int64: %d overflows int64", u)
		}
		return int64(u), nil
	case reflect.Float32, reflect.Float64:
		f := rv.Float()
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return 0, fmt.Errorf("conv.Int64: cannot convert %v to int", f)
		}
		if f >= float64(1<<63) || f < -float64(1<<63) {
			return 0, fmt.Errorf("conv.Int64: %v overflows int64", f)
		}
		return int64(f), nil
	case reflect.String:
		s := strings.TrimSpace(rv.String())
		if s == "" {
			return 0, fmt.Errorf("conv.Int64: cannot convert empty string to int")
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("conv.Int64: invalid literal for int(): %q", s)
		}
		return n, nil
	}
	return 0, fmt.Errorf("conv.Int64: cannot convert %s to int", pyval.TypeName(v))
}

// Uint64 等价于 Python 的非负整数转换，负数会返回 error。
func Uint64(v interface{}) (uint64, error) {
	if v == nil {
		return 0, fmt.Errorf("conv.Uint64: cannot convert NoneType to int")
	}
	if data, ok := v.([]byte); ok {
		v = string(data)
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return 0, fmt.Errorf("conv.Uint64: cannot convert NoneType to int")
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := rv.Int()
		if n < 0 {
			return 0, fmt.Errorf("conv.Uint64: cannot convert negative number %d", n)
		}
		return uint64(n), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint(), nil
	case reflect.Float32, reflect.Float64:
		f := rv.Float()
		if math.IsNaN(f) || math.IsInf(f, 0) || f < 0 {
			return 0, fmt.Errorf("conv.Uint64: cannot convert %v", f)
		}
		if f >= float64(1<<64) {
			return 0, fmt.Errorf("conv.Uint64: %v overflows uint64", f)
		}
		return uint64(f), nil
	case reflect.String:
		s := strings.TrimSpace(rv.String())
		if s == "" {
			return 0, fmt.Errorf("conv.Uint64: cannot convert empty string to int")
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("conv.Uint64: invalid literal for int(): %q", s)
		}
		return n, nil
	}
	return 0, fmt.Errorf("conv.Uint64: cannot convert %s to int", pyval.TypeName(v))
}

// Float 等价于 Python 的 float(v)。
func Float(v interface{}) (float64, error) {
	if v == nil {
		return 0, fmt.Errorf("conv.Float: cannot convert NoneType to float")
	}
	if data, ok := v.([]byte); ok {
		v = string(data)
	}
	if num, ok := v.(json.Number); ok {
		f, err := num.Float64()
		if err != nil {
			return 0, fmt.Errorf("conv.Float: invalid json.Number %q: %w", num, err)
		}
		return f, nil
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return 0, fmt.Errorf("conv.Float: cannot convert NoneType to float")
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil
	case reflect.String:
		s := strings.TrimSpace(rv.String())
		if s == "" {
			return 0, fmt.Errorf("conv.Float: could not convert string to float: ''")
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("conv.Float: could not convert string to float: %q", s)
		}
		return f, nil
	}
	return 0, fmt.Errorf("conv.Float: cannot convert %s to float", pyval.TypeName(v))
}

// Bool 把常见表示转换为 bool：支持 "1"/"true"/"yes"/"on" 等字符串，
// 数字非零为 true，nil 为 false。其他容器按 Python 真值规则处理。
func Bool(v interface{}) (bool, error) {
	if v == nil {
		return false, nil
	}
	if data, ok := v.([]byte); ok {
		v = string(data)
	}
	if s, ok := v.(string); ok {
		return parseBoolString(strings.TrimSpace(s))
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return false, nil
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Bool:
		return rv.Bool(), nil
	case reflect.String:
		return parseBoolString(strings.TrimSpace(rv.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0, nil
	}
	return pyval.Truthy(v), nil
}

func parseBoolString(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "1", "t", "true", "yes", "y", "on":
		return true, nil
	case "0", "f", "false", "no", "n", "off", "":
		return false, nil
	}
	return false, fmt.Errorf("conv.Bool: cannot parse %q as bool", s)
}

// String 等价于 Python 的 str(v)：nil 为 "None"，bool 为 "True"/"False"，
// 容器使用 repr 形式。
func String(v interface{}) string {
	return pyval.Str(v)
}

// IntDefault 转换失败时返回 def。
func IntDefault(v interface{}, def int) int {
	if n, err := Int(v); err == nil {
		return n
	}
	return def
}

// FloatDefault 转换失败时返回 def。
func FloatDefault(v interface{}, def float64) float64 {
	if f, err := Float(v); err == nil {
		return f
	}
	return def
}

// BoolDefault 转换失败时返回 def。
func BoolDefault(v interface{}, def bool) bool {
	if b, err := Bool(v); err == nil {
		return b
	}
	return def
}

// Strings 把可迭代对象转换为 []string。
func Strings(v interface{}) ([]string, error) {
	items, err := pyval.Iter(v)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = String(item)
	}
	return out, nil
}

// Ints 把可迭代对象转换为 []int。
func Ints(v interface{}) ([]int, error) {
	items, err := pyval.Iter(v)
	if err != nil {
		return nil, err
	}
	out := make([]int, len(items))
	for i, item := range items {
		n, err := Int(item)
		if err != nil {
			return nil, fmt.Errorf("conv.Ints: element %d: %w", i, err)
		}
		out[i] = n
	}
	return out, nil
}

// Floats 把可迭代对象转换为 []float64。
func Floats(v interface{}) ([]float64, error) {
	items, err := pyval.Iter(v)
	if err != nil {
		return nil, err
	}
	out := make([]float64, len(items))
	for i, item := range items {
		f, err := Float(item)
		if err != nil {
			return nil, fmt.Errorf("conv.Floats: element %d: %w", i, err)
		}
		out[i] = f
	}
	return out, nil
}

// InterfaceSlice 等价于 Python 的 list(v)。
func InterfaceSlice(v interface{}) ([]interface{}, error) {
	return pyval.Iter(v)
}

// ToSlice 是 InterfaceSlice 的别名。
func ToSlice(v interface{}) []interface{} {
	items, err := pyval.Iter(v)
	if err != nil {
		return nil
	}
	return items
}

// ToJSON 等价于 Python 的 json.dumps(v, ensure_ascii=False)，不转义 HTML。
func ToJSON(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return "", err
	}
	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// ToJSONIndent 等价于 Python 的 json.dumps(v, indent=...)。
func ToJSONIndent(v interface{}, indent string) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", indent)
	if err := encoder.Encode(v); err != nil {
		return "", err
	}
	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// FromJSON 等价于 Python 的 json.loads()，支持 string 和 []byte 输入。
func FromJSON(data interface{}, out interface{}) error {
	var raw []byte
	switch x := data.(type) {
	case string:
		raw = []byte(x)
	case []byte:
		raw = x
	default:
		return fmt.Errorf("conv.FromJSON: expected string or []byte, got %s", pyval.TypeName(data))
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("conv.FromJSON: %w", err)
	}
	return nil
}
