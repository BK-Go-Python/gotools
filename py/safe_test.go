package py

import (
	"errors"
	"strings"
	"testing"
)

func TestTryAndCatch(t *testing.T) {
	if value, err := Try(func() int { return 7 }); err != nil || value != 7 {
		t.Fatalf("Try = (%d,%v)", value, err)
	}
	if value, err := Try(func() int { panic("boom") }); err == nil || value != 0 || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("Try panic = (%d,%v)", value, err)
	}
	if err := Catch(func() { panic(ErrType) }); !errors.Is(err, ErrType) {
		t.Fatalf("Catch = %v", err)
	}
	if err := Catch(nil); err != nil {
		t.Fatalf("Catch nil = %v", err)
	}
	if a, b, err := Try2(func() (int, string) { return 1, "x" }); err != nil || a != 1 || b != "x" {
		t.Fatalf("Try2 = (%d,%q,%v)", a, b, err)
	}
	if a, b, c, err := Try3(func() (int, int, int) { return 1, 2, 3 }); err != nil || a+b+c != 6 {
		t.Fatalf("Try3 = (%d,%d,%d,%v)", a, b, c, err)
	}
}

func TestSafeArithmeticAndRange(t *testing.T) {
	if got, err := RangeE(3); err != nil || len(got) != 3 {
		t.Fatalf("RangeE = (%v,%v)", got, err)
	}
	if _, err := RangeE(1, 2, 3, 4); err == nil {
		t.Fatal("RangeE extra args should error")
	}
	if _, err := RangeE(0, 5, 0); err == nil {
		t.Fatal("RangeE zero step should error")
	}
	if _, err := FloorDivE(1, 0); !errors.Is(err, ErrZeroDivision) {
		t.Fatalf("FloorDivE = %v", err)
	}
	if _, err := ModE(1, 0); !errors.Is(err, ErrZeroDivision) {
		t.Fatalf("ModE = %v", err)
	}
	if _, _, err := DivModE(1, 0); !errors.Is(err, ErrZeroDivision) {
		t.Fatalf("DivModE = %v", err)
	}
	if q, r, err := DivModE(-7, 3); err != nil || q != -3 || r != 2 {
		t.Fatalf("DivModE = (%d,%d,%v)", q, r, err)
	}
}

func TestSafeDynamicAPIs(t *testing.T) {
	invalid := 42
	if _, err := SumAny(invalid); err == nil {
		t.Fatal("SumAny invalid should error")
	}
	if _, err := MinAny(invalid); err == nil {
		t.Fatal("MinAny invalid should error")
	}
	if _, err := MaxAny(invalid); err == nil {
		t.Fatal("MaxAny invalid should error")
	}
	if _, err := AnyE(invalid); err == nil {
		t.Fatal("AnyE invalid should error")
	}
	if _, err := AllE(invalid); err == nil {
		t.Fatal("AllE invalid should error")
	}
	if _, err := EnumerateAnyE(invalid); err == nil {
		t.Fatal("EnumerateAnyE invalid should error")
	}
	if _, err := ZipE(invalid); err == nil {
		t.Fatal("ZipE invalid should error")
	}
	if _, err := SortedAnyE(invalid); err == nil {
		t.Fatal("SortedAnyE invalid should error")
	}
	if _, err := SortedByE(invalid, func(v interface{}) interface{} { return v }); err == nil {
		t.Fatal("SortedByE invalid should error")
	}
	if _, err := ReversedAnyE(invalid); err == nil {
		t.Fatal("ReversedAnyE invalid should error")
	}
	if _, err := MapAnyE(func(v interface{}) interface{} { return v }, invalid); err == nil {
		t.Fatal("MapAnyE invalid should error")
	}
	if _, err := FilterAnyE(func(v interface{}) bool { return true }, invalid); err == nil {
		t.Fatal("FilterAnyE invalid should error")
	}

	if _, err := ReduceE(func(acc, item interface{}) interface{} { return acc }, []int{}); !errors.Is(err, ErrEmpty) {
		t.Fatalf("ReduceE empty = %v", err)
	}
	if _, err := ReduceE(func(acc, item interface{}) interface{} { panic("boom") }, []int{1, 2}); err == nil {
		t.Fatal("ReduceE callback panic should error")
	}
}

func TestSafeVariantsValid(t *testing.T) {
	if got, err := ZipE([]int{1, 2}, []string{"a"}); err != nil || len(got) != 1 {
		t.Fatalf("ZipE = (%v,%v)", got, err)
	}
	if got, err := SortedAnyE([]int{3, 1, 2}); err != nil || got.String() != "[1, 2, 3]" {
		t.Fatalf("SortedAnyE = (%v,%v)", got, err)
	}
	if got, err := ReversedAnyE([]int{1, 2, 3}); err != nil || got.String() != "[3, 2, 1]" {
		t.Fatalf("ReversedAnyE = (%v,%v)", got, err)
	}
	if got, err := MapAnyE(func(v interface{}) interface{} { return v.(int) * 2 }, []int{1, 2}); err != nil || got.String() != "[2, 4]" {
		t.Fatalf("MapAnyE = (%v,%v)", got, err)
	}
	if got, err := FilterAnyE(func(v interface{}) bool { return v.(int)%2 == 0 }, []int{1, 2, 3, 4}); err != nil || got.Len() != 2 {
		t.Fatalf("FilterAnyE = (%v,%v)", got, err)
	}
	if got, err := EnumerateAnyE([]string{"a"}, 5); err != nil || len(got) != 1 || got[0].Index != 5 {
		t.Fatalf("EnumerateAnyE = (%v,%v)", got, err)
	}
}
