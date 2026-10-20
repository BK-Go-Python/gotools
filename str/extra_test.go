package str

import "testing"

func TestCasefold(t *testing.T) {
	tests := []struct{ in, want string }{
		{"ABC", "abc"},
		{"ß", "ss"},
		{"ẞ", "ss"},
		{"İ", "i\u0307"},
		{"你好", "你好"},
	}
	for _, tt := range tests {
		if got := Casefold(tt.in); got != tt.want {
			t.Fatalf("Casefold(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestExpandtabs(t *testing.T) {
	tests := []struct {
		in   string
		size []int
		want string
	}{
		{"a\tb", []int{4}, "a   b"},
		{"\t", []int{4}, "    "},
		{"a\n\tb", []int{4}, "a\n    b"},
		{"a\tb", []int{0}, "ab"},
		{"a\tb", nil, "a       b"},
	}
	for _, tt := range tests {
		if got := Expandtabs(tt.in, tt.size...); got != tt.want {
			t.Fatalf("Expandtabs(%q,%v) = %q, want %q", tt.in, tt.size, got, tt.want)
		}
	}
}

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"_x", true},
		{"x1", true},
		{"变量", true},
		{"1x", false},
		{"", false},
		{"a-b", false},
	}
	for _, tt := range tests {
		if got := IsIdentifier(tt.in); got != tt.want {
			t.Fatalf("IsIdentifier(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"abc", true},
		{"你好", true},
		{"a b", true},
		{"a\nb", false},
	}
	for _, tt := range tests {
		if got := IsPrintable(tt.in); got != tt.want {
			t.Fatalf("IsPrintable(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSplitLinesKeepends(t *testing.T) {
	got := SplitLinesKeepends("a\r\nb\rc\nd\n")
	want := []string{"a\r\n", "b\r", "c\n", "d\n"}
	if len(got) != len(want) {
		t.Fatalf("SplitLinesKeepends = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SplitLinesKeepends[%d] = %q, want %q", i, got[i], want[i])
		}
	}
	if len(SplitLinesKeepends("")) != 0 {
		t.Fatal("SplitLinesKeepends empty should return empty")
	}
}
