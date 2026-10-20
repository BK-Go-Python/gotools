// Package shutil 提供 Python shutil 风格的文件树工具。
package shutil

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Copy 等价于 Python 的 shutil.copy(src, dst)：复制文件内容与权限。
// dst 是目录时复制到目录下。
func Copy(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("shutil.Copy: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("shutil.Copy: %s is a directory", src)
	}
	target, err := targetPath(src, dst)
	if err != nil {
		return err
	}
	return copyFile(src, target, info.Mode(), false)
}

// Copy2 等价于 Python 的 shutil.copy2(src, dst)：额外保留 mtime。
func Copy2(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("shutil.Copy2: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("shutil.Copy2: %s is a directory", src)
	}
	target, err := targetPath(src, dst)
	if err != nil {
		return err
	}
	return copyFile(src, target, info.Mode(), true)
}

// Copytree 等价于 Python 的 shutil.copytree(src, dst, dirs_exist_ok=False)。
func Copytree(src, dst string, dirsExistOK ...bool) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("shutil.Copytree: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("shutil.Copytree: %s is not a directory", src)
	}
	allowExisting := len(dirsExistOK) > 0 && dirsExistOK[0]
	if _, err := os.Stat(dst); err == nil {
		if !allowExisting {
			return fmt.Errorf("shutil.Copytree: destination exists: %s", dst)
		}
	} else if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return fmt.Errorf("shutil.Copytree: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("shutil.Copytree: %w", err)
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := Copytree(srcPath, dstPath, true); err != nil {
				return err
			}
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("shutil.Copytree: %w", err)
		}
		if err := copyFile(srcPath, dstPath, info.Mode(), true); err != nil {
			return fmt.Errorf("shutil.Copytree: %w", err)
		}
	}
	return nil
}

// Rmtree 等价于 Python 的 shutil.rmtree(path)；ignoreErrors=true 时忽略错误。
func Rmtree(path string, ignoreErrors ...bool) error {
	err := os.RemoveAll(path)
	if err != nil && !(len(ignoreErrors) > 0 && ignoreErrors[0]) {
		return fmt.Errorf("shutil.Rmtree: %w", err)
	}
	return nil
}

// Move 等价于 Python 的 shutil.move(src, dst)：优先 rename，跨设备时复制后删除。
func Move(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("shutil.Move: %w", err)
	}
	target, err := targetPath(src, dst)
	if err != nil {
		return err
	}
	if err := os.Rename(src, target); err == nil {
		return nil
	}
	if info.IsDir() {
		if err := Copytree(src, target); err != nil {
			return err
		}
	} else if err := copyFile(src, target, info.Mode(), true); err != nil {
		return fmt.Errorf("shutil.Move: %w", err)
	}
	if err := os.RemoveAll(src); err != nil {
		return fmt.Errorf("shutil.Move: %w", err)
	}
	return nil
}

// Which 等价于 Python 的 shutil.which(cmd)。
func Which(cmd string) (string, error) {
	return exec.LookPath(cmd)
}

func targetPath(src, dst string) (string, error) {
	info, err := os.Stat(dst)
	if err == nil && info.IsDir() {
		return filepath.Join(dst, filepath.Base(src)), nil
	}
	return dst, nil
}

func copyFile(src, dst string, mode os.FileMode, preserveTime bool) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	if preserveTime {
		if info, err := os.Stat(src); err == nil {
			_ = os.Chtimes(dst, info.ModTime(), info.ModTime())
		}
	}
	return nil
}
