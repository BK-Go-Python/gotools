package collections

import "testing"

func TestPythonSpellingAliases(t *testing.T) {
	c := NewCounter("a", "a", "b")
	got := c.Most_common(1)
	if len(got) != 1 || got[0].Value != "a" || got[0].Count != 2 {
		t.Fatalf("Most_common = %v", got)
	}

	d := NewDeque(2)
	d.Appendleft(1)
	d.Extendleft(0)
	if got, ok := d.Popleft(); !ok || got != 0 {
		t.Fatalf("Popleft = (%v,%v)", got, ok)
	}
	if got, ok := d.Peekleft(); !ok || got != 1 {
		t.Fatalf("Peekleft = (%v,%v)", got, ok)
	}

	od := NewOrderedDict[string, int]()
	od.Set("a", 1).Set("b", 2)
	item, ok := od.Popitem()
	if !ok || item.Key != "b" {
		t.Fatalf("Popitem = (%v,%v)", item, ok)
	}
	if !od.Move_to_end("a") {
		t.Fatal("Move_to_end failed")
	}
}
