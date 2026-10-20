package bytearray

import "github.com/BK-Go-Python/gotools/conv"

// Maketrans 等价于 Python 的 bytearray.maketrans(from, to)。
func Maketrans(from, to []byte) ([256]byte, error) {
	return conv.Bytes_maketrans(from, to)
}

// Strip 等价于 Python 的 bytearray.strip()。
func (b ByteArray) Strip() ByteArray { return ByteArray(conv.Bytes_strip([]byte(b))) }

// Lstrip 等价于 Python 的 bytearray.lstrip()。
func (b ByteArray) Lstrip() ByteArray { return ByteArray(conv.Bytes_lstrip([]byte(b))) }

// Rstrip 等价于 Python 的 bytearray.rstrip()。
func (b ByteArray) Rstrip() ByteArray { return ByteArray(conv.Bytes_rstrip([]byte(b))) }

// StripChars 等价于 Python 的 bytearray.strip(chars)。
func (b ByteArray) StripChars(cutset []byte) ByteArray {
	return ByteArray(conv.Bytes_strip_chars([]byte(b), cutset))
}

// LstripChars 等价于 Python 的 bytearray.lstrip(chars)。
func (b ByteArray) LstripChars(cutset []byte) ByteArray {
	return ByteArray(conv.Bytes_lstrip_chars([]byte(b), cutset))
}

// RstripChars 等价于 Python 的 bytearray.rstrip(chars)。
func (b ByteArray) RstripChars(cutset []byte) ByteArray {
	return ByteArray(conv.Bytes_rstrip_chars([]byte(b), cutset))
}

// Startswith 等价于 Python 的 bytearray.startswith(prefixes)。
func (b ByteArray) Startswith(prefixes ...[]byte) bool {
	return conv.Bytes_startswith([]byte(b), prefixes...)
}

// Endswith 等价于 Python 的 bytearray.endswith(suffixes)。
func (b ByteArray) Endswith(suffixes ...[]byte) bool {
	return conv.Bytes_endswith([]byte(b), suffixes...)
}

// Find 等价于 Python 的 bytearray.find(sub)。
func (b ByteArray) Find(sub []byte) int { return conv.Bytes_find([]byte(b), sub) }

// Rfind 等价于 Python 的 bytearray.rfind(sub)。
func (b ByteArray) Rfind(sub []byte) int { return conv.Bytes_rfind([]byte(b), sub) }

// Rindex 等价于 Python 的 bytearray.rindex(sub)。
func (b ByteArray) Rindex(sub []byte) (int, error) {
	return conv.Bytes_rindex([]byte(b), sub)
}

// Split 等价于 Python 的 bytearray.split(sep)。
func (b ByteArray) Split(sep []byte) ([]ByteArray, error) {
	parts, err := conv.Bytes_split([]byte(b), sep)
	if err != nil {
		return nil, err
	}
	return byteArrays(parts), nil
}

// SplitN 等价于 Python 的 bytearray.split(sep, maxsplit)。
func (b ByteArray) SplitN(sep []byte, maxsplit int) ([]ByteArray, error) {
	parts, err := conv.Bytes_split_n([]byte(b), sep, maxsplit)
	if err != nil {
		return nil, err
	}
	return byteArrays(parts), nil
}

// RSplit 等价于 Python 的 bytearray.rsplit(sep, maxsplit)。
func (b ByteArray) RSplit(sep []byte, maxsplit int) ([]ByteArray, error) {
	parts, err := conv.Bytes_rsplit([]byte(b), sep, maxsplit)
	if err != nil {
		return nil, err
	}
	return byteArrays(parts), nil
}

// Splitlines 等价于 Python 的 bytearray.splitlines(keepends=False)。
func (b ByteArray) Splitlines(keepends ...bool) []ByteArray {
	keep := len(keepends) > 0 && keepends[0]
	var parts [][]byte
	if keep {
		parts = conv.Bytes_splitlines_keepends([]byte(b))
	} else {
		parts = conv.Bytes_splitlines([]byte(b))
	}
	return byteArrays(parts)
}

// Join 等价于 Python 的 bytearray.join(parts)。
func (b ByteArray) Join(parts [][]byte) ByteArray {
	return ByteArray(conv.Bytes_join([]byte(b), parts))
}

