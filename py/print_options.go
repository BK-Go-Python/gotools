package py

import (
	"io"
	"os"
)

// Print_options 对应 Python print() 的 sep/end/file/flush 选项。
//
// Go 没有关键字参数，因此用结构体表达；零值不是 Python 默认值，
// 请用 Default_print_options() 获取 " "、"\\n"、stdout 的默认配置。
type Print_options struct {
	Sep   string
	End   string
	File  io.Writer
	Flush bool
}

// Default_print_options 返回 Python print 的默认选项：
// sep=" "、end="\n"、file=os.Stdout、flush=false。
func Default_print_options() Print_options {
	return Print_options{Sep: " ", End: "\n", File: os.Stdout}
}

// Print_with 等价于 Python 的 print(*args, sep=..., end=..., file=..., flush=...)。
//
// File 为 nil 时写 os.Stdout；Flush 为 true 且 File 实现 Flush() error 时调用它。
func Print_with(opts Print_options, args ...any) error {
	w := opts.File
	if w == nil {
		w = os.Stdout
	}
	if err := PrintTo(w, opts.Sep, opts.End, args...); err != nil {
		return err
	}
	if opts.Flush {
		if flusher, ok := w.(interface{ Flush() error }); ok {
			return flusher.Flush()
		}
	}
	return nil
}
