package file

import (
	"fmt"
	"os"
)

// Open 近似等价于 Python 的 open(path, mode)。
//
// mode 默认为 "r"，支持 r/w/a 与 +、b、t 的组合。
//
// Python: open(path, mode="r")
// DIFF: 不支持 encoding、newline、buffering 参数；文件对象是 *os.File。
func Open(path string, mode ...string) (*os.File, error) {
	selected := "r"
	if len(mode) > 0 {
		selected = mode[0]
	}
	flag, err := parseOpenMode(selected)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(path, flag, 0o644)
}

func parseOpenMode(mode string) (int, error) {
	if mode == "" {
		return 0, fmt.Errorf("file.Open: empty mode")
	}

	var flag int
	switch mode[0] {
	case 'r':
		flag = os.O_RDONLY
	case 'w':
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case 'a':
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	default:
		return 0, fmt.Errorf("file.Open: invalid mode %q", mode)
	}

	for _, char := range mode[1:] {
		switch char {
		case '+':
			flag = (flag &^ os.O_WRONLY) | os.O_RDWR
		case 'b', 't':
			// 二进制/文本模式在 Go 中没有区别。
		default:
			return 0, fmt.Errorf("file.Open: invalid mode %q", mode)
		}
	}
	return flag, nil
}
