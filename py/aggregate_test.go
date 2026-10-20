package py

import (
	"errors"
	"slices"
	"testing"
)

type namedIntSlice []int

func TestSumGeneric(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Fatalf("Sum([]int) = %d", got)
	}
	if got := Sum([]float64{1.5, 2.5}); got != 4 {
		t.Fatalf("Sum([]float64) = %v", got)
	}
	if got := Sum(namedIntSlice{4, 5}); got != 9 {
		t.Fatalf("Sum(named slice) = %d", got)
	}
	if got := Sum([]int{}); got != 0 {
		t.Fatalf("Sum(empty) = %d", got)
	}
	if got, err := SumAny([]any{1, 2.5}); err != nil || got != 3.5 {
		t.Fatalf("SumAny hetero = (%v,%v)", got, err)
	}
}

func TestSumOf(t *testing.T) {
	if got := SumOf([]int{1, 2, 3}); got != 6 {
		t.Fatalf("SumOf ints = %d, want 6", got)
	}
	if got := SumOf([]float64{1.5, 2.5}); got != 4 {
		t.Fatalf("SumOf floats = %v, want 4", got)
	}
	if got := SumOf([]int{}); got != 0 {
		t.Fatalf("SumOf empty = %d, want 0", got)
	}
	// Python: sum([127, 1]) 不会溢出；Go int8 会回绕。
	if got := SumOf([]int8{127, 1}); got != -128 {
		t.Fatalf("SumOf int8 overflow = %d, want -128", got)
	}
}

func TestMinMaxOf(t *testing.T) {
	if got, err := MinOf([]int{3, 1, 2}); err != nil || got != 1 {
		t.Fatalf("MinOf = (%d,%v), want (1,nil)", got, err)
	}
	if got, err := MaxOf([]int{3, 1, 2}); err != nil || got != 3 {
		t.Fatalf("MaxOf = (%d,%v), want (3,nil)", got, err)
	}
	if got, err := MinOf([]string{"b", "a"}); err != nil || got != "a" {
		t.Fatalf("MinOf strings = (%q,%v), want (a,nil)", got, err)
	}
	if _, err := MinOf([]int{}); !errors.Is(err, ErrEmpty) {
		t.Fatalf("MinOf empty err = %v, want ErrEmpty", err)
	}
	if _, err := MaxOf([]int{}); !errors.Is(err, ErrEmpty) {
		t.Fatalf("MaxOf empty err = %v, want ErrEmpty", err)
	}
}

func TestSortedOf(t *testing.T) {
	input := []int{3, 1, 2, 1}
	got := SortedOf(input)
	if !slices.Equal(got, []int{1, 1, 2, 3}) {
		t.Fatalf("SortedOf = %v, want [1 1 2 3]", got)
	}
	if !slices.Equal(input, []int{3, 1, 2, 1}) {
		t.Fatalf("SortedOf mutated input: %v", input)
	}
	if got := SortedOf([]int(nil)); got == nil || len(got) != 0 {
		t.Fatalf("SortedOf nil = %v, want empty non-nil", got)
	}
	if got := SortedOf([]string{"b", "a"}); !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("SortedOf strings = %v", got)
	}
}

func TestExistingAggregates(t *testing.T) {
	if !All([]interface{}{1, true, "x"}) {
		t.Fatal("All returned false")
	}
	if Any([]interface{}{0, "", nil}) {
		t.Fatal("Any returned true")
	}
	if got := Len("你好"); got != 2 {
		t.Fatalf("Len = %d, want 2", got)
	}
	if got := SortedAny([]int{3, 1, 2}).Join(","); got != "1,2,3" {
		t.Fatalf("Sorted = %q, want 1,2,3", got)
	}
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Fatalf("Sum = %v, want 6", got)
	}
	if got, err := SumAny([]int{1, 2, 3}); err != nil || got != 6 {
		t.Fatalf("SumAny = (%v,%v), want (6,nil)", got, err)
	}
	if got, err := Min([]int{3, 1, 2}); err != nil || got != 1 {
		t.Fatalf("Min = (%v,%v), want (1,nil)", got, err)
	}
	if got, err := Max([]int{3, 1, 2}); err != nil || got != 3 {
		t.Fatalf("Max = (%v,%v), want (3,nil)", got, err)
	}
}

type namedOrdered []int

func TestMinMaxSortedGeneric(t *testing.T) {
	if got, err := Min(namedOrdered{3, 1, 2}); err != nil || got != 1 {
		t.Fatalf("Min named = (%v,%v)", got, err)
	}
	if got, err := Max(namedOrdered{3, 1, 2}); err != nil || got != 3 {
		t.Fatalf("Max named = (%v,%v)", got, err)
	}
	if got := Sorted(namedOrdered{3, 1, 2}); !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("Sorted named = %v", got)
	}
	if got := Sorted([]int(nil)); got == nil || len(got) != 0 {
		t.Fatalf("Sorted nil = %v", got)
	}
	if _, err := Min([]int{}); !errors.Is(err, ErrEmpty) {
		t.Fatalf("Min empty err = %v", err)
	}
	if _, err := MaxAny([]interface{}{3, 1, 2}); err != nil {
		t.Fatalf("MaxAny = %v", err)
	}
}
