package list

import (
	"slices"
	"testing"
)

func TestOf(t *testing.T) {
	input := []int{1, 2, 3}
	got := Of(input...)
	got[0] = 9
	if input[0] != 1 {
		t.Fatal("Of did not copy input")
	}
	if !slices.Equal(got, []int{9, 2, 3}) {
		t.Fatalf("Of = %v", got)
	}
}
