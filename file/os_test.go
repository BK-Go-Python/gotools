package file

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
)

func TestStatChmodRename(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	if err := Write(path, "hello"); err != nil {
		t.Fatal(err)
	}
	info, err := Stat(path)
	if err != nil || info.Size() != 5 {
		t.Fatalf("Stat = (%v,%v)", info, err)
	}
	if err := Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		info, err = Stat(path)
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatalf("Chmod mode = (%v,%v)", info.Mode(), err)
		}
	}

	renamed := filepath.Join(dir, "b.txt")
	if err := Rename(path, renamed); err != nil {
		t.Fatal(err)
	}
	if Exists(path) || !Exists(renamed) {
		t.Fatal("Rename failed")
	}
	if runtime.GOOS != "windows" {
		if err := Chown(renamed, os.Getuid(), os.Getgid()); err != nil {
			t.Logf("Chown skipped: %v", err)
		}
	}
}

func TestWalk(t *testing.T) {
	dir := t.TempDir()
	if err := MkdirP(filepath.Join(dir, "sub")); err != nil {
		t.Fatal(err)
	}
	if err := Write(filepath.Join(dir, "a.txt"), "a"); err != nil {
		t.Fatal(err)
	}
	if err := Write(filepath.Join(dir, "sub", "b.txt"), "b"); err != nil {
		t.Fatal(err)
	}

	var walked []string
	if err := Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		walked = append(walked, info.Name())
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	sort.Strings(walked)
	want := []string{"a.txt", "b.txt", "sub", filepath.Base(dir)}
	sort.Strings(want)
	if len(walked) != len(want) {
		t.Fatalf("Walk = %v", walked)
	}

	var dirs int
	if err := WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			dirs++
		}
		return nil
	}); err != nil || dirs != 2 {
		t.Fatalf("WalkDir dirs = (%d,%v)", dirs, err)
	}
}

func TestSymlinkSameFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink may require privileges on Windows")
	}
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	link := filepath.Join(dir, "link.txt")
	if err := Write(target, "x"); err != nil {
		t.Fatal(err)
	}
	if err := Symlink(target, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if !IsSymlink(link) {
		t.Fatal("IsSymlink failed")
	}
	if got, err := Readlink(link); err != nil || got != target {
		t.Fatalf("Readlink = (%q,%v)", got, err)
	}
	if !SameFile(target, link) {
		t.Fatal("SameFile failed")
	}
	if _, err := Lstat(link); err != nil {
		t.Fatal(err)
	}
}
