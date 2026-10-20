package py

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// Print 等价于 Python 的 print(*args)，默认以空格分隔并换行。
//
// Python: print(*args)
func Print(args ...any) error {
	return PrintTo(os.Stdout, " ", "\n", args...)
}

// PrintTo 等价于 Python 的 print(*args, sep=..., end=..., file=...)。
//
// Python: print(*args, sep=" ", end="\n", file=w)
// DIFF: 使用 io.Writer 代替 Python 的文件对象。
func PrintTo(w io.Writer, sep, end string, args ...any) error {
	if w == nil {
		return fmt.Errorf("py.PrintTo: nil writer: %w", ErrType)
	}
	for i, arg := range args {
		if i > 0 {
			if _, err := io.WriteString(w, sep); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, pyval.Str(arg)); err != nil {
			return err
		}
	}
	_, err := io.WriteString(w, end)
	return err
}

// Input 等价于 Python 的 input(prompt)，从标准输入读取一行。
//
// Python: input(prompt)
// DIFF: 不支持 readline 行编辑；EOF 时返回 error（Python 抛 EOFError）。
func Input(prompt string) (string, error) {
	return inputFrom(os.Stdin, os.Stdout, prompt)
}

func inputFrom(r io.Reader, w io.Writer, prompt string) (string, error) {
	if prompt != "" && w != nil {
		if _, err := io.WriteString(w, prompt); err != nil {
			return "", err
		}
	}
	reader := bufio.NewReaderSize(r, 1<<20)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		if err == io.EOF {
			return "", fmt.Errorf("py.Input: EOF when reading a line: %w", err)
		}
		return "", err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	return line, nil
}
