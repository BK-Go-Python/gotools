package set

import (
	"slices"
	"testing"
)

func sortedInts(s Set[int]) []int {
	items := s.Items()
	slices.Sort(items)
	return items
}

func TestSetBasics(t *testing.T) {
	s := New(3, 1, 2, 3)
	if s.Len() != 3 || !s.Has(1) || s.Has(9) {
		t.Fatalf("unexpected set: %v", s)
	}
	if got := sortedInts(s); !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("Items = %v", got)
	}

	s = s.Add(4).Remove(1)
	if !s.Has(4) || s.Has(1) {
		t.Fatalf("after add/remove: %v", s)
	}
}

func TestSetZeroValue(t *testing.T) {
	var s Set[int]
	s = s.Add(1)
	if !s.Has(1) {
		t.Fatal("zero value Add failed")
	}
	s.Remove(1)
	if s.Len() != 0 {
		t.Fatal("zero value Remove failed")
	}
}

func TestSetOperations(t *testing.T) {
	a := New(1, 2, 3)
	b := New(3, 4)

	if got := sortedInts(a.Union(b)); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Fatalf("Union = %v", got)
	}
	if got := sortedInts(a.Intersection(b)); !slices.Equal(got, []int{3}) {
		t.Fatalf("Intersection = %v", got)
	}
	if got := sortedInts(a.Difference(b)); !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("Difference = %v", got)
	}
	if !New(1, 2).IsSubset(a) {
		t.Fatal("IsSubset returned false")
	}
	if New(1, 9).IsSubset(a) {
		t.Fatal("IsSubset returned true for non-subset")
	}
	if !New(1, 2).Equal(New(2, 1)) {
		t.Fatal("Equal returned false")
	}

	clone := a.Clone()
	clone.Add(9)
	if a.Has(9) {
		t.Fatal("Clone shares storage")
	}
}

func TestFrozen(t *testing.T) {
	f := NewFrozen(1, 2, 3)
	if f.Len() != 3 || !f.Has(2) {
		t.Fatalf("unexpected frozen: %v", f)
	}
	g := NewFrozen(3, 4)
	if got := f.Union(g).Len(); got != 4 {
		t.Fatalf("Union len = %d", got)
	}
	if got := f.Intersection(g).Len(); got != 1 {
		t.Fatalf("Intersection len = %d", got)
	}
	if !NewFrozen(1, 2).IsSubset(f) {
		t.Fatal("Frozen IsSubset failed")
	}
	if !f.Equal(f.Clone()) {
		t.Fatal("Frozen Equal/Clone failed")
	}
}
