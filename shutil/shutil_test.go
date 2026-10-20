package shutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.txt")
	dst := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Copy(src, dst); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(dst)
	if err != nil || string(data) != "hello" {
		t.Fatalf("Copy = (%q,%v)", data, err)
	}

	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Copy2(src, sub); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(sub, "a.txt")); err != nil {
		t.Fatal(err)
	}
	if err := Copy(dir, dst); err == nil {
		t.Fatal("copy directory should error")
	}
}

func TestCopytreeMoveRmtree(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	if err := os.MkdirAll(filepath.Join(src, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "nested", "b.txt"), []byte("b"), 0o644); err != nil {
		t.Fatal(err)
	}

	dst := filepath.Join(dir, "dst")
	if err := Copytree(src, dst, true); err != nil {
		t.Fatal(err)
	}
	if err := Copytree(src, dst); err == nil {
		t.Fatal("existing destination should error")
	}
	if err := Copytree(src, dst, true); err != nil {
		t.Fatalf("dirsExistOK = %v", err)
	}
	if data, err := os.ReadFile(filepath.Join(dst, "nested", "b.txt")); err != nil || string(data) != "b" {
		t.Fatalf("copied file = (%q,%v)", data, err)
	}

	moved := filepath.Join(dir, "moved")
	if err := Move(dst, moved); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Fatal("Move left source")
	}
	if err := Rmtree(moved); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(moved); !os.IsNotExist(err) {
		t.Fatal("Rmtree failed")
	}
}

func TestWhichDiskUsage(t *testing.T) {
	if _, err := Which("definitely-not-a-command-xyz"); err == nil {
		t.Fatal("missing command should error")
	}
	if runtime.GOOS != "windows" {
		path, err := Which("go")
		if err != nil || path == "" {
			t.Fatalf("Which go = (%q,%v)", path, err)
		}
		total, used, free, err := DiskUsage(t.TempDir())
		if err != nil || total == 0 || used > total || free > total {
			t.Fatalf("DiskUsage = (%d,%d,%d,%v)", total, used, free, err)
		}
	}
}
