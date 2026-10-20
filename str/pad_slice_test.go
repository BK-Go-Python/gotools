package str

import "testing"

// Expected values in this file were checked against CPython 3.13.15.

type padCase struct {
	name, in string
	width    int
	fill     []rune
	want     string
}

func runPadCases(t *testing.T, fnName string, fn func(string, int, ...rune) string, cases []padCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := fn(tc.in, tc.width, tc.fill...); got != tc.want {
				t.Fatalf("%s(%q, %d, %v) = %q, want %q", fnName, tc.in, tc.width, tc.fill, got, tc.want)
			}
		})
	}
}

func TestCenter(t *testing.T) {
	runPadCases(t, "Center", Center, []padCase{
		{"even width", "ab", 6, nil, "  ab  "},         // Python: "ab".center(6)
		{"odd width", "ab", 5, nil, "  ab "},           // Python: "ab".center(5)
		{"odd width seven", "ab", 7, nil, "   ab  "},   // Python: "ab".center(7)
		{"odd content", "abc", 7, nil, "  abc  "},      // Python: "abc".center(7)
		{"odd padding", "abc", 6, nil, " abc  "},       // Python: "abc".center(6)
		{"chinese", "你好", 5, nil, "  你好 "},             // Python: "你好".center(5)
		{"same width", "abc", 3, nil, "abc"},           // Python: "abc".center(3)
		{"narrower", "abc", 2, nil, "abc"},             // Python: "abc".center(2)
		{"empty", "", 3, nil, "   "},                   // Python: "".center(3) == "   "
		{"one char even", "a", 4, nil, " a  "},         // Python: "a".center(4)
		{"one char odd", "a", 5, nil, "  a  "},         // Python: "a".center(5)
		{"emoji", "😀", 3, nil, " 😀 "},                  // Python: "😀".center(3)
		{"negative", "ab", -1, nil, "ab"},              // Python: "ab".center(-1) == "ab"
		{"star fill", "ab", 6, []rune{'*'}, "**ab**"},  // Python: "ab".center(6, "*")
		{"dash fill", "ab", 5, []rune{'-'}, "--ab-"},   // Python: "ab".center(5, "-")
		{"dot fill", "ab", 7, []rune{'.'}, "...ab.."},  // Python: "ab".center(7, ".")
		{"empty with fill", "", 3, []rune{'x'}, "xxx"}, // Python: "".center(3, "x")
		{"unicode fill", "你", 4, []rune{'·'}, "·你··"},  // Python: "你".center(4, "·")
		{"chinese fill", "你好", 4, []rune{'*'}, "*你好*"}, // Python: "你好".center(4, "*")
	})
}

func TestLjust(t *testing.T) {
	runPadCases(t, "Ljust", Ljust, []padCase{
		{"basic", "ab", 5, nil, "ab   "},               // Python: "ab".ljust(5)
		{"same width", "ab", 2, nil, "ab"},             // Python: "ab".ljust(2)
		{"chinese", "你好", 4, nil, "你好  "},              // Python: "你好".ljust(4)
		{"empty", "", 3, nil, "   "},                   // Python: "".ljust(3) == "   "
		{"negative", "ab", -1, nil, "ab"},              // Python: "ab".ljust(-1) == "ab"
		{"dash fill", "ab", 5, []rune{'-'}, "ab---"},   // Python: "ab".ljust(5, "-")
		{"star fill", "你好", 4, []rune{'*'}, "你好**"},    // Python: "你好".ljust(4, "*")
		{"emoji fill", "ab", 6, []rune{'😀'}, "ab😀😀😀😀"}, // Python: "ab".ljust(6, "😀")
		{"zero fill", "-42", 5, []rune{'0'}, "-4200"},  // Python: "-42".ljust(5, "0")
	})
}

