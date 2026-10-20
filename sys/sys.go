// Package sys 提供 Python sys 模块风格的常用脚本接口。
package sys

import (
	"os"
	"runtime"
)

// 标准流，与 Python 的 sys.stdin/stdout/stderr 对应。
var (
	Stdin  = os.Stdin
	Stdout = os.Stdout
	Stderr = os.Stderr
)

// Args 是命令行参数，等价于 Python 的 sys.argv。
// 注意：它是在包初始化时从 os.Args 复制的；运行中修改 os.Args 不会同步。
var Args = append([]string(nil), os.Args...)

// Arg 返回第 i 个命令行参数；越界返回空字符串。
func Arg(i int) string {
	if i < 0 || i >= len(Args) {
		return ""
	}
	return Args[i]
}

// Exit 等价于 Python 的 sys.exit([code])。
func Exit(code int) {
	os.Exit(code)
}

// Getpid 等价于 Python 的 os.getpid()。
func Getpid() int { return os.Getpid() }

// Getppid 等价于 Python 的 os.getppid()。
func Getppid() int { return os.Getppid() }

// Platform 等价于 Python 的 sys.platform（Go 的 GOOS）。
func Platform() string { return runtime.GOOS }

// Version 等价于 Python 的 sys.version 的 Go 运行时版本部分。
func Version() string { return runtime.Version() }

// Maxsize 等价于 Python 的 sys.maxsize。
func Maxsize() int { return int(^uint(0) >> 1) }

// Getcwd 等价于 Python 的 os.getcwd()。
func Getcwd() (string, error) { return os.Getwd() }

// Chdir 等价于 Python 的 os.chdir(path)。
func Chdir(path string) error { return os.Chdir(path) }

// Executable 等价于 Python 的 sys.executable（当前进程可执行文件）。
func Executable() (string, error) { return os.Executable() }

// Hostname 等价于 Python 的 socket.gethostname()。
func Hostname() (string, error) { return os.Hostname() }
