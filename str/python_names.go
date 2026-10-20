package str

// 本文件提供 Python 拼写别名：Go 导出标识符必须大写开头，其余部分保持
// Python 的写法（例如 isdigit -> Isdigit、splitlines -> Splitlines）。

// Isascii 是 IsASCII 的 Python 拼写别名。
func Isascii(s string) bool { return IsASCII(s) }

// Isdigit 是 IsDigit 的 Python 拼写别名。
func Isdigit(s string) bool { return IsDigit(s) }

// Isdecimal 是 IsDecimal 的 Python 拼写别名。
func Isdecimal(s string) bool { return IsDecimal(s) }

// Isnumeric 是 IsNumeric 的 Python 拼写别名。
func Isnumeric(s string) bool { return IsNumeric(s) }

// Isalpha 是 IsAlpha 的 Python 拼写别名。
func Isalpha(s string) bool { return IsAlpha(s) }

// Isalnum 是 IsAlnum 的 Python 拼写别名。
func Isalnum(s string) bool { return IsAlnum(s) }

// Isspace 是 IsSpace 的 Python 拼写别名。
func Isspace(s string) bool { return IsSpace(s) }

// Isupper 是 IsUpper 的 Python 拼写别名。
func Isupper(s string) bool { return IsUpper(s) }

// Islower 是 IsLower 的 Python 拼写别名。
func Islower(s string) bool { return IsLower(s) }

// Istitle 是 IsTitle 的 Python 拼写别名。
func Istitle(s string) bool { return IsTitle(s) }

// Isidentifier 是 IsIdentifier 的 Python 拼写别名。
func Isidentifier(s string) bool { return IsIdentifier(s) }

// Isprintable 是 IsPrintable 的 Python 拼写别名。
func Isprintable(s string) bool { return IsPrintable(s) }

// Format_map 是 FormatMap 的 Python 拼写别名。
func Format_map(template string, values map[string]interface{}) string {
	return FormatMap(template, values)
}

// Splitlines 是 SplitLines 的 Python 拼写别名。
func Splitlines(s string) []string { return SplitLines(s) }

// Splitlines_keepends 是 SplitLinesKeepends 的 Python 拼写别名。
func Splitlines_keepends(s string) []string { return SplitLinesKeepends(s) }

// Removeprefix 是 RemovePrefix 的 Python 拼写别名。
func Removeprefix(s, prefix string) string { return RemovePrefix(s, prefix) }

// Removesuffix 是 RemoveSuffix 的 Python 拼写别名。
func Removesuffix(s, suffix string) string { return RemoveSuffix(s, suffix) }
