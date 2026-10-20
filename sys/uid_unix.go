//go:build !windows

package sys

import "os"

// Getuid 等价于 Python 的 os.getuid()。
func Getuid() int { return os.Getuid() }

// Geteuid 等价于 Python 的 os.geteuid()。
func Geteuid() int { return os.Geteuid() }

// Getgid 等价于 Python 的 os.getgid()。
func Getgid() int { return os.Getgid() }

// Getegid 等价于 Python 的 os.getegid()。
func Getegid() int { return os.Getegid() }
