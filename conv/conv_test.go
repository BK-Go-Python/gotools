package conv

import (
	"encoding/json"
	"math"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// floatEqual reports whether a and b are the same IEEE-754 value, treating
// NaN as equal to NaN so tables can assert Python's float("nan").
func floatEqual(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

func equalFloatSlice(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !floatEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestInt(t *testing.T) {
	n := 8
	tests := []struct {
		name    string
		in      interface{}
		want    int
		wantErr bool
	}{
		// Python: int("42") == 42; int() strips surrounding whitespace.
		{"decimal string", "42", 42, false},
		{"whitespace and sign", "  -42\t", -42, false},
		{"leading plus", "+7", 7, false},
		{"leading zeros", "007", 7, false},
		{"int64 max string", "9223372036854775807", math.MaxInt64, false},
		{"int64 min string", "-9223372036854775808", math.MinInt64, false},

		// Python: int(3.9) == 3, int(-3.9) == -3 (truncation toward zero).
		{"float truncates", 3.9, 3, false},
		{"negative float truncates toward zero", -3.9, -3, false},
		{"negative fraction truncates to zero", -0.9, 0, false},
		{"float32 truncates", float32(2.75), 2, false},

		// Python: int(True) == 1, int(False) == 0.
		{"true", true, 1, false},
		{"false", false, 0, false},

		{"int", 11, 11, false},
		{"int64", int64(-12), -12, false},
		{"uint", uint(13), 13, false},
		{"uint64 fits", uint64(math.MaxInt64), math.MaxInt64, false},
		{"pointer to int", &n, 8, false},
		// Python: int(b"123") == 123.
		{"bytes", []byte("123"), 123, false},
		// Go rule: json.Number is decoded according to its literal.
		{"json number", json.Number("456"), 456, false},

		// Python: int("") / int("abc") / int("1.5") raise ValueError.
		{"empty string", "", 0, true},
		{"blank string", " \t\n", 0, true},
		{"invalid literal", "abc", 0, true},
		{"float literal string", "1.5", 0, true},
		{"exponent literal string", "1e3", 0, true},
		{"hex literal string", "0x10", 0, true},
		{"trailing garbage", "12abc", 0, true},
		{"int64 overflow string", "9223372036854775808", 0, true},
		{"int64 underflow string", "-9223372036854775809", 0, true},
		{"nil", nil, 0, true},
		{"slice", []int{1}, 0, true},
		{"uint64 overflow", uint64(math.MaxUint64), 0, true},
		{"json number with fraction", json.Number("1.5"), 0, true},
		// Python: int(nan) / int(inf) raise ValueError/OverflowError.
		{"nan", math.NaN(), 0, true},
		{"positive infinity", math.Inf(1), 0, true},
		{"negative infinity", math.Inf(-1), 0, true},
		{"float overflow", math.Pow(2, 63), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Int(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Int(%#v) = %d, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Int(%#v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("Int(%#v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    int64
		wantErr bool
	}{
		{"max", "9223372036854775807", math.MaxInt64, false},
		{"min", "-9223372036854775808", math.MinInt64, false},
		{"negative decimal", "-5", -5, false},
		{"float truncates", 3.9, 3, false},
		{"negative float truncates", -2.9, -2, false},
		{"float32 truncates", float32(-7.9), -7, false},
		// Python: int(-2.0**63) == -9223372036854775808.
		{"min float", math.Pow(-2, 63), math.MinInt64, false},
		{"bool true", true, 1, false},
		{"bool false", false, 0, false},
		{"bytes", []byte("-12"), -12, false},
		{"json number", json.Number("99"), 99, false},

		{"string overflow", "9223372036854775808", 0, true},
		{"string underflow", "-9223372036854775809", 0, true},
		{"empty", "", 0, true},
		{"blank", "   ", 0, true},
		{"invalid", "12x", 0, true},
		{"float string", "1.0", 0, true},
		{"nil", nil, 0, true},
		{"json number fraction", json.Number("99.0"), 0, true},
		{"uint64 overflow", uint64(math.MaxUint64), 0, true},
		{"uint64 just over max", uint64(math.MaxInt64) + 1, 0, true},
		{"nan", math.NaN(), 0, true},
		{"inf", math.Inf(1), 0, true},
		{"float overflow", math.Pow(2, 63), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Int64(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Int64(%#v) = %d, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Int64(%#v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("Int64(%#v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    uint64
		wantErr bool
	}{
		{"zero", "0", 0, false},
		{"max", "18446744073709551615", math.MaxUint64, false},
		{"whitespace", " 42 ", 42, false},
		{"direct max", uint64(math.MaxUint64), math.MaxUint64, false},
		{"int zero", int64(0), 0, false},
		{"bool true", true, 1, false},
		{"bool false", false, 0, false},
		{"float truncates", 3.9, 3, false},
		// Python: int(-0.0) == 0; uint64 rejects only values < 0.
		{"negative zero", math.Copysign(0, -1), 0, false},
		{"bytes", []byte("7"), 7, false},
		{"json number string kind", json.Number("42"), 42, false},

		{"string overflow", "18446744073709551616", 0, true},
		{"negative string", "-1", 0, true},
		{"negative int", -1, 0, true},
		{"negative int64", int64(-42), 0, true},
		{"negative float", -0.5, 0, true},
		{"empty", "", 0, true},
		{"invalid", "abc", 0, true},
		{"nil", nil, 0, true},
		{"slice", []int{}, 0, true},
		{"nan", math.NaN(), 0, true},
		{"inf", math.Inf(1), 0, true},
		{"float overflow", math.Pow(2, 64), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Uint64(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Uint64(%#v) = %d, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Uint64(%#v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("Uint64(%#v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestFloat(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    float64
		wantErr bool
	}{
		{"decimal", "3.5", 3.5, false},
		{"whitespace and sign", "  -2.25 ", -2.25, false},
		{"exponent", "1e3", 1000, false},
		{"leading dot", ".5", 0.5, false},
		{"trailing dot", "5.", 5, false},
		// Python: float("inf"), float("Infinity"), float("-inf").
		{"inf literal", "inf", math.Inf(1), false},
		{"infinity literal", "Infinity", math.Inf(1), false},
		{"negative inf literal", "-inf", math.Inf(-1), false},
		// Python: float("nan") is nan.
		{"nan literal", "nan", math.NaN(), false},
		{"NaN literal", "NaN", math.NaN(), false},
		{"inf value", math.Inf(-1), math.Inf(-1), false},
		{"nan value", math.NaN(), math.NaN(), false},
		{"bool true", true, 1, false},
		{"bool false", false, 0, false},
		{"int", -3, -3, false},
		{"uint64", uint64(math.MaxUint64), float64(math.MaxUint64), false},
		{"bytes", []byte("2.5"), 2.5, false},
		{"json number", json.Number("2.5"), 2.5, false},

		{"empty", "", 0, true},
		{"blank", "  \t", 0, true},
		{"invalid", "abc", 0, true},
		{"comma decimal", "1,5", 0, true},
		{"double dot", "1.2.3", 0, true},
		{"hex literal without exponent", "0x10", 0, true},
		{"nil", nil, 0, true},
		{"slice", []int{1}, 0, true},
		{"complex", complex(1, 2), 0, true},
		{"bad json number", json.Number("bad"), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Float(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Float(%#v) = %v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Float(%#v) unexpected error: %v", tt.in, err)
			}
			if !floatEqual(got, tt.want) {
				t.Fatalf("Float(%#v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	yes := true
	tests := []struct {
		name    string
		in      interface{}
		want    bool
		wantErr bool
	}{
		{"bool true", true, true, false},
		{"bool false", false, false, false},
		{"true string", "true", true, false},
		{"TRUE string", "TRUE", true, false},
		{"false string with spaces", " false ", false, false},
		{"yes", "yes", true, false},
		{"YES", "YES", true, false},
		{"no", "no", false, false},
		{"on", "on", true, false},
		{"off", "off", false, false},
		{"one", "1", true, false},
		{"zero", "0", false, false},
		{"y", "y", true, false},
		{"n", "n", false, false},
		{"t", "t", true, false},
		{"f", "f", false, false},
		// Python: bool("") is False.
		{"empty string", "", false, false},
		{"bytes yes", []byte("yes"), true, false},
		{"bytes empty", []byte(""), false, false},
		// Python: bool(None) is False.
		{"nil", nil, false, false},
		{"nil pointer", (*bool)(nil), false, false},
		{"pointer to true", &yes, true, false},
		{"int zero", 0, false, false},
		{"int one", 1, true, false},
		{"negative int", -1, true, false},
		{"float zero", 0.0, false, false},
		{"negative zero", math.Copysign(0, -1), false, false},
		{"non-zero float", 3.14, true, false},
		// Python: bool(float("nan")) is True.
		{"nan", math.NaN(), true, false},
		{"inf", math.Inf(-1), true, false},
		{"uint64 max", uint64(math.MaxUint64), true, false},
		// Python truthiness for containers: empty is False, non-empty is True.
		{"empty slice", []int{}, false, false},
		{"nil slice", []int(nil), false, false},
		{"non-empty slice", []int{0}, true, false},
		{"empty any slice", []interface{}{}, false, false},
		{"any slice with falsey items", []interface{}{nil, 0, ""}, true, false},
		{"empty map", map[string]int{}, false, false},
		{"non-empty map", map[string]int{"a": 0}, true, false},
		{"empty struct", struct{}{}, true, false},

		{"unknown string", "maybe", false, true},
		{"binary-looking string", "2", false, true},
		{"hex string", "0x1", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Bool(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Bool(%#v) = %v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Bool(%#v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("Bool(%#v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	ptr := 5
	tests := []struct {
		name string
		in   interface{}
		want string
	}{
		// Python: str(None) == "None".
		{"nil", nil, "None"},
		{"string unchanged", "hello", "hello"},
		{"single quote in string", "it's", "it's"},
		{"chinese", "中文", "中文"},
		// Python: str(True) == "True".
		{"bool true", true, "True"},
		{"bool false", false, "False"},
		{"int", 42, "42"},
		{"negative int64", int64(-7), "-7"},
		{"uint64 max", uint64(math.MaxUint64), "18446744073709551615"},
		// Python: str(1.0) == "1.0"; repr floats keep a decimal point.
		{"float integral", 1.0, "1.0"},
		{"float", 1.5, "1.5"},
		{"negative zero", math.Copysign(0, -1), "-0.0"},
		{"nan", math.NaN(), "nan"},
		{"inf", math.Inf(1), "inf"},
		{"negative inf", math.Inf(-1), "-inf"},
		// Python: str(1e16) == "1e+16", str(1e-5) == "1e-05".
		{"large float", 1e16, "1e+16"},
		{"small float", 1e-5, "1e-05"},
		{"large fixed float", 1e15, "1000000000000000.0"},
		// Python: str(b"abc") == "b'abc'".
		{"bytes", []byte("abc"), "b'abc'"},
		// Python repr of list/dict.
		{"string slice", []string{"a", "b"}, "['a', 'b']"},
		{"mixed list", []interface{}{1, "a", true, nil}, "[1, 'a', True, None]"},
		{"nested containers", []interface{}{[]int{1, 2}, map[string]int{"a": 1}}, "[[1, 2], {'a': 1}]"},
		{"dict", map[string]int{"a": 1, "b": 2}, "{'a': 1, 'b': 2}"},
		{"pointer", &ptr, "5"},
		{"nil pointer", (*int)(nil), "None"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.in); got != tt.want {
				t.Fatalf("String(%#v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestDefaults(t *testing.T) {
	t.Run("IntDefault", func(t *testing.T) {
		tests := []struct {
			name string
			in   interface{}
			def  int
			want int
		}{
			{"valid string", "42", -1, 42},
			{"valid float", 3.9, 0, 3},
			{"invalid string", "bad", -1, -1},
			{"nil", nil, 7, 7},
			{"overflow", "9223372036854775808", 5, 5},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := IntDefault(tt.in, tt.def); got != tt.want {
					t.Fatalf("IntDefault(%#v, %d) = %d, want %d", tt.in, tt.def, got, tt.want)
				}
			})
		}
	})

	t.Run("FloatDefault", func(t *testing.T) {
		tests := []struct {
			name string
			in   interface{}
			def  float64
			want float64
		}{
			{"valid string", "2.5", -1, 2.5},
			{"valid bool", true, 0, 1},
			{"invalid string", "bad", -2.5, -2.5},
			{"nil", nil, 7.5, 7.5},
			{"nan succeeds", "nan", 0, math.NaN()},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := FloatDefault(tt.in, tt.def); !floatEqual(got, tt.want) {
					t.Fatalf("FloatDefault(%#v, %v) = %v, want %v", tt.in, tt.def, got, tt.want)
				}
			})
		}
	})

	t.Run("BoolDefault", func(t *testing.T) {
		tests := []struct {
			name string
			in   interface{}
			def  bool
			want bool
		}{
			{"yes", "yes", false, true},
			{"off", "off", true, false},
			{"invalid string", "bogus", true, true},
			{"nil uses false not default", nil, true, false},
			{"zero int", 0, true, false},
			{"non-zero int", 3, false, true},
			{"empty container", []int{}, true, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := BoolDefault(tt.in, tt.def); got != tt.want {
					t.Fatalf("BoolDefault(%#v, %v) = %v, want %v", tt.in, tt.def, got, tt.want)
				}
			})
		}
	})
}

func TestStrings(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    []string
		wantErr bool
	}{
		{"string slice", []string{"a", "b"}, []string{"a", "b"}, false},
		{"int slice", []int{1, 2}, []string{"1", "2"}, false},
		{"mixed any slice", []interface{}{"a", 1, true, nil}, []string{"a", "1", "True", "None"}, false},
		// Python: [str(c) for c in "abc"] == ["a", "b", "c"].
		{"string iterates by character", "abc", []string{"a", "b", "c"}, false},
		{"unicode string iterates by rune", "中a文", []string{"中", "a", "文"}, false},
		// Python: list(b"AB") == [65, 66], then str each.
		{"bytes iterate as ints", []byte("AB"), []string{"65", "66"}, false},
		// Go rule: map iteration is sorted and yields keys.
		{"map keys", map[string]int{"b": 2, "a": 1}, []string{"a", "b"}, false},
		{"int map keys", map[int]string{2: "b", 1: "a"}, []string{"1", "2"}, false},
		{"empty slice", []interface{}{}, []string{}, false},
		{"nil slice", []string(nil), []string{}, false},
		{"nil element", []interface{}{nil}, []string{"None"}, false},

		// Python: list(None) / list(42) raise TypeError.
		{"nil", nil, nil, true},
		{"not iterable", 42, nil, true},
		{"struct not iterable", struct{}{}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Strings(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Strings(%#v) = %#v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Strings(%#v) unexpected error: %v", tt.in, err)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Strings(%#v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestInts(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    []int
		wantErr bool
	}{
		{"string slice", []string{"1", "-2", "+3"}, []int{1, -2, 3}, false},
		{"int slice", []int{4, 5}, []int{4, 5}, false},
		{"mixed any slice", []interface{}{1, 2.9, -2.9, true, false, "6", []byte("7")}, []int{1, 2, -2, 1, 0, 6, 7}, false},
		{"float slice truncates", []float32{3.99, -3.99}, []int{3, -3}, false},
		// Python: [int(c) for c in "123"] == [1, 2, 3].
		{"string iterates by character", "123", []int{1, 2, 3}, false},
		{"json number", []interface{}{json.Number("8")}, []int{8}, false},
		{"empty", []int{}, []int{}, false},
		{"nil slice", []int(nil), []int{}, false},

		{"invalid element", []interface{}{1, "x"}, nil, true},
		{"nil element", []interface{}{nil}, nil, true},
		{"uint64 overflow element", []interface{}{uint64(math.MaxUint64)}, nil, true},
		{"nan element", []interface{}{1.0, math.NaN()}, nil, true},
		{"map keys are not ints", map[string]int{"a": 1}, nil, true},
		{"nil", nil, nil, true},
		{"not iterable", 42, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ints(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Ints(%#v) = %v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Ints(%#v) unexpected error: %v", tt.in, err)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Ints(%#v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}

	// Python reports the offending element; expose its index in the error.
	if _, err := Ints([]interface{}{1, "bad"}); err == nil || !strings.Contains(err.Error(), "element 1") {
		t.Fatalf("Ints element error = %v, want mention of element 1", err)
	}
}

func TestFloats(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    []float64
		wantErr bool
	}{
		{"string slice", []string{"1.5", "-2"}, []float64{1.5, -2}, false},
		{"int slice", []int{1, 2}, []float64{1, 2}, false},
		{"mixed any slice", []interface{}{1.5, 2, true, "3.25", []byte("4.5")}, []float64{1.5, 2, 1, 3.25, 4.5}, false},
		{"inf strings", []string{"inf", "-inf"}, []float64{math.Inf(1), math.Inf(-1)}, false},
		{"nan string", []string{"nan"}, []float64{math.NaN()}, false},
		// Python: [float(c) for c in "12"] == [1.0, 2.0].
		{"string iterates by character", "12", []float64{1, 2}, false},
		{"empty", []int{}, []float64{}, false},

		{"invalid element", []interface{}{"bad"}, nil, true},
		{"invalid second element", []interface{}{1, "bad"}, nil, true},
		{"nil element", []interface{}{nil}, nil, true},
		{"nil", nil, nil, true},
		{"not iterable", 42, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Floats(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Floats(%#v) = %v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Floats(%#v) unexpected error: %v", tt.in, err)
			}
			if !equalFloatSlice(got, tt.want) {
				t.Fatalf("Floats(%#v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}

	if _, err := Floats([]interface{}{1, "bad"}); err == nil || !strings.Contains(err.Error(), "element 1") {
		t.Fatalf("Floats element error = %v, want mention of element 1", err)
	}
}

func TestInterfaceSlice(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    []interface{}
		wantErr bool
	}{
		{"string slice", []string{"a", "b"}, []interface{}{"a", "b"}, false},
		{"int slice", []int{1, 2}, []interface{}{1, 2}, false},
		{"any slice", []interface{}{1, "x", nil}, []interface{}{1, "x", nil}, false},
		{"string iterates by rune", "中a", []interface{}{"中", "a"}, false},
		{"map keys sorted", map[string]int{"b": 2, "a": 1}, []interface{}{"a", "b"}, false},
		{"bytes iterate as ints", []byte("AB"), []interface{}{byte('A'), byte('B')}, false},
		{"empty", []interface{}{}, []interface{}{}, false},

		{"nil", nil, nil, true},
		{"not iterable", 42, nil, true},
		{"struct not iterable", struct{}{}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InterfaceSlice(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("InterfaceSlice(%#v) = %#v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("InterfaceSlice(%#v) unexpected error: %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("InterfaceSlice(%#v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		want []interface{}
	}{
		{"int slice", []int{1, 2}, []interface{}{1, 2}},
		{"string iterates by rune", "ab", []interface{}{"a", "b"}},
		{"map keys sorted", map[string]int{"b": 2, "a": 1}, []interface{}{"a", "b"}},
		{"bytes iterate as ints", []byte("AB"), []interface{}{byte('A'), byte('B')}},
		{"empty", []int{}, []interface{}{}},
		// ToSlice has no error return: non-iterable values yield nil.
		{"not iterable", 42, nil},
		{"nil", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSlice(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToSlice(%#v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	tests := []struct {
		name    string
		in      interface{}
		want    string
		wantErr bool
	}{
		// Python: json.dumps("中文", ensure_ascii=False) == '"中文"'.
		{"chinese not escaped", "中文", `"中文"`, false},
		// Python: ensure_ascii=False also leaves HTML characters alone.
		{"html not escaped", `<a href="x">&`, `"<a href=\"x\">&"`, false},
		{"newline escaped", "line\nbreak", `"line\nbreak"`, false},
		// Python: json.dumps(None) == "null".
		{"nil", nil, "null", false},
		{"bool", true, "true", false},
		{"float", 3.5, "3.5", false},
		// Go/JSON: object keys are sorted (Python preserves insertion order).
		{"map sorted keys", map[string]interface{}{"b": 2, "a": "中文"}, `{"a":"中文","b":2}`, false},
		{"slice", []interface{}{1, "x", nil, true}, `[1,"x",null,true]`, false},
		{"struct tags", user{Name: "汤姆", Age: 3}, `{"name":"汤姆","age":3}`, false},
		{"nested", map[string]interface{}{"a": []int{1, 2}}, `{"a":[1,2]}`, false},

		// Go/JSON: NaN and Inf are not valid JSON numbers.
		{"nan", math.NaN(), "", true},
		{"inf", math.Inf(1), "", true},
		{"func", func() {}, "", true},
		{"channel", make(chan int), "", true},
		{"complex", complex(1, 2), "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJSON(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ToJSON(%#v) = %q, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToJSON(%#v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("ToJSON(%#v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestToJSONIndent(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		indent  string
		want    string
		wantErr bool
	}{
		{
			"two spaces",
			map[string]interface{}{"b": []int{2, 3}, "a": 1},
			"  ",
			"{\n  \"a\": 1,\n  \"b\": [\n    2,\n    3\n  ]\n}",
			false,
		},
		{
			"tab",
			[]interface{}{1, []int{2, 3}},
			"\t",
			"[\n\t1,\n\t[\n\t\t2,\n\t\t3\n\t]\n]",
			false,
		},
		{
			"chinese not escaped",
			map[string]string{"msg": "中文"},
			"  ",
			"{\n  \"msg\": \"中文\"\n}",
			false,
		},
		// Empty indent means compact output.
		{"empty indent", map[string]int{"a": 1}, "", `{"a":1}`, false},
		{"scalar has no indentation", "中文", "  ", `"中文"`, false},
		{"nil", nil, "  ", "null", false},
		{"nan", math.NaN(), "  ", "", true},
		{"func", func() {}, "  ", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJSONIndent(tt.in, tt.indent)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ToJSONIndent(%#v, %q) = %q, want error", tt.in, tt.indent, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToJSONIndent(%#v, %q) unexpected error: %v", tt.in, tt.indent, err)
			}
			if got != tt.want {
				t.Fatalf("ToJSONIndent(%#v, %q) = %q, want %q", tt.in, tt.indent, got, tt.want)
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("string input into struct", func(t *testing.T) {
		var got user
		if err := FromJSON(`{"name":"中文","age":42}`, &got); err != nil {
			t.Fatalf("FromJSON error: %v", err)
		}
		if got.Name != "中文" || got.Age != 42 {
			t.Fatalf("FromJSON = %#v, want {Name:中文 Age:42}", got)
		}
	})

	t.Run("bytes input into map", func(t *testing.T) {
		var got map[string]interface{}
		if err := FromJSON([]byte(`{"a":1,"b":"中文"}`), &got); err != nil {
			t.Fatalf("FromJSON error: %v", err)
		}
		// Go/JSON: decoding into interface{} yields float64 for numbers.
		if n, ok := got["a"].(float64); !ok || n != 1 {
			t.Fatalf("got[a] = %#v, want float64(1)", got["a"])
		}
		if got["b"] != "中文" {
			t.Fatalf("got[b] = %#v, want 中文", got["b"])
		}
	})

	t.Run("array", func(t *testing.T) {
		var got []int
		if err := FromJSON(`[1,2,3]`, &got); err != nil {
			t.Fatalf("FromJSON error: %v", err)
		}
		if !slices.Equal(got, []int{1, 2, 3}) {
			t.Fatalf("FromJSON = %v, want [1 2 3]", got)
		}
	})

	t.Run("escaped unicode", func(t *testing.T) {
		var got map[string]string
		if err := FromJSON(`{"msg":"\u4e2d\u6587"}`, &got); err != nil {
			t.Fatalf("FromJSON error: %v", err)
		}
		if got["msg"] != "中文" {
			t.Fatalf("msg = %q, want 中文", got["msg"])
		}
	})

	t.Run("null leaves struct zero", func(t *testing.T) {
		got := user{Name: "keep", Age: 9}
		if err := FromJSON("null", &got); err != nil {
			t.Fatalf("FromJSON error: %v", err)
		}
		if got.Name != "keep" || got.Age != 9 {
			t.Fatalf("FromJSON(null) = %#v, want unchanged", got)
		}
	})

	t.Run("null sets pointer nil", func(t *testing.T) {
		p := &user{Name: "x"}
		if err := FromJSON("null", &p); err != nil {
			t.Fatalf("FromJSON error: %v", err)
		}
		if p != nil {
			t.Fatalf("FromJSON(null) = %#v, want nil", p)
		}
	})

	t.Run("missing and unknown fields", func(t *testing.T) {
		var got user
		if err := FromJSON(`{"name":"Tom","extra":true}`, &got); err != nil {
			t.Fatalf("FromJSON error: %v", err)
		}
		if got.Name != "Tom" || got.Age != 0 {
			t.Fatalf("FromJSON = %#v, want {Name:Tom Age:0}", got)
		}
	})

	invalid := []struct {
		name string
		data interface{}
		out  interface{}
	}{
		{"empty", "", &user{}},
		{"truncated object", "{", &user{}},
		{"trailing garbage", `{"name":"x"} garbage`, &user{}},
		{"array into struct", `[]`, &user{}},
		{"wrong field type", `{"age":"not-an-int"}`, &user{}},
		{"nil data", nil, &user{}},
		{"wrong data type", 42, &user{}},
		{"nil out", `{"name":"x"}`, nil},
		{"non-pointer out", `{"name":"x"}`, user{}},
	}
	for _, tt := range invalid {
		t.Run("invalid "+tt.name, func(t *testing.T) {
			if err := FromJSON(tt.data, tt.out); err == nil {
				t.Fatalf("FromJSON(%#v, %#v) = nil error, want error", tt.data, tt.out)
			}
		})
	}
}

type testAddress struct {
	City string `json:"city"`
	Zip  string `json:"zip,omitempty"`
}

type TestBase struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type testProfile struct {
	TestBase
	Name    string         `json:"name"`
	Age     int            `json:"age"`
	Score   float64        `json:"score"`
	Empty   string         `json:"empty,omitempty"`
	Count   int            `json:"count,omitempty"`
	Hidden  string         `json:"-"`
	Address *testAddress   `json:"address"`
	Tags    []string       `json:"tags"`
	Meta    map[string]int `json:"meta"`
	Nick    *string        `json:"nick"`
	NoTag   string
	hidden  string
}

type testOptional struct {
	Ptr  *int `json:"ptr"`
	Omit *int `json:"omit,omitempty"`
}

type testEmbeddedPtr struct {
	*TestBase
	Label string `json:"label"`
}

type testNested struct {
	Addr testAddress   `json:"addr"`
	Team []testAddress `json:"team"`
}

func TestStructToMap(t *testing.T) {
	createdAt := time.Date(2024, time.March, 5, 6, 7, 8, 0, time.UTC)
	nick := "neo"
	profile := testProfile{
		TestBase: TestBase{ID: 7, CreatedAt: createdAt},
		Name:     "汤姆",
		Age:      18,
		Score:    95.5,
		Empty:    "",
		Count:    0,
		Hidden:   "should be ignored",
		Address:  &testAddress{City: "北京"},
		Tags:     []string{"go", "py"},
		Meta:     map[string]int{"a": 1},
		Nick:     &nick,
		NoTag:    "keep",
		hidden:   "unexported",
	}

	want := map[string]interface{}{
		"id":         7,
		"created_at": createdAt,
		"name":       "汤姆",
		"age":        18,
		"score":      95.5,
		"address":    map[string]interface{}{"city": "北京"},
		"tags":       []interface{}{"go", "py"},
		"meta":       map[string]interface{}{"a": 1},
		"nick":       "neo",
		"NoTag":      "keep",
	}

	got, err := StructToMap(profile)
	if err != nil {
		t.Fatalf("StructToMap error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("StructToMap = %#v, want %#v", got, want)
	}

	// Numbers keep their original Go types instead of becoming float64.
	if _, ok := got["age"].(int); !ok {
		t.Fatalf("age has type %T, want int", got["age"])
	}
	if _, ok := got["score"].(float64); !ok {
		t.Fatalf("score has type %T, want float64", got["score"])
	}
	if _, ok := got["created_at"].(time.Time); !ok {
		t.Fatalf("created_at has type %T, want time.Time", got["created_at"])
	}

	t.Run("pointer input", func(t *testing.T) {
		got, err := StructToMap(&profile)
		if err != nil {
			t.Fatalf("StructToMap error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("StructToMap(&profile) = %#v, want %#v", got, want)
		}
	})

	t.Run("nil pointers", func(t *testing.T) {
		got, err := StructToMap(testOptional{})
		if err != nil {
			t.Fatalf("StructToMap error: %v", err)
		}
		want := map[string]interface{}{"ptr": nil}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("StructToMap(testOptional{}) = %#v, want %#v", got, want)
		}
	})

	t.Run("anonymous pointer expansion", func(t *testing.T) {
		got, err := StructToMap(testEmbeddedPtr{
			TestBase: &TestBase{ID: 9, CreatedAt: createdAt},
			Label:    "x",
		})
		if err != nil {
			t.Fatalf("StructToMap error: %v", err)
		}
		want := map[string]interface{}{
			"id":         9,
			"created_at": createdAt,
			"label":      "x",
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("StructToMap = %#v, want %#v", got, want)
		}
	})

	t.Run("nested struct and slice", func(t *testing.T) {
		got, err := StructToMap(testNested{
			Addr: testAddress{City: "上海"},
			Team: []testAddress{{City: "深圳", Zip: "518000"}},
		})
		if err != nil {
			t.Fatalf("StructToMap error: %v", err)
		}
		want := map[string]interface{}{
			"addr": map[string]interface{}{"city": "上海"},
			"team": []interface{}{
				map[string]interface{}{"city": "深圳", "zip": "518000"},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("StructToMap = %#v, want %#v", got, want)
		}
	})

	// encoding/json marshals exported fields promoted from an anonymous field
	// of an unexported struct type (see encoding/json field selection rules).
	// StructToMap documents the same rules, but structToMap currently skips the
	// whole embedded field when StructField.PkgPath != "". This case is kept
	// failing on purpose: it is an implementation bug, not a wrong expectation.
	t.Run("unexported embedded struct", func(t *testing.T) {
		type inner struct {
			ID int `json:"id"`
		}
		type outer struct {
			inner
			Name string `json:"name"`
		}
		got, err := StructToMap(outer{inner: inner{ID: 1}, Name: "x"})
		if err != nil {
			t.Fatalf("StructToMap error: %v", err)
		}
		want := map[string]interface{}{"id": 1, "name": "x"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("StructToMap(unexported embed) = %#v, want %#v", got, want)
		}
	})

	t.Run("errors", func(t *testing.T) {
		n := 1
		inputs := []struct {
			name string
			in   interface{}
		}{
			{"nil", nil},
			{"int", 42},
			{"string", "x"},
			{"slice", []int{1}},
			{"nil struct pointer", (*testAddress)(nil)},
			{"nil int pointer", (*int)(nil)},
			{"pointer to non-struct", &n},
			{"func", func() {}},
		}
		for _, tt := range inputs {
			t.Run(tt.name, func(t *testing.T) {
				if got, err := StructToMap(tt.in); err == nil {
					t.Fatalf("StructToMap(%#v) = %#v, want error", tt.in, got)
				}
			})
		}
	})
}

func TestFromStruct(t *testing.T) {
	createdAt := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	profile := testProfile{
		TestBase: TestBase{ID: 1, CreatedAt: createdAt},
		Name:     "x",
		Age:      2,
		Address:  &testAddress{City: "广州"},
	}
	want, err := StructToMap(profile)
	if err != nil {
		t.Fatalf("StructToMap error: %v", err)
	}
	got, err := FromStruct(profile)
	if err != nil {
		t.Fatalf("FromStruct error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FromStruct = %#v, want StructToMap result %#v", got, want)
	}

	if _, err := FromStruct(nil); err == nil {
		t.Fatal("FromStruct(nil) = nil error, want error")
	}
}

func TestMapToStruct(t *testing.T) {
	createdAt := time.Date(2024, time.March, 5, 6, 7, 8, 0, time.UTC)
	nick := "neo"
	original := testProfile{
		TestBase: TestBase{ID: 7, CreatedAt: createdAt},
		Name:     "汤姆",
		Age:      18,
		Score:    95.5,
		Hidden:   "json:- ignored",
		Address:  &testAddress{City: "北京"},
		Tags:     []string{"go", "py"},
		Meta:     map[string]int{"a": 1},
		Nick:     &nick,
		NoTag:    "keep",
		hidden:   "unexported",
	}

	m, err := StructToMap(original)
	if err != nil {
		t.Fatalf("StructToMap error: %v", err)
	}
	var round testProfile
	if err := MapToStruct(m, &round); err != nil {
		t.Fatalf("MapToStruct error: %v", err)
	}
	if round.ID != original.ID || !round.CreatedAt.Equal(original.CreatedAt) {
		t.Fatalf("roundtrip embedded base = %#v, want %#v", round.TestBase, original.TestBase)
	}
	if round.Name != original.Name || round.Age != original.Age || round.Score != original.Score {
		t.Fatalf("roundtrip scalar fields = %#v, want %#v", round, original)
	}
	if !reflect.DeepEqual(round.Address, original.Address) {
		t.Fatalf("roundtrip address = %#v, want %#v", round.Address, original.Address)
	}
	if !slices.Equal(round.Tags, original.Tags) {
		t.Fatalf("roundtrip tags = %v, want %v", round.Tags, original.Tags)
	}
	if !reflect.DeepEqual(round.Meta, original.Meta) {
		t.Fatalf("roundtrip meta = %v, want %v", round.Meta, original.Meta)
	}
	if round.Nick == nil || *round.Nick != *original.Nick {
		t.Fatalf("roundtrip nick = %v, want %v", round.Nick, original.Nick)
	}
	if round.NoTag != original.NoTag {
		t.Fatalf("roundtrip NoTag = %q, want %q", round.NoTag, original.NoTag)
	}
	// json:"-" and unexported fields must not be restored.
	if round.Hidden != "" || round.hidden != "" {
		t.Fatalf("roundtrip restored ignored fields: %#v", round)
	}
	if round.Empty != "" || round.Count != 0 {
		t.Fatalf("roundtrip restored omitted fields: %#v", round)
	}

	t.Run("nested", func(t *testing.T) {
		var got testNested
		err := MapToStruct(map[string]interface{}{
			"addr": map[string]interface{}{"city": "上海"},
			"team": []interface{}{
				map[string]interface{}{"city": "深圳", "zip": "518000"},
			},
		}, &got)
		if err != nil {
			t.Fatalf("MapToStruct error: %v", err)
		}
		want := testNested{
			Addr: testAddress{City: "上海"},
			Team: []testAddress{{City: "深圳", Zip: "518000"}},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("MapToStruct = %#v, want %#v", got, want)
		}
	})

	t.Run("missing and unknown fields", func(t *testing.T) {
		var got testAddress
		err := MapToStruct(map[string]interface{}{"city": "广州", "unknown": true}, &got)
		if err != nil {
			t.Fatalf("MapToStruct error: %v", err)
		}
		if got.City != "广州" || got.Zip != "" {
			t.Fatalf("MapToStruct = %#v, want {City:广州 Zip:}", got)
		}
	})

	t.Run("nil map is JSON null", func(t *testing.T) {
		got := testAddress{City: "keep"}
		if err := MapToStruct(nil, &got); err != nil {
			t.Fatalf("MapToStruct(nil) error: %v", err)
		}
		if got.City != "keep" {
			t.Fatalf("MapToStruct(nil) changed out: %#v", got)
		}
	})
}

func TestMapToStructErrors(t *testing.T) {
	var nilAddr *testAddress
	nonPointer := testAddress{}
	n := 1
	slice := []int{}
	tests := []struct {
		name string
		m    map[string]interface{}
		out  interface{}
	}{
		{"nil out", map[string]interface{}{"city": "x"}, nil},
		{"non-pointer out", map[string]interface{}{"city": "x"}, nonPointer},
		{"nil pointer out", map[string]interface{}{"city": "x"}, nilAddr},
		{"pointer to int", map[string]interface{}{"city": "x"}, &n},
		{"pointer to slice", map[string]interface{}{"city": "x"}, &slice},
		{"incompatible field", map[string]interface{}{"age": "not-an-int"}, &testProfile{}},
		{"unserializable value", map[string]interface{}{"f": func() {}}, &testProfile{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MapToStruct(tt.m, tt.out); err == nil {
				t.Fatalf("MapToStruct(%#v, %#v) = nil error, want error", tt.m, tt.out)
			}
		})
	}
}
