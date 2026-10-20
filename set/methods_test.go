package set

import (
	"slices"
	"strings"
	"testing"
)

func sortedIntItems(items []int) []int {
	slices.Sort(items)
	return items
}

func TestSetDiscard(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		discard []int
		want    []int
	}{
		{"existing", []int{1, 2, 3}, []int{2}, []int{1, 3}},
		{"multiple existing", []int{1, 2, 3, 4}, []int{2, 4}, []int{1, 3}},
		{"missing", []int{1, 2, 3}, []int{9}, []int{1, 2, 3}},
		{"existing and missing", []int{1, 2, 3}, []int{9, 2}, []int{1, 3}},
		{"remove all", []int{1, 2}, []int{1, 2}, nil},
		{"empty set", nil, []int{1}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)
			got := s.Discard(tt.discard...)

			if got.Len() != len(tt.want) {
				t.Fatalf("Discard(%v) Len = %d, want %d", tt.discard, got.Len(), len(tt.want))
			}
			if items := sortedIntItems(got.Items()); !slices.Equal(items, tt.want) {
				t.Fatalf("Discard(%v) = %v, want %v", tt.discard, items, tt.want)
			}

			// 返回值就是接收者本身：修改链式结果在原集合上可见。
			got.Add(-1)
			if !s.Has(-1) {
				t.Fatal("Discard did not return the receiver")
			}
		})
	}
}

func TestSetDiscardNil(t *testing.T) {
	var s Set[int]
	got := s.Discard(1, 2)
	if got != nil {
		t.Fatalf("nil.Discard = %v, want nil receiver", got)
	}
	if got.Len() != 0 || got.Has(1) {
		t.Fatalf("nil.Discard Len = %d, want 0", got.Len())
	}
}

func TestSetRemoveOne(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		item      int
		zeroValue bool
		want      []int
		wantErr   bool
	}{
		{"existing", []int{1, 2, 3}, 2, false, []int{1, 3}, false},
		{"missing", []int{1, 2, 3}, 9, false, []int{1, 2, 3}, true},
		{"only element", []int{1}, 1, false, nil, false},
		{"zero value", nil, 1, true, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Set[int]
			if !tt.zeroValue {
				s = New(tt.initial...)
			}

			err := s.RemoveOne(tt.item)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("RemoveOne(%d) error = nil, want error", tt.item)
				}
				if !strings.Contains(err.Error(), "not found") {
					t.Fatalf("RemoveOne(%d) error = %q, want it to mention not found", tt.item, err)
				}
			} else if err != nil {
				t.Fatalf("RemoveOne(%d) error = %v, want nil", tt.item, err)
			}

			if items := sortedIntItems(s.Items()); !slices.Equal(items, tt.want) {
				t.Fatalf("after RemoveOne(%d) = %v, want %v", tt.item, items, tt.want)
			}
		})
	}
}

func TestSetIssuperset(t *testing.T) {
	tests := []struct {
		name  string
		s     []int
		other []int
		want  bool
	}{
		{"proper superset", []int{1, 2, 3}, []int{1, 2}, true},
		{"equal sets", []int{1, 2, 3}, []int{3, 1, 2}, true},
		{"empty other", []int{1, 2, 3}, nil, true},
		{"both empty", nil, nil, true},
		{"missing element", []int{1, 2, 3}, []int{1, 4}, false},
		{"empty self", nil, []int{1}, false},
		{"other larger", []int{1, 2}, []int{1, 2, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.s...).Issuperset(New(tt.other...)); got != tt.want {
				t.Fatalf("New(%v).Issuperset(New(%v)) = %v, want %v", tt.s, tt.other, got, tt.want)
			}
		})
	}
}

func TestSetIssupersetNil(t *testing.T) {
	var s Set[int]
	if !s.Issuperset(New[int]()) {
		t.Fatal("nil.Issuperset(empty) = false, want true")
	}
	if s.Issuperset(New(1)) {
		t.Fatal("nil.Issuperset({1}) = true, want false")
	}
}

func TestFrozenIssuperset(t *testing.T) {
	tests := []struct {
		name  string
		f     []int
		other []int
		want  bool
	}{
		{"proper superset", []int{1, 2, 3}, []int{1, 2}, true},
		{"equal sets", []int{1, 2, 3}, []int{3, 1, 2}, true},
		{"empty other", []int{1, 2, 3}, nil, true},
		{"both empty", nil, nil, true},
		{"missing element", []int{1, 2, 3}, []int{1, 4}, false},
		{"empty self", nil, []int{1}, false},
		{"other larger", []int{1, 2}, []int{1, 2, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFrozen(tt.f...)
			other := NewFrozen(tt.other...)
			fLen, otherLen := f.Len(), other.Len()

			if got := f.Issuperset(other); got != tt.want {
				t.Fatalf("NewFrozen(%v).Issuperset(NewFrozen(%v)) = %v, want %v", tt.f, tt.other, got, tt.want)
			}
			// Frozen 只读：调用不得改变任何一方。
			if f.Len() != fLen || other.Len() != otherLen {
				t.Fatal("Frozen.Issuperset mutated its operands")
			}
		})
	}
}

func TestFrozenIssupersetZeroValue(t *testing.T) {
	var f Frozen[int]
	if !f.Issuperset(NewFrozen[int]()) {
		t.Fatal("zero Frozen.Issuperset(empty) = false, want true")
	}
	if f.Issuperset(NewFrozen(1)) {
		t.Fatal("zero Frozen.Issuperset({1}) = true, want false")
	}
}
