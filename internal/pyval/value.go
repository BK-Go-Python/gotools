// Package pyval implements the Python value semantics shared by the public
// packages of gotools: type names, str()/repr(), truthiness, length,
// iteration, comparison and equality.
//
// The package is internal. The behaviour is documented on the public
// packages that expose it (str, list, dict, py, conv).
package pyval

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// TypeName returns the Python type name of v, for example "int", "str",
// "list", "dict" or "NoneType".
func TypeName(v interface{}) string {
	if v == nil {
		return "NoneType"
	}
	switch v.(type) {
	case bool:
		return "bool"
	case string:
		return "str"
	case []byte:
		return "bytes"
	case int, int8, int16, int32, int64:
		return "int"
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return "int"
	case float32, float64:
		return "float"
	case complex64, complex128:
		return "complex"
	case time.Time:
		return "datetime"
	case time.Duration:
		return "timedelta"
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return "NoneType"
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		return "str"
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Complex64, reflect.Complex128:
		return "complex"
	case reflect.Slice, reflect.Array:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return "bytes"
		}
		return "list"
	case reflect.Map:
		return "dict"
	case reflect.Func:
		return "function"
	case reflect.Struct:
		if name := rv.Type().Name(); name != "" {
			return name
		}
		return "object"
	case reflect.Chan:
		return "channel"
	}
	return rv.Kind().String()
}

// Str returns Python's str(v). Strings are returned unchanged, everything
// else uses Repr.
func Str(v interface{}) string {
	if v == nil {
		return "None"
	}
	if s, ok := v.(string); ok {
		return s
	}
	return Repr(v)
}

// Repr returns a Python-style representation of v. It is intentionally close
// to, but not byte-for-byte identical with, CPython's repr for exotic types.
func Repr(v interface{}) string {
	if v == nil {
		return "None"
	}
	switch x := v.(type) {
	case string:
		return QuoteString(x)
	case bool:
		if x {
			return "True"
		}
		return "False"
	case []byte:
		return bytesRepr(x)
	case int:
		return strconv.FormatInt(int64(x), 10)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case uintptr:
		return strconv.FormatUint(uint64(x), 10)
	case float32:
		return FormatFloat(float64(x), 32)
	case float64:
		return FormatFloat(x, 64)
	case time.Duration:
		return fmt.Sprintf("timedelta(seconds=%s)", FormatFloat(x.Seconds(), 64))
	case time.Time:
		return fmt.Sprintf("%v", x)
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return "None"
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return QuoteString(rv.String())
	case reflect.Bool:
		if rv.Bool() {
			return "True"
		}
		return "False"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32:
		return FormatFloat(rv.Float(), 32)
	case reflect.Float64:
		return FormatFloat(rv.Float(), 64)
	case reflect.Slice, reflect.Array:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return bytesRepr(rv.Bytes())
		}
		if rv.Kind() == reflect.Slice && rv.IsNil() {
			return "[]"
		}
		parts := make([]string, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			parts[i] = Repr(rv.Index(i).Interface())
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case reflect.Map:
		if rv.IsNil() {
			return "{}"
		}
		keys := rv.MapKeys()
		values := make([]interface{}, len(keys))
		for i, k := range keys {
			values[i] = k.Interface()
		}
		sort.SliceStable(values, func(i, j int) bool { return Less(values[i], values[j]) })
		parts := make([]string, 0, len(values))
		for _, key := range values {
			kv := reflect.ValueOf(key)
			parts = append(parts, Repr(key)+": "+Repr(rv.MapIndex(kv).Interface()))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case reflect.Func:
		return "<function>"
	}
	return fmt.Sprintf("%v", v)
}

// FormatFloat formats a float the way Python's repr does for the common
// cases: integral values keep a trailing ".0" and very large or very small
// values use scientific notation. bits must be 32 or 64.
func FormatFloat(f float64, bits int) string {
	switch {
	case math.IsNaN(f):
		return "nan"
	case math.IsInf(f, 1):
		return "inf"
	case math.IsInf(f, -1):
		return "-inf"
	case f == 0:
		if math.Signbit(f) {
			return "-0.0"
		}
		return "0.0"
	}

	abs := math.Abs(f)
	if abs >= 1e16 || abs < 1e-4 {
		return strconv.FormatFloat(f, 'e', -1, bits)
	}
	s := strconv.FormatFloat(f, 'f', -1, bits)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

// QuoteString returns s as a Python string literal, preferring single quotes
// unless that would require escaping many of them.
func QuoteString(s string) string {
	quote := byte('\'')
	if strings.ContainsRune(s, '\'') && !strings.ContainsRune(s, '"') {
		quote = '"'
	}
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte(quote)
	for _, r := range s {
		if r == rune(quote) {
			b.WriteByte('\\')
			b.WriteByte(quote)
			continue
		}
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 || r == 0x7f {
				fmt.Fprintf(&b, `\x%02x`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte(quote)
	return b.String()
}

func bytesRepr(data []byte) string {
	var b strings.Builder
	b.WriteString("b'")
	for _, c := range data {
		switch c {
		case '\'':
			b.WriteString(`\'`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if c >= 0x20 && c < 0x7f {
				b.WriteByte(c)
			} else {
				fmt.Fprintf(&b, `\x%02x`, c)
			}
		}
	}
	b.WriteString("'")
	return b.String()
}

// Indirect follows pointers and interfaces, returning nil for nil.
func Indirect(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	return rv.Interface()
}

// IsNone reports whether v is nil or a nil pointer/interface/function/channel.
func IsNone(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return rv.IsNil()
	}
	return false
}

// runeLen returns the number of runes in s.
func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

// IsSpaceRune reports whether r is whitespace as understood by Python's
// str.strip()/str.split().
func IsSpaceRune(r rune) bool {
	// Python 的 str.isspace() 把 U+001C..U+001F 也算作空白；
	// Go 的 unicode.IsSpace 基于 White_Space，不包含这四个字符。
	if r >= 0x1c && r <= 0x1f {
		return true
	}
	return unicode.IsSpace(r)
}