func TestRjust(t *testing.T) {
	runPadCases(t, "Rjust", Rjust, []padCase{
		{"basic", "ab", 5, nil, "   ab"},               // Python: "ab".rjust(5)
		{"same width", "ab", 2, nil, "ab"},             // Python: "ab".rjust(2)
		{"chinese", "你好", 4, nil, "  你好"},              // Python: "你好".rjust(4)
		{"empty", "", 3, nil, "   "},                   // Python: "".rjust(3) == "   "
		{"negative", "ab", -1, nil, "ab"},              // Python: "ab".rjust(-1) == "ab"
		{"dash fill", "ab", 5, []rune{'-'}, "---ab"},   // Python: "ab".rjust(5, "-")
		{"star fill", "你好", 4, []rune{'*'}, "**你好"},    // Python: "你好".rjust(4, "*")
		{"emoji fill", "ab", 6, []rune{'😀'}, "😀😀😀😀ab"}, // Python: "ab".rjust(6, "😀")
		{"zero fill", "-42", 5, []rune{'0'}, "00-42"},  // Python: "-42".rjust(5, "0")
	})
}

func TestZfill(t *testing.T) {
	tests := []struct {
		name, in string
		width    int
		want     string
	}{
		{"positive", "42", 5, "00042"},             // Python: "42".zfill(5)
		{"negative", "-42", 5, "-0042"},            // Python: "-42".zfill(5)
		{"plus", "+42", 5, "+0042"},                // Python: "+42".zfill(5)
		{"same width", "-42", 3, "-42"},            // Python: "-42".zfill(3)
		{"narrower", "-42", 2, "-42"},              // Python: "-42".zfill(2)
		{"one zero", "-42", 4, "-042"},             // Python: "-42".zfill(4)
		{"empty", "", 3, "000"},                    // Python: "".zfill(3) == "000"
		{"letters", "abc", 5, "00abc"},             // Python: "abc".zfill(5)
		{"lone minus", "-", 3, "-00"},              // Python: "-".zfill(3)
		{"lone plus", "+", 3, "+00"},               // Python: "+".zfill(3)
		{"emoji", "😀", 3, "00😀"},                   // Python: "😀".zfill(3)
		{"chinese after sign", "-你好", 4, "-0你好"},   // Python: "-你好".zfill(4)
		{"already wide", "+42", 1, "+42"},          // Python: "+42".zfill(1)
		{"leading zeros", "000", 5, "00000"},       // Python: "000".zfill(5)
		{"negative zero", "-0", 4, "-000"},         // Python: "-0".zfill(4)
		{"space is not sign", " -42", 6, "00 -42"}, // Python: " -42".zfill(6)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Zfill(tt.in, tt.width); got != tt.want {
				t.Fatalf("Zfill(%q, %d) = %q, want %q", tt.in, tt.width, got, tt.want)
			}
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name, in       string
		wantLen, bytes int
	}{
		{"empty", "", 0, 0},            // Python: len("") == 0
		{"ascii", "abc", 3, 3},         // Python: len("abc") == 3
		{"chinese", "你好", 2, 6},        // Python: len("你好") == 2, len(utf8) == 6
		{"emoji", "😀", 1, 4},           // Python: len("😀") == 1, len(utf8) == 4
		{"combining", "e\u0301", 2, 3}, // Python: len("e\u0301") == 2
		{"mixed", "a😀b", 3, 6},         // Python: len("a😀b") == 3
		{"accented", "é", 1, 2},        // Python: len("é") == 1
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Len(tt.in); got != tt.wantLen {
				t.Errorf("Len(%q) = %d, want %d", tt.in, got, tt.wantLen)
			}
			if got := LenBytes(tt.in); got != tt.bytes {
				t.Errorf("LenBytes(%q) = %d, want %d", tt.in, got, tt.bytes)
			}
		})
	}
}

