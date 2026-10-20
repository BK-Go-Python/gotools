package str

import (
	"fmt"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// FormatValue 等价于 Python 内置函数 format(v, spec)。
//
// 与 str.Format（模板替换）不同，这里的 spec 是 Python 格式说明符：
// [[fill]align][width][.precision][type]，支持 s/d/f/e/g/x/X/o/b/%。
//
// Python: format(v, spec)
// DIFF: 不支持千位分隔符 ","、符号标志 "+"、以及自定义对象的 __format__。
func FormatValue(v any, spec string) (string, error) {
	if spec == "" {
		return pyval.Str(v), nil
	}
	parsed, ok := parseFormatSpec(spec)
	if !ok || !isSupportedFormatType(parsed.typ) {
		return "", fmt.Errorf("str.FormatValue: invalid format spec %q", spec)
	}
	switch parsed.typ {
	case 'd', 'x', 'X', 'o', 'b':
		if _, ok := toInt64(v); !ok {
			return "", fmt.Errorf("str.FormatValue: %q spec requires an integer, got %T", spec, v)
		}
	case 'f', 'F', 'e', 'E', 'g', 'G', '%':
		if _, ok := toFloat64(v); !ok {
			return "", fmt.Errorf("str.FormatValue: %q spec requires a number, got %T", spec, v)
		}
	}
	return formatValue(v, spec), nil
}

func isSupportedFormatType(typ byte) bool {
	switch typ {
	case 0, 's', 'd', 'x', 'X', 'o', 'b', 'f', 'F', 'e', 'E', 'g', 'G', '%':
		return true
	}
	return false
}
