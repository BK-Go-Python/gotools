package bytearray

import (
	"bytes"
	"testing"
)

func TestFromHexAndConcat(t *testing.T) {
	b, err := FromHex("de ad\nbe ef")
	if err != nil || !bytes.Equal(b, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Fatalf("FromHex = (%v,%v)", b, err)
	}
	if _, err := FromHex("zz"); err == nil {
		t.Fatal("FromHex invalid should error")
	}

	joined := New('a').Concat(New('b'), New('c'))
	if string(joined) != "abc" {
		t.Fatalf("Concat = %q", joined)
	}
	if got := New('a', 'b').Repeat(3); string(got) != "ababab" {
		t.Fatalf("Repeat = %q", got)
	}
	if got := New('a').Repeat(0); len(got) != 0 {
		t.Fatalf("Repeat(0) = %q", got)
	}
	if got := New('a').Repeat(-1); len(got) != 0 {
		t.Fatalf("Repeat(-1) = %q", got)
	}
}
