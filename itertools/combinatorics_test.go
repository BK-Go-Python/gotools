package itertools

import (
	"reflect"
	"strings"
	"testing"
)

func TestProduct(t *testing.T) {
	tests := []struct {
		name string
		sets [][]int
		want [][]int
	}{
		// Python 3.13: list(itertools.product()) == [()]
		{name: "zero sets", sets: nil, want: [][]int{{}}},
		// Python 3.13: list(itertools.product([1, 2])) == [(1,), (2,)]
		{name: "one set", sets: [][]int{{1, 2}}, want: [][]int{{1}, {2}}},
		// Python 3.13: list(itertools.product([1, 2], [3, 4]))
		{
			name: "two sets",
			sets: [][]int{{1, 2}, {3, 4}},
			want: [][]int{{1, 3}, {1, 4}, {2, 3}, {2, 4}},
		},
		// Python 3.13: list(itertools.product([1, 2], [3], [4, 5]))
		{
			name: "three sets",
			sets: [][]int{{1, 2}, {3}, {4, 5}},
			want: [][]int{{1, 3, 4}, {1, 3, 5}, {2, 3, 4}, {2, 3, 5}},
		},
		// Python 3.13: list(itertools.product([1, 2], [], [3])) == []
		{name: "empty set", sets: [][]int{{1, 2}, {}, {3}}, want: [][]int{}},
		// Python 3.13: list(itertools.product([])) == []
		{name: "only empty set", sets: [][]int{{}}, want: [][]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Product(tt.sets...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Product(%v) = %v, want %v", tt.sets, got, tt.want)
			}
		})
	}
}

func TestProductStrings(t *testing.T) {
	// Python 3.13: list(itertools.product(['ab', 'cd'], ['x', 'y']))
	got := Product([]string{"ab", "cd"}, []string{"x", "y"})
	want := [][]string{
		{"ab", "x"},
		{"ab", "y"},
		{"cd", "x"},
		{"cd", "y"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Product(strings) = %v, want %v", got, want)
	}
}

func TestCombinations(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		r     int
		want  [][]int
	}{
		// Python 3.13: list(itertools.combinations([1, 2, 3], 0)) == [()]
		{name: "r zero", items: []int{1, 2, 3}, r: 0, want: [][]int{{}}},
		// Python 3.13: list(itertools.combinations([1, 2, 3], 1))
		{
			name:  "r one",
			items: []int{1, 2, 3},
			r:     1,
			want:  [][]int{{1}, {2}, {3}},
		},
		// Python 3.13: list(itertools.combinations([1, 2, 3, 4], 2))
		{
			name:  "r two",
			items: []int{1, 2, 3, 4},
			r:     2,
			want:  [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}},
		},
		// Python 3.13: list(itertools.combinations([1, 2, 3], 3)) == [(1, 2, 3)]
		{name: "r equals len", items: []int{1, 2, 3}, r: 3, want: [][]int{{1, 2, 3}}},
		// Python 3.13: list(itertools.combinations([1, 2, 3], 4)) == []
		{name: "r greater than len", items: []int{1, 2, 3}, r: 4, want: [][]int{}},
		// Python 3.13: list(itertools.combinations([1, 2, 3], -1)) == []
		{name: "negative r", items: []int{1, 2, 3}, r: -1, want: [][]int{}},
		// Python 3.13: list(itertools.combinations([1, 1, 2], 2))
		{
			name:  "duplicate elements",
			items: []int{1, 1, 2},
			r:     2,
			want:  [][]int{{1, 1}, {1, 2}, {1, 2}},
		},
		// Python 3.13: list(itertools.combinations([1, 1], 2)) == [(1, 1)]
		{name: "all duplicate elements", items: []int{1, 1}, r: 2, want: [][]int{{1, 1}}},
		// Python 3.13: list(itertools.combinations([], 0)) == [()]
		{name: "empty items r zero", items: []int{}, r: 0, want: [][]int{{}}},
		// Python 3.13: list(itertools.combinations([], 1)) == []
		{name: "empty items r one", items: []int{}, r: 1, want: [][]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Combinations(tt.items, tt.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Combinations(%v, %d) = %v, want %v", tt.items, tt.r, got, tt.want)
			}
		})
	}
}

