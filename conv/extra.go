package conv

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

// BytesHex 等价于 Python 的 bytes.hex()，返回小写十六进制字符串。
//
// Python: b.hex()
func BytesHex(data []byte) string {
	return hex.EncodeToString(data)
}

// FromHex 等价于 Python 的 bytes.fromhex(s)，忽略输入中的空白。
//
// Python: bytes.fromhex(s)
func FromHex(s string) ([]byte, error) {
	cleaned := strings.Map(func(r rune) rune {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return -1
		}
		return r
	}, s)
	data, err := hex.DecodeString(cleaned)
	if err != nil {
		return nil, fmt.Errorf("conv.FromHex: invalid hex string: %w", err)
	}
	return data, nil
}

// IsInteger 等价于 Python 的 float.is_integer()。
//
// Python: f.is_integer()
func IsInteger(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0) && math.Trunc(f) == f
}
