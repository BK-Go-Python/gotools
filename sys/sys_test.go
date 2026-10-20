package sys

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestArgsAndMeta(t *testing.T) {
	old := Args
	defer func() { Args = old }()
	Args = []string{"prog", "a", "b"}
	if Arg(1) != "a" || Arg(3) != "" || Arg(-1) != "" {
		t.Fatalf("Arg = %q %q", Arg(1), Arg(3))
	}
	if Platform() != runtime.GOOS || Version() != runtime.Version() {
		t.Fatal("Platform/Version mismatch")
	}
	if Maxsize() <= 0 || Getpid() <= 0 || Getppid() < 0 {
		t.Fatal("pid/maxsize invalid")
	}
	if runtime.GOOS != "windows" && (Getuid() < 0 || Getgid() < 0) {
		t.Fatal("uid/gid invalid")
	}
}

func TestCwdChdir(t *testing.T) {
	original, err := Getcwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Chdir(original); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	}()

	dir := t.TempDir()
	if err := Chdir(dir); err != nil {
		t.Fatal(err)
	}
	current, err := Getcwd()
	if err != nil {
		t.Fatal(err)
	}
	resolvedCurrent, err := filepath.EvalSymlinks(current)
	if err != nil {
		t.Fatal(err)
	}
	resolvedDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	if resolvedCurrent != resolvedDir {
		t.Fatalf("cwd = %q, want %q", current, dir)
	}
}

func TestExecutableHostname(t *testing.T) {
	if _, err := Executable(); err != nil {
		t.Fatalf("Executable: %v", err)
	}
	if _, err := Hostname(); err != nil {
		t.Fatalf("Hostname: %v", err)
	}
	if Stdin != os.Stdin || Stdout != os.Stdout || Stderr != os.Stderr {
		t.Fatal("stdio mismatch")
	}
}
