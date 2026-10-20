package pyval

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

// pyvalTestStruct pins down TypeName's rule for named Go structs: the Go type
// name is reported, analogous to Python's type(obj).__name__ for a class.
type pyvalTestStruct struct{ N int }

// pyvalTestString and pyvalTestBytes exercise the reflection fallback in
// TypeName/Repr for named Go types.
type pyvalTestString string

type pyvalTestBytes []byte

func TestTruthy(t *testing.T) {
	one := 1
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		// Python: bool(None) is False.
		{"nil", nil, false},
		{"false", false, false},
		{"true", true, true},
		// Python: bool(0), bool(0.0), bool(-0.0) are False.
		{"int zero", 0, false},
		{"int one", 1, true},
		{"int negative", -1, true},
		{"int8 zero", int8(0), false},
		{"uint zero", uint(0), false},
		{"float zero", 0.0, false},
		{"negative zero float", math.Copysign(0, -1), false},
		{"float nonzero", 0.5, true},
		{"nan is truthy", math.NaN(), true}, // Python: bool(float("nan")) is True.
		{"complex zero", complex(0, 0), false},
		{"complex nonzero", complex(1, 2), true},
		// Python: bool("") is False, any other string is True.
		{"empty string", "", false},
		{"space string", " ", true},
		{"chinese string", "中文", true},
		// Python: list/tuple truthiness is based on length only.
		{"empty slice", []int{}, false},
		{"nil slice", []int(nil), false},
		{"slice with zero", []int{0}, true},
		{"empty array", [0]int{}, false},
		{"array with zero", [1]int{0}, true},
		{"empty map", map[string]int{}, false},
		{"nil map", map[string]int(nil), false},
		{"map with zero value", map[string]int{"a": 0}, true},
		// Python has no pointers; nil pointer models None, a non-nil pointer is
		// a live object reference and therefore truthy.
		{"nil pointer", (*int)(nil), false},
		{"non-nil pointer", &one, true},
		{"empty struct", struct{}{}, true},
		{"time", time.Time{}, true}, // Python: datetime objects are always truthy.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Truthy(tt.value); got != tt.want {
				t.Errorf("Truthy(%#v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestLen(t *testing.T) {
	valid := []struct {
		name  string
		value any
		want  int
	}{
		// Python: len(str) counts code points (runes for valid UTF-8).
		{"ascii", "hello", 5},
		{"empty string", "", 0},
		{"chinese", "中文", 2},
		{"accented", "héllo", 5},
		{"emoji", "🙂a", 2},
		{"combining mark", "e\u0301", 2}, // e + U+0301, two code points.
		// Python: len(b"...") counts bytes, not runes.
		{"bytes", []byte("中文"), 6},
		{"slice", []int{1, 2, 3}, 3},
		{"nil slice", []int(nil), 0},
		{"array", [4]string{"a", "b", "c", "d"}, 4},
		{"map", map[string]int{"a": 1, "b": 2}, 2},
		{"nil map", map[string]int(nil), 0},
	}
	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Len(tt.value)
			if err != nil {
				t.Fatalf("Len(%#v) error = %v, want nil", tt.value, err)
			}
			if got != tt.want {
				t.Fatalf("Len(%#v) = %d, want %d", tt.value, got, tt.want)
			}
		})
	}

	// Python: len(42), len(None), len(True) etc. raise TypeError.
	invalid := []struct {
		name  string
		value any
	}{
		{"nil", nil},
		{"int", 42},
		{"float", 3.14},
		{"bool", true},
		{"complex", complex(1, 2)},
		{"struct", struct{}{}},
		{"time", time.Time{}},
		{"nil pointer", (*int)(nil)},
	}
	for _, tt := range invalid {
		t.Run("error/"+tt.name, func(t *testing.T) {
			_, err := Len(tt.value)
			if err == nil {
				t.Fatalf("Len(%#v) error = nil, want non-nil", tt.value)
			}
			if !strings.Contains(err.Error(), "has no len()") {
				t.Fatalf("Len(%#v) error = %q, want it to mention has no len()", tt.value, err)
			}
		})
	}
}

