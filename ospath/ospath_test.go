package ospath

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestStringOps(t *testing.T) {
	if got := Join("a", "b", "c"); got != filepath.Join("a", "b", "c") {
		t.Fatalf("Join = %q", got)
	}
	if Base("/a/b.txt") != "b.txt" || Dir("/a/b.txt") != "/a" {
		t.Fatal("Base/Dir failed")
	}
	if Ext("a/b.txt") != ".txt" || Ext(".bashrc") != "" {
		t.Fatalf("Ext = %q %q", Ext("a/b.txt"), Ext(".bashrc"))
	}
	root, ext := Splitext("a/b.txt")
	if root != "a/b" || ext != ".txt" {
		t.Fatalf("Splitext = (%q,%q)", root, ext)
	}
	dir, file := Split("a/b.txt")
	if dir != "a/" || file != "b.txt" {
		t.Fatalf("Split = (%q,%q)", dir, file)
	}
	if !IsAbs("/a") || IsAbs("a") {
		t.Fatal("IsAbs failed")
	}
	if Normpath("a/../b") != "b" {
		t.Fatalf("Normpath = %q", Normpath("a/../b"))
	}
}

func TestFilesystem(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !Exists(path) || !IsFile(path) || IsDir(path) {
		t.Fatal("predicates failed")
	}
	abs, err := Abs(path)
	if err != nil || !IsAbs(abs) {
		t.Fatalf("Abs = (%q,%v)", abs, err)
	}
	rel, err := Rel(dir, path)
	if err != nil || rel != "a.txt" {
		t.Fatalf("Rel = (%q,%v)", rel, err)
	}
	if !SameFile(path, path) || SameFile(path, filepath.Join(dir, "missing")) {
		t.Fatal("SameFile failed")
	}
}

func TestExpand(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip(err)
	}
	if got, err := Expanduser("~/x"); err != nil || got != filepath.Join(home, "x") {
		t.Fatalf("Expanduser = (%q,%v)", got, err)
	}
	if _, err := Expanduser("~nobody/x"); err == nil {
		t.Fatal("~user should error")
	}
	t.Setenv("GOTOOLS_OSPATH_TEST", "value")
	if got := Expandvars("$GOTOOLS_OSPATH_TEST/x"); got != "value/x" {
		t.Fatalf("Expandvars = %q", got)
	}
	if runtime.GOOS == "windows" && strings.Contains(Expandvars("$X"), "\\") {
		t.Skip("windows env syntax")
	}
}
