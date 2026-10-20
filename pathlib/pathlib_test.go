package pathlib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathParts(t *testing.T) {
	p := New("tmp", "dir", "file.txt")
	if p.Name() != "file.txt" || p.Suffix() != ".txt" || p.Stem() != "file" {
		t.Fatalf("parts = %q %q %q", p.Name(), p.Suffix(), p.Stem())
	}
	if p.Parent() != Path(filepath.Join("tmp", "dir")) {
		t.Fatalf("Parent = %q", p.Parent())
	}
	if got := p.WithSuffix(".md"); got != Path(filepath.Join("tmp", "dir", "file.md")) {
		t.Fatalf("WithSuffix = %q", got)
	}
	if New().String() != "." {
		t.Fatalf("New() = %q", New())
	}
	if !From("/tmp").IsAbsolute() {
		t.Fatal("IsAbsolute failed")
	}
}

func TestPathIO(t *testing.T) {
	dir := t.TempDir()
	p := New(dir, "file.txt")
	if err := p.WriteText("hello"); err != nil {
		t.Fatal(err)
	}
	if !p.Exists() || !p.IsFile() || p.IsDir() {
		t.Fatal("predicates failed")
	}
	if got, err := p.ReadText(); err != nil || got != "hello" {
		t.Fatalf("ReadText = (%q,%v)", got, err)
	}
	if err := p.WriteBytes([]byte{1, 2}); err != nil {
		t.Fatal(err)
	}
	if got, err := p.ReadBytes(); err != nil || len(got) != 2 {
		t.Fatalf("ReadBytes = (%v,%v)", got, err)
	}
	if size, err := p.Size(); err != nil || size != 2 {
		t.Fatalf("Size = (%d,%v)", size, err)
	}
	if err := p.Touch(); err != nil {
		t.Fatal(err)
	}
}

func TestPathDirectories(t *testing.T) {
	dir := t.TempDir()
	nested := New(dir, "a", "b")
	if err := nested.Mkdir(true); err != nil {
		t.Fatal(err)
	}
	if !nested.IsDir() {
		t.Fatal("Mkdir failed")
	}
	if err := New(dir, "a").Mkdir(); err == nil {
		t.Fatal("Mkdir existing should error")
	}
	if err := New(dir, "a", "b", "c.txt").WriteText("x"); err != nil {
		t.Fatal(err)
	}
	items, err := nested.Iterdir()
	if err != nil || len(items) != 1 || items[0].Name() != "c.txt" {
		t.Fatalf("Iterdir = (%v,%v)", items, err)
	}
	matches, err := nested.Glob("*.txt")
	if err != nil || len(matches) != 1 {
		t.Fatalf("Glob = (%v,%v)", matches, err)
	}
}

func TestPathRenameUnlink(t *testing.T) {
	dir := t.TempDir()
	oldPath := New(dir, "old.txt")
	newPath := New(dir, "new.txt")
	if err := oldPath.WriteText("x"); err != nil {
		t.Fatal(err)
	}
	renamed, err := oldPath.Rename(newPath)
	if err != nil || renamed != newPath {
		t.Fatalf("Rename = (%q,%v)", renamed, err)
	}
	if oldPath.Exists() || !newPath.Exists() {
		t.Fatal("Rename file state wrong")
	}
	if err := newPath.Unlink(); err != nil {
		t.Fatal(err)
	}
	if newPath.Exists() {
		t.Fatal("Unlink failed")
	}
	emptyDir := New(dir, "empty")
	if err := emptyDir.Mkdir(); err != nil {
		t.Fatal(err)
	}
	if err := emptyDir.Rmdir(); err != nil {
		t.Fatal(err)
	}
}

func TestPathAbsAndErrors(t *testing.T) {
	dir := t.TempDir()
	p := New(dir, "missing.txt")
	if _, err := p.ReadText(); os.IsNotExist(err) == false {
		t.Fatalf("ReadText missing err = %v", err)
	}
	abs, err := p.Abs()
	if err != nil || !abs.IsAbsolute() {
		t.Fatalf("Abs = (%q,%v)", abs, err)
	}
}
