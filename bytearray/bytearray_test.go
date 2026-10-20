package bytearray

import (
	"bytes"
	"testing"
)

func TestConstruction(t *testing.T) {
	b := New('a', 'b')
	if string(b) != "ab" {
		t.Fatalf("New = %q", b)
	}
	from, err := From("hi")
	if err != nil || string(from) != "hi" {
		t.Fatalf("From string = (%q,%v)", from, err)
	}
	zeros, err := From(3)
	if err != nil || !bytes.Equal(zeros, []byte{0, 0, 0}) {
		t.Fatalf("From int = (%v,%v)", zeros, err)
	}
	if _, err := From([]int{-1}); err == nil {
		t.Fatal("From negative byte should error")
	}
}

func TestMutation(t *testing.T) {
	b := New(1, 2, 3)
	b.Append(4).Insert(0, 0)
	if !bytes.Equal(b, []byte{0, 1, 2, 3, 4}) {
		t.Fatalf("Append/Insert = %v", b)
	}
	if err := b.Extend([]int{5, 6}); err != nil || !bytes.Equal(b, []byte{0, 1, 2, 3, 4, 5, 6}) {
		t.Fatalf("Extend = (%v,%v)", b, err)
	}
	if value, ok := b.Pop(); !ok || value != 6 {
		t.Fatalf("Pop = (%d,%v)", value, ok)
	}
	if err := b.Remove(0); err != nil || b[0] != 1 {
		t.Fatalf("Remove = (%v,%v)", b, err)
	}
	if err := b.Remove(99); err == nil {
		t.Fatal("Remove missing should error")
	}
	if err := b.Set(-1, 9); err != nil || b[len(b)-1] != 9 {
		t.Fatalf("Set = (%v,%v)", b, err)
	}
	if got, err := b.Get(-1); err != nil || got != 9 {
		t.Fatalf("Get = (%d,%v)", got, err)
	}
	b.Reverse()
	if b[0] != 9 {
		t.Fatalf("Reverse = %v", b)
	}
	b.Clear()
	if b.Len() != 0 {
		t.Fatalf("Clear = %v", b)
	}
}

func TestSearchAndSlice(t *testing.T) {
	b := New('b', 'a', 'n', 'a', 'n', 'a')
	if b.Count('a') != 3 {
		t.Fatalf("Count = %d", b.Count('a'))
	}
	if idx, err := b.Index('n'); err != nil || idx != 2 {
		t.Fatalf("Index = (%d,%v)", idx, err)
	}
	if _, err := b.Index('z'); err == nil {
		t.Fatal("Index missing should error")
	}
	if got := b.Slice(1, -1); string(got) != "anan" {
		t.Fatalf("Slice = %q", got)
	}
	if got := b.Slice(-3, 100); string(got) != "ana" {
		t.Fatalf("Slice negative = %q", got)
	}
}

func TestTransforms(t *testing.T) {
	b := New('a', 'B', 'c')
	if got := b.Upper(); string(got) != "ABC" {
		t.Fatalf("Upper = %q", got)
	}
	if got := b.Lower(); string(got) != "abc" {
		t.Fatalf("Lower = %q", got)
	}
	b.Replace([]byte("B"), []byte("x"), 1)
	if string(b) != "axc" {
		t.Fatalf("Replace = %q", b)
	}
	decoded, err := New('h', 'i').Decode()
	if err != nil || decoded != "hi" {
		t.Fatalf("Decode = (%q,%v)", decoded, err)
	}
	if _, err := New(0xff).Decode(); err == nil {
		t.Fatal("Decode invalid UTF-8 should error")
	}
	original := New(1, 2)
	copied := original.Copy()
	copied[0] = 9
	if original[0] != 1 || !original.Equal(New(1, 2)) {
		t.Fatal("Copy/Equal failed")
	}
	if got := original.Repr(); got != "bytearray(b'\\x01\\x02')" {
		t.Fatalf("Repr = %q", got)
	}
}
