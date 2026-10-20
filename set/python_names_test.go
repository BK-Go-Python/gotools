package set

import "testing"

func TestPythonSpellingNames(t *testing.T) {
	a, b := New(1, 2), New(2, 3)
	if !New(1, 2).Issubset(a) || !a.Isdisjoint(New(4)) {
		t.Fatal("subset/disjoint alias failed")
	}
	if got := a.Symmetric_difference(b).Len(); got != 2 {
		t.Fatalf("Symmetric_difference len = %d", got)
	}
	if got := New(1, 2, 3).Intersection_update(New(2, 3)).Len(); got != 2 {
		t.Fatalf("Intersection_update len = %d", got)
	}
	if got := New(1, 2, 3).Difference_update(New(2)).Len(); got != 2 {
		t.Fatalf("Difference_update len = %d", got)
	}
	if got := New(1, 2).Symmetric_difference_update(New(2, 3)).Len(); got != 2 {
		t.Fatalf("Symmetric_difference_update len = %d", got)
	}
}
