// Package subprocess 提供 Python subprocess 风格的高频命令执行封装。
//
// 与 Python 一样，命令默认不经 shell；需要 shell 时显式使用 RunShell。
package subprocess

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Options 是运行命令的可选配置。
type Options struct {
	Timeout time.Duration
	Dir     string
	Env     []string
	Stdin   string
}

// Result 是命令执行结果，对应 Python 的 CompletedProcess。
type Result struct {
	Args       []string
	Returncode int
	Stdout     string
	Stderr     string
}

// Success 判断退出码是否为 0。
func (r Result) Success() bool { return r.Returncode == 0 }

// Run 执行命令并捕获 stdout/stderr；非零退出码不会返回 error。
func Run(name string, args ...string) (Result, error) {
	return RunWithOptions(Options{}, name, args...)
}

// RunWithOptions 执行命令并应用超时、工作目录、环境变量和 stdin。
func RunWithOptions(opts Options, name string, args ...string) (Result, error) {
	result := Result{Args: append([]string{name}, args...)}
	ctx := context.Background()
	cancel := func() {}
	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if opts.Stdin != "" {
		cmd.Stdin = strings.NewReader(opts.Stdin)
	}
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}
	if opts.Env != nil {
		cmd.Env = opts.Env
	}

	err := cmd.Run()
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()
	if cmd.ProcessState != nil {
		result.Returncode = cmd.ProcessState.ExitCode()
	}
	if err == nil {
		return result, nil
	}
	if ctx.Err() == context.DeadlineExceeded {
		return result, fmt.Errorf("subprocess: command timed out after %s", opts.Timeout)
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return result, nil
	}
	return result, fmt.Errorf("subprocess: %w", err)
}

// CheckOutput 执行命令并在非零退出码时返回 error；成功时返回 stdout。
func CheckOutput(name string, args ...string) (string, error) {
	result, err := Run(name, args...)
	if err != nil {
		return "", err
	}
	if result.Returncode != 0 {
		return "", fmt.Errorf("subprocess: %s exited with %d: %s", name, result.Returncode, strings.TrimSpace(result.Stderr))
	}
	return result.Stdout, nil
}

// CheckCall 执行命令并在非零退出码时返回 error。
func CheckCall(name string, args ...string) error {
	_, err := CheckOutput(name, args...)
	return err
}

// RunShell 通过系统 shell 执行命令字符串。
//
// Unix 使用 sh -c，Windows 使用 cmd /C。仅用于可信输入。
func RunShell(command string) (Result, error) {
	if runtime.GOOS == "windows" {
		return Run("cmd", "/C", command)
	}
	return Run("sh", "-c", command)
}
