//go:build windows

package shutil

import "fmt"

// DiskUsage 在 Windows 上未实现。
func DiskUsage(path string) (total, used, free uint64, err error) {
	return 0, 0, 0, fmt.Errorf("shutil.DiskUsage: not implemented on windows")
}
