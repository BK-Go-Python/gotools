package collections

import "testing"

func TestCounterBoundaries(t *testing.T) {
	var c Counter[string]
	c = c.Add("a")
	c = c.Update("a", "b")
	c = c.Subtract("a")
	if c.Count("a") != 1 || c.Count("b") != 1 {
		t.Fatalf("Counter after nil updates = %v", c)
	}
	if got := c.Most_common(0); len(got) != 0 {
		t.Fatalf("Most_common(0) = %v", got)
	}
	if got := c.Most_common(-1); len(got) != 2 {
		t.Fatalf("Most_common(-1) = %v", got)
	}
	if c.Total() != 2 {
		t.Fatalf("Total = %d", c.Total())
	}
	if !c.Equal(c.Copy()) {
		t.Fatal("Counter Copy/Equal failed")
	}

	var zero Counter[int]
	if zero.Count(1) != 0 || zero.Total() != 0 || len(zero.Most_common()) != 0 {
		t.Fatal("zero Counter read failed")
	}
	if zero.Equal(Counter[int]{1: 0}) {
		t.Fatal("zero Counter should not equal explicit zero count")
	}
}

func TestDequeBoundaries(t *testing.T) {
	var d Deque[int]
	d.Append(1, 2, 3)
	d.AppendLeft(4, 5)
	if got := d.Items(); len(got) != 5 || got[0] != 5 || got[4] != 3 {
		t.Fatalf("Items = %v", got)
	}
	d.Rotate(2)
	if got := d.Items(); len(got) != 5 || got[0] != 2 {
		t.Fatalf("Rotate(2) = %v", got)
	}
	d.Rotate(-7)
	if d.Len() != 5 {
		t.Fatalf("Rotate(-7) len = %d", d.Len())
	}

	for d.Len() > 0 {
		if _, ok := d.PopLeft(); !ok {
			t.Fatal("PopLeft failed")
		}
	}
	if _, ok := d.Pop(); ok {
		t.Fatal("Pop empty should be false")
	}

	var nilDeque *Deque[int]
	nilDeque.Append(1)
	nilDeque.Rotate(1)
	if nilDeque.Len() != 0 || len(nilDeque.Items()) != 0 {
		t.Fatal("nil Deque should be empty")
	}
	if _, ok := nilDeque.Pop(); ok {
		t.Fatal("nil Deque Pop should be false")
	}
}

func TestOrderedDictBoundaries(t *testing.T) {
	var od OrderedDict[string, int]
	od.Set("a", 1).Set("b", 2).Set("c", 3)
	od.Delete("b")
	od.Set("b", 20)
	if got := od.Keys(); len(got) != 3 || got[2] != "b" {
		t.Fatalf("Keys after delete/reinsert = %v", got)
	}
	item, ok := od.Popitem()
	if !ok || item.Key != "b" || item.Value != 20 {
		t.Fatalf("Popitem LIFO = (%v,%v)", item, ok)
	}
	item, ok = od.Popitem(false)
	if !ok || item.Key != "a" {
		t.Fatalf("Popitem FIFO = (%v,%v)", item, ok)
	}
	if od.Len() != 1 {
		t.Fatalf("Len = %d", od.Len())
	}

	other := NewOrderedDict[string, int]()
	other.Set("c", 3)
	if !od.Equal(other) {
		t.Fatal("OrderedDict Equal failed")
	}
	other.Set("d", 4)
	other.Delete("c")
	other.Set("c", 3)
	if od.Equal(other) {
		t.Fatal("OrderedDict Equal should be order sensitive")
	}

	var nilDict *OrderedDict[string, int]
	if nilDict.Len() != 0 || len(nilDict.Keys()) != 0 || nilDict.Has("x") {
		t.Fatal("nil OrderedDict read failed")
	}
	nilDict.Set("x", 1)
	if nilDict.Len() != 0 {
		t.Fatal("nil OrderedDict Set should be no-op")
	}
}
