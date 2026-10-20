package py

import "github.com/BK-Go-Python/gotools/str"

// Format 等价于 Python 内置函数 format(v, spec)。
//
// Python: format(v, spec)
// DIFF: 仅支持格式说明符子集，详见 str.FormatValue。
func Format(v any, spec string) (string, error) {
	return str.FormatValue(v, spec)
}
