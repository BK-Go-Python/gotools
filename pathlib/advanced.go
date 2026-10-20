package pathlib

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Resolve 等价于 Python 的 Path.resolve()：返回绝对路径并解析符号链接。
// 路径不存在时（strict=False 语义）返回绝对路径而不报错。
func (p Path) Resolve() (Path, error) {
	absolute, err := filepath.Abs(string(p))
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		if os.IsNotExist(err) {
			return Path(absolute), nil
		}
		return "", err
	}
	return Path(resolved), nil
}

// Relative_to 等价于 Python 的 Path.relative_to(other)。
func (p Path) Relative_to(other Path) (Path, error) {
	relative, err := filepath.Rel(string(other), string(p))
	if err != nil {
		return "", err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("pathlib: %q is not relative to %q", p, other)
	}
	return Path(relative), nil
}

// Is_relative_to 等价于 Python 的 Path.is_relative_to(other)。
func (p Path) Is_relative_to(other Path) bool {
	_, err := p.Relative_to(other)
	return err == nil
}

// With_name 等价于 Python 的 Path.with_name(name)：替换最后一段。
func (p Path) With_name(name string) (Path, error) {
	if name == "" || strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("pathlib.With_name: invalid name %q", name)
	}
	return p.Parent().Join(name), nil
}

// Match 等价于 Python 的 Path.match(pattern)：默认匹配文件名，
// pattern 含路径分隔符时匹配完整路径。
func (p Path) Match(pattern string) (bool, error) {
	target := p.Name()
	if strings.ContainsAny(pattern, `/\`) {
		target = string(p)
	}
	return filepath.Match(pattern, target)
}

// Rglob 等价于 Python 的 Path.rglob(pattern)：递归匹配。
func (p Path) Rglob(pattern string) ([]Path, error) {
	var out []Path
	err := filepath.WalkDir(string(p), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == string(p) {
			return nil
		}
		relative, err := filepath.Rel(string(p), path)
		if err != nil {
			return err
		}
		matched, err := filepath.Match(pattern, relative)
		if err != nil {
			return err
		}
		if !matched {
			matched, err = filepath.Match(pattern, entry.Name())
			if err != nil {
				return err
			}
		}
		if matched {
			out = append(out, Path(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