func TestIter(t *testing.T) {
	valid := []struct {
		name  string
		value any
		want  []any
	}{
		// Python: iter(str) yields one-character strings, by code point.
		{"ascii string", "abc", []any{"a", "b", "c"}},
		{"empty string", "", []any{}},
		{"chinese string", "中文", []any{"中", "文"}},
		{"emoji string", "a🙂", []any{"a", "🙂"}},
		{"slice", []int{3, 1, 2}, []any{3, 1, 2}},
		{"nil slice", []int(nil), []any{}}, // nil slice is an empty list, not None.
		{"array", [2]string{"x", "y"}, []any{"x", "y"}},
		{"empty map", map[string]int{}, []any{}},
		// Python maps iterate in insertion order; Go maps are unordered, so the
		// package rule is to sort keys, which makes the result deterministic.
		{"string keys sorted", map[string]int{"b": 2, "a": 1, "c": 0}, []any{"a", "b", "c"}},
		{"int keys sorted", map[int]string{10: "ten", 2: "two", 1: "one"}, []any{1, 2, 10}},
	}
	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Iter(tt.value)
			if err != nil {
				t.Fatalf("Iter(%#v) error = %v, want nil", tt.value, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Iter(%#v) = %#v, want %#v", tt.value, got, tt.want)
			}
		})
	}

	// Python: iter(None) and iter(42) raise TypeError.
	invalid := []struct {
		name  string
		value any
	}{
		{"nil", nil},
		{"int", 42},
		{"float", 3.14},
		{"bool", true},
		{"struct", struct{}{}},
		{"time", time.Time{}},
		{"nil pointer", (*int)(nil)},
	}
	for _, tt := range invalid {
		t.Run("error/"+tt.name, func(t *testing.T) {
			_, err := Iter(tt.value)
			if err == nil {
				t.Fatalf("Iter(%#v) error = nil, want non-nil", tt.value)
			}
			if !strings.Contains(err.Error(), "is not iterable") {
				t.Fatalf("Iter(%#v) error = %q, want it to mention is not iterable", tt.value, err)
			}
		})
	}
}

