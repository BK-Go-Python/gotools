package conv

// Bytes_isascii 等价于 Python 的 bytes.isascii()：空 bytes 返回 true。
func Bytes_isascii(data []byte) bool {
	for _, b := range data {
		if b >= 0x80 {
			return false
		}
	}
	return true
}

// Bytes_isalpha 等价于 Python 的 bytes.isalpha()：非空且全部是 ASCII 字母。
func Bytes_isalpha(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	for _, b := range data {
		if !isASCIILetter(b) {
			return false
		}
	}
	return true
}

// Bytes_isalnum 等价于 Python 的 bytes.isalnum()：非空且全部是 ASCII 字母或数字。
func Bytes_isalnum(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	for _, b := range data {
		if !isASCIILetter(b) && !(b >= '0' && b <= '9') {
			return false
		}
	}
	return true
}

// Bytes_isdigit 等价于 Python 的 bytes.isdigit()：非空且全部是 ASCII 数字。
func Bytes_isdigit(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// Bytes_isspace 等价于 Python 的 bytes.isspace()：非空且全部是 ASCII 空白。
func Bytes_isspace(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	for _, b := range data {
		if !isASCIIWhitespace(b) {
			return false
		}
	}
	return true
}

// Bytes_islower 等价于 Python 的 bytes.islower()：至少有一个 ASCII 小写
// 字母，且没有 ASCII 大写字母；非 ASCII 字节不参与判断。
func Bytes_islower(data []byte) bool {
	hasLower := false
	for _, b := range data {
		if b >= 'A' && b <= 'Z' {
			return false
		}
		if b >= 'a' && b <= 'z' {
			hasLower = true
		}
	}
	return hasLower
}

// Bytes_isupper 等价于 Python 的 bytes.isupper()：至少有一个 ASCII 大写
// 字母，且没有 ASCII 小写字母；非 ASCII 字节不参与判断。
func Bytes_isupper(data []byte) bool {
	hasUpper := false
	for _, b := range data {
		if b >= 'a' && b <= 'z' {
			return false
		}
		if b >= 'A' && b <= 'Z' {
			hasUpper = true
		}
	}
	return hasUpper
}

// Bytes_istitle 等价于 Python 的 bytes.istitle()：至少有一个大小写字母，
// 每个单词首字母大写、其余字母小写；非 ASCII 字节视为非字母。
func Bytes_istitle(data []byte) bool {
	hasCased := false
	prevCased := false
	for _, b := range data {
		switch {
		case b >= 'A' && b <= 'Z':
			if prevCased {
				return false
			}
			hasCased = true
			prevCased = true
		case b >= 'a' && b <= 'z':
			if !prevCased {
				return false
			}
			hasCased = true
			prevCased = true
		default:
			prevCased = false
		}
	}
	return hasCased
}
