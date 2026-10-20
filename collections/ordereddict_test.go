package collections

import (
	"reflect"
	"testing"
)

// odTestCheckEqual 比较 got 与 want，用于表驱动测试。
func odTestCheckEqual(t *testing.T, label string, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %#v, want %#v", label, got, want)
	}
}

// odTestBaseDict 对应 Python 的
// OrderedDict([("a", 1), ("b", 2), ("c", 3)])。
func odTestBaseDict() *OrderedDict[string, int] {
	d := NewOrderedDict[string, int]()
	d.Set("a", 1).Set("b", 2).Set("c", 3)
	return d
}

// TestOrderedDictSetOrder 覆盖插入顺序与更新语义。期望值与 Python 3.13 一致：
//
//	od = OrderedDict([("a", 1), ("b", 2), ("c", 3)])
//	od["b"] = 20
//	list(od.items())  # [("a", 1), ("b", 20), ("c", 3)]
func TestOrderedDictSetOrder(t *testing.T) {
	tests := []struct {
		name       string
		build      func(d *OrderedDict[string, int])
		wantKeys   []string
		wantValues []int
		wantItems  []OrderedItem[string, int]
	}{
		{
			name: "insert appends in order",
			build: func(d *OrderedDict[string, int]) {
				d.Set("a", 1).Set("b", 2).Set("c", 3)
			},
			wantKeys:   []string{"a", "b", "c"},
			wantValues: []int{1, 2, 3},
			wantItems: []OrderedItem[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
			},
		},
		{
			name: "update existing keeps position",
			build: func(d *OrderedDict[string, int]) {
				d.Set("a", 1).Set("b", 2).Set("c", 3)
				d.Set("b", 20)
			},
			wantKeys:   []string{"a", "b", "c"},
			wantValues: []int{1, 20, 3},
			wantItems: []OrderedItem[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 20},
				{Key: "c", Value: 3},
			},
		},
		{
			name: "delete then reinsert appends to end",
			build: func(d *OrderedDict[string, int]) {
				d.Set("a", 1).Set("b", 2).Set("c", 3)
				d.Delete("b")
				d.Set("b", 20)
			},
			wantKeys:   []string{"a", "c", "b"},
			wantValues: []int{1, 3, 20},
			wantItems: []OrderedItem[string, int]{
				{Key: "a", Value: 1},
				{Key: "c", Value: 3},
				{Key: "b", Value: 20},
			},
		},
		{
			name:       "empty",
			build:      func(*OrderedDict[string, int]) {},
			wantKeys:   []string{},
			wantValues: []int{},
			wantItems:  []OrderedItem[string, int]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewOrderedDict[string, int]()
			tt.build(d)
			odTestCheckEqual(t, "Keys", d.Keys(), tt.wantKeys)
			odTestCheckEqual(t, "Values", d.Values(), tt.wantValues)
			odTestCheckEqual(t, "Items", d.Items(), tt.wantItems)
			odTestCheckEqual(t, "Len", d.Len(), len(tt.wantKeys))
		})
	}
}

// TestOrderedDictSetReturnsReceiver 验证 Set 返回自身以支持链式调用。
func TestOrderedDictSetReturnsReceiver(t *testing.T) {
	t.Run("chain", func(t *testing.T) {
		d := NewOrderedDict[string, int]()
		got := d.Set("a", 1).Set("b", 2)
		if got != d {
			t.Fatalf("Set returned %p, want receiver %p", got, d)
		}
		odTestCheckEqual(t, "Keys", d.Keys(), []string{"a", "b"})
		odTestCheckEqual(t, "Values", d.Values(), []int{1, 2})
	})
}