func TestTypeName(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"nil", nil, "NoneType"},
		{"nil pointer", (*int)(nil), "NoneType"}, // nil pointer models None.
		{"bool", true, "bool"},
		{"string", "abc", "str"},
		{"named string", pyvalTestString("abc"), "str"},
		{"bytes", []byte("abc"), "bytes"},
		{"named bytes", pyvalTestBytes{1, 2}, "bytes"},
		{"int", int(1), "int"},
		{"int8", int8(1), "int"},
		{"int16", int16(1), "int"},
		{"int32", int32(1), "int"},
		{"int64", int64(1), "int"},
		{"uint", uint(1), "int"},
		{"uint8", uint8(1), "int"},
		{"uint16", uint16(1), "int"},
		{"uint32", uint32(1), "int"},
		{"uint64", uint64(1), "int"},
		{"uintptr", uintptr(1), "int"},
		{"float32", float32(1), "float"},
		{"float64", float64(1), "float"},
		{"complex64", complex64(complex(1, 2)), "complex"},
		{"complex128", complex(1, 2), "complex"},
		{"slice", []int{1, 2}, "list"},
		{"array", [2]int{1, 2}, "list"},
		{"map", map[string]int{"a": 1}, "dict"},
		// Explicit rule: a named struct reports its Go type name, like a Python
		// class's __name__; anonymous structs have no name and report "object".
		{"named struct", pyvalTestStruct{N: 1}, "pyvalTestStruct"},
		{"anonymous struct", struct{ N int }{N: 1}, "object"},
		{"time", time.Time{}, "datetime"},
		{"duration", 1500 * time.Millisecond, "timedelta"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypeName(tt.value); got != tt.want {
				t.Errorf("TypeName(%#v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestStr(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		// Python: str(None) == "None", str(True) == "True".
		{"none", nil, "None"},
		{"true", true, "True"},
		{"false", false, "False"},
		{"int", 42, "42"},
		{"float", 1.5, "1.5"},
		{"float one", 1.0, "1.0"},
		{"float large", 1e16, "1e+16"},
		{"negative zero", math.Copysign(0, -1), "-0.0"},
		// Python: str(str) is the string itself, without quotes or escapes.
		{"string unquoted", "hello", "hello"},
		{"string with newline", "a\nb", "a\nb"},
		{"string with quotes", "it's", "it's"},
		{"list", []any{1, "a", true}, "[1, 'a', True]"},
		// Go maps are unordered; the package sorts keys for determinism.
		{"dict", map[string]int{"b": 2, "a": 1}, "{'a': 1, 'b': 2}"},
		{"bytes", []byte("hi"), "b'hi'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Str(tt.value); got != tt.want {
				t.Errorf("Str(%#v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestRepr(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		// Python: repr(None), repr(True), repr(False).
		{"none", nil, "None"},
		{"true", true, "True"},
		{"false", false, "False"},
		{"int", 42, "42"},
		// Python 3.13 float repr: integral values keep ".0"; 1e16 and 1e-5
		// use exponent notation; -0.0 keeps its sign.
		{"float one", 1.0, "1.0"},
		{"float large", 1e16, "1e+16"},
		{"float small", 1e-5, "1e-05"},
		{"float negative zero", math.Copysign(0, -1), "-0.0"},
		{"float32", float32(0.1), "0.1"},
		{"nan", math.NaN(), "nan"},
		{"inf", math.Inf(1), "inf"},
		{"negative inf", math.Inf(-1), "-inf"},
		// CPython picks double quotes iff the string contains a single quote
		// but no double quote; otherwise it uses single quotes.
		{"string", "hello", "'hello'"},
		{"empty string", "", "''"},
		{"apostrophe", "it's", `"it's"`},
		{"both quotes", `a'b"c`, `'a\'b"c'`},
		{"escapes", "a\nb\tc\\d\r", `'a\nb\tc\\d\r'`},
		{"controls", "a\x00b\x1fc\x7fd", `'a\x00b\x1fc\x7fd'`},
		{"chinese", "中文", "'中文'"},
		{"list", []any{1, "a", true, nil}, "[1, 'a', True, None]"},
		{"nested list", []any{[]int{1, 2}, map[string]int{"b": 2, "a": 1}}, "[[1, 2], {'a': 1, 'b': 2}]"},
		// Go maps are unordered; nested dicts sort their keys so the repr is
		// deterministic. Lists keep their order.
		{
			"nested dict",
			map[string]any{"z": []int{3, 2}, "a": map[string]int{"y": 2, "x": 1}},
			"{'a': {'x': 1, 'y': 2}, 'z': [3, 2]}",
		},
		{"int keys sorted numerically", map[int]string{10: "ten", 2: "two", 1: "one"}, "{1: 'one', 2: 'two', 10: 'ten'}"},
		{"nil map", map[string]int(nil), "{}"},
		{"nil slice", []int(nil), "[]"},
		{"none value", map[string]any{"a": nil}, "{'a': None}"},
		{"bytes", []byte("hi"), "b'hi'"},
		{"bytes escapes", []byte{0x00, 0x7f, 0xff}, `b'\x00\x7f\xff'`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Repr(tt.value); got != tt.want {
				t.Errorf("Repr(%#v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestQuoteString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", "''"},
		{"plain", "hello", "'hello'"},
		{"contains apostrophe", "it's", `"it's"`},
		{"contains double quote", `say "hi"`, `'say "hi"'`},
		{"both quotes", `a'b"c`, `'a\'b"c'`},
		{"only apostrophe", "'", `"'"`},
		{"only double quote", `"`, `'"'`},
		{"newline", "a\nb", `'a\nb'`},
		{"tab", "a\tb", `'a\tb'`},
		{"backslash", `a\b`, `'a\\b'`},
		{"carriage return", "a\rb", `'a\rb'`},
		{"controls", "\x00\x1f\x7f", `'\x00\x1f\x7f'`},
		{"chinese", "中文", "'中文'"},
		{"accented", "é", "'é'"},
		{"emoji", "🙂", "'🙂'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QuoteString(tt.in); got != tt.want {
				t.Errorf("QuoteString(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		bits int
		want string
	}{
		// Expected values match Python 3.13 repr(float).
		{"zero", 0, 64, "0.0"},
		{"negative zero", math.Copysign(0, -1), 64, "-0.0"},
		{"one", 1, 64, "1.0"},
		{"negative one", -1, 64, "-1.0"},
		{"one point five", 1.5, 64, "1.5"},
		{"1e15 stays fixed", 1e15, 64, "1000000000000000.0"},
		{"1e16 uses exponent", 1e16, 64, "1e+16"},
		{"just below 1e16", 9999999999999998.0, 64, "9999999999999998.0"},
		{"1e-5 uses exponent", 1e-5, 64, "1e-05"},
		{"negative 1e-5", -1e-5, 64, "-1e-05"},
		{"1e-4 stays fixed", 1e-4, 64, "0.0001"},
		{"nan", math.NaN(), 64, "nan"},
		{"positive inf", math.Inf(1), 64, "inf"},
		{"negative inf", math.Inf(-1), 64, "-inf"},
		{"1e100", 1e100, 64, "1e+100"},
		{"smallest subnormal", math.SmallestNonzeroFloat64, 64, "5e-324"},
		// Shortest repr at 32-bit precision: float32(0.1) is not 0.1, but
		// "0.1" is the shortest decimal that round-trips to the same float32.
		{"float32 0.1", float64(float32(0.1)), 32, "0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFloat(tt.f, tt.bits); got != tt.want {
				t.Errorf("FormatFloat(%v, %d) = %q, want %q", tt.f, tt.bits, got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	t1 := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	zone := time.FixedZone("UTC+8", 8*60*60)

	tests := []struct {
		name    string
		a, b    any
		want    int
		wantErr bool
	}{
		{"int less", 1, 2, -1, false},
		{"int greater", 2, 1, 1, false},
		{"int equal", 1, 1, 0, false},
		// Python: 1 == 1.0 and 1 < 1.5; numeric types compare by value.
		{"int vs float equal", 1, 1.0, 0, false},
		{"float vs int equal", 1.0, 1, 0, false},
		{"int vs float less", 1, 1.5, -1, false},
		{"float vs int greater", 1.5, 1, 1, false},
		{"negative float vs int", -1.5, -1, -1, false},
		// Python: True == 1, False == 0; bool is an int subclass.
		{"true vs int", true, 1, 0, false},
		{"false vs float", false, 0.0, 0, false},
		{"true vs two", true, 2, -1, false},
		{"false vs negative", false, -1, 1, false},
		{"strings", "a", "b", -1, false},
		{"string prefix", "abc", "ab", 1, false},
		{"string case", "a", "A", 1, false},
		{"chinese strings", "中", "文", -1, false},
		{"empty string", "", "a", -1, false},
		// Python compares sequences lexicographically.
		{"slice less", []int{1, 2}, []int{1, 3}, -1, false},
		{"slice shorter prefix", []int{1}, []int{1, 0}, -1, false},
		{"slice longer prefix", []int{1, 2}, []int{1}, 1, false},
		{"slice equal", []int{1, 2}, []int{1, 2}, 0, false},
		{"array less", [2]int{1, 2}, [2]int{1, 3}, -1, false},
		// Package rule: arrays are treated as Python lists.
		{"array vs slice", [2]int{1, 2}, []int{1, 2}, 0, false},
		{"slice numeric coercion", []any{1}, []any{1.0}, 0, false},
		{"bytes lexicographic", []byte{1, 2}, []byte{1, 3}, -1, false},
		{"bytes prefix", []byte{1}, []byte{1, 0}, -1, false},
		{"time before", t1, t2, -1, false},
		{"time after", t2, t1, 1, false},
		{"time equal", t1, t1, 0, false},
		{"time same instant different zone", t1, t1.In(zone), 0, false},
		{"positive infinity", math.Inf(1), 1e308, 1, false},
		{"negative infinity", math.Inf(-1), -1e308, -1, false},
		// Python compares int and float exactly: 2**53+1 > 2.0**53.
		{"large int vs float ordering", (int64(1) << 53) + 1, float64(int64(1) << 53), 1, false},

		// Python raises TypeError for these pairs.
		{"struct", pyvalTestStruct{N: 1}, pyvalTestStruct{N: 2}, 0, true},
		{"map", map[string]int{"a": 1}, map[string]int{"b": 2}, 0, true},
		{"int vs string", 1, "1", 0, true},
		{"bool vs string", true, "true", 0, true},
		{"slice vs string", []int{1}, "1", 0, true},
		{"bytes vs list", []byte{1}, []int{1}, 0, true},
		{"string vs bytes", "a", []byte("a"), 0, true},
		{"nil vs int", nil, 1, 0, true},
		{"time vs string", t1, t1.String(), 0, true},
		{"complex", complex(1, 2), complex(2, 3), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compare(tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Compare(%#v, %#v) error = nil, want non-nil", tt.a, tt.b)
				}
				if !strings.Contains(err.Error(), "cannot compare") {
					t.Fatalf("Compare(%#v, %#v) error = %q, want it to mention cannot compare", tt.a, tt.b, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Compare(%#v, %#v) error = %v, want nil", tt.a, tt.b, err)
			}
			if got != tt.want {
				t.Errorf("Compare(%#v, %#v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	t1 := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	zone := time.FixedZone("UTC+8", 8*60*60)

	tests := []struct {
		name string
		a, b any
		want bool
	}{
		{"none none", nil, nil, true},
		{"none zero", nil, 0, false},
		{"int", 1, 1, true},
		// Python: 1 == 1.0 is True.
		{"int float", 1, 1.0, true},
		{"float int", 1.0, 1, true},
		{"int float unequal", 1, 1.5, false},
		// Python: True == 1 and False == 0.0.
		{"true int", true, 1, true},
		{"false float", false, 0.0, true},
		{"true two", true, 2, false},
		{"zero negative zero", 0.0, math.Copysign(0, -1), true},
		{"strings", "a", "a", true},
		{"strings unequal", "a", "b", false},
		{"string int", "1", 1, false},
		{"slice", []int{1, 2}, []int{1, 2}, true},
		{"slice unequal", []int{1, 2}, []int{1, 3}, false},
		{"slice length", []int{1}, []int{1, 2}, false},
		// Package rule: a nil slice models the empty list, so [] == [] is True.
		{"nil slice empty slice", []int(nil), []int{}, true},
		{"slice numeric coercion", []int{1}, []float64{1.0}, true},
		{"nested slices", []any{1, "a"}, []any{1.0, "a"}, true},
		// Package rule: arrays are treated as Python lists.
		{"array vs slice", [2]int{1, 2}, []int{1, 2}, true},
		{"bytes", []byte("ab"), []byte("ab"), true},
		// Python: b"ab" != "ab".
		{"bytes vs string", []byte("ab"), "ab", false},
		{"dict", map[string]int{"a": 1, "b": 2}, map[string]int{"b": 2, "a": 1}, true},
		{"dict value", map[string]int{"a": 1}, map[string]int{"a": 2}, false},
		{"dict length", map[string]int{"a": 1}, map[string]int{"a": 1, "b": 2}, false},
		{
			"nested dict",
			map[string]any{"a": []int{1, 2}},
			map[string]any{"a": []int{1, 2}},
			true,
		},
		{"time", t1, t1, true},
		{"time unequal", t1, t2, false},
		{"time zones", t1, t1.In(zone), true},
		{"nil pointer nil", (*int)(nil), nil, true},
		{"nil pointer zero", (*int)(nil), 0, false},
		// Python: float("nan") != float("nan").
		{"nan", math.NaN(), math.NaN(), false},
		{"inf", math.Inf(1), math.Inf(1), true},
		// Python compares int and float exactly. float64(1<<53) is
		// 9007199254740992.0, so it must not equal the next integer up.
		{"large int vs float", (int64(1) << 53) + 1, float64(int64(1) << 53), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.a, tt.b); got != tt.want {
				t.Errorf("Equal(%#v, %#v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLess(t *testing.T) {
	t1 := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	t2 := t1.Add(time.Hour)

	tests := []struct {
		name string
		a, b any
		want bool
	}{
		{"int", 1, 2, true},
		{"int reversed", 2, 1, false},
		{"int equal", 1, 1, false},
		// Python: 1 == 1.0, so 1 < 1.0 is False.
		{"int float equal", 1, 1.0, false},
		{"int float", 1, 1.5, true},
		{"float int", 1.5, 1, false},
		{"true two", true, 2, true},
		{"false zero", false, 0.0, false},
		{"strings", "a", "b", true},
		{"strings reversed", "b", "a", false},
		{"slice", []int{1, 2}, []int{1, 3}, true},
		{"slice reversed", []int{1, 3}, []int{1, 2}, false},
		{"slice prefix", []int{1}, []int{1, 0}, true},
		{"bytes", []byte{1}, []byte{2}, true},
		{"time", t1, t2, true},
		{"time reversed", t2, t1, false},
		// Package rule: for incomparable values Less falls back to Repr. A
		// string repr starts with a quote, which sorts before '-' and digits.
		{"fallback string before int", "a", 1, true},
		{"fallback int after string", 1, "a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Less(tt.a, tt.b); got != tt.want {
				t.Errorf("Less(%#v, %#v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsNone(t *testing.T) {
	var nilPtr *int
	var nilFunc func()
	var nilChan chan int
	var nilAny any
	x := 1

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"nil", nil, true},
		{"nil interface", nilAny, true},
		{"nil pointer", nilPtr, true},
		{"non-nil pointer", &x, false},
		// Documented: nil functions and channels are also None.
		{"nil func", nilFunc, true},
		{"nil chan", nilChan, true},
		{"zero", 0, false},
		{"empty string", "", false},
		// Python: empty list/dict are not None.
		{"nil slice", []int(nil), false},
		{"nil map", map[string]int(nil), false},
		{"struct", struct{}{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNone(tt.value); got != tt.want {
				t.Errorf("IsNone(%#v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIndirect(t *testing.T) {
	x := 7
	p := &x
	pp := &p
	var nilPtr *int
	nestedNil := &nilPtr
	var nilAny any
	anyPtr := any(p)
	pi := any(p)
	ppi := &pi

	tests := []struct {
		name string
		in   any
		want any
	}{
		{"nil", nil, nil},
		{"nil interface", nilAny, nil},
		{"nil pointer", nilPtr, nil},
		{"pointer", p, 7},
		{"double pointer", pp, 7},
		{"pointer to nil pointer", nestedNil, nil},
		{"interface holding pointer", anyPtr, 7},
		{"pointer to interface", ppi, 7},
		{"value", 42, 42},
		{"string", "abc", "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Indirect(tt.in)
			if tt.want == nil {
				if got != nil {
					t.Errorf("Indirect(%#v) = %#v, want nil", tt.in, got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Indirect(%#v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestSortValues(t *testing.T) {
	tests := []struct {
		name string
		in   []any
		want []any
	}{
		{"empty", []any{}, []any{}},
		{"already sorted", []any{1, 2, 3}, []any{1, 2, 3}},
		{"numbers", []any{3, 1, 2}, []any{1, 2, 3}},
		{"mixed numeric", []any{3.5, 1, 2.0, -1}, []any{-1, 1, 2.0, 3.5}},
		{"strings", []any{"b", "a", "c"}, []any{"a", "b", "c"}},
		// Python sorts str by code point; UTF-8 byte order agrees.
		{"unicode strings", []any{"文", "中", "a"}, []any{"a", "中", "文"}},
		// Python 3.13 raises TypeError for str < int. pyval's documented rule
		// is a Repr fallback: a string repr starts with a quote, which sorts
		// before '-' and digits, so strings precede numbers; numbers among
		// themselves still sort numerically, and the sort is stable.
		{"strings and ints", []any{"b", 3, "a", 1}, []any{"a", "b", 1, 3}},
		{"strings and negative ints", []any{-2, "a", 5}, []any{"a", -2, 5}},
		// 1 and 1.0 compare equal; stable sort preserves their input order.
		{"equal numeric values stay stable", []any{1, 1.0, 0}, []any{0, 1, 1.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]any, 0, len(tt.in))
			got = append(got, tt.in...)
			SortValues(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SortValues(%#v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}
