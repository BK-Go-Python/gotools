package json

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoads(t *testing.T) {
	value, err := Loads(`{"a":1,"b":[true,null,"x"],"c":1.5}`)
	if err != nil {
		t.Fatal(err)
	}
	object, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("type = %T", value)
	}
	if object["a"] != int64(1) || object["c"] != 1.5 {
		t.Fatalf("object = %#v", object)
	}
	array := object["b"].([]any)
	if array[0] != true || array[1] != nil || array[2] != "x" {
		t.Fatalf("array = %#v", array)
	}

	for _, input := range []string{`{`, `1 2`, ``} {
		if _, err := Loads(input); err == nil {
			t.Fatalf("Loads(%q) expected error", input)
		}
	}
}

func TestLoadsInto(t *testing.T) {
	var target struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	if err := LoadsInto(`{"name":"Tom","age":18}`, &target); err != nil {
		t.Fatal(err)
	}
	if target.Name != "Tom" || target.Age != 18 {
		t.Fatalf("target = %#v", target)
	}
	if err := LoadsInto(`{`, &target); err == nil {
		t.Fatal("invalid JSON should error")
	}
}

func TestDumpsAndIndent(t *testing.T) {
	text, err := Dumps(map[string]any{"a": 1, "b": "你好"})
	if err != nil {
		t.Fatal(err)
	}
	if text != `{"a":1,"b":"你好"}` {
		t.Fatalf("Dumps = %q", text)
	}
	pretty, err := DumpsIndent(map[string]any{"a": 1}, "  ")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pretty, "\n  \"a\": 1") {
		t.Fatalf("DumpsIndent = %q", pretty)
	}
	indented, err := Indent(`{"a":1}`, "  ")
	if err != nil || !strings.Contains(indented, "\n  \"a\": 1") {
		t.Fatalf("Indent = (%q,%v)", indented, err)
	}
}

func TestLoadDump(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	if err := Dump(path, map[string]any{"a": 1}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"a":1}` {
		t.Fatalf("file = %q", data)
	}
	value, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if value.(map[string]any)["a"] != int64(1) {
		t.Fatalf("Load = %#v", value)
	}
	if _, err := Load(filepath.Join(dir, "missing.json")); err == nil {
		t.Fatal("missing file should error")
	}
}
