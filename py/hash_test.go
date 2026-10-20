package py

import (
	"errors"
	"testing"
)

func TestHash(t *testing.T) {
	a, err := Hash("abc")
	if err != nil {
		t.Fatal(err)
	}
	b, err := Hash("abc")
	if err != nil || a != b {
		t.Fatalf("Hash stable = (%d,%d,%v)", a, b, err)
	}
	c, err := Hash("abd")
	if err != nil || a == c {
		t.Fatalf("Hash collision for abc/abd: %d %d %v", a, c, err)
	}
	if _, err := Hash(nil); err != nil {
		t.Fatalf("Hash(nil) = %v", err)
	}
	one, _ := Hash(1)
	oneFloat, _ := Hash(1.0)
	oneBool, _ := Hash(true)
	if one != oneFloat || one != oneBool {
		t.Fatalf("hash equality contract broken: 1=%d 1.0=%d true=%d", one, oneFloat, oneBool)
	}
	if _, err := Hash(struct{ A int }{1}); err != nil {
		t.Fatalf("Hash(struct) = %v", err)
	}

	for _, v := range []any{[]int{1}, map[string]int{"a": 1}, func() {}} {
		if _, err := Hash(v); !errors.Is(err, ErrType) {
			t.Fatalf("Hash(%T) err = %v, want ErrType", v, err)
		}
	}
}
