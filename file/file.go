// Package file 提供 Python os / os.path / shutil 风格的常用文件操作：
// 读写文本、按行读写、判断存在性、递归建目录、复制、移动、列目录等。
package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	stdtime "time"
)

// Read 等价于 Python 的 open(path).read()，读取整个文件为字符串。
func Read(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadBytes 读取整个文件的字节内容。
func ReadBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Write 等价于 Python 的 open(path, "w").write(content)，会覆盖已有内容。
// 文件权限为 0644，父目录必须存在。
func Write(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

// WriteBytes 写入字节内容，会覆盖已有内容。
func WriteBytes(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}

// Append 等价于 Python 的 open(path, "a").write(content)。
func Append(path, content string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// ReadLines 等价于 Python 的 readlines()，返回按行切分、不含行尾的结果。
func ReadLines(path string) ([]string, error) {
	content, err := Read(path)
	if err != nil {
		return nil, err
	}
	return splitLines(content), nil
}

// WriteLines 按行写入文件，每行之间以及文件末尾都会补上 "\n"。
func WriteLines(path string, lines []string) error {
	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}
	return Write(path, content)
}

// Exists 等价于 Python 的 os.path.exists(path)。
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

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

// MkdirP 等价于 Python 的 os.makedirs(path, exist_ok=True)。
func MkdirP(path string) error {
	return os.MkdirAll(path, 0o755)
}

// ListDir 等价于 Python 的 os.listdir(path)，返回排序后的文件名（不含路径）。
func ListDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names, nil
}

// Glob 等价于 Python 的 glob.glob(pattern)。
func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// Copy 等价于 Python 的 shutil.copy(src, dst)，会保留源文件权限并覆盖目标。
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("file.Copy: %q is a directory", src)
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// Move 等价于 Python 的 shutil.move；跨文件系统时退化为复制后删除。
func Move(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	if err := Copy(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

// Remove 等价于 Python 的 os.remove(path)，只能删除文件或空目录。
func Remove(path string) error {
	return os.Remove(path)
}

// RemoveAll 等价于 Python 的 shutil.rmtree(path)，不存在时返回 nil。
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Size 等价于 Python 的 os.path.getsize(path)。
func Size(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ModTime 返回文件的修改时间。
func ModTime(path string) (stdtime.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return stdtime.Time{}, err
	}
	return info.ModTime(), nil
}

// Abs 等价于 Python 的 os.path.abspath(path)。
func Abs(path string) (string, error) {
	return filepath.Abs(path)
}

// Cwd 等价于 Python 的 os.getcwd()。
func Cwd() (string, error) {
	return os.Getwd()
}

// Home 返回当前用户的主目录。
func Home() (string, error) {
	return os.UserHomeDir()
}

// TempFile 创建临时文件，dir 为空时使用系统临时目录。
func TempFile(dir, pattern string) (*os.File, error) {
	return os.CreateTemp(dir, pattern)
}

// TempDir 创建临时目录，dir 为空时使用系统临时目录。
func TempDir(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

// Touch 等价于 Python 的 pathlib.Path.touch()：不存在则创建，存在则更新修改时间。
func Touch(path string) error {
	now := stdtime.Now()
	if _, err := os.Stat(path); err == nil {
		return os.Chtimes(path, now, now)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	return f.Close()
}

func splitLines(content string) []string {
	if content == "" {
		return []string{}
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