// TestOrderedDictGetPop 覆盖 Get/Pop 的默认值语义。Python 中取值和
// 删除分别对应 od.get(key, default) 与 od.pop(key, default)。
func TestOrderedDictGetPop(t *testing.T) {
	tests := []struct {
		name     string
		op       func(d *OrderedDict[string, int]) int
		want     int
		wantKeys []string
	}{
		{
			name:     "get existing",
			op:       func(d *OrderedDict[string, int]) int { return d.Get("b") },
			want:     2,
			wantKeys: []string{"a", "b", "c"},
		},
		{
			name:     "get missing returns zero value",
			op:       func(d *OrderedDict[string, int]) int { return d.Get("z") },
			want:     0,
			wantKeys: []string{"a", "b", "c"},
		},
		{
			name:     "get missing returns default",
			op:       func(d *OrderedDict[string, int]) int { return d.Get("z", 99) },
			want:     99,
			wantKeys: []string{"a", "b", "c"},
		},
		{
			name:     "get existing ignores default",
			op:       func(d *OrderedDict[string, int]) int { return d.Get("a", 99) },
			want:     1,
			wantKeys: []string{"a", "b", "c"},
		},
		{
			name:     "pop existing returns value and removes",
			op:       func(d *OrderedDict[string, int]) int { return d.Pop("b") },
			want:     2,
			wantKeys: []string{"a", "c"},
		},
		{
			name:     "pop missing returns zero value",
			op:       func(d *OrderedDict[string, int]) int { return d.Pop("z") },
			want:     0,
			wantKeys: []string{"a", "b", "c"},
		},
		{
			name:     "pop missing returns default",
			op:       func(d *OrderedDict[string, int]) int { return d.Pop("z", 99) },
			want:     99,
			wantKeys: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := odTestBaseDict()
			odTestCheckEqual(t, "result", tt.op(d), tt.want)
			odTestCheckEqual(t, "Keys after op", d.Keys(), tt.wantKeys)
			odTestCheckEqual(t, "Len after op", d.Len(), len(tt.wantKeys))
		})
	}
}

// TestOrderedDictHasDelete 覆盖成员判断和删除返回值。
func TestOrderedDictHasDelete(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		wantHas    bool
		wantDelete bool
		wantKeys   []string
	}{
		{"first key", "a", true, true, []string{"b", "c"}},
		{"middle key", "b", true, true, []string{"a", "c"}},
		{"last key", "c", true, true, []string{"a", "b"}},
		{"missing key", "z", false, false, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := odTestBaseDict()
			odTestCheckEqual(t, "Has before delete", d.Has(tt.key), tt.wantHas)
			odTestCheckEqual(t, "Delete", d.Delete(tt.key), tt.wantDelete)
			odTestCheckEqual(t, "Has after delete", d.Has(tt.key), false)
			odTestCheckEqual(t, "second Delete", d.Delete(tt.key), false)
			odTestCheckEqual(t, "Keys", d.Keys(), tt.wantKeys)
		})
	}
}

// TestOrderedDictMoveToEnd 覆盖 move_to_end 的位置变化。期望值与
// Python 3.13 的 od.move_to_end(key) 一致。
func TestOrderedDictMoveToEnd(t *testing.T) {
	tests := []struct {
		name       string
		empty      bool
		key        string
		wantMoved  bool
		wantKeys   []string
		wantValues []int
	}{
		{"first key", false, "a", true, []string{"b", "c", "a"}, []int{2, 3, 1}},
		{"middle key", false, "b", true, []string{"a", "c", "b"}, []int{1, 3, 2}},
		{"last key stays in place", false, "c", true, []string{"a", "b", "c"}, []int{1, 2, 3}},
		{"missing key", false, "z", false, []string{"a", "b", "c"}, []int{1, 2, 3}},
		{"empty dict", true, "a", false, []string{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewOrderedDict[string, int]()
			if !tt.empty {
				d = odTestBaseDict()
			}
			odTestCheckEqual(t, "MoveToEnd", d.MoveToEnd(tt.key), tt.wantMoved)
			odTestCheckEqual(t, "Keys", d.Keys(), tt.wantKeys)
			odTestCheckEqual(t, "Values", d.Values(), tt.wantValues)
		})
	}
}

