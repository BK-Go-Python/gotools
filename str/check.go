package str

import "unicode"

// IsDigit 近似等价于 Python 的 s.isdigit()：非空且所有字符都是十进制数字。
//
// 实现基于 Unicode Nd 类，Python 对少数上标数字（如 "²"）的 isdigit
// 返回 true，这里为 false；绝大多数场景行为一致。
func IsDigit(s string) bool {
	return allRunes(s, unicode.IsDigit)
}

// IsDecimal 等价于 Python 的 s.isdecimal()：非空且所有字符都是 Unicode Nd。
func IsDecimal(s string) bool {
	return allRunes(s, func(r rune) bool { return unicode.Is(unicode.Nd, r) })
}

// IsNumeric 等价于 Python 的 s.isnumeric()：非空且所有字符都是数字。
func IsNumeric(s string) bool {
	return allRunes(s, unicode.IsNumber)
}

// IsAlpha 等价于 Python 的 s.isalpha()：非空且所有字符都是字母。
func IsAlpha(s string) bool {
	return allRunes(s, unicode.IsLetter)
}

// IsAlnum 等价于 Python 的 s.isalnum()：非空且所有字符都是字母或数字。
func IsAlnum(s string) bool {
	return allRunes(s, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsNumber(r)
	})
}

// IsSpace 等价于 Python 的 s.isspace()：非空且所有字符都是空白。
func IsSpace(s string) bool {
	return allRunes(s, isPythonSpace)
}

// IsUpper 等价于 Python 的 s.isupper()：至少有一个带大小写的字符，
// 且这些字符全部为大写。
func IsUpper(s string) bool {
	hasCased := false
	for _, r := range s {
		if !isCased(r) {
			continue
		}
		hasCased = true
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return hasCased
}

// IsLower 等价于 Python 的 s.islower()：至少有一个带大小写的字符，
// 且这些字符全部为小写。
func IsLower(s string) bool {
	hasCased := false
	for _, r := range s {
		if !isCased(r) {
			continue
		}
		hasCased = true
		if !unicode.IsLower(r) {
			return false
		}
	}
	return hasCased
}

// IsTitle 等价于 Python 的 s.istitle()：每个单词首字母大写，其余字母小写。
func IsTitle(s string) bool {
	hasCased := false
	prevCased := false
	for _, r := range s {
		if !isCased(r) {
			prevCased = false
			continue
		}
		hasCased = true
		if prevCased {
			if !unicode.IsLower(r) {
				return false
			}
		} else if !unicode.IsUpper(r) && !unicode.IsTitle(r) {
			return false
		}
		prevCased = true
	}
	return hasCased
}

// IsASCII 等价于 Python 的 s.isascii()：所有字符都在 ASCII 范围内。
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

func allRunes(s string, fn func(rune) bool) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !fn(r) {
			return false
		}
	}
	return true
}

func isCased(r rune) bool {
	return unicode.IsUpper(r) || unicode.IsLower(r) || unicode.IsTitle(r)
}

// isPythonSpace 对应 Python str.isspace() 的空白集合。
// Go 的 unicode.IsSpace 不包含 U+001C..U+001F。
func isPythonSpace(r rune) bool {
	if r >= 0x1c && r <= 0x1f {
		return true
	}
	return unicode.IsSpace(r)
}
