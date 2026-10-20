package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenModes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	f, err := Open(path, "w")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("hello"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	f, err = Open(path, "a")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(" world"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Fatalf("content = %q", data)
	}

	f, err = Open(path, "r+")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("HELLO"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(path)
	if string(data) != "HELLO world" {
		t.Fatalf("r+ content = %q", data)
	}
}

func TestOpenErrors(t *testing.T) {
	dir := t.TempDir()
	if _, err := Open(filepath.Join(dir, "missing.txt")); err == nil {
		t.Fatal("expected error for missing file")
	}
	for _, mode := range []string{"", "x", "rw", "+"} {
		if _, err := Open(filepath.Join(dir, "x"), mode); err == nil {
			t.Fatalf("expected error for mode %q", mode)
		}
	}
}
