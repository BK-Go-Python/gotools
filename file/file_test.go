package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadWriteAppendLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")

	if err := Write(path, "hello"); err != nil {
		t.Fatal(err)
	}
	if got, err := Read(path); err != nil || got != "hello" {
		t.Fatalf("Read = (%q,%v)", got, err)
	}
	if err := Append(path, " world"); err != nil {
		t.Fatal(err)
	}
	if got, _ := Read(path); got != "hello world" {
		t.Fatalf("after Append = %q", got)
	}

	lines := []string{"a", "b", "c"}
	if err := WriteLines(path, lines); err != nil {
		t.Fatal(err)
	}
	got, err := ReadLines(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Fatalf("ReadLines = %v", got)
	}
	if err := WriteLines(path, nil); err != nil {
		t.Fatal(err)
	}
	if got, _ := Read(path); got != "" {
		t.Fatalf("WriteLines(nil) content = %q", got)
	}
}

func TestPathOperations(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "x", "y")
	if err := MkdirP(nested); err != nil {
		t.Fatal(err)
	}
	if !Exists(nested) || !IsDir(nested) || IsFile(nested) {
		t.Fatal("directory predicates mismatch")
	}

	path := filepath.Join(nested, "f.txt")
	if err := Write(path, "data"); err != nil {
		t.Fatal(err)
	}
	if !IsFile(path) {
		t.Fatal("IsFile failed")
	}
	if size, err := Size(path); err != nil || size != 4 {
		t.Fatalf("Size = (%d,%v)", size, err)
	}
	if _, err := ModTime(path); err != nil {
		t.Fatal(err)
	}
	if abs, err := Abs(path); err != nil || !filepath.IsAbs(abs) {
		t.Fatalf("Abs = (%q,%v)", abs, err)
	}
	if cwd, err := Cwd(); err != nil || cwd == "" {
		t.Fatalf("Cwd = (%q,%v)", cwd, err)
	}
	if home, err := Home(); err != nil || home == "" {
		t.Fatalf("Home = (%q,%v)", home, err)
	}

	names, err := ListDir(nested)
	if err != nil || len(names) != 1 || names[0] != "f.txt" {
		t.Fatalf("ListDir = (%v,%v)", names, err)
	}
	matches, err := Glob(filepath.Join(nested, "*.txt"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("Glob = (%v,%v)", matches, err)
	}
}

func TestCopyMoveRemove(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	moved := filepath.Join(dir, "moved.txt")

	if err := Write(src, "payload"); err != nil {
		t.Fatal(err)
	}
	if err := Copy(src, dst); err != nil {
		t.Fatal(err)
	}
	if got, _ := Read(dst); got != "payload" {
		t.Fatalf("Copy content = %q", got)
	}
	if err := Move(dst, moved); err != nil {
		t.Fatal(err)
	}
	if Exists(dst) || !Exists(moved) {
		t.Fatal("Move did not relocate file")
	}
	if err := Remove(moved); err != nil {
		t.Fatal(err)
	}
	if Exists(moved) {
		t.Fatal("Remove failed")
	}
	if err := RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	if Exists(dir) {
		t.Fatal("RemoveAll failed")
	}
}

func TestTempAndTouch(t *testing.T) {
	dir := t.TempDir()
	f, err := TempFile(dir, "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	f.Close()
	if !Exists(name) {
		t.Fatal("TempFile not created")
	}
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}

	tempDir, err := TempDir(dir, "dir-*")
	if err != nil || !IsDir(tempDir) {
		t.Fatalf("TempDir = (%q,%v)", tempDir, err)
	}

	path := filepath.Join(dir, "touch.txt")
	if err := Touch(path); err != nil {
		t.Fatal(err)
	}
	if !Exists(path) {
		t.Fatal("Touch did not create file")
	}
	if err := Touch(path); err != nil {
		t.Fatal(err)
	}
}
