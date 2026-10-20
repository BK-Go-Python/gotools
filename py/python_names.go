package py

import "io"

// 本文件提供 Python 拼写别名：Go 导出标识符必须大写开头，因此保留
// Python 名字的其余拼写（例如 bit_count -> Bit_count）。

// Bit_length 是 BitLength 的 Python 拼写别名。
func Bit_length[T Integer](x T) int { return BitLength(x) }

// Bit_count 是 BitCount 的 Python 拼写别名。
func Bit_count[T Integer](x T) int { return BitCount(x) }

// Print_to 是 PrintTo 的 Python 拼写别名。
func Print_to(w io.Writer, sep, end string, args ...any) error {
	return PrintTo(w, sep, end, args...)
}

// Isinstance 是 IsInstance 的 Python 拼写别名。
func Isinstance(v any, types ...any) bool { return IsInstance(v, types...) }

// Issubclass 是 IsSubclass 的 Python 拼写别名。
func Issubclass(child, parent any) bool { return IsSubclass(child, parent) }

// Getattr 是 GetAttr 的 Python 拼写别名。
func Getattr(v any, name string, def ...any) (any, error) {
	return GetAttr(v, name, def...)
}

// Setattr 是 SetAttr 的 Python 拼写别名。
func Setattr(v any, name string, value any) error { return SetAttr(v, name, value) }

// Hasattr 是 HasAttr 的 Python 拼写别名。
func Hasattr(v any, name string) bool { return HasAttr(v, name) }

// Delattr 是 DelAttr 的 Python 拼写别名。
func Delattr(v any, name string) error { return DelAttr(v, name) }

// Memoryview 是 MemoryView 的 Python 拼写别名。
func Memoryview(b []byte) []byte { return MemoryView(b) }