// Center 等价于 Python 的 bytearray.center(width[, fill])。
func (b ByteArray) Center(width int, fill ...byte) ByteArray {
	return ByteArray(conv.Bytes_center([]byte(b), width, fill...))
}

// Ljust 等价于 Python 的 bytearray.ljust(width[, fill])。
func (b ByteArray) Ljust(width int, fill ...byte) ByteArray {
	return ByteArray(conv.Bytes_ljust([]byte(b), width, fill...))
}

// Rjust 等价于 Python 的 bytearray.rjust(width[, fill])。
func (b ByteArray) Rjust(width int, fill ...byte) ByteArray {
	return ByteArray(conv.Bytes_rjust([]byte(b), width, fill...))
}

// Zfill 等价于 Python 的 bytearray.zfill(width)。
func (b ByteArray) Zfill(width int) ByteArray {
	return ByteArray(conv.Bytes_zfill([]byte(b), width))
}

// Expandtabs 等价于 Python 的 bytearray.expandtabs(tabsize=8)。
func (b ByteArray) Expandtabs(tabsize ...int) ByteArray {
	return ByteArray(conv.Bytes_expandtabs([]byte(b), tabsize...))
}

// Capitalize 等价于 Python 的 bytearray.capitalize()。
func (b ByteArray) Capitalize() ByteArray {
	return ByteArray(conv.Bytes_capitalize([]byte(b)))
}

// Title 等价于 Python 的 bytearray.title()。
func (b ByteArray) Title() ByteArray { return ByteArray(conv.Bytes_title([]byte(b))) }

// Swapcase 等价于 Python 的 bytearray.swapcase()。
func (b ByteArray) Swapcase() ByteArray {
	return ByteArray(conv.Bytes_swapcase([]byte(b)))
}

// RemovePrefix 等价于 Python 3.9+ 的 bytearray.removeprefix(prefix)。
func (b ByteArray) RemovePrefix(prefix []byte) ByteArray {
	return ByteArray(conv.Bytes_removeprefix([]byte(b), prefix))
}

// Removesuffix 等价于 Python 3.9+ 的 bytearray.removesuffix(suffix)。
func (b ByteArray) Removesuffix(suffix []byte) ByteArray {
	return ByteArray(conv.Bytes_removesuffix([]byte(b), suffix))
}

// Translate 等价于 Python 的 bytearray.translate(table, delete=b"")。
func (b ByteArray) Translate(table [256]byte, delete ...[]byte) ByteArray {
	return ByteArray(conv.Bytes_translate([]byte(b), table, delete...))
}

// Hex 等价于 Python 的 bytearray.hex(sep, bytes_per_sep)。
func (b ByteArray) Hex(sep string, bytesPerSep ...int) (string, error) {
	return conv.Bytes_hex_sep([]byte(b), sep, bytesPerSep...)
}

// Isalnum 等价于 Python 的 bytearray.isalnum()。
func (b ByteArray) Isalnum() bool { return conv.Bytes_isalnum([]byte(b)) }

// Isalpha 等价于 Python 的 bytearray.isalpha()。
func (b ByteArray) Isalpha() bool { return conv.Bytes_isalpha([]byte(b)) }

// Isascii 等价于 Python 的 bytearray.isascii()。
func (b ByteArray) Isascii() bool { return conv.Bytes_isascii([]byte(b)) }

// Isdigit 等价于 Python 的 bytearray.isdigit()。
func (b ByteArray) Isdigit() bool { return conv.Bytes_isdigit([]byte(b)) }

// Islower 等价于 Python 的 bytearray.islower()。
func (b ByteArray) Islower() bool { return conv.Bytes_islower([]byte(b)) }

// Isspace 等价于 Python 的 bytearray.isspace()。
func (b ByteArray) Isspace() bool { return conv.Bytes_isspace([]byte(b)) }

// Istitle 等价于 Python 的 bytearray.istitle()。
func (b ByteArray) Istitle() bool { return conv.Bytes_istitle([]byte(b)) }

// Isupper 等价于 Python 的 bytearray.isupper()。
func (b ByteArray) Isupper() bool { return conv.Bytes_isupper([]byte(b)) }

func byteArrays(parts [][]byte) []ByteArray {
	out := make([]ByteArray, len(parts))
	for i, part := range parts {
		out[i] = ByteArray(part)
	}
	return out
}
