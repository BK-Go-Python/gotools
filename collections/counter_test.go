package collections

import (
	"slices"
	"testing"
)

func TestCounterNew(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  Counter[string]
	}{
		{name: "empty", items: nil, want: Counter[string]{}},
		{name: "single", items: []string{"x"}, want: Counter[string]{"x": 1}},
		{name: "duplicates", items: []string{"a", "b", "a"}, want: Counter[string]{"a": 2, "b": 1}},
		{name: "all same", items: []string{"x", "x", "x"}, want: Counter[string]{"x": 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCounter(tt.items...)
			if !got.Equal(tt.want) {
				t.Fatalf("NewCounter(%v) = %v, want %v", tt.items, got, tt.want)
			}
		})
	}
}

func TestCounterZeroValue(t *testing.T) {
	t.Run("queries", func(t *testing.T) {
		var c Counter[string]
		if got := c.Count("missing"); got != 0 {
			t.Fatalf("Count(missing) = %d, want 0", got)
		}
		if got := c.Total(); got != 0 {
			t.Fatalf("Total() = %d, want 0", got)
		}
		if got := c.MostCommon(); len(got) != 0 {
			t.Fatalf("MostCommon() = %v, want empty", got)
		}
		if got := c.Elements(); len(got) != 0 {
			t.Fatalf("Elements() = %v, want empty", got)
		}
		if !c.Equal(Counter[string]{}) {
			t.Fatal("zero Counter should equal empty Counter")
		}
	})

	t.Run("add allocates", func(t *testing.T) {
		var c Counter[string]
		c = c.Add("a")
		if c == nil {
			t.Fatal("Add on nil Counter should allocate a map")
		}
		if got := c.Count("a"); got != 1 {
			t.Fatalf("Count(a) = %d, want 1", got)
		}
	})
}

func TestCounterAdd(t *testing.T) {
	tests := []struct {
		name  string
		start Counter[string]
		item  string
		n     []int
		want  Counter[string]
	}{
		{
			name:  "default increment on nil",
			start: nil,
			item:  "a",
			want:  Counter[string]{"a": 1},
		},
		{
			name:  "explicit positive",
			start: Counter[string]{"a": 1},
			item:  "a",
			n:     []int{4},
			want:  Counter[string]{"a": 5},
		},
		{
			name:  "negative",
			start: Counter[string]{"a": 2},
			item:  "a",
			n:     []int{-5},
			want:  Counter[string]{"a": -3},
		},
		{
			name:  "negative on nil",
			start: nil,
			item:  "z",
			n:     []int{-2},
			want:  Counter[string]{"z": -2},
		},
		{
			name:  "zero creates key",
			start: NewCounter("a"),
			item:  "b",
			n:     []int{0},
			want:  Counter[string]{"a": 1, "b": 0},
		},
		{
			name:  "only first n used",
			start: nil,
			item:  "a",
			n:     []int{2, 100},
			want:  Counter[string]{"a": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.start.Add(tt.item, tt.n...)
			if !got.Equal(tt.want) {
				t.Fatalf("Add(%q, %v) = %v, want %v", tt.item, tt.n, got, tt.want)
			}
		})
	}
}

func TestCounterCount(t *testing.T) {
	tests := []struct {
		name string
		c    Counter[string]
		item string
		want int
	}{
		{name: "present", c: NewCounter("a", "b", "a"), item: "a", want: 2},
		{name: "missing", c: NewCounter("a"), item: "z", want: 0},
		{name: "nil receiver", c: nil, item: "a", want: 0},
		{name: "negative", c: Counter[string]{"x": -3}, item: "x", want: -3},
		{name: "zero count", c: Counter[string]{"x": 0}, item: "x", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Count(tt.item); got != tt.want {
				t.Fatalf("Count(%q) = %d, want %d", tt.item, got, tt.want)
			}
		})
	}
}

