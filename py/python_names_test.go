package py

import (
	"bytes"
	"testing"
)

func TestPythonSpellingNames(t *testing.T) {
	if got := Bit_count(7); got != 3 {
		t.Fatalf("Bit_count = %d", got)
	}
	if got := Bit_length(8); got != 4 {
		t.Fatalf("Bit_length = %d", got)
	}

	var buf bytes.Buffer
	if err := Print_to(&buf, "-", "!", 1, "a"); err != nil || buf.String() != "1-a!" {
		t.Fatalf("Print_to = (%q,%v)", buf.String(), err)
	}

	if !Isinstance(1, 0) || !Issubclass((*bytes.Buffer)(nil), (*bytes.Buffer)(nil)) {
		t.Fatalf("type predicates failed")
	}

	m := map[string]any{"a": 1}
	if got, err := Getattr(m, "a"); err != nil || got != 1 {
		t.Fatalf("Getattr = (%v,%v)", got, err)
	}
	if err := Setattr(m, "b", 2); err != nil || m["b"] != 2 {
		t.Fatalf("Setattr = (%v,%v)", m, err)
	}
	if !Hasattr(m, "b") || Hasattr(m, "c") {
		t.Fatal("Hasattr mismatch")
	}
	if err := Delattr(m, "b"); err != nil || Hasattr(m, "b") {
		t.Fatalf("Delattr = %v", err)
	}
	if got := Memoryview([]byte{1}); len(got) != 1 {
		t.Fatalf("Memoryview = %v", got)
	}
}
