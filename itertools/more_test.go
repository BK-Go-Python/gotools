package itertools

import (
	"slices"
	"testing"
)

func TestFilterfalseCompress(t *testing.T) {
	got := Filterfalse([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	if !slices.Equal(got, []int{1, 3}) {
		t.Fatalf("Filterfalse = %v", got)
	}
	if got := Filterfalse([]int{1}, nil); len(got) != 0 {
		t.Fatalf("Filterfalse nil pred = %v", got)
	}
	compressed := Compress([]string{"a", "b", "c"}, []bool{true, false, true})
	if !slices.Equal(compressed, []string{"a", "c"}) {
		t.Fatalf("Compress = %v", compressed)
	}
	if got := Compress([]string{"a"}, nil); len(got) != 0 {
		t.Fatalf("Compress nil selectors = %v", got)
	}
}

func TestStarmapCount(t *testing.T) {
	pairs := []Pair2[int, int]{{1, 2}, {3, 4}}
	got := Starmap(pairs, func(a, b int) int { return a + b })
	if !slices.Equal(got, []int{3, 7}) {
		t.Fatalf("Starmap = %v", got)
	}
	if got := Starmap[int, int, int](pairs, nil); len(got) != 0 {
		t.Fatalf("Starmap nil fn = %v", got)
	}
	if got := Count(5, 2, 3); !slices.Equal(got, []int{5, 7, 9}) {
		t.Fatalf("Count = %v", got)
	}
	if got := Count(0, 1, 0); len(got) != 0 {
		t.Fatalf("Count 0 = %v", got)
	}
}
