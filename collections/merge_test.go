package collections

import "testing"

func TestCounterMergeOperations(t *testing.T) {
	a := NewCounter("a", "a", "a", "b")
	b := NewCounter("a", "b", "b")

	if got := a.AddCounter(b); got.Count("a") != 4 || got.Count("b") != 3 {
		t.Fatalf("AddCounter = %v", got)
	}
	if got := a.SubCounter(b); got.Count("a") != 2 || got.Count("b") != 0 {
		t.Fatalf("SubCounter = %v", got)
	}
	if got := a.Union(b); got.Count("a") != 3 || got.Count("b") != 2 {
		t.Fatalf("Union = %v", got)
	}
	if got := a.Intersection(b); got.Count("a") != 1 || got.Count("b") != 1 {
		t.Fatalf("Intersection = %v", got)
	}
	negative := NewCounter[string]()
	negative["x"] = -2
	negative["y"] = 1
	if got := negative.Negate(); got.Count("x") != 2 || got.Count("y") != 0 {
		t.Fatalf("Negate = %v", got)
	}

	var nilCounter Counter[string]
	if got := nilCounter.UpdateCounter(a); got.Count("a") != 3 {
		t.Fatalf("UpdateCounter nil = %v", got)
	}
	if got := a.SubtractCounter(b); got.Count("a") != 2 || got.Count("b") != -1 {
		t.Fatalf("SubtractCounter = %v", got)
	}
}

func TestOrderedDictUpdateAndMerge(t *testing.T) {
	d := NewOrderedDict[string, int]()
	d.Set("a", 1).Set("b", 2)
	if got := d.Setdefault("a", 9); got != 1 {
		t.Fatalf("Setdefault existing = %v", got)
	}
	if got := d.Setdefault("c", 3); got != 3 || d.Keys()[2] != "c" {
		t.Fatalf("Setdefault missing = (%v,%v)", got, d.Keys())
	}

	other := NewOrderedDict[string, int]()
	other.Set("b", 20).Set("d", 4)
	d.Update(other)
	if got := d.Keys(); len(got) != 4 || got[0] != "a" || got[1] != "b" || got[3] != "d" {
		t.Fatalf("Update keys = %v", got)
	}
	if value := d.Get("b"); value != 20 {
		t.Fatalf("Update b = %d", value)
	}

	merged := d.Merge(NewOrderedDict[string, int]())
	if !merged.Equal(d) {
		t.Fatal("Merge empty failed")
	}
	extra := NewOrderedDict[string, int]()
	extra.Set("e", 5).Set("a", 10)
	merged = d.Merge(extra)
	if merged.Len() != 5 || merged.Keys()[4] != "e" {
		t.Fatalf("Merge keys = %v", merged.Keys())
	}
	if value := merged.Get("a"); value != 10 {
		t.Fatalf("Merge overwrite a = %d", value)
	}
	if d.Get("a") != 1 {
		t.Fatal("Merge mutated source")
	}
}
