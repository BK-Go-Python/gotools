package file

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Stat 等价于 Python 的 os.stat(path)。
func Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Lstat 等价于 Python 的 os.lstat(path)，不跟随符号链接。
func Lstat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

// Chmod 等价于 Python 的 os.chmod(path, mode)。
func Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// Chown 等价于 Python 的 os.chown(path, uid, gid)。
//
// 在 Windows 上会返回错误；普通用户通常只能修改自己的文件。
func Chown(path string, uid, gid int) error {
	return os.Chown(path, uid, gid)
}

// Rename 等价于 Python 的 os.rename(src, dst)。
func Rename(src, dst string) error {
	return os.Rename(src, dst)
}

// Walk 等价于 Python 的 os.walk 的 Go 回调形式，递归遍历 root。
//
// fn 收到路径、FileInfo 和可能的错误；返回非 nil error 会终止遍历。
func Walk(root string, fn func(path string, info os.FileInfo, err error) error) error {
	return filepath.Walk(root, fn)
}

// WalkDir 与 Walk 类似，但使用 fs.DirEntry，避免额外 stat。
func WalkDir(root string, fn func(path string, d fs.DirEntry, err error) error) error {
	return filepath.WalkDir(root, fn)
}

// IsSymlink 判断 path 是否是符号链接（不跟随）。
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

// Readlink 等价于 Python 的 os.readlink(path)。
func Readlink(path string) (string, error) {
	return os.Readlink(path)
}

// Symlink 等价于 Python 的 os.symlink(oldname, newname)。
func Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

// SameFile 等价于 Python 的 os.path.samefile(a, b)。
func SameFile(a, b string) bool {
	infoA, err := os.Stat(a)
	if err != nil {
		return false
	}
	infoB, err := os.Stat(b)
	if err != nil {
		return false
	}
	return os.SameFile(infoA, infoB)
}
