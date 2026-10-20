package dict

import "testing"

func TestPythonSpellingNames(t *testing.T) {
	d := Fromkeys([]string{"a", "b"}, 0)
	if len(d) != 2 || d["a"] != 0 {
		t.Fatalf("Fromkeys = %v", d)
	}
	if got := d.Setdefault("a", 9); got != 0 {
		t.Fatalf("Setdefault existing = %v", got)
	}
	if got := d.Setdefault("c", 9); got != 9 || d["c"] != 9 {
		t.Fatalf("Setdefault missing = %v", got)
	}
	item, err := d.Popitem()
	if err != nil || item.Key == "" {
		t.Fatalf("Popitem = (%v,%v)", item, err)
	}
}
