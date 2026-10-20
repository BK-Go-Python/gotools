package conv

import (
	"bytes"
	"testing"
)

func TestPythonSpellingNames(t *testing.T) {
	b, err := Int_to_bytes(1024, 2, "big")
	if err != nil || !bytes.Equal(b, []byte{0x04, 0x00}) {
		t.Fatalf("Int_to_bytes = (%v,%v)", b, err)
	}
	if n, err := Int_from_bytes(b, "big"); err != nil || n != 1024 {
		t.Fatalf("Int_from_bytes = (%d,%v)", n, err)
	}
	if got := Float_hex(1.0); got != "0x1.0000000000000p+0" {
		t.Fatalf("Float_hex = %q", got)
	}
	if f, err := Float_fromhex("0x1.8p+1"); err != nil || f != 3 {
		t.Fatalf("Float_fromhex = (%v,%v)", f, err)
	}
	if got := Bytes_hex([]byte{0xde, 0xad}); got != "dead" {
		t.Fatalf("Bytes_hex = %q", got)
	}
	if got, err := Bytes_fromhex("de ad"); err != nil || !bytes.Equal(got, []byte{0xde, 0xad}) {
		t.Fatalf("Bytes_fromhex = (%v,%v)", got, err)
	}
	if got := Str_encode("hi"); string(got) != "hi" {
		t.Fatalf("Str_encode = %q", got)
	}
	if got, err := Bytes_decode([]byte("hi")); err != nil || got != "hi" {
		t.Fatalf("Bytes_decode = (%q,%v)", got, err)
	}
}
