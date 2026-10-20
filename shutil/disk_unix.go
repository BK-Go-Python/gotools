//go:build !windows

package shutil

import "syscall"

// DiskUsage 返回 path 所在文件系统的总容量、已用、可用字节数。
func DiskUsage(path string) (total, used, free uint64, err error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0, 0, err
	}
	blockSize := uint64(stat.Bsize)
	total = stat.Blocks * blockSize
	free = stat.Bavail * blockSize
	used = (stat.Blocks - stat.Bfree) * blockSize
	return total, used, free, nil
}
