package conv

import (
	"bytes"
	"math"
	"testing"
)

func TestBytesHex(t *testing.T) {
	if got := BytesHex([]byte{0xde, 0xad, 0xbe, 0xef}); got != "deadbeef" {
		t.Fatalf("BytesHex = %q", got)
	}
	if got := BytesHex(nil); got != "" {
		t.Fatalf("BytesHex(nil) = %q", got)
	}
}

func TestFromHex(t *testing.T) {
	got, err := FromHex("de ad\nbe ef")
	if err != nil || !bytes.Equal(got, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Fatalf("FromHex = (%v,%v)", got, err)
	}
	for _, s := range []string{"abc", "zz"} {
		if _, err := FromHex(s); err == nil {
			t.Fatalf("FromHex(%q) expected error", s)
		}
	}
}

func TestIsInteger(t *testing.T) {
	tests := []struct {
		f    float64
		want bool
	}{
		{3, true},
		{3.5, false},
		{math.NaN(), false},
		{math.Inf(1), false},
		{-2, true},
	}
	for _, tt := range tests {
		if got := IsInteger(tt.f); got != tt.want {
			t.Fatalf("IsInteger(%v) = %v, want %v", tt.f, got, tt.want)
		}
	}
}
