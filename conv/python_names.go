package conv

import "math/big"

// 本文件提供 Python 拼写别名。

// Int_to_bytes 是 IntToBytes 的 Python 拼写别名。
func Int_to_bytes(x int64, length int, order string, signed ...bool) ([]byte, error) {
	return IntToBytes(x, length, order, signed...)
}

// Int_from_bytes 是 IntFromBytes 的 Python 拼写别名。
func Int_from_bytes(data []byte, order string, signed ...bool) (int64, error) {
	return IntFromBytes(data, order, signed...)
}

// Float_hex 是 FloatHex 的 Python 拼写别名。
func Float_hex(f float64) string { return FloatHex(f) }

// Float_fromhex 是 FromFloatHex 的 Python 拼写别名。
func Float_fromhex(s string) (float64, error) { return FromFloatHex(s) }

// Float_as_integer_ratio 是 FloatAsIntegerRatio 的 Python 拼写别名。
func Float_as_integer_ratio(f float64) (num, den *big.Int, err error) {
	return FloatAsIntegerRatio(f)
}

// Bytes_hex 是 BytesHex 的 Python 拼写别名。
func Bytes_hex(data []byte) string { return BytesHex(data) }

// Bytes_fromhex 是 FromHex 的 Python 拼写别名。
func Bytes_fromhex(s string) ([]byte, error) { return FromHex(s) }

// Bytes_decode 是 Decode 的 Python 拼写别名。
func Bytes_decode(data []byte) (string, error) { return Decode(data) }

// Str_encode 是 Encode 的 Python 拼写别名。
func Str_encode(s string) []byte { return Encode(s) }
