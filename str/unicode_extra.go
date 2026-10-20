package str

import "strings"

var upperSpecialReplacer = strings.NewReplacer(
	"ß", "SS",
	"ﬀ", "FF",
	"ﬁ", "FI",
	"ﬂ", "FL",
	"ﬃ", "FFI",
	"ﬄ", "FFL",
	"ﬅ", "ST",
	"ﬆ", "ST",
)

var casefoldSpecialReplacer = strings.NewReplacer(
	"ß", "ss",
	"ẞ", "ss",
	"ﬀ", "ff",
	"ﬁ", "fi",
	"ﬂ", "fl",
	"ﬃ", "ffi",
	"ﬄ", "ffl",
	"ﬅ", "st",
	"ﬆ", "st",
)

func upperSpecial(s string) string {
	return upperSpecialReplacer.Replace(s)
}

func casefoldSpecial(s string) string {
	return casefoldSpecialReplacer.Replace(s)
}

// titleLigature 返回连字在 title/capitalize 下的展开形式。
// Python 只展开首个字符，其余部分按 lower 处理，例如 ﬀ -> "Ff"。
func titleLigature(r rune) (string, bool) {
	switch r {
	case 'ﬀ':
		return "Ff", true
	case 'ﬁ':
		return "Fi", true
	case 'ﬂ':
		return "Fl", true
	case 'ﬃ':
		return "Ffi", true
	case 'ﬄ':
		return "Ffl", true
	case 'ﬅ', 'ﬆ':
		return "St", true
	}
	return "", false
}

func swapcaseSpecial(r rune) (string, bool) {
	switch r {
	case 'ß':
		return "SS", true
	case 'İ':
		return "i\u0307", true
	case 'ﬀ':
		return "FF", true
	case 'ﬁ':
		return "FI", true
	case 'ﬂ':
		return "FL", true
	case 'ﬃ':
		return "FFI", true
	case 'ﬄ':
		return "FFL", true
	case 'ﬅ', 'ﬆ':
		return "ST", true
	}
	return "", false
}