func TestCounterTotal(t *testing.T) {
	tests := []struct {
		name string
		c    Counter[string]
		want int
	}{
		{name: "nil", c: nil, want: 0},
		{name: "empty", c: Counter[string]{}, want: 0},
		{name: "duplicates", c: NewCounter("a", "b", "a"), want: 3},
		{name: "zero count", c: Counter[string]{"a": 0, "b": 2}, want: 2},
		{name: "negative counts", c: Counter[string]{"a": 2, "b": -1}, want: 1},
		{name: "all negative", c: Counter[string]{"a": -2, "b": -3}, want: -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Total(); got != tt.want {
				t.Fatalf("Total() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCounterMostCommon(t *testing.T) {
	// Counter("abracadabra") in Python 3.13:
	// [('a', 5), ('b', 2), ('r', 2), ('c', 1), ('d', 1)]
	letters := Counter[string]{"a": 5, "b": 2, "r": 2, "c": 1, "d": 1}
	allLetters := []CountEntry[string]{
		{Value: "a", Count: 5},
		{Value: "b", Count: 2},
		{Value: "r", Count: 2},
		{Value: "c", Count: 1},
		{Value: "d", Count: 1},
	}

	tests := []struct {
		name string
		c    Counter[string]
		n    []int
		want []CountEntry[string]
	}{
		{name: "all", c: letters, want: allLetters},
		{name: "top 1", c: letters, n: []int{1}, want: allLetters[:1]},
		{name: "top 3", c: letters, n: []int{3}, want: allLetters[:3]},
		{name: "zero", c: letters, n: []int{0}, want: []CountEntry[string]{}},
		{name: "negative means all", c: letters, n: []int{-1}, want: allLetters},
		{name: "n larger than keys", c: letters, n: []int{99}, want: allLetters},
		{name: "only first n used", c: letters, n: []int{2, 99}, want: allLetters[:2]},
		{name: "empty", c: Counter[string]{}, want: []CountEntry[string]{}},
		{name: "nil", c: nil, want: []CountEntry[string]{}},
		{
			name: "ties sort by fmt.Sprint",
			c:    NewCounter("pear", "apple", "fig"),
			want: []CountEntry[string]{
				{Value: "apple", Count: 1},
				{Value: "fig", Count: 1},
				{Value: "pear", Count: 1},
			},
		},
		{
			name: "negative counts still sorted descending",
			c:    Counter[string]{"a": -1, "b": 2, "c": 0, "d": -2},
			want: []CountEntry[string]{
				{Value: "b", Count: 2},
				{Value: "c", Count: 0},
				{Value: "a", Count: -1},
				{Value: "d", Count: -2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.MostCommon(tt.n...)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("MostCommon(%v) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}

	t.Run("integer ties use fmt.Sprint lexicographic order", func(t *testing.T) {
		c := Counter[int]{2: 1, 10: 1}
		want := []CountEntry[int]{
			{Value: 10, Count: 1},
			{Value: 2, Count: 1},
		}
		if got := c.MostCommon(); !slices.Equal(got, want) {
			t.Fatalf("MostCommon() = %v, want %v", got, want)
		}
	})
}

func TestCounterElements(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		tests := []struct {
			name string
			c    Counter[string]
			want []string
		}{
			{name: "nil", c: nil, want: []string{}},
			{name: "empty", c: Counter[string]{}, want: []string{}},
			{name: "duplicates", c: NewCounter("a", "b", "a"), want: []string{"a", "a", "b"}},
			{
				name: "zero and negative counts skipped",
				c:    Counter[string]{"a": 2, "b": 1, "c": 0, "d": -1},
				want: []string{"a", "a", "b"},
			},
			{
				name: "sorted by fmt.Sprint",
				c:    Counter[string]{"z": 1, "a": 2, "m": 1},
				want: []string{"a", "a", "m", "z"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.c.Elements()
				if !slices.Equal(got, tt.want) {
					t.Fatalf("Elements() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("ints", func(t *testing.T) {
		tests := []struct {
			name string
			c    Counter[int]
			want []int
		}{
			{name: "positive", c: Counter[int]{3: 1, 1: 2, 2: 0}, want: []int{1, 1, 3}},
			{name: "fmt.Sprint is lexicographic", c: Counter[int]{2: 1, 10: 1}, want: []int{10, 2}},
			{name: "negative counts", c: Counter[int]{1: -3, 2: 2}, want: []int{2, 2}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.c.Elements()
				if !slices.Equal(got, tt.want) {
					t.Fatalf("Elements() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestCounterUpdate(t *testing.T) {
	tests := []struct {
		name  string
		start Counter[string]
		items []string
		want  Counter[string]
	}{
		{
			name:  "nil allocates",
			start: nil,
			items: []string{"a", "b", "a"},
			want:  Counter[string]{"a": 2, "b": 1},
		},
		{
			name:  "existing counts increase",
			start: Counter[string]{"a": 2, "b": -1},
			items: []string{"a", "b"},
			want:  Counter[string]{"a": 3, "b": 0},
		},
		{
			name:  "no items leaves counts unchanged",
			start: NewCounter("a", "a"),
			want:  Counter[string]{"a": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.start.Update(tt.items...)
			if !got.Equal(tt.want) {
				t.Fatalf("Update(%v) = %v, want %v", tt.items, got, tt.want)
			}
		})
	}

	t.Run("nil with no items still allocates", func(t *testing.T) {
		var c Counter[string]
		got := c.Update()
		if got == nil {
			t.Fatal("Update on nil Counter should allocate a map")
		}
		if len(got) != 0 {
			t.Fatalf("Update() = %v, want empty", got)
		}
	})
}

func TestCounterSubtract(t *testing.T) {
	tests := []struct {
		name  string
		start Counter[string]
		items []string
		want  Counter[string]
	}{
		{
			name:  "nil allocates negative counts",
			start: nil,
			items: []string{"x", "x"},
			want:  Counter[string]{"x": -2},
		},
		{
			name:  "existing counts decrease",
			start: NewCounter("a", "a", "b"),
			items: []string{"a", "a", "a", "b"},
			want:  Counter[string]{"a": -1, "b": 0},
		},
		{
			name:  "negative counts keep decreasing",
			start: Counter[string]{"a": -2},
			items: []string{"a"},
			want:  Counter[string]{"a": -3},
		},
		{
			name:  "no items leaves counts unchanged",
			start: NewCounter("a"),
			want:  Counter[string]{"a": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.start.Subtract(tt.items...)
			if !got.Equal(tt.want) {
				t.Fatalf("Subtract(%v) = %v, want %v", tt.items, got, tt.want)
			}
		})
	}

	t.Run("total includes negative counts", func(t *testing.T) {
		got := NewCounter("a", "a", "b").Subtract("a", "a", "a")
		if got.Count("a") != -1 || got.Count("b") != 1 {
			t.Fatalf("after Subtract: %v", got)
		}
		if got.Total() != 0 {
			t.Fatalf("Total() = %d, want 0", got.Total())
		}
	})
}

func TestCounterCopy(t *testing.T) {
	t.Run("independent maps", func(t *testing.T) {
		original := NewCounter("a", "a", "b")
		clone := original.Copy()
		if !clone.Equal(original) {
			t.Fatalf("Copy() = %v, want %v", clone, original)
		}

		clone = clone.Add("a", 5).Add("c")
		if original.Count("a") != 2 || original.Count("c") != 0 {
			t.Fatalf("modifying copy changed original: %v", original)
		}

		original = original.Subtract("b")
		if clone.Count("b") != 1 {
			t.Fatalf("modifying original changed copy: %v", clone)
		}
	})

	t.Run("nil", func(t *testing.T) {
		var c Counter[int]
		clone := c.Copy()
		if clone == nil {
			t.Fatal("Copy of nil Counter should be a usable empty map")
		}
		if clone.Total() != 0 {
			t.Fatalf("Copy of nil Total() = %d, want 0", clone.Total())
		}
	})
}

func TestCounterEqual(t *testing.T) {
	tests := []struct {
		name string
		a    Counter[string]
		b    Counter[string]
		want bool
	}{
		{name: "both nil", a: nil, b: nil, want: true},
		{name: "nil and empty", a: nil, b: Counter[string]{}, want: true},
		{name: "same counts", a: NewCounter("a", "b", "a"), b: Counter[string]{"a": 2, "b": 1}, want: true},
		{name: "different counts", a: Counter[string]{"a": 2}, b: Counter[string]{"a": 3}, want: false},
		{name: "different keys", a: Counter[string]{"a": 1}, b: Counter[string]{"b": 1}, want: false},
		{name: "negative counts equal", a: Counter[string]{"a": -1, "b": 2}, b: Counter[string]{"a": -1, "b": 2}, want: true},
		{name: "zero count versus missing key", a: Counter[string]{"a": 0}, b: Counter[string]{}, want: false},
		{name: "subset", a: Counter[string]{"a": 1}, b: Counter[string]{"a": 1, "b": 0}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.b); got != tt.want {
				t.Fatalf("%v.Equal(%v) = %t, want %t", tt.a, tt.b, got, tt.want)
			}
			if got := tt.b.Equal(tt.a); got != tt.want {
				t.Fatalf("%v.Equal(%v) = %t, want %t (symmetry)", tt.b, tt.a, got, tt.want)
			}
		})
	}
}
