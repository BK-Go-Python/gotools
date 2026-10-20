package collections

import (
	"slices"
	"testing"
)

func TestDefaultDict(t *testing.T) {
	calls := 0
	d := NewDefaultDict[string](func() int { calls++; return 7 })
	if got := d.Get("a"); got != 7 || calls != 1 {
		t.Fatalf("Get missing = (%d,%d)", got, calls)
	}
	if got := d.Get("a"); got != 7 || calls != 1 {
		t.Fatalf("Get cached = (%d,%d)", got, calls)
	}
	d.Set("b", 2)
	if !d.Has("b") || d.Len() != 2 {
		t.Fatalf("Set/Has/Len = %v %d", d.Keys(), d.Len())
	}
	if got := d.Keys(); !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("Keys = %v", got)
	}
	items := d.Items()
	if len(items) != 2 || items[0].Key != "a" || items[0].Value != 7 {
		t.Fatalf("Items = %v", items)
	}
	copied := d.Copy()
	copied.Set("c", 3)
	if d.Has("c") {
		t.Fatal("Copy shares storage")
	}
	if !copied.Equal(copied.Copy()) {
		t.Fatal("Equal failed")
	}
	if d.Delete("a") != true || d.Delete("a") != false {
		t.Fatal("Delete mismatch")
	}
	d.Clear()
	if d.Len() != 0 {
		t.Fatal("Clear failed")
	}
}

func TestDefaultDictNil(t *testing.T) {
	var d *DefaultDict[string, int]
	if d.Get("x") != 0 || d.Has("x") || d.Len() != 0 || d.Delete("x") {
		t.Fatal("nil DefaultDict read failed")
	}
	d.Set("x", 1)
	if d.Len() != 0 {
		t.Fatal("nil Set should be no-op")
	}
}
