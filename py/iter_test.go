package py

import (
	"errors"
	"iter"
	"slices"
	"testing"
)

func collectSeq[T any](seq iter.Seq[T]) []T {
	var out []T
	for value := range seq {
		out = append(out, value)
	}
	return out
}

func TestIter(t *testing.T) {
	it := Iter(slices.Values([]int{1, 2, 3}))
	for want := 1; want <= 3; want++ {
		got, ok := it.Next()
		if !ok || got != want {
			t.Fatalf("Next() = (%d,%v), want (%d,true)", got, ok, want)
		}
	}
	if got, ok := it.Next(); ok || got != 0 {
		t.Fatalf("exhausted Next() = (%d,%v), want (0,false)", got, ok)
	}
	it.Stop()
	it.Stop()
	if _, ok := it.Next(); ok {
		t.Fatal("Next() after Stop returned ok")
	}
}

func TestIterOf(t *testing.T) {
	items := []int{4, 5}
	it := IterOf(items)
	items[0] = 9
	if got, ok := it.Next(); !ok || got != 9 {
		t.Fatalf("Next() = (%d,%v), want (9,true)", got, ok)
	}
	if got, ok := Next(it); !ok || got != 5 {
		t.Fatalf("Next() = (%d,%v), want (5,true)", got, ok)
	}
	if _, ok := Next(it); ok {
		t.Fatal("exhausted iterator returned ok")
	}
	if _, ok := Next[int](nil); ok {
		t.Fatal("Next(nil) returned ok")
	}
}

func TestRangeSeq(t *testing.T) {
	tests := []struct {
		name              string
		start, stop, step int
		want              []int
		wantErr           error
	}{
		{"positive", 0, 5, 1, []int{0, 1, 2, 3, 4}, nil},
		{"step 2", 0, 10, 2, []int{0, 2, 4, 6, 8}, nil},
		{"negative", 5, 0, -1, []int{5, 4, 3, 2, 1}, nil},
		{"empty positive", 5, 0, 1, nil, nil},
		{"empty negative", 0, 5, -1, nil, nil},
		{"zero step", 0, 5, 0, nil, ErrZeroDivision},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, err := RangeSeq(tt.start, tt.stop, tt.step)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			got := collectSeq(seq)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("RangeSeq = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZip(t *testing.T) {
	got2 := Zip2([]int{1, 2, 3}, []string{"a", "b"})
	want2 := []Pair2[int, string]{{1, "a"}, {2, "b"}}
	if !slices.Equal(got2, want2) {
		t.Fatalf("Zip2 = %v, want %v", got2, want2)
	}

	got3 := Zip3([]int{1, 2}, []string{"a"}, []bool{true, false})
	want3 := []Pair3[int, string, bool]{{1, "a", true}}
	if !slices.Equal(got3, want3) {
		t.Fatalf("Zip3 = %v, want %v", got3, want3)
	}
}

func TestTransformations(t *testing.T) {
	pairs := EnumerateOf([]string{"a", "b"}, 10)
	if len(pairs) != 2 || pairs[0] != (PairOf[string]{10, "a"}) || pairs[1] != (PairOf[string]{11, "b"}) {
		t.Fatalf("EnumerateOf = %v", pairs)
	}

	doubled := MapOf([]int{1, 2, 3}, func(v int) int { return v * 2 })
	if !slices.Equal(doubled, []int{2, 4, 6}) {
		t.Fatalf("MapOf = %v", doubled)
	}

	even := FilterOf([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	if !slices.Equal(even, []int{2, 4}) {
		t.Fatalf("FilterOf = %v", even)
	}

	input := []int{1, 2, 3}
	reversed := ReversedOf(input)
	if !slices.Equal(reversed, []int{3, 2, 1}) {
		t.Fatalf("ReversedOf = %v", reversed)
	}
	if !slices.Equal(input, []int{1, 2, 3}) {
		t.Fatalf("ReversedOf mutated input: %v", input)
	}
}
