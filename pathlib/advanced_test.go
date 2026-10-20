package pathlib

import (
	"path/filepath"
	"testing"
)

func TestResolveRelative(t *testing.T) {
	dir := t.TempDir()
	file := New(dir, "a.txt")
	if err := file.WriteText("x"); err != nil {
		t.Fatal(err)
	}
	resolved, err := file.Resolve()
	if err != nil || !resolved.IsAbsolute() {
		t.Fatalf("Resolve = (%q,%v)", resolved, err)
	}
	relative, err := file.Relative_to(From(dir))
	if err != nil || relative != Path("a.txt") {
		t.Fatalf("Relative_to = (%q,%v)", relative, err)
	}
	if !file.Is_relative_to(From(dir)) {
		t.Fatal("Is_relative_to failed")
	}
	outside, err := From(filepath.Dir(dir)).Relative_to(From(dir))
	if err == nil {
		t.Fatalf("outside relative should error: %q", outside)
	}
}

func TestWithNameMatch(t *testing.T) {
	p := New("a", "b.txt")
	renamed, err := p.With_name("c.md")
	if err != nil || renamed != Path(filepath.Join("a", "c.md")) {
		t.Fatalf("With_name = (%q,%v)", renamed, err)
	}
	if _, err := p.With_name("a/b"); err == nil {
		t.Fatal("separator name should error")
	}
	if matched, err := p.Match("*.txt"); err != nil || !matched {
		t.Fatalf("Match = (%v,%v)", matched, err)
	}
	if matched, err := p.Match("*.md"); err != nil || matched {
		t.Fatalf("Match miss = (%v,%v)", matched, err)
	}
}

func TestRglob(t *testing.T) {
	dir := t.TempDir()
	if err := New(dir, "a.txt").WriteText("a"); err != nil {
		t.Fatal(err)
	}
	if err := New(dir, "sub").Mkdir(); err != nil {
		t.Fatal(err)
	}
	if err := New(dir, "sub", "b.txt").WriteText("b"); err != nil {
		t.Fatal(err)
	}
	if err := New(dir, "sub", "c.md").WriteText("c"); err != nil {
		t.Fatal(err)
	}
	matches, err := New(dir).Rglob("*.txt")
	if err != nil || len(matches) != 2 {
		t.Fatalf("Rglob = (%v,%v)", matches, err)
	}
	matches, err = New(dir).Rglob("sub/*.md")
	if err != nil || len(matches) != 1 || matches[0].Name() != "c.md" {
		t.Fatalf("Rglob nested = (%v,%v)", matches, err)
	}
}
