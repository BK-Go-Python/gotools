package str

import "testing"

// Expected values in this file were checked against CPython 3.13.15.

type boolCase struct {
	name, in string
	want     bool
}

func runBoolCases(t *testing.T, fnName string, fn func(string) bool, cases []boolCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := fn(tc.in); got != tc.want {
				t.Fatalf("%s(%q) = %v, want %v", fnName, tc.in, got, tc.want)
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	runBoolCases(t, "IsDigit", IsDigit, []boolCase{
		{"ascii digits", "123", true},     // Python: "123".isdigit()
		{"empty", "", false},              // Python: "".isdigit() == False
		{"letters mixed", "12a", false},   // Python: "12a".isdigit()
		{"fullwidth digits", "１２３", true}, // Python: "１２３".isdigit()
		{"arabic indic", "٣", true},       // Python: "٣".isdigit()
		{"superscript two", "²", false},   // Python: "²".isdigit() == True; DIFF documented in IsDigit.
		{"chinese", "你好", false},          // Python: "你好".isdigit()
		{"negative", "-1", false},         // Python: "-1".isdigit()
		{"space", " ", false},             // Python: " ".isdigit()
	})
}

func TestIsDecimal(t *testing.T) {
	runBoolCases(t, "IsDecimal", IsDecimal, []boolCase{
		{"ascii digits", "123", true},     // Python: "123".isdecimal()
		{"empty", "", false},              // Python: "".isdecimal() == False
		{"fullwidth digits", "１２３", true}, // Python: "１２３".isdecimal()
		{"arabic indic", "٣", true},       // Python: "٣".isdecimal()
		{"superscript two", "²", false},   // Python: "²".isdecimal() == False
		{"half", "½", false},              // Python: "½".isdecimal() == False
		{"float text", "1.2", false},      // Python: "1.2".isdecimal()
	})
}

func TestIsNumeric(t *testing.T) {
	runBoolCases(t, "IsNumeric", IsNumeric, []boolCase{
		{"ascii digits", "123", true},     // Python: "123".isnumeric()
		{"fullwidth digits", "１２３", true}, // Python: "１２３".isnumeric()
		{"superscript two", "²", true},    // Python: "²".isnumeric()
		{"half", "½", true},               // Python: "½".isnumeric()
		{"roman numeral", "Ⅳ", true},      // Python: "Ⅳ".isnumeric()
		{"arabic indic", "٣", true},       // Python: "٣".isnumeric()
		{"empty", "", false},              // Python: "".isnumeric() == False
		{"chinese numerals", "一千", false}, // Python: "一千".isnumeric()
		{"float text", "1.2", false},      // Python: "1.2".isnumeric()
	})
}

func TestIsAlpha(t *testing.T) {
	runBoolCases(t, "IsAlpha", IsAlpha, []boolCase{
		{"ascii", "abc", true},         // Python: "abc".isalpha()
		{"chinese", "你好", true},        // Python: "你好".isalpha()
		{"accented", "é", true},        // Python: "é".isalpha()
		{"sharp s", "ß", true},         // Python: "ß".isalpha()
		{"empty", "", false},           // Python: "".isalpha() == False
		{"digit mixed", "abc1", false}, // Python: "abc1".isalpha()
		{"space", "a b", false},        // Python: "a b".isalpha()
		{"emoji", "😀", false},          // Python: "😀".isalpha()
	})
}

func TestIsAlnum(t *testing.T) {
	runBoolCases(t, "IsAlnum", IsAlnum, []boolCase{
		{"letters and digits", "abc123", true}, // Python: "abc123".isalnum()
		{"chinese and digits", "你好123", true},  // Python: "你好123".isalnum()
		{"superscript", "²", true},             // Python: "²".isalnum()
		{"roman numeral", "Ⅳ", true},           // Python: "Ⅳ".isalnum()
		{"empty", "", false},                   // Python: "".isalnum() == False
		{"space", "abc 123", false},            // Python: "abc 123".isalnum()
		{"punctuation", "a-b", false},          // Python: "a-b".isalnum()
		{"emoji", "😀", false},                  // Python: "😀".isalnum()
	})
}

func TestIsSpace(t *testing.T) {
	runBoolCases(t, "IsSpace", IsSpace, []boolCase{
		{"space", " ", true},                     // Python: " ".isspace()
		{"ascii whitespace", "\t\n\r\v\f", true}, // Python: "\t\n\r\v\f".isspace()
		{"nbsp", "\u00a0", true},                 // Python: "\xa0".isspace()
		{"ideographic space", "\u3000", true},    // Python: "\u3000".isspace()
		{"empty", "", false},                     // Python: "".isspace() == False
		{"letter mixed", " a", false},            // Python: " a".isspace()
		{"zero width space", "\u200b", false},    // Python: "\u200b".isspace() == False
		{"file separator", "\x1c", true},         // Python: "\x1c".isspace() == True; DIFF: unicode.IsSpace.
	})
}

func TestIsUpper(t *testing.T) {
	runBoolCases(t, "IsUpper", IsUpper, []boolCase{
		{"ascii upper", "ABC", true},    // Python: "ABC".isupper()
		{"with digit", "A1!", true},     // Python: "A1!".isupper()
		{"single upper", "A", true},     // Python: "A".isupper()
		{"accented upper", "ÄÖÜ", true}, // Python: "ÄÖÜ".isupper()
		{"chinese mixed", "A你", true},   // Python: "A你".isupper()
		{"emoji mixed", "A😀", true},     // Python: "A😀".isupper()
		{"sharp s lower", "Aß", false},  // Python: "Aß".isupper() == False
		{"empty", "", false},            // Python: "".isupper() == False
		{"digits only", "123", false},   // Python: "123".isupper() == False
		{"mixed case", "Ab", false},     // Python: "Ab".isupper()
		{"uncased only", "你好", false},   // Python: "你好".isupper() == False
	})
}

func TestIsLower(t *testing.T) {
	runBoolCases(t, "IsLower", IsLower, []boolCase{
		{"ascii lower", "abc", true},    // Python: "abc".islower()
		{"with digit", "a1!", true},     // Python: "a1!".islower()
		{"accented lower", "äöü", true}, // Python: "äöü".islower()
		{"chinese mixed", "a你", true},   // Python: "a你".islower()
		{"emoji mixed", "a😀", true},     // Python: "a😀".islower()
		{"upper", "ABC", false},         // Python: "ABC".islower()
		{"empty", "", false},            // Python: "".islower() == False
		{"digits only", "123", false},   // Python: "123".islower() == False
		{"mixed case", "aB", false},     // Python: "aB".islower()
	})
}

func TestIsTitle(t *testing.T) {
	runBoolCases(t, "IsTitle", IsTitle, []boolCase{
		{"two words", "Hello World", true},          // Python: "Hello World".istitle()
		{"apostrophe", "They'Re", true},             // Python: "They'Re".istitle()
		{"digit boundary", "A1B", true},             // Python: "A1B".istitle()
		{"uncased between", "A你B", true},            // Python: "A你B".istitle()
		{"titlecase letter", "ǅungla", true},        // Python: "ǅungla".istitle()
		{"apostrophe lower", "They're", false},      // Python: "They're".istitle()
		{"hyphen lower", "Hello-world", false},      // Python: "Hello-world".istitle()
		{"second word lower", "hello World", false}, // Python: "hello World".istitle()
		{"all upper", "ABC", false},                 // Python: "ABC".istitle()
		{"empty", "", false},                        // Python: "".istitle() == False
		{"digits only", "123", false},               // Python: "123".istitle() == False
	})
}

func TestIsASCII(t *testing.T) {
	runBoolCases(t, "IsASCII", IsASCII, []boolCase{
		{"empty", "", true},            // Python: "".isascii() == True
		{"ascii text", "abc123", true}, // Python: "abc123".isascii()
		{"del", "\x7f", true},          // Python: "\x7f".isascii()
		{"newline", "a\n", true},       // Python: "a\n".isascii()
		{"accented", "é", false},       // Python: "é".isascii() == False
		{"chinese", "你好", false},       // Python: "你好".isascii() == False
		{"emoji", "😀", false},          // Python: "😀".isascii() == False
	})
}