// TestOrderedDictPopItem 覆盖 popitem 的 LIFO/FIFO 与空字典语义。
func TestOrderedDictPopItem(t *testing.T) {
	tests := []struct {
		name     string
		last     []bool
		wantItem OrderedItem[string, int]
		wantKeys []string
	}{
		{
			name:     "default removes last",
			last:     nil,
			wantItem: OrderedItem[string, int]{Key: "c", Value: 3},
			wantKeys: []string{"a", "b"},
		},
		{
			name:     "last true removes last",
			last:     []bool{true},
			wantItem: OrderedItem[string, int]{Key: "c", Value: 3},
			wantKeys: []string{"a", "b"},
		},
		{
			name:     "last false removes first",
			last:     []bool{false},
			wantItem: OrderedItem[string, int]{Key: "a", Value: 1},
			wantKeys: []string{"b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := odTestBaseDict()
			item, ok := d.PopItem(tt.last...)
			if !ok {
				t.Fatal("PopItem returned ok = false, want true")
			}
			odTestCheckEqual(t, "item", item, tt.wantItem)
			odTestCheckEqual(t, "Keys", d.Keys(), tt.wantKeys)
			odTestCheckEqual(t, "Len", d.Len(), 2)
		})
	}

	t.Run("empty returns zero and false", func(t *testing.T) {
		d := NewOrderedDict[string, int]()
		item, ok := d.PopItem()
		if ok {
			t.Fatal("PopItem on empty dict returned ok = true, want false")
		}
		odTestCheckEqual(t, "item", item, OrderedItem[string, int]{})
		odTestCheckEqual(t, "Len", d.Len(), 0)
	})

	t.Run("LIFO drain", func(t *testing.T) {
		d := odTestBaseDict()
		var got []OrderedItem[string, int]
		for {
			item, ok := d.PopItem()
			if !ok {
				break
			}
			got = append(got, item)
		}
		want := []OrderedItem[string, int]{
			{Key: "c", Value: 3},
			{Key: "b", Value: 2},
			{Key: "a", Value: 1},
		}
		odTestCheckEqual(t, "drained items", got, want)
		odTestCheckEqual(t, "Len", d.Len(), 0)
	})

	t.Run("FIFO drain", func(t *testing.T) {
		d := odTestBaseDict()
		var got []OrderedItem[string, int]
		for {
			item, ok := d.PopItem(false)
			if !ok {
				break
			}
			got = append(got, item)
		}
		want := []OrderedItem[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
			{Key: "c", Value: 3},
		}
		odTestCheckEqual(t, "drained items", got, want)
		odTestCheckEqual(t, "Len", d.Len(), 0)
	})
}

// TestOrderedDictViews 覆盖 Keys/Values/Items 的顺序和快照语义。
func TestOrderedDictViews(t *testing.T) {
	t.Run("insertion order with update", func(t *testing.T) {
		d := odTestBaseDict()
		d.Set("b", 20)
		odTestCheckEqual(t, "Keys", d.Keys(), []string{"a", "b", "c"})
		odTestCheckEqual(t, "Values", d.Values(), []int{1, 20, 3})
		odTestCheckEqual(t, "Items", d.Items(), []OrderedItem[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 20},
			{Key: "c", Value: 3},
		})
	})

	t.Run("empty views", func(t *testing.T) {
		d := NewOrderedDict[string, int]()
		odTestCheckEqual(t, "Keys", d.Keys(), []string{})
		odTestCheckEqual(t, "Values", d.Values(), []int{})
		odTestCheckEqual(t, "Items", d.Items(), []OrderedItem[string, int]{})
	})

	t.Run("returned slices are copies", func(t *testing.T) {
		d := odTestBaseDict()
		keys := d.Keys()
		keys[0] = "z"
		values := d.Values()
		values[0] = 99
		items := d.Items()
		items[0].Key = "z"
		items[0].Value = 99

		odTestCheckEqual(t, "Keys after mutating view", d.Keys(), []string{"a", "b", "c"})
		odTestCheckEqual(t, "Values after mutating view", d.Values(), []int{1, 2, 3})
		odTestCheckEqual(t, "Items after mutating view", d.Items(), []OrderedItem[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
			{Key: "c", Value: 3},
		})
	})
}

// TestOrderedDictLenClear 覆盖长度跟踪和 Clear 后继续使用。
func TestOrderedDictLenClear(t *testing.T) {
	t.Run("length tracks operations", func(t *testing.T) {
		d := NewOrderedDict[string, int]()
		odTestCheckEqual(t, "empty Len", d.Len(), 0)
		d.Set("a", 1).Set("b", 2)
		odTestCheckEqual(t, "Len after Set", d.Len(), 2)
		d.Set("a", 10)
		odTestCheckEqual(t, "Len after update", d.Len(), 2)
		d.Delete("a")
		odTestCheckEqual(t, "Len after Delete", d.Len(), 1)
		d.PopItem()
		odTestCheckEqual(t, "Len after PopItem", d.Len(), 0)
	})

	t.Run("clear empties and dict stays usable", func(t *testing.T) {
		d := odTestBaseDict()
		keysBefore := d.Keys()
		d.Clear()
		odTestCheckEqual(t, "Len", d.Len(), 0)
		odTestCheckEqual(t, "Keys", d.Keys(), []string{})
		odTestCheckEqual(t, "Values", d.Values(), []int{})
		odTestCheckEqual(t, "Items", d.Items(), []OrderedItem[string, int]{})
		odTestCheckEqual(t, "Has", d.Has("a"), false)
		odTestCheckEqual(t, "Get", d.Get("a"), 0)
		odTestCheckEqual(t, "Keys before clear unchanged", keysBefore, []string{"a", "b", "c"})

		d.Set("x", 24)
		odTestCheckEqual(t, "Keys after reuse", d.Keys(), []string{"x"})
		odTestCheckEqual(t, "Values after reuse", d.Values(), []int{24})
	})

	t.Run("clear empty dict", func(t *testing.T) {
		d := NewOrderedDict[string, int]()
		d.Clear()
		d.Clear()
		odTestCheckEqual(t, "Len", d.Len(), 0)
	})
}

