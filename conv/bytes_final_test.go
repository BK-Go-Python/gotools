package conv

import (
	"bytes"
	"testing"
)

func TestBytesRindex(t *testing.T) {
	idx, err := Bytes_rindex([]byte("banana"), []byte("na"))
	if err != nil || idx != 4 {
		t.Fatalf("rindex = (%d,%v)", idx, err)
	}
	if _, err := Bytes_rindex([]byte("abc"), []byte("x")); err == nil {
		t.Fatal("rindex missing should error")
	}
}

func TestBytesSplitlinesKeepends(t *testing.T) {
	got := Bytes_splitlines_keepends([]byte("a\r\nb\rc\nd\n"))
	want := [][]byte{[]byte("a\r\n"), []byte("b\r"), []byte("c\n"), []byte("d\n")}
	if len(got) != len(want) {
		t.Fatalf("splitlines = %q", got)
	}
	for i := range want {
		if !bytes.Equal(got[i], want[i]) {
			t.Fatalf("splitlines[%d] = %q, want %q", i, got[i], want[i])
		}
	}
	if len(Bytes_splitlines_keepends(nil)) != 0 {
		t.Fatal("nil splitlines should be empty")
	}
}

func TestBytesHexSep(t *testing.T) {
	tests := []struct {
		data []byte
		sep  string
		n    int
		want string
	}{
		{[]byte{0xde, 0xad, 0xbe, 0xef}, " ", 1, "de ad be ef"},
		{[]byte{0xde, 0xad, 0xbe, 0xef}, ":", 2, "dead:beef"},
		{[]byte("12345"), "-", -2, "3132-3334-35"},
		{[]byte{0xde, 0xad}, "", 1, "dead"},
	}
	for _, tt := range tests {
		got, err := Bytes_hex_sep(tt.data, tt.sep, tt.n)
		if err != nil || got != tt.want {
			t.Fatalf("hex_sep(%v,%q,%d) = (%q,%v), want %q", tt.data, tt.sep, tt.n, got, err, tt.want)
		}
	}
	if _, err := Bytes_hex_sep([]byte{1}, "-", 0); err == nil {
		t.Fatal("zero bytes_per_sep should error")
	}
}
