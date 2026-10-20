package subprocess

import (
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunEcho(t *testing.T) {
	result, err := Run("echo", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if result.Returncode != 0 || !result.Success() {
		t.Fatalf("result = %#v", result)
	}
	if strings.TrimSpace(result.Stdout) != "hello" {
		t.Fatalf("stdout = %q", result.Stdout)
	}
}

func TestNonZeroExit(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sh not available")
	}
	result, err := Run("sh", "-c", "echo oops >&2; exit 3")
	if err != nil {
		t.Fatal(err)
	}
	if result.Returncode != 3 || result.Success() {
		t.Fatalf("result = %#v", result)
	}
	if !strings.Contains(result.Stderr, "oops") {
		t.Fatalf("stderr = %q", result.Stderr)
	}
	if _, err := CheckOutput("sh", "-c", "exit 4"); err == nil {
		t.Fatal("CheckOutput should error")
	}
}

func TestOptions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sh not available")
	}
	result, err := RunWithOptions(Options{Stdin: "hi"}, "cat")
	if err != nil || result.Stdout != "hi" {
		t.Fatalf("stdin = (%#v,%v)", result, err)
	}

	result, err = RunWithOptions(Options{Env: []string{"FOO=bar"}}, "sh", "-c", `printf %s "$FOO"`)
	if err != nil || result.Stdout != "bar" {
		t.Fatalf("env = (%#v,%v)", result, err)
	}

	if _, err := RunWithOptions(Options{Dir: t.TempDir()}, "pwd"); err != nil {
		t.Fatalf("dir = %v", err)
	}

	start := time.Now()
	_, err = RunWithOptions(Options{Timeout: 50 * time.Millisecond}, "sh", "-c", "sleep 1")
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("timeout err = %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatal("timeout did not stop command")
	}
}

func TestShellAndErrors(t *testing.T) {
	result, err := RunShell("echo shell")
	if err != nil || !strings.Contains(result.Stdout, "shell") {
		t.Fatalf("RunShell = (%#v,%v)", result, err)
	}
	if _, err := Run("definitely-not-a-command-xyz"); err == nil {
		t.Fatal("missing command should error")
	}
}
