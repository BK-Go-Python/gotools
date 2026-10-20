package itertools

import (
	"reflect"
	"testing"
)

// streamCheckNonNil 断言结果是非 nil 切片（空结果也必须是非 nil）。
func streamCheckNonNil[T any](t *testing.T, name string, got []T) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s returned nil slice, want non-nil empty slice", name)
	}
}

func streamCheckEqual[T any](t *testing.T, name string, got, want []T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s = %#v, want %#v", name, got, want)
	}
}

func TestChain(t *testing.T) {
	tests := []struct {
		name string
		seqs [][]int
		want []int
	}{
		{"no args", nil, []int{}},
		{"nil slice", [][]int{nil}, []int{}},
		{"empty slice", [][]int{{}}, []int{}},
		{"single", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
		{"multiple with nil", [][]int{{1, 2}, {}, nil, {3, 4}}, []int{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chain(tt.seqs...)
			streamCheckNonNil(t, "Chain", got)
			streamCheckEqual(t, "Chain", got, tt.want)
		})
	}

	t.Run("strings", func(t *testing.T) {
		got := Chain([]string{"ab", "cd"}, nil, []string{"ef"})
		streamCheckNonNil(t, "Chain", got)
		streamCheckEqual(t, "Chain", got, []string{"ab", "cd", "ef"})
	})
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name  string
		value int
		n     int
		want  []int
	}{
		{"zero", 7, 0, []int{}},
		{"negative", 7, -3, []int{}},
		{"one", 7, 1, []int{7}},
		{"many", 3, 4, []int{3, 3, 3, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Repeat(tt.value, tt.n)
			streamCheckNonNil(t, "Repeat", got)
			streamCheckEqual(t, "Repeat", got, tt.want)
		})
	}

	t.Run("strings", func(t *testing.T) {
		got := Repeat("ab", 3)
		streamCheckNonNil(t, "Repeat", got)
		streamCheckEqual(t, "Repeat", got, []string{"ab", "ab", "ab"})
	})
}

func TestCycle(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		n     int
		want  []int
	}{
		{"empty items", []int{}, 5, []int{}},
		{"nil items", nil, 5, []int{}},
		{"zero n", []int{1, 2, 3}, 0, []int{}},
		{"negative n", []int{1, 2, 3}, -1, []int{}},
		{"less than one cycle", []int{1, 2, 3}, 2, []int{1, 2}},
		{"exact cycle", []int{1, 2, 3}, 3, []int{1, 2, 3}},
		{"more than one cycle", []int{1, 2, 3}, 8, []int{1, 2, 3, 1, 2, 3, 1, 2}},
		{"many cycles", []int{1, 2, 3}, 10, []int{1, 2, 3, 1, 2, 3, 1, 2, 3, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cycle(tt.items, tt.n)
			streamCheckNonNil(t, "Cycle", got)
			streamCheckEqual(t, "Cycle", got, tt.want)
		})
	}

	t.Run("strings", func(t *testing.T) {
		got := Cycle([]string{"ab", "cd"}, 5)
		streamCheckNonNil(t, "Cycle", got)
		streamCheckEqual(t, "Cycle", got, []string{"ab", "cd", "ab", "cd", "ab"})
	})
}

func TestISlice(t *testing.T) {
	ten := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	tests := []struct {
		name        string
		items       []int
		start, stop int
		step        int
		want        []int
	}{
		{"empty input", []int{}, 0, 5, 1, []int{}},
		{"normal", ten, 2, 6, 1, []int{2, 3, 4, 5}},
		{"with step", ten, 1, 8, 3, []int{1, 4, 7}},
		{"step two", ten, 0, 10, 2, []int{0, 2, 4, 6, 8}},
		{"empty range", ten, 4, 4, 1, []int{}},
		{"stop before start", ten, 5, 3, 1, []int{}},
		{"stop beyond end", ten, 3, 100, 2, []int{3, 5, 7, 9}},
		{"start beyond end", ten, 20, 30, 1, []int{}},
		{"step beyond range", ten, 2, 10, 100, []int{2}},
		{"single element", []int{42}, 0, 1, 1, []int{42}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Islice(tt.items, tt.start, tt.stop, tt.step)
			if err != nil {
				t.Fatalf("Islice returned error: %v", err)
			}
			streamCheckNonNil(t, "Islice", got)
			streamCheckEqual(t, "Islice", got, tt.want)
		})
	}

	t.Run("invalid arguments", func(t *testing.T) {
		tests := []struct {
			name        string
			start, stop int
			step        int
		}{
			{"negative start", -1, 5, 1},
			{"negative stop", 0, -1, 1},
			{"zero step", 0, 5, 0},
			{"negative step", 0, 5, -2},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Islice(ten, tt.start, tt.stop, tt.step)
				if err == nil {
					t.Fatalf("Islice(%d, %d, %d) error = nil, want non-nil", tt.start, tt.stop, tt.step)
				}
				streamCheckNonNil(t, "Islice", got)
				if len(got) != 0 {
					t.Fatalf("Islice on invalid arguments = %v, want empty", got)
				}
			})
		}
	})
}

