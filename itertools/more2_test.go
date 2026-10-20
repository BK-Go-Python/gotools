package itertools

import (
	"slices"
	"testing"
)

func TestTee(t *testing.T) {
	input := []int{1, 2, 3}
	parts := Tee(input, 3)
	if len(parts) != 3 {
		t.Fatalf("Tee = %v", parts)
	}
	parts[0][0] = 9
	if input[0] != 1 || parts[1][0] != 1 {
		t.Fatal("Tee did not copy independently")
	}
	if got := Tee(input, 0); len(got) != 0 {
		t.Fatalf("Tee(0) = %v", got)
	}
}

func TestCombinationsWithReplacement(t *testing.T) {
	got := Combinations_with_replacement([]int{1, 2}, 2)
	want := [][]int{{1, 1}, {1, 2}, {2, 2}}
	if len(got) != len(want) {
		t.Fatalf("combinations = %v", got)
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Fatalf("combinations[%d] = %v, want %v", i, got[i], want[i])
		}
	}
	if got := Combinations_with_replacement([]int{1}, 0); len(got) != 1 || len(got[0]) != 0 {
		t.Fatalf("r=0 = %v", got)
	}
	if got := Combinations_with_replacement([]int{}, 1); len(got) != 0 {
		t.Fatalf("empty = %v", got)
	}
}