// TestOrderedDictCopy 覆盖浅拷贝的独立性和顺序保持。
func TestOrderedDictCopy(t *testing.T) {
	t.Run("preserves order and is equal", func(t *testing.T) {
		d := odTestBaseDict()
		cp := d.Copy()
		if cp == d {
			t.Fatal("Copy returned the same pointer")
		}
		if !cp.Equal(d) || !d.Equal(cp) {
			t.Fatal("Copy should be Equal to the original")
		}
		odTestCheckEqual(t, "Keys", cp.Keys(), d.Keys())
		odTestCheckEqual(t, "Values", cp.Values(), d.Values())
	})

	t.Run("independent from original", func(t *testing.T) {
		d := odTestBaseDict()
		cp := d.Copy()

		cp.Set("b", 20)
		cp.Set("d", 4)
		odTestCheckEqual(t, "original after copy changes", d.Items(), []OrderedItem[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
			{Key: "c", Value: 3},
		})

		d.Delete("a")
		d.Set("e", 5)
		odTestCheckEqual(t, "copy after original changes", cp.Items(), []OrderedItem[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 20},
			{Key: "c", Value: 3},
			{Key: "d", Value: 4},
		})
	})

	t.Run("empty copy is usable", func(t *testing.T) {
		d := NewOrderedDict[string, int]()
		cp := d.Copy()
		if cp == nil || cp == d {
			t.Fatalf("Copy returned %p, want a new non-nil dict", cp)
		}
		cp.Set("x", 1)
		odTestCheckEqual(t, "original Len", d.Len(), 0)
		odTestCheckEqual(t, "copy Keys", cp.Keys(), []string{"x"})
	})
}

// TestOrderedDictEqual 覆盖有序相等比较。Python 3.13 中
// OrderedDict([("a", 1), ("b", 2)]) == OrderedDict([("b", 2), ("a", 1)])
// 为 False，而普通 dict 的同类比较为 True。
func TestOrderedDictEqual(t *testing.T) {
	pair := func(key string, value any) OrderedItem[string, any] {
		return OrderedItem[string, any]{Key: key, Value: value}
	}
	newDict := func(items ...OrderedItem[string, any]) *OrderedDict[string, any] {
		d := NewOrderedDict[string, any]()
		for _, item := range items {
			d.Set(item.Key, item.Value)
		}
		return d
	}

	ab := newDict(pair("a", 1), pair("b", 2))
	ba := newDict(pair("b", 2), pair("a", 1))

	tests := []struct {
		name  string
		left  *OrderedDict[string, any]
		right *OrderedDict[string, any]
		want  bool
	}{
		{"same pairs in same order", ab, newDict(pair("a", 1), pair("b", 2)), true},
		{"same pairs in different order", ab, ba, false},
		{"different value", ab, newDict(pair("a", 1), pair("b", 3)), false},
		{"different key", ab, newDict(pair("a", 1), pair("c", 2)), false},
		{"different length", ab, newDict(pair("a", 1)), false},
		{"longer prefix", newDict(pair("a", 1), pair("b", 2), pair("c", 3)), ab, false},
		{"empty dicts", NewOrderedDict[string, any](), NewOrderedDict[string, any](), true},
		{"nil and nil", nil, nil, true},
		{"nil and empty", nil, NewOrderedDict[string, any](), true},
		{"empty and nil", NewOrderedDict[string, any](), nil, true},
		{"nil and non-empty", nil, ab, false},
		{"non-empty and nil", ab, nil, false},
		{
			name:  "deep equal slice values",
			left:  newDict(pair("a", []int{1, 2})),
			right: newDict(pair("a", []int{1, 2})),
			want:  true,
		},
		{
			name:  "different slice values",
			left:  newDict(pair("a", []int{1, 2})),
			right: newDict(pair("a", []int{1, 3})),
			want:  false,
		},
		{
			name:  "same slice values in different order",
			left:  newDict(pair("a", []int{1, 2}), pair("b", []int{3})),
			right: newDict(pair("b", []int{3}), pair("a", []int{1, 2})),
			want:  false,
		},
		{
			name:  "numeric cross type equality",
			left:  newDict(pair("a", 1)),
			right: newDict(pair("a", 1.0)),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.left.Equal(tt.right); got != tt.want {
				t.Errorf("left.Equal(right) = %v, want %v", got, tt.want)
			}
			if got := tt.right.Equal(tt.left); got != tt.want {
				t.Errorf("right.Equal(left) = %v, want %v", got, tt.want)
			}
			if tt.want {
				if !tt.left.Equal(tt.left) {
					t.Error("dict should equal itself")
				}
				if !tt.right.Equal(tt.right) {
					t.Error("dict should equal itself")
				}
			}
		})
	}
}

