package set

import (
	"slices"
	"testing"
)

func TestSetExtraOperations(t *testing.T) {
	a := New(1, 2, 3)
	b := New(3, 4)

	if got := sortedInts(a.SymmetricDifference(b)); !slices.Equal(got, []int{1, 2, 4}) {
		t.Fatalf("SymmetricDifference = %v", got)
	}
	if !New(1, 2).IsDisjoint(New(3, 4)) || New(1, 2).IsDisjoint(New(2, 3)) {
		t.Fatal("IsDisjoint mismatch")
	}

	a = a.Update(New(5))
	if !a.Has(5) {
		t.Fatal("Update failed")
	}
	a = a.IntersectionUpdate(New(2, 5, 9))
	if a.Len() != 2 || !a.Has(2) || !a.Has(5) {
		t.Fatalf("IntersectionUpdate = %v", a)
	}
	a = a.DifferenceUpdate(New(2))
	if a.Len() != 1 || !a.Has(5) {
		t.Fatalf("DifferenceUpdate = %v", a)
	}
	a = a.SymmetricDifferenceUpdate(New(5, 6))
	if a.Len() != 1 || !a.Has(6) {
		t.Fatalf("SymmetricDifferenceUpdate = %v", a)
	}
	if got, ok := a.Pop(); !ok || got != 6 || a.Len() != 0 {
		t.Fatalf("Pop = (%v,%v), len=%d", got, ok, a.Len())
	}
	if _, ok := a.Pop(); ok {
		t.Fatal("Pop empty returned ok")
	}

	copied := New(1, 2).Copy()
	copied.Add(3)
	if New(1, 2).Has(3) {
		t.Fatal("Copy shares storage")
	}
	if New(1).Clear().Len() != 0 {
		t.Fatal("Clear failed")
	}
}

func TestFrozenExtraOperations(t *testing.T) {
	a := NewFrozen(1, 2)
	b := NewFrozen(2, 3)
	if got := a.SymmetricDifference(b).Len(); got != 2 {
		t.Fatalf("Frozen SymmetricDifference len = %d", got)
	}
	if !a.IsDisjoint(NewFrozen(3, 4)) || a.IsDisjoint(b) {
		t.Fatal("Frozen IsDisjoint mismatch")
	}
	if !a.Copy().Equal(a) {
		t.Fatal("Frozen Copy failed")
	}
}
