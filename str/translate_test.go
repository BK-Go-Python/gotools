package str

import (
	"reflect"
	"testing"
)

// Expected values in this file were checked against CPython 3.13.15.

func TestMaketrans(t *testing.T) {
	tests := []struct {
		name     string
		from, to string
		want     map[rune]string
	}{
		{"ascii", "abc", "xyz", map[rune]string{'a': "x", 'b': "y", 'c': "z"}},
		{"chinese and emoji", "你好😀", "我她🚀", map[rune]string{'你': "我", '好': "她", '😀': "🚀"}},
		{"empty", "", "", map[rune]string{}},
		{"duplicate from rune last wins", "aba", "xyz", map[rune]string{'a': "z", 'b': "y"}},
		{"multi rune values", "ab", "你好", map[rune]string{'a': "你", 'b': "好"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Maketrans(tt.from, tt.to)
			if err != nil {
				t.Fatalf("Maketrans(%q, %q) returned error: %v", tt.from, tt.to, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Maketrans(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestMaketransLengthMismatch(t *testing.T) {
	tests := []struct {
		name     string
		from, to string
	}{
		{"ascii", "abc", "xy"},
		{"to longer", "ab", "xyz"},
		{"chinese", "你好", "好"},
		{"emoji", "😀", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := Maketrans(tt.from, tt.to)
			if err == nil {
				t.Fatalf("Maketrans(%q, %q) = %v, want error", tt.from, tt.to, table)
			}
			if table != nil {
				t.Fatalf("Maketrans(%q, %q) table = %v on error, want nil", tt.from, tt.to, table)
			}
		})
	}
}

func TestTranslate(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		table map[rune]string
		want  string
	}{
		{"map and delete", "abc", map[rune]string{'a': "A", 'b': ""}, "Ac"},
		{"unmapped kept", "abc", map[rune]string{'b': "B"}, "aBc"},
		{"nil table", "你好😀", nil, "你好😀"},
		{"empty table", "abc", map[rune]string{}, "abc"},
		{"empty input", "", map[rune]string{'a': "A"}, ""},
		{"delete all", "abc", map[rune]string{'a': "", 'b': "", 'c': ""}, ""},
		{"multi rune replacement", "abc", map[rune]string{'a': "XYZ"}, "XYZbc"},
		{"replace with same", "abc", map[rune]string{'a': "a"}, "abc"},
		{"chinese delete", "你好世界", map[rune]string{'你': "我", '好': ""}, "我世界"},
		{"emoji", "a😀b", map[rune]string{'😀': "👋"}, "a👋b"},
		{"emoji multi rune replacement", "你", map[rune]string{'你': "😀🚀"}, "😀🚀"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Translate(tt.s, tt.table); got != tt.want {
				t.Fatalf("Translate(%q, %v) = %q, want %q", tt.s, tt.table, got, tt.want)
			}
		})
	}
}