func TestTakeWhile(t *testing.T) {
	isEven := func(x int) bool { return x%2 == 0 }
	tests := []struct {
		name  string
		items []int
		pred  func(int) bool
		want  []int
	}{
		{"empty", []int{}, isEven, []int{}},
		{"all true", []int{2, 4, 6}, isEven, []int{2, 4, 6}},
		{"all false", []int{1, 3, 5}, isEven, []int{}},
		{"stops in middle", []int{2, 4, 5, 6}, isEven, []int{2, 4}},
		{"first false", []int{1, 2, 4}, isEven, []int{}},
		{"single true", []int{2}, isEven, []int{2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Takewhile(tt.items, tt.pred)
			streamCheckNonNil(t, "Takewhile", got)
			streamCheckEqual(t, "Takewhile", got, tt.want)
		})
	}

	t.Run("strings", func(t *testing.T) {
		got := Takewhile([]string{"a", "bb", "ccc", "d"}, func(s string) bool { return len(s) < 3 })
		streamCheckNonNil(t, "Takewhile", got)
		streamCheckEqual(t, "Takewhile", got, []string{"a", "bb"})
	})

	t.Run("nil predicate", func(t *testing.T) {
		got := Takewhile([]int{1, 2, 3}, nil)
		streamCheckNonNil(t, "Takewhile", got)
		streamCheckEqual(t, "Takewhile", got, []int{})
	})
}

func TestDropWhile(t *testing.T) {
	isEven := func(x int) bool { return x%2 == 0 }
	tests := []struct {
		name  string
		items []int
		pred  func(int) bool
		want  []int
	}{
		{"empty", []int{}, isEven, []int{}},
		{"all true", []int{2, 4, 6}, isEven, []int{}},
		{"all false", []int{1, 3, 5}, isEven, []int{1, 3, 5}},
		{"stops in middle", []int{2, 4, 5, 6}, isEven, []int{5, 6}},
		{"first false", []int{1, 2, 4}, isEven, []int{1, 2, 4}},
		{"single true", []int{2}, isEven, []int{}},
		{"single false", []int{3}, isEven, []int{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Dropwhile(tt.items, tt.pred)
			streamCheckNonNil(t, "Dropwhile", got)
			streamCheckEqual(t, "Dropwhile", got, tt.want)
		})
	}

	t.Run("strings", func(t *testing.T) {
		got := Dropwhile([]string{"a", "bb", "ccc", "d"}, func(s string) bool { return len(s) < 3 })
		streamCheckNonNil(t, "Dropwhile", got)
		streamCheckEqual(t, "Dropwhile", got, []string{"ccc", "d"})
	})

	t.Run("nil predicate", func(t *testing.T) {
		got := Dropwhile([]int{1, 2, 3}, nil)
		streamCheckNonNil(t, "Dropwhile", got)
		streamCheckEqual(t, "Dropwhile", got, []int{1, 2, 3})
	})
}

func TestAccumulate(t *testing.T) {
	add := func(acc, x int) int { return acc + x }
	mul := func(acc, x int) int { return acc * x }
	tests := []struct {
		name  string
		items []int
		fn    func(acc, x int) int
		want  []int
	}{
		{"empty", []int{}, add, []int{}},
		{"single", []int{5}, add, []int{5}},
		{"addition", []int{1, 2, 3, 4}, add, []int{1, 3, 6, 10}},
		{"multiplication", []int{1, 2, 3, 4}, mul, []int{1, 2, 6, 24}},
		{"nil fn", []int{1, 2, 3}, nil, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Accumulate(tt.items, tt.fn)
			streamCheckNonNil(t, "Accumulate", got)
			streamCheckEqual(t, "Accumulate", got, tt.want)
		})
	}

	t.Run("strings", func(t *testing.T) {
		got := Accumulate([]string{"a", "b", "c"}, func(acc, x string) string { return acc + x })
		streamCheckNonNil(t, "Accumulate", got)
		streamCheckEqual(t, "Accumulate", got, []string{"a", "ab", "abc"})
	})
}

