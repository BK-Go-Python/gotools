package str

import "testing"

func TestASCII(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{
		{"ascii string", "abc", "'abc'"},
		{"latin1", "é", "'\\xe9'"},
		{"bmp", "你", "'\\u4f60'"},
		{"number", 123, "123"},
		{"nil", nil, "None"},
		{"bool", true, "True"},
		{"newline", "\n", "'\\n'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ASCII(tt.in); got != tt.want {
				t.Fatalf("ASCII(%#v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
