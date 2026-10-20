// Package pathlib 提供 Python pathlib.Path 风格的路径操作。
//
// Path 是 string 的具名类型，方法覆盖读写文本/字节、目录遍历、glob、
// 创建/删除/重命名等常用操作；不改变当前工作目录。
package pathlib

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BK-Go-Python/gotools/file"
)

// Path 是文件路径。
type Path string

// New 等价于 Python 的 Path(*parts)，用 filepath.Join 拼接。
func New(parts ...string) Path {
	if len(parts) == 0 {
		return Path(".")
	}
	return Path(filepath.Join(parts...))
}

// From 把字符串转换为 Path。
func From(path string) Path {
	return Path(path)
}

// String 返回路径字符串。
func (p Path) String() string {
	return string(p)
}

// Join 等价于 Python 的 p / "a" / "b"。
func (p Path) Join(parts ...string) Path {
	all := make([]string, 0, len(parts)+1)
	all = append(all, string(p))
	all = append(all, parts...)
	return Path(filepath.Join(all...))
}

// Parent 等价于 Python 的 p.parent。
func (p Path) Parent() Path {
	return Path(filepath.Dir(string(p)))
}

// Name 等价于 Python 的 p.name。
func (p Path) Name() string {
	return filepath.Base(string(p))
}

// Suffix 等价于 Python 的 p.suffix；点文件（如 .bashrc）返回空字符串。
func (p Path) Suffix() string {
	name := p.Name()
	ext := filepath.Ext(name)
	if ext == name {
		return ""
	}
	return ext
}

// Stem 等价于 Python 的 p.stem。
func (p Path) Stem() string {
	name := p.Name()
	suffix := p.Suffix()
	if suffix == "" {
		return name
	}
	return strings.TrimSuffix(name, suffix)
}

// WithSuffix 等价于 Python 的 p.with_suffix(suffix)。
func (p Path) WithSuffix(suffix string) Path {
	stem := strings.TrimSuffix(string(p), p.Suffix())
	return Path(stem + suffix)
}

// IsAbsolute 等价于 Python 的 p.is_absolute()。
func (p Path) IsAbsolute() bool {
	return filepath.IsAbs(string(p))
}

// Abs 等价于 Python 的 p.resolve() / absolute()。
func (p Path) Abs() (Path, error) {
	abs, err := filepath.Abs(string(p))
	if err != nil {
		return "", err
	}
	return Path(abs), nil
}

// Exists 等价于 Python 的 p.exists()。
func (p Path) Exists() bool { return file.Exists(string(p)) }

// IsFile 等价于 Python 的 p.is_file()。
func (p Path) IsFile() bool { return file.IsFile(string(p)) }

// IsDir 等价于 Python 的 p.is_dir()。
func (p Path) IsDir() bool { return file.IsDir(string(p)) }

// ReadText 等价于 Python 的 p.read_text()（UTF-8）。
func (p Path) ReadText() (string, error) { return file.Read(string(p)) }

// WriteText 等价于 Python 的 p.write_text(text)（UTF-8，覆盖写入）。
func (p Path) WriteText(text string) error { return file.Write(string(p), text) }

// ReadBytes 等价于 Python 的 p.read_bytes()。
func (p Path) ReadBytes() ([]byte, error) { return file.ReadBytes(string(p)) }

// WriteBytes 等价于 Python 的 p.write_bytes(data)。
func (p Path) WriteBytes(data []byte) error { return file.WriteBytes(string(p), data) }

// Mkdir 等价于 Python 的 p.mkdir()；parents=true 时等价于 mkdir(parents=True)。
func (p Path) Mkdir(parents ...bool) error {
	if len(parents) > 0 && parents[0] {
		return os.MkdirAll(string(p), 0o755)
	}
	return os.Mkdir(string(p), 0o755)
}

// Iterdir 等价于 Python 的 p.iterdir()，返回目录下的直接子项。
func (p Path) Iterdir() ([]Path, error) {
	names, err := file.ListDir(string(p))
	if err != nil {
		return nil, err
	}
	out := make([]Path, len(names))
	for i, name := range names {
		out[i] = p.Join(name)
	}
	return out, nil
}

// Glob 等价于 Python 的 p.glob(pattern)。
func (p Path) Glob(pattern string) ([]Path, error) {
	matches, err := filepath.Glob(filepath.Join(string(p), pattern))
	if err != nil {
		return nil, err
	}
	out := make([]Path, len(matches))
	for i, match := range matches {
		out[i] = Path(match)
	}
	return out, nil
}

// Size 等价于 Python 的 p.stat().st_size。
func (p Path) Size() (int64, error) { return file.Size(string(p)) }

// Touch 等价于 Python 的 p.touch()。
func (p Path) Touch() error { return file.Touch(string(p)) }

// Rename 等价于 Python 的 p.rename(target)，返回新路径。
func (p Path) Rename(target Path) (Path, error) {
	if err := file.Move(string(p), string(target)); err != nil {
		return "", err
	}
	return target, nil
}

// Unlink 等价于 Python 的 p.unlink()，删除文件。
func (p Path) Unlink() error { return os.Remove(string(p)) }

// Rmdir 等价于 Python 的 p.rmdir()，删除空目录。
func (p Path) Rmdir() error { return os.Remove(string(p)) }

// Equal 比较两个路径字符串。
func (p Path) Equal(other Path) bool { return p == other }