func TestPairwise(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  [][2]int
	}{
		{"empty", []int{}, [][2]int{}},
		{"nil", nil, [][2]int{}},
		{"single", []int{1}, [][2]int{}},
		{"two", []int{1, 2}, [][2]int{{1, 2}}},
		{"three", []int{1, 2, 3}, [][2]int{{1, 2}, {2, 3}}},
		{"many", []int{1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pairwise(tt.items)
			streamCheckNonNil(t, "Pairwise", got)
			streamCheckEqual(t, "Pairwise", got, tt.want)
		})
	}

	t.Run("strings", func(t *testing.T) {
		got := Pairwise([]string{"a", "b", "c"})
		streamCheckNonNil(t, "Pairwise", got)
		streamCheckEqual(t, "Pairwise", got, [][2]string{{"a", "b"}, {"b", "c"}})
	})
}

func TestZipLongest(t *testing.T) {
	tests := []struct {
		name string
		fill int
		seqs [][]int
		want [][]int
	}{
		{"no args", 0, nil, [][]int{}},
		{"one empty seq", 0, [][]int{{}}, [][]int{}},
		{"all empty", 0, [][]int{{}, nil}, [][]int{}},
		{"single seq", 0, [][]int{{1, 2}}, [][]int{{1}, {2}}},
		{"equal lengths", 0, [][]int{{1, 2}, {3, 4}}, [][]int{{1, 3}, {2, 4}}},
		{"first longer", 0, [][]int{{1, 2, 3}, {4, 5}}, [][]int{{1, 4}, {2, 5}, {3, 0}}},
		{"second longer", -1, [][]int{{1}, {2, 3}}, [][]int{{1, 2}, {-1, 3}}},
		{"three seqs", -1, [][]int{{1, 2}, {3}, {4, 5, 6}}, [][]int{{1, 3, 4}, {2, -1, 5}, {-1, -1, 6}}},
		{"empty and nonempty", 9, [][]int{{}, {1, 2}}, [][]int{{9, 1}, {9, 2}}},
		{"nil and nonempty", 7, [][]int{nil, {1}}, [][]int{{7, 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Zip_longest(tt.fill, tt.seqs...)
			streamCheckNonNil(t, "Zip_longest", got)
			streamCheckEqual(t, "Zip_longest", got, tt.want)
		})
	}

	t.Run("strings", func(t *testing.T) {
		got := Zip_longest("?", []string{"a", "b"}, []string{"x"})
		streamCheckNonNil(t, "Zip_longest", got)
		streamCheckEqual(t, "Zip_longest", got, [][]string{{"a", "x"}, {"b", "?"}})
	})
}

func TestStreamFunctionsDoNotModifyInputs(t *testing.T) {
	t.Run("Chain", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := Chain(in, []int{4})
		out[0] = 99
		streamCheckEqual(t, "input", in, []int{1, 2, 3})
	})

	t.Run("Cycle", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := Cycle(in, 4)
		out[0] = 99
		streamCheckEqual(t, "input", in, []int{1, 2, 3})
	})

	t.Run("Islice", func(t *testing.T) {
		in := []int{1, 2, 3}
		out, err := Islice(in, 0, 2, 1)
		if err != nil {
			t.Fatalf("Islice returned error: %v", err)
		}
		out[0] = 99
		streamCheckEqual(t, "input", in, []int{1, 2, 3})
	})

	t.Run("Takewhile", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := Takewhile(in, func(int) bool { return true })
		out[0] = 99
		streamCheckEqual(t, "input", in, []int{1, 2, 3})
	})

	t.Run("Dropwhile", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := Dropwhile(in, func(int) bool { return false })
		out[0] = 99
		streamCheckEqual(t, "input", in, []int{1, 2, 3})
	})

	t.Run("Accumulate", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := Accumulate(in, func(acc, x int) int { return acc + x })
		out[0] = 99
		streamCheckEqual(t, "input", in, []int{1, 2, 3})
	})

	t.Run("Pairwise", func(t *testing.T) {
		in := []int{1, 2, 3}
		out := Pairwise(in)
		out[0][0] = 99
		streamCheckEqual(t, "input", in, []int{1, 2, 3})
	})

	t.Run("Zip_longest", func(t *testing.T) {
		a := []int{1, 2}
		b := []int{3}
		out := Zip_longest(0, a, b)
		out[0][0] = 99
		streamCheckEqual(t, "first input", a, []int{1, 2})
		streamCheckEqual(t, "second input", b, []int{3})
	})
}
