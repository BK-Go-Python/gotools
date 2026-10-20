package py

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Abs 返回 x 的绝对值。
//
// Python: abs(x)
// DIFF: 仅支持实数；complex 请使用 math/cmplx.Abs。对最小整数取负会溢出。
func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Bin 等价于 Python 的 bin(x)，返回带 "0b" 前缀的二进制字符串。
//
// Python: bin(x)
func Bin[T Integer](x T) string {
	return formatRadix(x, 2, "0b")
}

// Hex 等价于 Python 的 hex(x)，返回带 "0x" 前缀的十六进制字符串。
//
// Python: hex(x)
func Hex[T Integer](x T) string {
	return formatRadix(x, 16, "0x")
}

// Oct 等价于 Python 的 oct(x)，返回带 "0o" 前缀的八进制字符串。
//
// Python: oct(x)
func Oct[T Integer](x T) string {
	return formatRadix(x, 8, "0o")
}

// DivModOf 等价于 Python 的 divmod(a, b)，使用向下取整的除法语义。
//
// Deprecated: 使用 Divmod。
//
// Python: divmod(a, b)
// DIFF: 仅支持整数；b 为 0 时返回 ErrZeroDivision（Python 抛 ZeroDivisionError）。
func DivModOf[T Integer](a, b T) (q, r T, err error) {
	if b == 0 {
		return 0, 0, fmt.Errorf("py.DivModOf: %w", ErrZeroDivision)
	}
	q = a / b
	r = a % b
	if r != 0 && (r < 0) != (b < 0) {
		q--
		r += b
	}
	return q, r, nil
}

// PowMod 等价于 Python 的三参数 pow(base, exp, mod)。
//
// Python: pow(base, exp, mod)
// DIFF: 只接受 int64；mod 必须为正数，结果超出 int64 时返回 ErrOutOfRange。
func PowMod(base, exp, mod int64) (int64, error) {
	if mod <= 0 {
		return 0, fmt.Errorf("py.PowMod: mod must be positive: %d: %w", mod, ErrOutOfRange)
	}
	if exp == math.MinInt64 {
		return 0, fmt.Errorf("py.PowMod: exponent overflows int64: %d: %w", exp, ErrOutOfRange)
	}

	b := big.NewInt(base)
	m := big.NewInt(mod)
	b.Mod(b, m)

	if exp < 0 {
		inverse := new(big.Int).ModInverse(b, m)
		if inverse == nil {
			return 0, fmt.Errorf("py.PowMod: base %d is not invertible modulo %d: %w", base, mod, ErrType)
		}
		b = inverse
		exp = -exp
	}

	result := new(big.Int).Exp(b, big.NewInt(exp), m)
	if !result.IsInt64() {
		return 0, fmt.Errorf("py.PowMod: result overflows int64: %w", ErrOutOfRange)
	}
	return result.Int64(), nil
}

func formatRadix[T Integer](x T, base int, prefix string) string {
	switch v := any(x).(type) {
	case int:
		return signedRadix(int64(v), base, prefix)
	case int8:
		return signedRadix(int64(v), base, prefix)
	case int16:
		return signedRadix(int64(v), base, prefix)
	case int32:
		return signedRadix(int64(v), base, prefix)
	case int64:
		return signedRadix(v, base, prefix)
	case uint:
		return prefix + strconv.FormatUint(uint64(v), base)
	case uint8:
		return prefix + strconv.FormatUint(uint64(v), base)
	case uint16:
		return prefix + strconv.FormatUint(uint64(v), base)
	case uint32:
		return prefix + strconv.FormatUint(uint64(v), base)
	case uint64:
		return prefix + strconv.FormatUint(v, base)
	case uintptr:
		return prefix + strconv.FormatUint(uint64(v), base)
	}
	return prefix
}

func signedRadix(n int64, base int, prefix string) string {
	s := strconv.FormatInt(n, base)
	if strings.HasPrefix(s, "-") {
		return "-" + prefix + s[1:]
	}
	return prefix + s
}

// Divmod 等价于 Python 的 divmod(a, b)，使用向下取整的除法语义，
// 并保留元素类型 T。
//
// Python: divmod(a, b)
// DIFF: 仅支持整数；b 为 0 时返回 ErrZeroDivision。
func Divmod[T Integer](a, b T) (q, r T, err error) {
	return DivModOf(a, b)
}
