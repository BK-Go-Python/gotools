package conv

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// FloatHex 等价于 Python 的 float.hex()，返回 C99 风格的十六进制浮点
// 字符串，例如 1.0 -> "0x1.0000000000000p+0"，0.1 ->
// "0x1.999999999999ap-4"。
//
// 尾数固定为 13 位十六进制数字；次正规数写作 "0x0.<尾数>p-1022"；
// +0.0 与 -0.0 分别返回 "0x0.0p+0" 与 "-0x0.0p+0"；±Inf 与 NaN 返回
// "inf"、"-inf" 与 "nan"。返回值可被 FromFloatHex 精确还原。
//
// Python: f.hex()
func FloatHex(f float64) string {
	switch {
	case math.IsNaN(f):
		return "nan"
	case math.IsInf(f, 1):
		return "inf"
	case math.IsInf(f, -1):
		return "-inf"
	case f == 0:
		if math.Signbit(f) {
			return "-0x0.0p+0"
		}
		return "0x0.0p+0"
	}

	bits := math.Float64bits(f)
	sign := ""
	if bits>>63 != 0 {
		sign = "-"
	}
	exponentBits := int((bits >> 52) & 0x7ff)
	mantissa := bits & ((1 << 52) - 1)
	leading := 1
	exponent := exponentBits - 1023
	if exponentBits == 0 {
		// 次正规数：0x0.<52 位尾数>p-1022。
		leading = 0
		exponent = -1022
	}
	return fmt.Sprintf("%s0x%d.%013xp%+d", sign, leading, mantissa, exponent)
}

// FromFloatHex 解析 FloatHex（Python float.hex()）的输出和常见的十六进制
// 浮点写法，如 "0x1.8p+1"、"0X1P-2"；同时兼容 Python float.fromhex 允许
// 省略的 "p" 指数（"0x1.8"）或 "0x" 前缀（"1.8p+1"），并忽略首尾空白。
//
// 换算委托给 strconv.ParseFloat，因此也可解析十进制写法；"inf"、"-inf"、
// "nan" 等特殊值遵循 ParseFloat 的语义。空串、非法写法或溢出 float64 时
// 返回 error。
//
// Python: float.fromhex(s)
// DIFF: 无 "0x" 前缀且无 "p" 指数的字符串（如 "1.8"）按十进制而非十六
// 进制解析；"-nan" 按 strconv.ParseFloat 的规则视为非法。
func FromFloatHex(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("conv.FromFloatHex: empty string")
	}

	original := s
	s = normalizeHexFloat(s)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("conv.FromFloatHex: invalid floating-point string %q: %w", original, err)
	}
	return f, nil
}

// normalizeHexFloat 把 Python 允许但 Go 语法不接受的十六进制浮点写法补全
// 为 strconv.ParseFloat 可解析的形式。
func normalizeHexFloat(s string) string {
	sign := ""
	body := s
	if body != "" && (body[0] == '+' || body[0] == '-') {
		sign, body = body[:1], body[1:]
	}

	switch strings.ToLower(body) {
	case "inf", "infinity", "nan":
		return s
	}

	hasPrefix := len(body) >= 2 && body[0] == '0' && (body[1] == 'x' || body[1] == 'X')
	hasExponent := strings.ContainsAny(body, "pP")
	switch {
	case hasPrefix && !hasExponent:
		// Python 允许 "0x1.8"；Go 要求必须有 p 指数。
		return sign + body + "p0"
	case !hasPrefix && hasExponent:
		// Python 允许 "1.8p+1"；Go 要求必须有 0x 前缀。
		return sign + "0x" + body
	}
	return s
}

// FloatAsIntegerRatio 等价于 Python 的 float.as_integer_ratio()，把 f 精确
// 表示为最简分数 num/den，其中 den > 0 且 |num| 与 den 互质。例如
// 0.5 -> (1, 2)，0.1 -> (3602879701896397, 36028797018963968)。
//
// NaN 与 ±Inf 无法表示为整数比，返回 error。
//
// Python: f.as_integer_ratio()
func FloatAsIntegerRatio(f float64) (num, den *big.Int, err error) {
	switch {
	case math.IsNaN(f):
		return nil, nil, fmt.Errorf("conv.FloatAsIntegerRatio: cannot convert NaN to integer ratio")
	case math.IsInf(f, 1), math.IsInf(f, -1):
		return nil, nil, fmt.Errorf("conv.FloatAsIntegerRatio: cannot convert Infinity to integer ratio")
	}
	if f == 0 {
		return new(big.Int), big.NewInt(1), nil
	}

	bits := math.Float64bits(f)
	negative := bits>>63 != 0
	exponentBits := int((bits >> 52) & 0x7ff)
	mantissa := bits & ((1 << 52) - 1)

	var exponent int
	if exponentBits == 0 {
		exponent = -1074
	} else {
		mantissa |= 1 << 52
		exponent = exponentBits - 1075
	}

	// 约去 2 的因子，使分子为奇数，从而得到最简分数。
	for mantissa&1 == 0 {
		mantissa >>= 1
		exponent++
	}

	num = new(big.Int).SetUint64(mantissa)
	den = new(big.Int)
	if exponent >= 0 {
		num.Lsh(num, uint(exponent))
		den.SetInt64(1)
	} else {
		den.Lsh(big.NewInt(1), uint(-exponent))
	}
	if negative {
		num.Neg(num)
	}
	return num, den, nil
}
