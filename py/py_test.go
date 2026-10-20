package py

import (
	"reflect"
	"testing"
)

func TestTruthyAndTypes(t *testing.T) {
	falseValues := []any{nil, false, 0, 0.0, "", []int{}, map[string]int{}}
	for _, value := range falseValues {
		if Truthy(value) {
			t.Fatalf("Truthy(%#v) = true, want false", value)
		}
	}
	trueValues := []any{1, -1, "x", []int{1}, map[string]int{"a": 1}, struct{}{}}
	for _, value := range trueValues {
		if !Truthy(value) {
			t.Fatalf("Truthy(%#v) = false, want true", value)
		}
	}

	if got := Repr("a"); got != "'a'" {
		t.Fatalf("Repr string = %q", got)
	}
	if got := Repr(nil); got != "None" {
		t.Fatalf("Repr(nil) = %q", got)
	}
	if got := Str(1); got != "1" {
		t.Fatalf("Str(1) = %q", got)
	}
	if got := Type([]int{}); got != "list" {
		t.Fatalf("Type slice = %q", got)
	}
	if got := Type(nil); got != "NoneType" {
		t.Fatalf("Type(nil) = %q", got)
	}

	var nilPointer *int
	if !IsNone(nil) || !IsNone(nilPointer) || IsNone(0) {
		t.Fatal("IsNone mismatch")
	}
}

func TestRangeAndZip(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want []int
	}{
		{"stop", []int{5}, []int{0, 1, 2, 3, 4}},
		{"start stop", []int{2, 5}, []int{2, 3, 4}},
		{"step", []int{0, 10, 2}, []int{0, 2, 4, 6, 8}},
		{"negative", []int{5, 0, -1}, []int{5, 4, 3, 2, 1}},
		{"empty", []int{5, 0, 1}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Range(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Range(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}

	got := Zip([]int{1, 2, 3}, []string{"a", "b"})
	want := [][]interface{}{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Zip = %v, want %v", got, want)
	}
}

func TestFunctionalHelpers(t *testing.T) {
	mapped := Map(func(v interface{}) interface{} { return v.(int) * 2 }, []int{1, 2, 3})
	if got := mapped.Join(","); got != "2,4,6" {
		t.Fatalf("Map = %q", got)
	}

	filtered := Filter(func(v interface{}) bool { return v.(int)%2 == 0 }, []int{1, 2, 3, 4})
	if got := filtered.Join(","); got != "2,4" {
		t.Fatalf("Filter = %q", got)
	}

	if got := Reduce(func(acc, v interface{}) interface{} { return acc.(int) + v.(int) }, []int{1, 2, 3}); got != 6 {
		t.Fatalf("Reduce = %v", got)
	}
	if got := Reduce(func(acc, v interface{}) interface{} { return acc.(int) + v.(int) }, []int{1, 2}, 10); got != 13 {
		t.Fatalf("Reduce initial = %v", got)
	}

	if got := Sorted([]int{3, 1, 2}).Join(","); got != "1,2,3" {
		t.Fatalf("Sorted = %q", got)
	}
	if got := Reversed([]int{1, 2, 3}).Join(","); got != "3,2,1" {
		t.Fatalf("Reversed = %q", got)
	}

	pairs := Enumerate([]string{"a", "b"}, 5)
	if len(pairs) != 2 || pairs[0].Index != 5 || pairs[1].Value != "b" {
		t.Fatalf("Enumerate = %v", pairs)
	}
}

func TestNumericHelpers(t *testing.T) {
	if got := FloorDiv(-7, 3); got != -3 {
		t.Fatalf("FloorDiv = %d", got)
	}
	if got := Mod(-7, 3); got != 2 {
		t.Fatalf("Mod = %d", got)
	}
	q, r := DivMod(7, -3)
	if q != -3 || r != -2 {
		t.Fatalf("DivMod = (%d,%d)", q, r)
	}
	if got := Pow(2, 10); got != 1024 {
		t.Fatalf("Pow = %v", got)
	}
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Fatalf("Sum = %v", got)
	}
	if got, err := SumAny([]int{1, 2, 3}); err != nil || got != 6 {
		t.Fatalf("SumAny = (%v,%v)", got, err)
	}
	if got, err := Min([]int{3, 1, 2}); err != nil || got != 1 {
		t.Fatalf("Min = (%v,%v)", got, err)
	}
	if got, err := Max([]int{3, 1, 2}); err != nil || got != 3 {
		t.Fatalf("Max = (%v,%v)", got, err)
	}
	if !Any([]interface{}{0, "x"}) || All([]interface{}{1, "x"}) != true {
		t.Fatal("Any/All mismatch")
	}
	if got := Len("你好"); got != 2 {
		t.Fatalf("Len = %d", got)
	}
	if !Bool("false") {
		t.Fatal("Bool string truthiness should be true")
	}
}
