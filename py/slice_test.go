package py

import (
	"errors"
	"testing"
)

func TestNewSliceIndices(t *testing.T) {
	tests := []struct {
		name                          string
		args                          []any
		n                             int
		wantStart, wantStop, wantStep int
	}{
		{"all none", nil, 5, 0, 5, 1},
		{"stop only", []any{3}, 5, 0, 3, 1},
		{"start stop", []any{1, 3}, 5, 1, 3, 1},
		{"start stop step", []any{1, 5, 2}, 10, 1, 5, 2},
		{"negative start", []any{-3, nil, 1}, 5, 2, 5, 1},
		{"negative stop", []any{1, -1, 1}, 5, 1, 4, 1},
		{"reverse", []any{nil, nil, -1}, 5, 4, -1, -1},
		{"reverse start beyond", []any{5, 0, -1}, 5, 4, 0, -1},
		{"reverse explicit stop", []any{nil, -1, -1}, 5, 4, 4, -1},
		{"empty n", []any{1, 3}, 0, 0, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewSlice(tt.args...)
			if err != nil {
				t.Fatal(err)
			}
			start, stop, step := s.Indices(tt.n)
			if start != tt.wantStart || stop != tt.wantStop || step != tt.wantStep {
				t.Fatalf("Indices(%d) = (%d,%d,%d), want (%d,%d,%d)",
					tt.n, start, stop, step, tt.wantStart, tt.wantStop, tt.wantStep)
			}
		})
	}
}

func TestNewSliceErrors(t *testing.T) {
	if _, err := NewSlice(1, 2, 0); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("zero step err = %v, want ErrOutOfRange", err)
	}
	if _, err := NewSlice(1, 2, 3, 4); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("too many args err = %v, want ErrOutOfRange", err)
	}
	if _, err := NewSlice("x"); !errors.Is(err, ErrType) {
		t.Fatalf("bad type err = %v, want ErrType", err)
	}
}

func TestSliceDirectZeroStep(t *testing.T) {
	start, stop, step := Slice{Step: 0}.Indices(4)
	if start != 0 || stop != 4 || step != 1 {
		t.Fatalf("Indices = (%d,%d,%d), want (0,4,1)", start, stop, step)
	}
}
