//go:build windows

package sys

// Getuid 在 Windows 上返回 -1。
func Getuid() int { return -1 }

// Geteuid 在 Windows 上返回 -1。
func Geteuid() int { return -1 }

// Getgid 在 Windows 上返回 -1。
func Getgid() int { return -1 }

// Getegid 在 Windows 上返回 -1。
func Getegid() int { return -1 }
