package collections

import (
	"slices"
	"testing"
)

func TestNewDeque(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  []int
	}{
		{name: "empty", items: nil, want: []int{}},
		{name: "single", items: []int{42}, want: []int{42}},
		{name: "multiple", items: []int{1, 2, 3}, want: []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeque(tt.items...)
			if got := d.Len(); got != len(tt.want) {
				t.Fatalf("Len() = %d, want %d", got, len(tt.want))
			}
			if got := d.Items(); !slices.Equal(got, tt.want) {
				t.Fatalf("Items() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("zero value usable", func(t *testing.T) {
		var d Deque[int]
		d.Append(2)
		d.AppendLeft(1)
		d.Append(3)
		if got := d.Items(); !slices.Equal(got, []int{1, 2, 3}) {
			t.Fatalf("Items() = %v, want [1 2 3]", got)
		}
	})
}

func TestDequeAppend(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		left    bool
		items   []int
		want    []int
	}{
		{name: "append to empty", items: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{name: "append to existing", initial: []int{1, 2}, items: []int{3, 4}, want: []int{1, 2, 3, 4}},
		{name: "append no items", initial: []int{1, 2}, items: nil, want: []int{1, 2}},
		{
			name:  "append grows buffer",
			items: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			want:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{name: "appendleft to empty", left: true, items: []int{1, 2, 3}, want: []int{3, 2, 1}},
		{name: "appendleft to existing", initial: []int{1, 2}, left: true, items: []int{3, 4}, want: []int{4, 3, 1, 2}},
		{name: "appendleft no items", initial: []int{1, 2}, left: true, items: nil, want: []int{1, 2}},
		{
			name:  "appendleft grows buffer",
			left:  true,
			items: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			want:  []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeque(tt.initial...)
			if tt.left {
				d.AppendLeft(tt.items...)
			} else {
				d.Append(tt.items...)
			}
			if got := d.Items(); !slices.Equal(got, tt.want) {
				t.Fatalf("Items() = %v, want %v", got, tt.want)
			}
			if got := d.Len(); got != len(tt.want) {
				t.Fatalf("Len() = %d, want %d", got, len(tt.want))
			}
		})
	}
}

func TestDequeExtend(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		left    bool
		items   []int
		want    []int
	}{
		{name: "extend empty", items: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{name: "extend existing", initial: []int{1, 2}, items: []int{3, 4}, want: []int{1, 2, 3, 4}},
		{name: "extend no items", initial: []int{1, 2}, items: nil, want: []int{1, 2}},
		{name: "extendleft empty", left: true, items: []int{1, 2, 3}, want: []int{3, 2, 1}},
		{name: "extendleft existing", initial: []int{1, 2}, left: true, items: []int{3, 4}, want: []int{4, 3, 1, 2}},
		{name: "extendleft no items", initial: []int{1, 2}, left: true, items: nil, want: []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeque(tt.initial...)
			if tt.left {
				d.ExtendLeft(tt.items...)
			} else {
				d.Extend(tt.items...)
			}
			if got := d.Items(); !slices.Equal(got, tt.want) {
				t.Fatalf("Items() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDequePopPeek(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		prep      func(*Deque[int])
		op        func(*Deque[int]) (int, bool)
		want      int
		wantOK    bool
		wantItems []int
	}{
		{name: "pop empty", initial: nil, op: (*Deque[int]).Pop, wantOK: false, wantItems: []int{}},
		{name: "popleft empty", initial: nil, op: (*Deque[int]).PopLeft, wantOK: false, wantItems: []int{}},
		{name: "peek empty", initial: nil, op: (*Deque[int]).Peek, wantOK: false, wantItems: []int{}},
		{name: "peekleft empty", initial: nil, op: (*Deque[int]).PeekLeft, wantOK: false, wantItems: []int{}},
		{name: "pop", initial: []int{1, 2, 3}, op: (*Deque[int]).Pop, want: 3, wantOK: true, wantItems: []int{1, 2}},
		{name: "popleft", initial: []int{1, 2, 3}, op: (*Deque[int]).PopLeft, want: 1, wantOK: true, wantItems: []int{2, 3}},
		{name: "peek", initial: []int{1, 2, 3}, op: (*Deque[int]).Peek, want: 3, wantOK: true, wantItems: []int{1, 2, 3}},
		{name: "peekleft", initial: []int{1, 2, 3}, op: (*Deque[int]).PeekLeft, want: 1, wantOK: true, wantItems: []int{1, 2, 3}},
		{name: "pop single", initial: []int{7}, op: (*Deque[int]).Pop, want: 7, wantOK: true, wantItems: []int{}},
		{name: "popleft single", initial: []int{7}, op: (*Deque[int]).PopLeft, want: 7, wantOK: true, wantItems: []int{}},
		{
			name:      "pop after popleft wraps head",
			initial:   []int{1, 2, 3, 4},
			prep:      func(d *Deque[int]) { d.PopLeft() },
			op:        (*Deque[int]).Pop,
			want:      4,
			wantOK:    true,
			wantItems: []int{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeque(tt.initial...)
			if tt.prep != nil {
				tt.prep(d)
			}
			got, ok := tt.op(d)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("op() = (%d, %v), want (%d, %v)", got, ok, tt.want, tt.wantOK)
			}
			if gotItems := d.Items(); !slices.Equal(gotItems, tt.wantItems) {
				t.Fatalf("Items() = %v, want %v", gotItems, tt.wantItems)
			}
		})
	}
}

func TestDequeDrain(t *testing.T) {
	tests := []struct {
		name    string
		popLeft bool
		want    []int
	}{
		{name: "pop right", popLeft: false, want: []int{5, 4, 3, 2, 1}},
		{name: "popleft", popLeft: true, want: []int{1, 2, 3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeque(1, 2, 3, 4, 5)
			for i, want := range tt.want {
				var got int
				var ok bool
				if tt.popLeft {
					got, ok = d.PopLeft()
				} else {
					got, ok = d.Pop()
				}
				if !ok || got != want {
					t.Fatalf("step %d: got (%d, %v), want (%d, true)", i, got, ok, want)
				}
				if gotLen := d.Len(); gotLen != len(tt.want)-i-1 {
					t.Fatalf("step %d: Len() = %d, want %d", i, gotLen, len(tt.want)-i-1)
				}
			}
			if got := d.Len(); got != 0 {
				t.Fatalf("Len() = %d, want 0", got)
			}
			if _, ok := d.Pop(); ok {
				t.Fatal("Pop() on drained deque returned true")
			}
			if _, ok := d.PopLeft(); ok {
				t.Fatal("PopLeft() on drained deque returned true")
			}
			if _, ok := d.Peek(); ok {
				t.Fatal("Peek() on drained deque returned true")
			}
			if _, ok := d.PeekLeft(); ok {
				t.Fatal("PeekLeft() on drained deque returned true")
			}

			// 清空后应可继续使用。
			d.Append(9)
			if got, ok := d.PeekLeft(); !ok || got != 9 {
				t.Fatalf("PeekLeft() after reuse = (%d, %v), want (9, true)", got, ok)
			}
		})
	}
}

func TestDequeRotate(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		prep    func(*Deque[int])
		n       int
		want    []int
	}{
		{name: "zero", initial: []int{1, 2, 3, 4, 5}, n: 0, want: []int{1, 2, 3, 4, 5}},
		{name: "right 1", initial: []int{1, 2, 3, 4, 5}, n: 1, want: []int{5, 1, 2, 3, 4}},
		{name: "right 2", initial: []int{1, 2, 3, 4, 5}, n: 2, want: []int{4, 5, 1, 2, 3}},
		{name: "right equal length", initial: []int{1, 2, 3, 4, 5}, n: 5, want: []int{1, 2, 3, 4, 5}},
		{name: "right greater than length", initial: []int{1, 2, 3, 4, 5}, n: 7, want: []int{4, 5, 1, 2, 3}},
		{name: "left 1", initial: []int{1, 2, 3, 4, 5}, n: -1, want: []int{2, 3, 4, 5, 1}},
		{name: "left 2", initial: []int{1, 2, 3, 4, 5}, n: -2, want: []int{3, 4, 5, 1, 2}},
		{name: "left equal length", initial: []int{1, 2, 3, 4, 5}, n: -5, want: []int{1, 2, 3, 4, 5}},
		{name: "left greater than length", initial: []int{1, 2, 3, 4, 5}, n: -7, want: []int{3, 4, 5, 1, 2}},
		{name: "empty positive", initial: nil, n: 3, want: []int{}},
		{name: "empty negative", initial: nil, n: -3, want: []int{}},
		{name: "single positive", initial: []int{42}, n: 3, want: []int{42}},
		{name: "single negative", initial: []int{42}, n: -9, want: []int{42}},
		{name: "two elements right", initial: []int{1, 2}, n: 1, want: []int{2, 1}},
		{name: "two elements left", initial: []int{1, 2}, n: -1, want: []int{2, 1}},
		{
			name:    "wrapped buffer",
			initial: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			prep: func(d *Deque[int]) {
				for i := 0; i < 5; i++ {
					d.PopLeft()
				}
				d.Append(10, 11, 12)
			},
			n:    3,
			want: []int{10, 11, 12, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeque(tt.initial...)
			if tt.prep != nil {
				tt.prep(d)
			}
			d.Rotate(tt.n)
			if got := d.Items(); !slices.Equal(got, tt.want) {
				t.Fatalf("Items() after Rotate(%d) = %v, want %v", tt.n, got, tt.want)
			}
			if got := d.Len(); got != len(tt.want) {
				t.Fatalf("Len() = %d, want %d", got, len(tt.want))
			}
		})
	}
}

func TestDequeItems(t *testing.T) {
	tests := []struct {
		name  string
		build func() *Deque[int]
		want  []int
	}{
		{
			name:  "empty",
			build: func() *Deque[int] { return NewDeque[int]() },
			want:  []int{},
		},
		{
			name: "zero value",
			build: func() *Deque[int] {
				var d Deque[int]
				d.Append(2, 3)
				d.AppendLeft(1)
				return &d
			},
			want: []int{1, 2, 3},
		},
		{
			name: "wrapped buffer",
			build: func() *Deque[int] {
				d := NewDeque(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
				for i := 0; i < 5; i++ {
					d.PopLeft()
				}
				d.Append(10, 11)
				return d
			},
			want: []int{5, 6, 7, 8, 9, 10, 11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.build()
			if got := d.Items(); !slices.Equal(got, tt.want) {
				t.Fatalf("Items() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("empty is non-nil", func(t *testing.T) {
		var nilDeque *Deque[int]
		for name, got := range map[string][]int{
			"nil deque":  nilDeque.Items(),
			"new deque":  NewDeque[int]().Items(),
			"zero value": (&Deque[int]{}).Items(),
		} {
			if got == nil {
				t.Fatalf("%s: Items() = nil, want non-nil empty slice", name)
			}
			if len(got) != 0 {
				t.Fatalf("%s: Items() = %v, want empty", name, got)
			}
		}
	})

	t.Run("returns independent slice", func(t *testing.T) {
		d := NewDeque(1, 2, 3)
		items := d.Items()
		items[0] = 99
		d.Append(4)

		if got := d.Items(); !slices.Equal(got, []int{1, 2, 3, 4}) {
			t.Fatalf("Items() = %v, want [1 2 3 4]", got)
		}
		if !slices.Equal(items, []int{99, 2, 3}) {
			t.Fatalf("snapshot = %v, want [99 2 3]", items)
		}
	})
}

func TestDequeClear(t *testing.T) {
	tests := []struct {
		name  string
		build func() *Deque[int]
	}{
		{name: "empty", build: func() *Deque[int] { return NewDeque[int]() }},
		{name: "with items", build: func() *Deque[int] { return NewDeque(1, 2, 3) }},
		{name: "zero value", build: func() *Deque[int] { return &Deque[int]{} }},
		{
			name: "wrapped buffer",
			build: func() *Deque[int] {
				d := NewDeque(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
				for i := 0; i < 4; i++ {
					d.PopLeft()
				}
				d.Append(10, 11)
				return d
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.build()
			d.Clear()
			if got := d.Len(); got != 0 {
				t.Fatalf("Len() = %d, want 0", got)
			}
			if got := d.Items(); len(got) != 0 {
				t.Fatalf("Items() = %v, want empty", got)
			}
			if _, ok := d.Peek(); ok {
				t.Fatal("Peek() after Clear returned true")
			}
			if _, ok := d.PeekLeft(); ok {
				t.Fatal("PeekLeft() after Clear returned true")
			}
			if _, ok := d.Pop(); ok {
				t.Fatal("Pop() after Clear returned true")
			}
			if _, ok := d.PopLeft(); ok {
				t.Fatal("PopLeft() after Clear returned true")
			}

			d.Clear() // 重复清空应无副作用。
			d.Append(7)
			d.AppendLeft(6)
			if got := d.Items(); !slices.Equal(got, []int{6, 7}) {
				t.Fatalf("Items() after reuse = %v, want [6 7]", got)
			}
		})
	}
}

func TestDequeCopy(t *testing.T) {
	tests := []struct {
		name  string
		build func() *Deque[int]
		want  []int
	}{
		{name: "empty", build: func() *Deque[int] { return NewDeque[int]() }, want: []int{}},
		{name: "with items", build: func() *Deque[int] { return NewDeque(1, 2, 3) }, want: []int{1, 2, 3}},
		{
			name: "wrapped buffer",
			build: func() *Deque[int] {
				d := NewDeque(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
				for i := 0; i < 5; i++ {
					d.PopLeft()
				}
				d.Append(100, 101)
				return d
			},
			want: []int{5, 6, 7, 8, 9, 100, 101},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.build()
			copied := original.Copy()
			if copied == nil {
				t.Fatal("Copy() = nil, want non-nil deque")
			}
			if copied == original {
				t.Fatal("Copy() returned the same pointer, want independent deque")
			}
			if got := copied.Items(); !slices.Equal(got, tt.want) {
				t.Fatalf("copied Items() = %v, want %v", got, tt.want)
			}
			if got := original.Items(); !slices.Equal(got, tt.want) {
				t.Fatalf("original Items() = %v, want %v", got, tt.want)
			}

			// 修改副本不影响原队列。
			copied.Append(901)
			copied.AppendLeft(900)
			if _, ok := copied.Pop(); !ok {
				t.Fatal("Pop() on copy returned false")
			}
			wantCopied := append([]int{900}, tt.want...)
			if got := copied.Items(); !slices.Equal(got, wantCopied) {
				t.Fatalf("copied Items() = %v, want %v", got, wantCopied)
			}
			if got := original.Items(); !slices.Equal(got, tt.want) {
				t.Fatalf("original Items() after copy mutation = %v, want %v", got, tt.want)
			}

			// 修改原队列也不影响副本。
			original.Append(902)
			original.Rotate(1)
			if got := copied.Items(); !slices.Equal(got, wantCopied) {
				t.Fatalf("copied Items() after original mutation = %v, want %v", got, wantCopied)
			}
		})
	}
}

func TestDequeNilSafety(t *testing.T) {
	var d *Deque[int]

	t.Run("read methods", func(t *testing.T) {
		if got := d.Len(); got != 0 {
			t.Fatalf("Len() = %d, want 0", got)
		}
		items := d.Items()
		if items == nil || len(items) != 0 {
			t.Fatalf("Items() = %v, want non-nil empty slice", items)
		}

		tests := []struct {
			name string
			op   func(*Deque[int]) (int, bool)
		}{
			{name: "Pop", op: (*Deque[int]).Pop},
			{name: "PopLeft", op: (*Deque[int]).PopLeft},
			{name: "Peek", op: (*Deque[int]).Peek},
			{name: "PeekLeft", op: (*Deque[int]).PeekLeft},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got, ok := tt.op(d); got != 0 || ok {
					t.Fatalf("%s() = (%d, %v), want (0, false)", tt.name, got, ok)
				}
			})
		}
	})

	t.Run("mutating methods are no-ops", func(t *testing.T) {
		tests := []struct {
			name string
			run  func()
		}{
			{name: "Append", run: func() { d.Append(1) }},
			{name: "AppendLeft", run: func() { d.AppendLeft(1) }},
			{name: "Extend", run: func() { d.Extend(1, 2) }},
			{name: "ExtendLeft", run: func() { d.ExtendLeft(1, 2) }},
			{name: "Rotate positive", run: func() { d.Rotate(5) }},
			{name: "Rotate negative", run: func() { d.Rotate(-5) }},
			{name: "Clear", run: func() { d.Clear() }},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.run()
				if got := d.Len(); got != 0 {
					t.Fatalf("Len() after %s = %d, want 0", tt.name, got)
				}
			})
		}
	})

	t.Run("Copy", func(t *testing.T) {
		copied := d.Copy()
		if copied == nil {
			t.Fatal("Copy() on nil deque = nil, want empty deque")
		}
		if got := copied.Len(); got != 0 {
			t.Fatalf("copied Len() = %d, want 0", got)
		}
		copied.Append(1)
		if got := copied.Items(); !slices.Equal(got, []int{1}) {
			t.Fatalf("copied Items() = %v, want [1]", got)
		}
		if got := d.Len(); got != 0 {
			t.Fatalf("nil deque Len() = %d, want 0", got)
		}
	})

	t.Run("zero value of element type", func(t *testing.T) {
		var s *Deque[string]
		if got, ok := s.Pop(); got != "" || ok {
			t.Fatalf("Pop() = (%q, %v), want (\"\", false)", got, ok)
		}
	})
}

func TestDequeRingStress(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T) *Deque[int]
		want func() []int
	}{
		{
			name: "grow shift and extend",
			run: func(t *testing.T) *Deque[int] {
				t.Helper()
				d := NewDeque[int]()
				for i := 0; i < 100; i++ {
					d.Append(i)
				}
				for i := 0; i < 50; i++ {
					if got, ok := d.PopLeft(); !ok || got != i {
						t.Fatalf("PopLeft() = (%d, %v), want (%d, true)", got, ok, i)
					}
				}
				d.AppendLeft(999)
				d.Extend(200, 201, 202, 203, 204, 205, 206, 207, 208, 209,
					210, 211, 212, 213, 214, 215, 216, 217, 218, 219,
					220, 221, 222, 223, 224, 225, 226, 227, 228, 229)
				return d
			},
			want: func() []int {
				want := []int{999}
				for i := 50; i < 100; i++ {
					want = append(want, i)
				}
				for i := 200; i < 230; i++ {
					want = append(want, i)
				}
				return want
			},
		},
		{
			name: "alternating ends",
			run: func(t *testing.T) *Deque[int] {
				t.Helper()
				d := NewDeque[int]()
				for i := 0; i < 10; i++ {
					d.Append(i)
				}
				for i := 0; i < 5; i++ {
					if _, ok := d.PopLeft(); !ok {
						t.Fatal("PopLeft() returned false")
					}
					d.Append(100 + i)
				}
				for i := 0; i < 3; i++ {
					if _, ok := d.Pop(); !ok {
						t.Fatal("Pop() returned false")
					}
					d.AppendLeft(200 + i)
				}
				return d
			},
			want: func() []int {
				return []int{202, 201, 200, 5, 6, 7, 8, 9, 100, 101}
			},
		},
		{
			name: "appendleft many and wrap",
			run: func(t *testing.T) *Deque[int] {
				t.Helper()
				d := NewDeque[int]()
				for i := 0; i < 100; i++ {
					d.AppendLeft(i)
				}
				for i := 0; i < 40; i++ {
					if got, ok := d.Pop(); !ok || got != i {
						t.Fatalf("Pop() = (%d, %v), want (%d, true)", got, ok, i)
					}
				}
				for i := 0; i < 20; i++ {
					d.AppendLeft(1000 + i)
				}
				return d
			},
			want: func() []int {
				want := make([]int, 0, 80)
				for i := 1019; i >= 1000; i-- {
					want = append(want, i)
				}
				for i := 99; i >= 40; i-- {
					want = append(want, i)
				}
				return want
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.run(t)
			want := tt.want()
			if got := d.Items(); !slices.Equal(got, want) {
				t.Fatalf("Items() = %v, want %v", got, want)
			}
			if got := d.Len(); got != len(want) {
				t.Fatalf("Len() = %d, want %d", got, len(want))
			}
		})
	}
}