func TestAt(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		index   int
		want    rune
		wantErr bool
	}{
		{"first", "abc", 0, 'a', false},               // Python: "abc"[0] == "a"
		{"last negative", "abc", -1, 'c', false},      // Python: "abc"[-1] == "c"
		{"first negative", "abc", -3, 'a', false},     // Python: "abc"[-3] == "a"
		{"out of range", "abc", 3, 0, true},           // Python: IndexError
		{"negative out of range", "abc", -4, 0, true}, // Python: IndexError
		{"empty", "", 0, 0, true},                     // Python: IndexError
		{"chinese", "你好", 1, '好', false},              // Python: "你好"[1] == "好"
		{"emoji first", "😀a", 0, '😀', false},          // Python: "😀a"[0] == "😀"
		{"emoji last", "a😀", -1, '😀', false},          // Python: "a😀"[-1] == "😀"
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := At(tt.s, tt.index)
			if (err != nil) != tt.wantErr {
				t.Fatalf("At(%q, %d) error = %v, wantErr %v", tt.s, tt.index, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("At(%q, %d) = %q, want %q", tt.s, tt.index, got, tt.want)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		start, end int
		want       string
	}{
		{"basic", "abcdef", 1, 4, "bcd"},                       // Python: "abcdef"[1:4]
		{"negative", "abcdef", -3, -1, "de"},                   // Python: "abcdef"[-3:-1]
		{"end beyond", "abcdef", 2, 100, "cdef"},               // Python: "abcdef"[2:100]
		{"start beyond negative", "abcdef", -100, 2, "ab"},     // Python: "abcdef"[-100:2]
		{"start after end", "abcdef", 4, 2, ""},                // Python: "abcdef"[4:2] == ""
		{"both negative first", "abcdef", -1, -3, ""},          // Python: "abcdef"[-1:-3] == ""
		{"empty range", "abcdef", 6, 6, ""},                    // Python: "abcdef"[6:6] == ""
		{"empty string", "", 0, 0, ""},                         // Python: ""[0:0] == ""
		{"chinese", "你好世界", 1, 3, "好世"},                        // Python: "你好世界"[1:3]
		{"emoji", "a😀b", 1, 2, "😀"},                            // Python: "a😀b"[1:2] == "😀"
		{"negative start beyond end", "abcdef", -2, 100, "ef"}, // Python: "abcdef"[-2:100]
		{"end negative beyond", "abcdef", 0, -100, ""},         // Python: "abcdef"[0:-100] == ""
		{"both beyond", "abcdef", 100, -100, ""},               // Python: "abcdef"[100:-100] == ""
		{"both clipped", "abcdef", -100, 100, "abcdef"},        // Python: "abcdef"[-100:100]
		{"start clipped negative", "abc", -5, -1, "ab"},        // Python: "abc"[-5:-1] == "ab"
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slice(tt.s, tt.start, tt.end); got != tt.want {
				t.Fatalf("Slice(%q, %d, %d) = %q, want %q", tt.s, tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"empty", "", ""},                   // Python: ""[::-1] == ""
		{"ascii", "abc", "cba"},             // Python: "abc"[::-1]
		{"chinese", "你好", "好你"},             // Python: "你好"[::-1]
		{"emoji", "a😀b", "b😀a"},             // Python: "a😀b"[::-1]
		{"combining", "e\u0301", "\u0301e"}, // Python: "e\u0301"[::-1] reverses code points.
		{"repeated emoji", "😀😀", "😀😀"},      // Python: "😀😀"[::-1]
		{"digits", "a1b", "b1a"},            // Python: "a1b"[::-1]
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.in); got != tt.want {
				t.Fatalf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestOrd(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    rune
		wantErr bool
	}{
		{"ascii", "a", 97, false},              // Python: ord("a") == 97
		{"chinese", "你", 20320, false},         // Python: ord("你") == 20320
		{"emoji", "😀", 128512, false},          // Python: ord("😀") == 128512
		{"nul", "\x00", 0, false},              // Python: ord("\x00") == 0
		{"empty", "", 0, true},                 // Python: TypeError: ord() expected a character
		{"two chars", "ab", 0, true},           // Python: TypeError
		{"combining pair", "e\u0301", 0, true}, // Python: TypeError, length 2
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ord(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Ord(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("Ord(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestChr(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		want    string
		wantErr bool
	}{
		{"ascii", 97, "a", false},                   // Python: chr(97) == "a"
		{"chinese", 20320, "你", false},              // Python: chr(20320) == "你"
		{"emoji", 128512, "😀", false},               // Python: chr(128512) == "😀"
		{"nul", 0, "\x00", false},                   // Python: chr(0) == "\x00"
		{"max rune", 0x10FFFF, "\U0010FFFF", false}, // Python: chr(0x10FFFF) == "\U0010ffff"
		{"negative", -1, "", true},                  // Python: ValueError: chr() arg not in range
		{"too large", 0x110000, "", true},           // Python: ValueError
		{"high surrogate", 0xD800, "", true},        // Python: chr(0xD800) yields a surrogate; DIFF: Go rejects it.
		{"low surrogate", 0xDFFF, "", true},         // Python: chr(0xDFFF) yields a surrogate; DIFF: Go rejects it.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Chr(tt.code)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Chr(%d) error = %v, wantErr %v", tt.code, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("Chr(%d) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name, in string
		n        int
		want     string
	}{
		{"three", "ab", 3, "ababab"}, // Python: "ab" * 3
		{"zero", "ab", 0, ""},        // Python: "ab" * 0 == ""
		{"negative", "ab", -2, ""},   // Python: "ab" * -2 == ""
		{"empty", "", 5, ""},         // Python: "" * 5 == ""
		{"chinese", "你好", 2, "你好你好"}, // Python: "你好" * 2
		{"one", "a", 1, "a"},         // Python: "a" * 1 == "a"
		{"emoji", "😀", 2, "😀😀"},      // Python: "😀" * 2
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Repeat(tt.in, tt.n); got != tt.want {
				t.Fatalf("Repeat(%q, %d) = %q, want %q", tt.in, tt.n, got, tt.want)
			}
		})
	}
}

type affixCase struct {
	name, in, affix, want string
}

func runAffixCases(t *testing.T, fnName string, fn func(string, string) string, cases []affixCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := fn(tc.in, tc.affix); got != tc.want {
				t.Fatalf("%s(%q, %q) = %q, want %q", fnName, tc.in, tc.affix, got, tc.want)
			}
		})
	}
}

func TestRemovePrefix(t *testing.T) {
	runAffixCases(t, "RemovePrefix", RemovePrefix, []affixCase{
		{"match", "abc", "ab", "c"},             // Python: "abc".removeprefix("ab")
		{"whole", "abc", "abc", ""},             // Python: "abc".removeprefix("abc")
		{"empty prefix", "abc", "", "abc"},      // Python: "abc".removeprefix("")
		{"no match", "abc", "x", "abc"},         // Python: "abc".removeprefix("x")
		{"longer prefix", "abc", "abcd", "abc"}, // Python: "abc".removeprefix("abcd")
		{"chinese", "你好世界", "你好", "世界"},         // Python: "你好世界".removeprefix("你好")
		{"emoji", "a😀b", "a😀", "b"},             // Python: "a😀b".removeprefix("a😀")
		{"only once", "aaab", "aa", "ab"},       // Python: "aaab".removeprefix("aa")
		{"mismatch", "abc", "abab", "abc"},      // Python: "abc".removeprefix("abab")
	})
}

func TestRemoveSuffix(t *testing.T) {
	runAffixCases(t, "RemoveSuffix", RemoveSuffix, []affixCase{
		{"match", "abc", "bc", "a"},             // Python: "abc".removesuffix("bc")
		{"whole", "abc", "abc", ""},             // Python: "abc".removesuffix("abc")
		{"empty suffix", "abc", "", "abc"},      // Python: "abc".removesuffix("")
		{"no match", "abc", "x", "abc"},         // Python: "abc".removesuffix("x")
		{"longer suffix", "abc", "zabc", "abc"}, // Python: "abc".removesuffix("zabc")
		{"chinese", "你好世界", "世界", "你好"},         // Python: "你好世界".removesuffix("世界")
		{"emoji", "a😀b", "😀b", "a"},             // Python: "a😀b".removesuffix("😀b")
		{"only once", "abaaa", "aa", "aba"},     // Python: "abaaa".removesuffix("aa")
		{"mismatch", "abc", "bcbc", "abc"},      // Python: "abc".removesuffix("bcbc")
	})
}
