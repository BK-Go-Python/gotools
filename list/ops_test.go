package list

import (
	"slices"
	"testing"
)

func TestConcat(t *testing.T) {
	tests := []struct {
		name   string
		base   List
		others []List
		want   List
	}{
		{"empty", nil, nil, List{}},
		{"no others", New(1, 2), nil, New(1, 2)},
		{"single other", New(1), []List{New(2, 3)}, New(1, 2, 3)},
		{"multiple others", New(1), []List{New(2), New(3, 4)}, New(1, 2, 3, 4)},
		{"empty others", New(1), []List{{}, New(2), nil}, New(1, 2)},
		{"all empty", List{}, []List{{}, nil}, List{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.Concat(tt.others...)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Concat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcatDoesNotModifyInputs(t *testing.T) {
	a := New(1, 2)
	b := New(3)
	c := New(4, 5)
	got := a.Concat(b, c)

	if !slices.Equal(a, New(1, 2)) || !slices.Equal(b, New(3)) || !slices.Equal(c, New(4, 5)) {
		t.Fatalf("Concat modified inputs: a=%v b=%v c=%v", a, b, c)
	}

	got[0] = 99
	if !slices.Equal(a, New(1, 2)) {
		t.Fatalf("Concat result aliases input: a=%v", a)
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name string
		base List
		n    int
		want List
	}{
		{"zero", New(1, 2), 0, List{}},
		{"one", New(1, 2), 1, New(1, 2)},
		{"three", New(1, 2), 3, New(1, 2, 1, 2, 1, 2)},
		{"negative", New(1, 2), -2, List{}},
		{"empty list", List{}, 5, List{}},
		{"nil list", nil, 5, List{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.Repeat(tt.n)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Repeat(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestRepeatDoesNotModifyInput(t *testing.T) {
	src := New(1, 2)
	for _, n := range []int{1, 3} {
		got := src.Repeat(n)
		if !slices.Equal(src, New(1, 2)) {
			t.Fatalf("Repeat(%d) modified input: %v", n, src)
		}
		got[0] = 99
		if !slices.Equal(src, New(1, 2)) {
			t.Fatalf("Repeat(%d) result aliases input: %v", n, src)
		}
	}
}
