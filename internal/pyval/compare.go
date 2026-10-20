package pyval

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// Truthy implements Python truthiness: nil, false, numeric zero, empty
// strings/containers and nil pointers are false, everything else is true.
func Truthy(v interface{}) bool {
	if v == nil {
		return false
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return false
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Bool:
		return rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		return rv.Complex() != 0
	case reflect.String:
		return rv.Len() > 0
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return rv.Len() > 0
	}
	return true
}

// Len returns Python's len(v): runes for strings, elements for slices,
// arrays, maps and channels.
func Len(v interface{}) (int, error) {
	if s, ok := v.(string); ok {
		return utf8.RuneCountInString(s), nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return 0, fmt.Errorf("pyval: object of type NoneType has no len()")
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		return utf8.RuneCountInString(rv.String()), nil
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return rv.Len(), nil
	}
	return 0, fmt.Errorf("pyval: object of type %s has no len()", TypeName(v))
}

// Iter converts strings, slices, arrays and maps into a slice of elements.
// Strings iterate by rune, maps iterate over their keys (sorted so that the
// result is deterministic, which Go maps otherwise are not).
func Iter(v interface{}) ([]interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("pyval: NoneType is not iterable")
	}
	if s, ok := v.(string); ok {
		runes := []rune(s)
		out := make([]interface{}, len(runes))
		for i, r := range runes {
			out[i] = string(r)
		}
		return out, nil
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil, fmt.Errorf("pyval: NoneType is not iterable")
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return Iter(rv.String())
	case reflect.Slice, reflect.Array:
		out := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out[i] = rv.Index(i).Interface()
		}
		return out, nil
	case reflect.Map:
		if rv.Len() == 0 {
			return []interface{}{}, nil
		}
		keys := rv.MapKeys()
		out := make([]interface{}, len(keys))
		for i, k := range keys {
			out[i] = k.Interface()
		}
		SortValues(out)
		return out, nil
	}
	return nil, fmt.Errorf("pyval: object of type %s is not iterable", TypeName(v))
}

// SortValues sorts values in place using Python-style ordering. Values that
// cannot be compared fall back to their repr, which keeps sorting total and
// deterministic instead of panicking.
func SortValues(values []interface{}) {
	sort.SliceStable(values, func(i, j int) bool { return Less(values[i], values[j]) })
}

// Less reports whether a sorts before b.
func Less(a, b interface{}) bool {
	c, err := Compare(a, b)
	if err != nil {
		return Repr(a) < Repr(b)
	}
	return c < 0
}

// Compare compares a and b using Python-style ordering. It returns a
// negative number when a < b, zero when they are equal and a positive number
// when a > b.
func Compare(a, b interface{}) (int, error) {
	a, b = Indirect(a), Indirect(b)
	if a == nil || b == nil {
		if a == nil && b == nil {
			return 0, nil
		}
		return 0, fmt.Errorf("pyval: cannot compare %s with None", TypeName(a))
	}

	if ai, ok := asInt64(a); ok {
		if bi, ok := asInt64(b); ok {
			return cmpInt64(ai, bi), nil
		}
	}
	if au, ok := asUint64(a); ok {
		if bu, ok := asUint64(b); ok {
			switch {
			case au < bu:
				return -1, nil
			case au > bu:
				return 1, nil
			}
			return 0, nil
		}
	}
	if ai, ok := asInt64(a); ok {
		if bf, ok := asFloatOnly(b); ok {
			return cmpIntFloat(ai, bf)
		}
	}
	if bi, ok := asInt64(b); ok {
		if af, ok := asFloatOnly(a); ok {
			c, err := cmpIntFloat(bi, af)
			return -c, err
		}
	}
	if au, ok := asUint64(a); ok {
		if bf, ok := asFloatOnly(b); ok {
			return cmpUintFloat(au, bf)
		}
	}
	if bu, ok := asUint64(b); ok {
		if af, ok := asFloatOnly(a); ok {
			c, err := cmpUintFloat(bu, af)
			return -c, err
		}
	}
	if af, ok := asFloat(a); ok {
		if bf, ok := asFloat(b); ok {
			if math.IsNaN(af) || math.IsNaN(bf) {
				return 0, fmt.Errorf("pyval: cannot compare nan")
			}
			switch {
			case af < bf:
				return -1, nil
			case af > bf:
				return 1, nil
			}
			return 0, nil
		}
	}

	if at, ok := a.(time.Time); ok {
		if bt, ok := b.(time.Time); ok {
			switch {
			case at.Before(bt):
				return -1, nil
			case at.After(bt):
				return 1, nil
			}
			return 0, nil
		}
	}

	ra, rb := reflect.ValueOf(a), reflect.ValueOf(b)
	if ra.Kind() == reflect.String && rb.Kind() == reflect.String {
		return strings.Compare(ra.String(), rb.String()), nil
	}
	if isSequence(ra) && isSequence(rb) {
		if isByteSequence(ra) != isByteSequence(rb) {
			return 0, fmt.Errorf("pyval: cannot compare %s with %s", TypeName(a), TypeName(b))
		}
		n := ra.Len()
		if rb.Len() < n {
			n = rb.Len()
		}
		for i := 0; i < n; i++ {
			c, err := Compare(ra.Index(i).Interface(), rb.Index(i).Interface())
			if err != nil {
				return 0, err
			}
			if c != 0 {
				return c, nil
			}
		}
		return cmpInt64(int64(ra.Len()), int64(rb.Len())), nil
	}
	return 0, fmt.Errorf("pyval: cannot compare %s with %s", TypeName(a), TypeName(b))
}

