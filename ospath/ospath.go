// Package ospath 提供 Python os.path 风格的路径字符串工具。
package ospath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Join 等价于 Python 的 os.path.join(*parts)。
func Join(parts ...string) string { return filepath.Join(parts...) }

// Base 等价于 Python 的 os.path.basename(path)。
func Base(path string) string { return filepath.Base(path) }

// Dir 等价于 Python 的 os.path.dirname(path)。
func Dir(path string) string { return filepath.Dir(path) }

// Ext 等价于 Python 的 os.path.splitext(path)[1]；点文件返回空。
func Ext(path string) string {
	ext := filepath.Ext(filepath.Base(path))
	if ext == filepath.Base(path) {
		return ""
	}
	return ext
}

// Split 等价于 Python 的 os.path.split(path)，返回目录和文件名。
func Split(path string) (dir, file string) { return filepath.Split(path) }

// Splitext 等价于 Python 的 os.path.splitext(path)，返回根和扩展名。
func Splitext(path string) (root, ext string) {
	ext = Ext(path)
	return strings.TrimSuffix(path, ext), ext
}

// IsAbs 等价于 Python 的 os.path.isabs(path)。
func IsAbs(path string) bool { return filepath.IsAbs(path) }

// Exists 等价于 Python 的 os.path.exists(path)。
func Exists(path string) bool { return fileExists(path) }

// IsFile 等价于 Python 的 os.path.isfile(path)。
func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// IsDir 等价于 Python 的 os.path.isdir(path)。
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// Abs 等价于 Python 的 os.path.abspath(path)。
func Abs(path string) (string, error) { return filepath.Abs(path) }

// Rel 等价于 Python 的 os.path.relpath(path, start)。
func Rel(base, target string) (string, error) { return filepath.Rel(base, target) }

// Normpath 等价于 Python 的 os.path.normpath(path)。
func Normpath(path string) string { return filepath.Clean(path) }

// Expanduser 等价于 Python 的 os.path.expanduser(path)：只支持 ~ 和 ~/。
func Expanduser(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		if strings.HasPrefix(path, "~") {
			return "", fmt.Errorf("ospath.Expanduser: ~user is not supported")
		}
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if path == "~" {
		return home, nil
	}
	return filepath.Join(home, path[2:]), nil
}

// Expandvars 等价于 Python 的 os.path.expandvars(path)。
func Expandvars(path string) string { return os.ExpandEnv(path) }

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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
