package str

import (
	"math"
	"testing"
)

// Expected values in this file were checked against CPython 3.13.15.
// Rows marked DIFF document behaviour that intentionally differs from Python.

func TestStr(t *testing.T) {
	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type hidden struct {
		Name   string
		secret int
	}

	tests := []struct {
		name string
		in   any
		want string
	}{
		{"nil", nil, ""}, // Python: str(None) == "None"; DIFF: Str is Go-style, nil -> "".
		{"string", "hello", "hello"},
		{"empty string", "", ""},
		{"bytes", []byte("hello"), "hello"},
		{"nil bytes", []byte(nil), ""},
		{"empty bytes", []byte{}, ""},
		{"int", -42, "-42"},
		{"int8", int8(-128), "-128"},
		{"int16", int16(32767), "32767"},
		{"int32", int32(-2147483648), "-2147483648"},
		{"int64", int64(1234567890123), "1234567890123"},
		{"uint", uint(42), "42"},
		{"uint8", uint8(255), "255"},
		{"uint16", uint16(65535), "65535"},
		{"uint32", uint32(4294967295), "4294967295"},
		{"uint64", uint64(18446744073709551615), "18446744073709551615"},
		{"uintptr", uintptr(7), "7"},
		{"float64", 3.14, "3.14"},
		{"float64 integral", 1.0, "1"}, // Python: str(1.0) == "1.0"; DIFF: Go-style float.
		{"float32 precision", float32(1.0 / 3.0), "0.33333334"},
		{"float32 tenth", float32(0.1), "0.1"},
		{"bool true", true, "true"},                                 // Python: str(True) == "True"; DIFF: Go-style bool.
		{"bool false", false, "false"},                              // Python: str(False) == "False"; DIFF.
		{"struct json", user{"Tom", 18}, `{"name":"Tom","age":18}`}, // Python: str(dict) uses single quotes.
		{"struct json pointer", &user{"Tom", 18}, `{"name":"Tom","age":18}`},
		{"struct unexported ignored", hidden{Name: "Tom", secret: 1}, `{"Name":"Tom"}`},
		{"map json sorted", map[string]int{"b": 2, "a": 1}, `{"a":1,"b":2}`},
		{"json failure inf", math.Inf(1), "+Inf"}, // json.Marshal fails; falls back to fmt.Sprint.
		{"json failure nan", math.NaN(), "NaN"},   // json.Marshal fails; falls back to fmt.Sprint.
		{"json failure complex", complex(1, 2), "(1+2i)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Str(tt.in); got != tt.want {
				t.Fatalf("Str(%#v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

type caseStringCase struct {
	name, in, want string
}

func TestCaseConversions(t *testing.T) {
	groups := []struct {
		name  string
		fn    func(string) string
		cases []caseStringCase
	}{
		{
			name: "Upper",
			fn:   Upper,
			cases: []caseStringCase{
				{"empty", "", ""},
				{"ascii", "hello", "HELLO"},
				{"already upper", "HELLO", "HELLO"},
				{"space", "hello world", "HELLO WORLD"},
				{"chinese", "他们", "他们"}, // Python: "他们".upper() == "他们"
				{"emoji", "😀", "😀"},     // Python: "😀".upper() == "😀"
				{"apostrophe", "they're", "THEY'RE"},
				{"digit", "a1b", "A1B"},
				{"punctuation", "a-b c_d", "A-B C_D"},
				{"sharp s", "ß", "SS"}, // Python: "ß".upper() == "SS"; DIFF: Go simple case mapping.
			},
		},
		{
			name: "Lower",
			fn:   Lower,
			cases: []caseStringCase{
				{"empty", "", ""},
				{"ascii", "HELLO", "hello"},
				{"already lower", "hello", "hello"},
				{"space", "HELLO WORLD", "hello world"},
				{"chinese", "他们", "他们"}, // Python: "他们".lower() == "他们"
				{"emoji", "😀", "😀"},     // Python: "😀".lower() == "😀"
				{"apostrophe", "they're", "they're"},
				{"digit", "a1b", "a1b"},
				{"punctuation", "A-B C_D", "a-b c_d"},
			},
		},
		{
			name: "Capitalize",
			fn:   Capitalize,
			cases: []caseStringCase{
				{"empty", "", ""},
				{"ascii", "hello", "Hello"},
				{"upper input", "HELLO", "Hello"},
				{"space", "hello world", "Hello world"},
				{"chinese", "他们", "他们"}, // Python: "他们".capitalize() == "他们"
				{"emoji", "😀", "😀"},     // Python: "😀".capitalize() == "😀"
				{"apostrophe", "they're", "They're"},
				{"digit", "a1b", "A1b"},
				{"punctuation", "a-b c_d", "A-b c_d"},
			},
		},
		{
			name: "Title",
			fn:   Title,
			cases: []caseStringCase{
				{"empty", "", ""},
				{"ascii", "hello", "Hello"},
				{"upper input", "HELLO", "Hello"},
				{"space", "hello world", "Hello World"},
				{"chinese", "他们", "他们"}, // Python: "他们".title() == "他们"
				{"emoji", "😀", "😀"},     // Python: "😀".title() == "😀"
				{"apostrophe", "they're", "They'Re"},
				{"already title", "They'Re", "They'Re"},
				{"digit", "a1b", "A1B"},
				{"punctuation", "a-b c_d", "A-B C_D"},
				{"titlecase letter", "ǅungla", "ǅungla"}, // Python: "ǅungla".title() == "ǅungla"; DIFF: Go uses ToUpper.
			},
		},
		{
			name: "Swapcase",
			fn:   Swapcase,
			cases: []caseStringCase{
				{"empty", "", ""},
				{"ascii", "hello", "HELLO"},
				{"upper input", "HELLO", "hello"},
				{"space", "hello world", "HELLO WORLD"},
				{"chinese", "他们", "他们"}, // Python: "他们".swapcase() == "他们"
				{"emoji", "😀", "😀"},     // Python: "😀".swapcase() == "😀"
				{"apostrophe", "they're", "THEY'RE"},
				{"digit", "a1b", "A1B"},
				{"punctuation", "a-b c_d", "A-B C_D"},
			},
		},
	}
	for _, group := range groups {
		t.Run(group.name, func(t *testing.T) {
			for _, tc := range group.cases {
				t.Run(tc.name, func(t *testing.T) {
					if got := group.fn(tc.in); got != tc.want {
						t.Fatalf("%s(%q) = %q, want %q", group.name, tc.in, got, tc.want)
					}
				})
			}
		})
	}
}

type trimCase struct {
	name, in, want string
}

func TestStrip(t *testing.T) {
	cases := []trimCase{
		{"spaces", "  hello  ", "hello"},                   // Python: "  hello  ".strip()
		{"mixed whitespace", "\t\n hello \r\v\f", "hello"}, // Python: .strip()
		{"empty", "", ""},                                  // Python: "".strip() == ""
		{"only spaces", "   ", ""},                         // Python: "   ".strip() == ""
		{"inner space kept", "你好 世界", "你好 世界"},             // Python: "你好 世界".strip()
		{"nbsp", "\u00a0x\u00a0", "x"},                     // Python: "\xa0x\xa0".strip() == "x"
		{"ideographic space", "\u3000x\u3000", "x"},        // Python: "\u3000x\u3000".strip()
		{"control separators", "\x1cfoo\x1f", "foo"},       // Python: "\x1cfoo\x1f".strip() == "foo"; DIFF: unicode.IsSpace.
		{"emoji is not space", "😀 x 😀", "😀 x 😀"},           // Python: "😀 x 😀".strip()
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Strip(tc.in); got != tc.want {
				t.Fatalf("Strip(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestLstrip(t *testing.T) {
	cases := []trimCase{
		{"spaces", "  hello  ", "hello  "},                        // Python: "  hello  ".lstrip()
		{"mixed whitespace", "\t\n hello \r\v\f", "hello \r\v\f"}, // Python: .lstrip()
		{"empty", "", ""},                                // Python: "".lstrip() == ""
		{"only spaces", "   ", ""},                       // Python: "   ".lstrip() == ""
		{"chinese", "你好 世界", "你好 世界"},                    // Python: no leading whitespace.
		{"emoji", "😀 x", "😀 x"},                          // Python: "😀 x".lstrip()
		{"control separators", "\x1cfoo\x1f", "foo\x1f"}, // Python: "\x1cfoo\x1f".lstrip() == "foo\x1f"; DIFF.
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Lstrip(tc.in); got != tc.want {
				t.Fatalf("Lstrip(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestRstrip(t *testing.T) {
	cases := []trimCase{
		{"spaces", "  hello  ", "  hello"},                      // Python: "  hello  ".rstrip()
		{"mixed whitespace", "\t\n hello \r\v\f", "\t\n hello"}, // Python: .rstrip()
		{"empty", "", ""},                                // Python: "".rstrip() == ""
		{"only spaces", "   ", ""},                       // Python: "   ".rstrip() == ""
		{"chinese", "你好 世界", "你好 世界"},                    // Python: no trailing whitespace.
		{"emoji", "x 😀", "x 😀"},                          // Python: "x 😀".rstrip()
		{"control separators", "\x1cfoo\x1f", "\x1cfoo"}, // Python: "\x1cfoo\x1f".rstrip() == "\x1cfoo"; DIFF.
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Rstrip(tc.in); got != tc.want {
				t.Fatalf("Rstrip(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

type trimCharsCase struct {
	name, in, cutset, want string
}

func TestStripChars(t *testing.T) {
	cases := []trimCharsCase{
		{"example", "www.example.com", "wcom.", "example"}, // Python: "www.example.com".strip("wcom.")
		{"xyz", "xyzHellozyx", "xyz", "Hello"},             // Python: "xyzHellozyx".strip("xyz")
		{"empty cutset", "abc", "", "abc"},                 // Python: "abc".strip("") == "abc"
		{"empty input", "", "abc", ""},                     // Python: "".strip("abc") == ""
		{"emoji cutset", "😀a😀", "😀", "a"},                  // Python: "😀a😀".strip("😀") == "a"
		{"chinese cutset", "他们是", "他是", "们"},               // Python: "他们是".strip("他是") == "们"
		{"strips all", "abc", "abc", ""},                   // Python: "abc".strip("abc") == ""
		{"minus", "-42-", "-", "42"},                       // Python: "-42-".strip("-") == "42"
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := StripChars(tc.in, tc.cutset); got != tc.want {
				t.Fatalf("StripChars(%q, %q) = %q, want %q", tc.in, tc.cutset, got, tc.want)
			}
		})
	}
}

func TestLstripChars(t *testing.T) {
	cases := []trimCharsCase{
		{"example", "www.example.com", "wcom.", "example.com"}, // Python: .lstrip("wcom.")
		{"xyz", "xyzHellozyx", "xyz", "Hellozyx"},              // Python: .lstrip("xyz")
		{"empty cutset", "abc", "", "abc"},                     // Python: "abc".lstrip("") == "abc"
		{"empty input", "", "abc", ""},                         // Python: "".lstrip("abc") == ""
		{"emoji cutset", "😀a😀", "😀", "a😀"},                     // Python: "😀a😀".lstrip("😀") == "a😀"
		{"chinese cutset", "他们是", "他是", "们是"},                  // Python: "他们是".lstrip("他是") == "们是"
		{"strips all", "abc", "abc", ""},                       // Python: "abc".lstrip("abc") == ""
		{"minus", "-42-", "-", "42-"},                          // Python: "-42-".lstrip("-") == "42-"
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := LstripChars(tc.in, tc.cutset); got != tc.want {
				t.Fatalf("LstripChars(%q, %q) = %q, want %q", tc.in, tc.cutset, got, tc.want)
			}
		})
	}
}

func TestRstripChars(t *testing.T) {
	cases := []trimCharsCase{
		{"example", "www.example.com", "wcom.", "www.example"}, // Python: .rstrip("wcom.")
		{"xyz", "xyzHellozyx", "xyz", "xyzHello"},              // Python: .rstrip("xyz")
		{"empty cutset", "abc", "", "abc"},                     // Python: "abc".rstrip("") == "abc"
		{"empty input", "", "abc", ""},                         // Python: "".rstrip("abc") == ""
		{"emoji cutset", "😀a😀", "😀", "😀a"},                     // Python: "😀a😀".rstrip("😀") == "😀a"
		{"chinese cutset", "他们是", "他是", "他们"},                  // Python: "他们是".rstrip("他是") == "他们"
		{"strips all", "abc", "abc", ""},                       // Python: "abc".rstrip("abc") == ""
		{"minus", "-42-", "-", "-42"},                          // Python: "-42-".rstrip("-") == "-42"
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := RstripChars(tc.in, tc.cutset); got != tc.want {
				t.Fatalf("RstripChars(%q, %q) = %q, want %q", tc.in, tc.cutset, got, tc.want)
			}
		})
	}
}