func TestPermutations(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		r     int
		want  [][]int
	}{
		// Python 3.13: list(itertools.permutations([1, 2, 3], 0)) == [()]
		{name: "r zero", items: []int{1, 2, 3}, r: 0, want: [][]int{{}}},
		// Python 3.13: list(itertools.permutations([1, 2, 3], 1))
		{
			name:  "r one",
			items: []int{1, 2, 3},
			r:     1,
			want:  [][]int{{1}, {2}, {3}},
		},
		// Python 3.13: list(itertools.permutations([1, 2, 3], 2))
		{
			name:  "r two",
			items: []int{1, 2, 3},
			r:     2,
			want:  [][]int{{1, 2}, {1, 3}, {2, 1}, {2, 3}, {3, 1}, {3, 2}},
		},
		// Python 3.13: list(itertools.permutations([1, 2, 3], 3))
		{
			name:  "r equals len",
			items: []int{1, 2, 3},
			r:     3,
			want: [][]int{
				{1, 2, 3}, {1, 3, 2},
				{2, 1, 3}, {2, 3, 1},
				{3, 1, 2}, {3, 2, 1},
			},
		},
		// Python 3.13: list(itertools.permutations([1, 2, 3], 4)) == []
		{name: "r greater than len", items: []int{1, 2, 3}, r: 4, want: [][]int{}},
		// Python 3.13: list(itertools.permutations([1, 2, 3], -1)) == []
		{name: "negative r", items: []int{1, 2, 3}, r: -1, want: [][]int{}},
		// Python 3.13: list(itertools.permutations([1, 1, 2], 2))
		{
			name:  "duplicate elements",
			items: []int{1, 1, 2},
			r:     2,
			want:  [][]int{{1, 1}, {1, 2}, {1, 1}, {1, 2}, {2, 1}, {2, 1}},
		},
		// Python 3.13: list(itertools.permutations([1, 1], 2))
		{name: "all duplicate elements", items: []int{1, 1}, r: 2, want: [][]int{{1, 1}, {1, 1}}},
		// Python 3.13: list(itertools.permutations([], 0)) == [()]
		{name: "empty items r zero", items: []int{}, r: 0, want: [][]int{{}}},
		// Python 3.13: list(itertools.permutations([], 1)) == []
		{name: "empty items r one", items: []int{}, r: 1, want: [][]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Permutations(tt.items, tt.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Permutations(%v, %d) = %v, want %v", tt.items, tt.r, got, tt.want)
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	identity := func(v int) int { return v }
	parity := func(v int) int { return v % 2 }
	tests := []struct {
		name      string
		items     []int
		key       func(int) int
		want      []Group[int, int]
		wantCalls int
	}{
		// Python 3.13: [(k, list(g)) for k, g in itertools.groupby([1, 1, 2, 2, 3])]
		{
			name:  "consecutive groups",
			items: []int{1, 1, 2, 2, 3},
			key:   identity,
			want: []Group[int, int]{
				{Key: 1, Values: []int{1, 1}},
				{Key: 2, Values: []int{2, 2}},
				{Key: 3, Values: []int{3}},
			},
			wantCalls: 5,
		},
		// Python 3.13: [(k, list(g)) for k, g in itertools.groupby([1, 2, 1])]
		{
			name:  "same key not adjacent",
			items: []int{1, 2, 1},
			key:   identity,
			want: []Group[int, int]{
				{Key: 1, Values: []int{1}},
				{Key: 2, Values: []int{2}},
				{Key: 1, Values: []int{1}},
			},
			wantCalls: 3,
		},
		// Python 3.13: [(k, list(g)) for k, g in itertools.groupby([])] == []
		{
			name:      "empty input",
			items:     []int{},
			key:       identity,
			want:      []Group[int, int]{},
			wantCalls: 0,
		},
		// Python 3.13: keys 1, 0, 0, 1; key called once per element.
		{
			name:  "key called once per element",
			items: []int{1, 2, 2, 3},
			key:   parity,
			want: []Group[int, int]{
				{Key: 1, Values: []int{1}},
				{Key: 0, Values: []int{2, 2}},
				{Key: 1, Values: []int{3}},
			},
			wantCalls: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			key := func(v int) int {
				calls++
				return tt.key(v)
			}
			got := Groupby(tt.items, key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Groupby(%v) = %v, want %v", tt.items, got, tt.want)
			}
			if calls != tt.wantCalls {
				t.Fatalf("key called %d times, want %d", calls, tt.wantCalls)
			}
		})
	}
}

func TestGroupByStrings(t *testing.T) {
	// Python 3.13: [(k, list(g)) for k, g in
	//     itertools.groupby(['a', 'A', 'b', 'B', 'b'], key=str.lower)]
	items := []string{"a", "A", "b", "B", "b"}
	got := Groupby(items, strings.ToLower)
	want := []Group[string, string]{
		{Key: "a", Values: []string{"a", "A"}},
		{Key: "b", Values: []string{"b", "B", "b"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Groupby(strings) = %v, want %v", got, want)
	}
}

func TestCombinatoricsNoAliasing(t *testing.T) {
	t.Run("Product", func(t *testing.T) {
		first := []int{1, 2}
		second := []int{3, 4}
		got := Product(first, second)
		got[0][0] = 99
		if !reflect.DeepEqual(first, []int{1, 2}) || !reflect.DeepEqual(second, []int{3, 4}) {
			t.Fatalf("Product modified input: first=%v second=%v", first, second)
		}
		if got[1][0] != 1 {
			t.Fatalf("Product rows share backing arrays: got[1][0] = %d, want 1", got[1][0])
		}
	})

	t.Run("Combinations", func(t *testing.T) {
		items := []int{1, 2, 3}
		got := Combinations(items, 2)
		got[0][0] = 99
		if !reflect.DeepEqual(items, []int{1, 2, 3}) {
			t.Fatalf("Combinations modified input: %v", items)
		}
		if got[1][0] != 1 {
			t.Fatalf("Combinations rows share backing arrays: got[1][0] = %d, want 1", got[1][0])
		}
	})

	t.Run("Permutations", func(t *testing.T) {
		items := []int{1, 2}
		got := Permutations(items, 2)
		got[0][0] = 99
		if !reflect.DeepEqual(items, []int{1, 2}) {
			t.Fatalf("Permutations modified input: %v", items)
		}
		if got[1][0] != 2 {
			t.Fatalf("Permutations rows share backing arrays: got[1][0] = %d, want 2", got[1][0])
		}
	})

	t.Run("Groupby", func(t *testing.T) {
		items := []int{1, 1, 2, 2}
		got := Groupby(items, func(v int) int { return v })
		got[0].Values[0] = 99
		if !reflect.DeepEqual(items, []int{1, 1, 2, 2}) {
			t.Fatalf("Groupby modified input: %v", items)
		}
	})
}