// Equal implements Python equality for the common cases, including numeric
// cross-type equality (1 == 1.0) and deep comparison of lists and dicts.
func Equal(a, b interface{}) bool {
	a, b = Indirect(a), Indirect(b)
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if isNumericValue(a) && isNumericValue(b) {
		c, err := Compare(a, b)
		return err == nil && c == 0
	}
	if at, ok := a.(time.Time); ok {
		if bt, ok := b.(time.Time); ok {
			return at.Equal(bt)
		}
	}

	ra, rb := reflect.ValueOf(a), reflect.ValueOf(b)
	if ra.Kind() == reflect.String && rb.Kind() == reflect.String {
		return ra.String() == rb.String()
	}
	if isSequence(ra) && isSequence(rb) {
		if isByteSequence(ra) != isByteSequence(rb) {
			return false
		}
		if ra.Len() != rb.Len() {
			return false
		}
		for i := 0; i < ra.Len(); i++ {
			if !Equal(ra.Index(i).Interface(), rb.Index(i).Interface()) {
				return false
			}
		}
		return true
	}
	if ra.Kind() == reflect.Map && rb.Kind() == reflect.Map {
		if ra.Len() != rb.Len() {
			return false
		}
		for _, ka := range ra.MapKeys() {
			found := false
			for _, kb := range rb.MapKeys() {
				if Equal(ka.Interface(), kb.Interface()) &&
					Equal(ra.MapIndex(ka).Interface(), rb.MapIndex(kb).Interface()) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	return reflect.DeepEqual(a, b)
}

func isSequence(rv reflect.Value) bool {
	return rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array
}

func isByteSequence(rv reflect.Value) bool {
	return isSequence(rv) && rv.Type().Elem().Kind() == reflect.Uint8
}

func asInt64(v interface{}) (int64, bool) {
	switch x := v.(type) {
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	case int:
		return int64(x), true
	case int8:
		return int64(x), true
	case int16:
		return int64(x), true
	case int32:
		return int64(x), true
	case int64:
		return x, true
	case uint:
		if uint64(x) <= math.MaxInt64 {
			return int64(x), true
		}
	case uint8:
		return int64(x), true
	case uint16:
		return int64(x), true
	case uint32:
		return int64(x), true
	case uint64:
		if x <= math.MaxInt64 {
			return int64(x), true
		}
	case uintptr:
		if uint64(x) <= math.MaxInt64 {
			return int64(x), true
		}
	}
	return 0, false
}

func asUint64(v interface{}) (uint64, bool) {
	switch x := v.(type) {
	case uint:
		return uint64(x), true
	case uint8:
		return uint64(x), true
	case uint16:
		return uint64(x), true
	case uint32:
		return uint64(x), true
	case uint64:
		return x, true
	case uintptr:
		return uint64(x), true
	}
	return 0, false
}

func asFloatOnly(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float32:
		return float64(x), true
	case float64:
		return x, true
	}
	return 0, false
}

func isNumericValue(v interface{}) bool {
	switch v.(type) {
	case bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64:
		return true
	}
	return false
}

func cmpIntFloat(i int64, f float64) (int, error) {
	if math.IsNaN(f) {
		return 0, fmt.Errorf("pyval: cannot compare nan")
	}
	if math.IsInf(f, 1) {
		return -1, nil
	}
	if math.IsInf(f, -1) {
		return 1, nil
	}
	left := new(big.Rat).SetInt64(i)
	right := new(big.Rat).SetFloat64(f)
	return left.Cmp(right), nil
}

func cmpUintFloat(u uint64, f float64) (int, error) {
	if math.IsNaN(f) {
		return 0, fmt.Errorf("pyval: cannot compare nan")
	}
	if math.IsInf(f, 1) {
		return -1, nil
	}
	if math.IsInf(f, -1) {
		return 1, nil
	}
	integer := new(big.Int).SetUint64(u)
	left := new(big.Rat).SetInt(integer)
	right := new(big.Rat).SetFloat64(f)
	return left.Cmp(right), nil
}

func asFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	case int:
		return float64(x), true
	case int8:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint8:
		return float64(x), true
	case uint16:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	case uintptr:
		return float64(x), true
	case float32:
		return float64(x), true
	case float64:
		return x, true
	}
	return 0, false
}

func cmpInt64(a, b int64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}
