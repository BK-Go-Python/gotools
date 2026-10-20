package str

import "testing"

// Expected values in this file were checked against CPython 3.13.15.

func TestReplace(t *testing.T) {
	tests := []struct {
		name        string
		s, old, new string
		count       []int
		want        string
	}{
		{"all by default", "banana", "a", "o", nil, "bonono"},
		{"count one", "banana", "a", "o", []int{1}, "bonana"},
		{"count zero", "banana", "a", "o", []int{0}, "banana"},
		{"count two", "banana", "a", "o", []int{2}, "bonona"},
		{"negative count replaces all", "banana", "a", "o", []int{-2}, "bonono"},
		{"count beyond occurrences", "banana", "a", "o", []int{99}, "bonono"},
		{"empty old inserts", "abc", "", "-", nil, "-a-b-c-"},
		{"empty old with count", "abc", "", "-", []int{2}, "-a-bc"},
		{"empty old with count zero", "abc", "", "-", []int{0}, "abc"},
		{"empty old on empty input", "", "", "-", nil, "-"},
		{"empty old at unicode boundaries", "你😀", "", "-", nil, "-你-😀-"},
		{"no match", "banana", "z", "x", nil, "banana"},
		{"empty input", "", "a", "b", nil, ""},
		{"delete occurrences", "banana", "a", "", nil, "bnn"},
		{"old equals new", "aaa", "a", "a", nil, "aaa"},
		{"non-overlapping matches", "aaaa", "aa", "b", nil, "bb"},
		{"chinese", "你好世界", "好", "棒", nil, "你棒世界"},
		{"emoji", "a😀b😀", "😀", "🚀", nil, "a🚀b🚀"},
		{"multi rune pattern with count", "xabcxabc", "abc", "z", []int{1}, "xzxabc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Replace(tt.s, tt.old, tt.new, tt.count...); got != tt.want {
				t.Fatalf("Replace(%q, %q, %q, %v) = %q, want %q",
					tt.s, tt.old, tt.new, tt.count, got, tt.want)
			}
		})
	}
}
