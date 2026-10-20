package str

import (
	"fmt"
	"strings"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// ASCII 等价于 Python 的 ascii(v)：返回只包含 ASCII 字符的表示。
//
// 非 ASCII 字符按范围转义为 \xNN、\uNNNN 或 \UNNNNNNNN。
//
// Python: ascii(v)
func ASCII(v any) string {
	repr := pyval.Repr(v)
	var b strings.Builder
	b.Grow(len(repr))
	for _, r := range repr {
		switch {
		case r <= 0x7F:
			b.WriteRune(r)
		case r <= 0xFF:
			fmt.Fprintf(&b, `\x%02x`, r)
		case r <= 0xFFFF:
			fmt.Fprintf(&b, `\u%04x`, r)
		default:
			fmt.Fprintf(&b, `\U%08x`, r)
		}
	}
	return b.String()
}
