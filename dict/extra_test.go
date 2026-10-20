package dict

import (
	"testing"

	"github.com/BK-Go-Python/gotools/list"
)

type namedMap map[string]int

func TestFromKeys(t *testing.T) {
	d := FromKeys([]string{"a", "b"}, 0)
	if len(d) != 2 || d["a"] != 0 || d["b"] != 0 {
		t.Fatalf("FromKeys = %v", d)
	}
	empty := FromKeys(nil)
	if len(empty) != 0 {
		t.Fatalf("FromKeys(nil) = %v", empty)
	}
}

func TestGetPathNamedContainers(t *testing.T) {
	d := Dict{
		"items": list.New(
			map[string]interface{}{"name": "a"},
			map[string]interface{}{"name": "b"},
		),
		"counts": namedMap{"x": 3},
	}
	if got := d.GetPath("items.1.name"); got != "b" {
		t.Fatalf("GetPath list = %v", got)
	}
	if got := d.GetPath("counts.x"); got != 3 {
		t.Fatalf("GetPath named map = %v", got)
	}
	if got := d.GetPath("items.9.name", "missing"); got != "missing" {
		t.Fatalf("GetPath missing = %v", got)
	}
}
