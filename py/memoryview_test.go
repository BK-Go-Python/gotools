package py

import "testing"

func TestMemoryView(t *testing.T) {
	base := []byte{1, 2, 3}
	view := MemoryView(base)
	if len(view) != 3 || cap(view) != 3 {
		t.Fatalf("len/cap = %d/%d, want 3/3", len(view), cap(view))
	}
	view[0] = 9
	if base[0] != 9 {
		t.Fatal("MemoryView does not share backing array")
	}
	if MemoryView(nil) != nil {
		t.Fatal("MemoryView(nil) != nil")
	}
	if got := MemoryView([]byte{}); len(got) != 0 {
		t.Fatalf("MemoryView(empty) len = %d", len(got))
	}
}