// TestOrderedDictNilReceiver 验证 nil 接收者上的方法不会 panic。
func TestOrderedDictNilReceiver(t *testing.T) {
	var d *OrderedDict[string, int]

	t.Run("read methods", func(t *testing.T) {
		odTestCheckEqual(t, "Len", d.Len(), 0)
		odTestCheckEqual(t, "Get", d.Get("a"), 0)
		odTestCheckEqual(t, "Get with default", d.Get("a", 42), 42)
		odTestCheckEqual(t, "Has", d.Has("a"), false)
		odTestCheckEqual(t, "Keys length", len(d.Keys()), 0)
		odTestCheckEqual(t, "Values length", len(d.Values()), 0)
		odTestCheckEqual(t, "Items length", len(d.Items()), 0)
		odTestCheckEqual(t, "Equal nil", d.Equal(nil), true)
		odTestCheckEqual(t, "Equal empty", d.Equal(NewOrderedDict[string, int]()), true)
		odTestCheckEqual(t, "Equal non-empty", d.Equal(odTestBaseDict()), false)
	})

	t.Run("delete methods", func(t *testing.T) {
		odTestCheckEqual(t, "Delete", d.Delete("a"), false)
		odTestCheckEqual(t, "Pop", d.Pop("a"), 0)
		odTestCheckEqual(t, "Pop with default", d.Pop("a", 42), 42)
		odTestCheckEqual(t, "MoveToEnd", d.MoveToEnd("a"), false)
		item, ok := d.PopItem()
		if ok {
			t.Fatal("PopItem on nil receiver returned ok = true, want false")
		}
		odTestCheckEqual(t, "PopItem item", item, OrderedItem[string, int]{})
	})

	t.Run("write methods are no-ops", func(t *testing.T) {
		if got := d.Set("a", 1); got != nil {
			t.Fatalf("Set on nil receiver = %p, want nil", got)
		}
		d.Clear()
	})

	t.Run("Copy returns usable empty dict", func(t *testing.T) {
		cp := d.Copy()
		if cp == nil {
			t.Fatal("Copy on nil receiver returned nil")
		}
		odTestCheckEqual(t, "Len", cp.Len(), 0)
		cp.Set("x", 9)
		odTestCheckEqual(t, "Keys", cp.Keys(), []string{"x"})
	})
}

// TestOrderedDictZeroValue 验证 OrderedDict 的零值可以直接使用。
func TestOrderedDictZeroValue(t *testing.T) {
	t.Run("usable without constructor", func(t *testing.T) {
		var d OrderedDict[string, int]
		d.Set("a", 1).Set("b", 2).Set("a", 3)
		odTestCheckEqual(t, "Keys", d.Keys(), []string{"a", "b"})
		odTestCheckEqual(t, "Values", d.Values(), []int{3, 2})
		odTestCheckEqual(t, "MoveToEnd", d.MoveToEnd("a"), true)
		odTestCheckEqual(t, "Keys after move", d.Keys(), []string{"b", "a"})
		d.Clear()
		odTestCheckEqual(t, "Len after Clear", d.Len(), 0)
		d.Set("c", 4)
		odTestCheckEqual(t, "Keys after reuse", d.Keys(), []string{"c"})
	})

	t.Run("equal to empty constructor result", func(t *testing.T) {
		var zero OrderedDict[string, int]
		if !zero.Equal(NewOrderedDict[string, int]()) {
			t.Error("zero value should equal an empty NewOrderedDict")
		}
	})
}
