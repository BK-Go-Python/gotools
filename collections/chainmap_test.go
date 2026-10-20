package collections

import (
	"slices"
	"testing"
)

func TestChainMapLookupAndWrite(t *testing.T) {
	first := map[string]int{"a": 1}
	second := map[string]int{"a": 2, "b": 3}
	cm := NewChainMap(first, second)

	if got := cm.Get("a"); got != 1 {
		t.Fatalf("Get a = %d", got)
	}
	if got := cm.Get("b"); got != 3 {
		t.Fatalf("Get b = %d", got)
	}
	if got := cm.Get("x", 9); got != 9 {
		t.Fatalf("Get default = %d", got)
	}
	if !cm.Has("b") || cm.Has("x") {
		t.Fatal("Has mismatch")
	}
	cm.Set("b", 30)
	if first["b"] != 30 || second["b"] != 3 {
		t.Fatalf("Set wrote wrong map: %v %v", first, second)
	}
	if !cm.Delete("a") || first["a"] != 0 || second["a"] != 2 {
		t.Fatalf("Delete = %v %v", first, second)
	}
	if cm.Delete("a") {
		t.Fatal("Delete missing should be false")
	}
	if got := cm.Len(); got != 2 {
		t.Fatalf("Len = %d", got)
	}
	if got := cm.Keys(); !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("Keys = %v", got)
	}
	items := cm.Items()
	if len(items) != 2 || items[0].Key != "a" || items[0].Value != 2 {
		t.Fatalf("Items = %v", items)
	}
}

func TestChainMapChildAndNil(t *testing.T) {
	base := NewChainMap(map[string]int{"a": 1})
	child := base.NewChild()
	child.Set("b", 2)
	if base.Has("b") {
		t.Fatal("child write leaked to base")
	}
	if child.Get("a") != 1 {
		t.Fatal("child lookup failed")
	}
	if len(child.Maps()) != 2 {
		t.Fatalf("Maps = %d", len(child.Maps()))
	}

	var nilCM *ChainMap[string, int]
	if nilCM.Get("a") != 0 || nilCM.Has("a") || nilCM.Len() != 0 {
		t.Fatal("nil ChainMap read failed")
	}
	nilCM.Set("a", 1)
	if nilCM.Delete("a") {
		t.Fatal("nil Delete should be false")
	}
}
