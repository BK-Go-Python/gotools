package conv

import (
	"bytes"
	"slices"
	"testing"
)

func TestBytes(t *testing.T) {
	if got, err := Bytes("hello"); err != nil || string(got) != "hello" {
		t.Fatalf("Bytes(string) = (%q,%v)", got, err)
	}

	source := []byte{1, 2, 3}
	got, err := Bytes(source)
	if err != nil {
		t.Fatal(err)
	}
	got[0] = 9
	if source[0] != 1 {
		t.Fatal("Bytes([]byte) did not copy input")
	}

	if got, err := Bytes([]int{65, 66}); err != nil || !bytes.Equal(got, []byte("AB")) {
		t.Fatalf("Bytes([]int) = (%v,%v)", got, err)
	}
	if got, err := Bytes([]any{65, 66}); err != nil || !bytes.Equal(got, []byte("AB")) {
		t.Fatalf("Bytes([]any) = (%v,%v)", got, err)
	}
	if got, err := Bytes(3); err != nil || !slices.Equal(got, []byte{0, 0, 0}) {
		t.Fatalf("Bytes(3) = (%v,%v)", got, err)
	}
	array, err := ByteArray(2)
	if err != nil || !slices.Equal(array, []byte{0, 0}) {
		t.Fatalf("ByteArray(2) = (%v,%v)", array, err)
	}

	for _, v := range []any{nil, []int{-1}, []int{256}, []any{65, "x"}, []string{"x"}} {
		if _, err := Bytes(v); err == nil {
			t.Fatalf("Bytes(%#v) expected error", v)
		}
	}
}
