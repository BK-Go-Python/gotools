package functools

import (
	"slices"
	"sync"
	"testing"
)

func TestMemoize(t *testing.T) {
	calls := 0
	fn := Memoize(func(key int) int {
		calls++
		return key * 2
	})
	if fn(2) != 4 || fn(2) != 4 || fn(3) != 6 {
		t.Fatal("Memoize result mismatch")
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if fn(10) != 20 {
				t.Error("concurrent Memoize mismatch")
			}
		}()
	}
	wg.Wait()
}

func TestLruCache(t *testing.T) {
	c := NewLruCache[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	if value, ok := c.Get("a"); !ok || value != 1 {
		t.Fatalf("Get a = (%d,%v)", value, ok)
	}
	c.Put("c", 3)
	if _, ok := c.Get("b"); ok {
		t.Fatal("b should have been evicted")
	}
	if got := c.Keys(); !slices.Equal(got, []string{"c", "a"}) {
		t.Fatalf("Keys = %v", got)
	}
	if c.Len() != 2 || c.Cap() != 2 {
		t.Fatalf("Len/Cap = %d/%d", c.Len(), c.Cap())
	}
	if !c.Remove("a") || c.Remove("a") {
		t.Fatal("Remove mismatch")
	}
	c.Clear()
	if c.Len() != 0 {
		t.Fatal("Clear failed")
	}

	zero := NewLruCache[int, int](0)
	zero.Put(1, 1)
	if _, ok := zero.Get(1); ok || zero.Len() != 0 {
		t.Fatal("zero-capacity cache should not store")
	}
}
