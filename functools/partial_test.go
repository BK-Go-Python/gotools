package functools

import "testing"

func TestPartialReflect(t *testing.T) {
	fn := func(a int, b string) string { return b + string(rune('0'+a)) }
	partial, err := NewPartial(fn, 1)
	if err != nil {
		t.Fatal(err)
	}
	result, err := partial.Call("x")
	if err != nil || len(result) != 1 || result[0] != "x1" {
		t.Fatalf("Call = (%v,%v)", result, err)
	}
	if _, err := partial.Call(); err == nil {
		t.Fatal("too few args should error")
	}
	if _, err := partial.Call("x", "extra"); err == nil {
		t.Fatal("too many args should error")
	}
	if _, err := NewPartial(42); err == nil {
		t.Fatal("non-function should error")
	}
}

func TestPartialVariadicAndPanic(t *testing.T) {
	sum := func(prefix string, values ...int) string {
		total := 0
		for _, value := range values {
			total += value
		}
		return prefix + string(rune('0'+total))
	}
	partial, err := NewPartial(sum, "n=")
	if err != nil {
		t.Fatal(err)
	}
	result, err := partial.Call(1, 2, 3)
	if err != nil || result[0] != "n=6" {
		t.Fatalf("variadic = (%v,%v)", result, err)
	}

	panicky, _ := NewPartial(func() { panic("boom") })
	if _, err := panicky.Call(); err == nil {
		t.Fatal("panic should become error")
	}
}

func TestTypedPartials(t *testing.T) {
	add := func(a, b int) int { return a + b }
	inc := Partial2(add, 1)
	if inc(2) != 3 {
		t.Fatal("Partial2 failed")
	}
	identity := func(a string) string { return a }
	bound := Partial1(identity, "x")
	if bound() != "x" {
		t.Fatal("Partial1 failed")
	}
}
