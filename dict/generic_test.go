package dict

import "testing"

func TestOf(t *testing.T) {
	input := map[string]int{"a": 1}
	got := Of(input)
	got["a"] = 9
	got["b"] = 2
	if input["a"] != 1 || len(input) != 1 {
		t.Fatal("Of did not copy map")
	}
}
