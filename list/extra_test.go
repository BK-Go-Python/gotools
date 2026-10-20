package list

import (
	"slices"
	"testing"
)

func TestIndexRange(t *testing.T) {
	l := New(1, 2, 3, 2)
	if got, err := l.IndexRange(2, 0, 4); err != nil || got != 1 {
		t.Fatalf("IndexRange full = (%d,%v)", got, err)
	}
	if got, err := l.IndexRange(2, 2, 4); err != nil || got != 3 {
		t.Fatalf("IndexRange start = (%d,%v)", got, err)
	}
	if got, err := l.IndexRange(2, -2, 4); err != nil || got != 3 {
		t.Fatalf("IndexRange negative start = (%d,%v)", got, err)
	}
	if _, err := l.IndexRange(9, 0, 4); err == nil {
		t.Fatal("IndexRange missing should return error")
	}
}

func TestSortKeyCallsKeyOnce(t *testing.T) {
	l := New(3, 1, 2)
	calls := 0
	l.SortKey(func(v interface{}) interface{} {
		calls++
		return v
	})
	if calls != 3 {
		t.Fatalf("key called %d times, want 3", calls)
	}
	if !slices.Equal([]interface{}(l), []interface{}{1, 2, 3}) {
		t.Fatalf("SortKey = %v", l)
	}

	l = New(3, 1, 2)
	l.SortKey(func(v interface{}) interface{} { return v }, true)
	if !slices.Equal([]interface{}(l), []interface{}{3, 2, 1}) {
		t.Fatalf("SortKey reverse = %v", l)
	}
}

func TestSortByUsesCachedKeys(t *testing.T) {
	l := New(3, 1, 2)
	calls := 0
	l.SortBy(func(v interface{}) interface{} {
		calls++
		return v
	})
	if calls != 3 {
		t.Fatalf("SortBy key called %d times, want 3", calls)
	}
}
