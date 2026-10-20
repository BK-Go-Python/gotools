package str

import "testing"

// Expected values in this file were checked against CPython 3.13.15.
// Rows marked DIFF document behaviour that intentionally differs from Python:
// missing fields are kept verbatim and invalid specs fall back to str(value)
// instead of raising ValueError/IndexError/KeyError.

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []interface{}
		want     string
	}{
		{"auto one", "{}", []interface{}{"hello"}, "hello"},
		{"auto two", "{} {}", []interface{}{"a", "b"}, "a b"},
		{"explicit repeat", "{0} {1} {0}", []interface{}{"a", "b"}, "a b a"},
		{"explicit second", "{1}", []interface{}{"a", "b"}, "b"},
		{"auto ignores extra", "{}", []interface{}{"a", "b"}, "a"},
		{"unicode template", "你好，{0}！", []interface{}{"世界"}, "你好，世界！"},
		{"plain", "plain", nil, "plain"},
		{"empty template", "", nil, ""},
		{"brace escape", "{{}}", nil, "{}"},
		{"brace escape inside", "a{{b}}c", nil, "a{b}c"},
		{"triple braces", "{{{}}}", []interface{}{"x"}, "{x}"},
		{"triple braces explicit", "{{{0}}}", []interface{}{"x"}, "{x}"},
		{"double open escape", "{{{{", nil, "{{"},
		{"double close escape", "}}}}", nil, "}}"},
		{"open escape", "{{", nil, "{"},
		{"close escape", "}}", nil, "}"},
		{"lone open", "{", nil, "{"},                                     // Python: ValueError; DIFF: kept verbatim.
		{"lone close", "}", nil, "}"},                                    // Python: ValueError; DIFF: kept verbatim.
		{"missing auto", "{}", nil, "{}"},                                // Python: IndexError; DIFF: kept verbatim.
		{"missing explicit", "{0}", nil, "{0}"},                          // Python: IndexError; DIFF: kept verbatim.
		{"missing out of range", "{2}", []interface{}{"a", "b"}, "{2}"},  // Python: IndexError; DIFF.
		{"missing name", "{missing}", nil, "{missing}"},                  // Python: KeyError; DIFF: kept verbatim.
		{"missing name with spec", "{missing:>5}", nil, "{missing:>5}"},  // Python: KeyError; DIFF.
		{"conversion unsupported", "{0!r}", []interface{}{"x"}, "{0!r}"}, // Python: "'x'"; DIFF: !r unsupported.
		{"named through args map", "{name}", []interface{}{map[string]interface{}{"name": "Go"}}, "Go"},
		{"named map and age", "{name} is {age}", []interface{}{map[string]interface{}{"name": "Tom", "age": 18}}, "Tom is 18"},
		{"bool", "{}", []interface{}{true}, "True"},           // Python: "{}".format(True)
		{"float", "{}", []interface{}{1.0}, "1.0"},            // Python: "{}".format(1.0)
		{"nil", "{}", []interface{}{nil}, "None"},             // Python: "{}".format(None)
		{"slice", "{}", []interface{}{[]int{1, 2}}, "[1, 2]"}, // Python: "{}".format([1, 2])
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.template, tt.args...); got != tt.want {
				t.Fatalf("Format(%q, %#v) = %q, want %q", tt.template, tt.args, got, tt.want)
			}
		})
	}
}

func TestFormatSpec(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []interface{}
		want     string
	}{
		{"right align string", "{:>5}", []interface{}{"ab"}, "   ab"},            // Python: format("ab", ">5")
		{"left align string", "{:<5}", []interface{}{"ab"}, "ab   "},             // Python: format("ab", "<5")
		{"center odd", "{:^5}", []interface{}{"ab"}, " ab  "},                    // Python: format("ab", "^5")
		{"center seven", "{:^7}", []interface{}{"ab"}, "  ab   "},                // Python: format("ab", "^7")
		{"default string align", "{:5}", []interface{}{"ab"}, "ab   "},           // Python: format("ab", "5")
		{"default number align", "{:5}", []interface{}{42}, "   42"},             // Python: format(42, "5")
		{"explicit left number", "{:<5}", []interface{}{42}, "42   "},            // Python: format(42, "<5")
		{"zero right string", "{:0>5}", []interface{}{"ab"}, "000ab"},            // Python: format("ab", "0>5")
		{"zero left string", "{:0<5}", []interface{}{"ab"}, "ab000"},             // Python: format("ab", "0<5")
		{"star center", "{:*^7}", []interface{}{"ab"}, "**ab***"},                // Python: format("ab", "*^7")
		{"unicode width", "{:>5}", []interface{}{"你好"}, "   你好"},                 // Python: format("你好", ">5")
		{"unicode center", "{:^6}", []interface{}{"你好"}, "  你好  "},               // Python: format("你好", "^6")
		{"unicode right", "{:>6}", []interface{}{"你好"}, "    你好"},                // Python: format("你好", ">6")
		{"unicode fill left", "{:你<4}", []interface{}{"ab"}, "ab你你"},             // Python: format("ab", "你<4"); DIFF: multibyte fill not parsed.
		{"emoji fill center", "{:😀^5}", []interface{}{"ab"}, "😀ab😀😀"},            // Python: format("ab", "😀^5"); DIFF: multibyte fill not parsed.
		{"precision string", "{:.2}", []interface{}{"abcdef"}, "ab"},             // Python: format("abcdef", ".2")
		{"precision zero", "{:.0}", []interface{}{"abc"}, ""},                    // Python: format("abc", ".0")
		{"precision unicode", "{:.2}", []interface{}{"你好世界"}, "你好"},              // Python: format("你好世界", ".2")
		{"precision explicit s", "{:.2s}", []interface{}{"abcdef"}, "ab"},        // Python: format("abcdef", ".2s")
		{"precision float", "{:.2f}", []interface{}{3.14159}, "3.14"},            // Python: format(3.14159, ".2f")
		{"precision half even", "{:.0f}", []interface{}{2.5}, "2"},               // Python: format(2.5, ".0f")
		{"width precision float", "{:8.2f}", []interface{}{3.14159}, "    3.14"}, // Python: format(3.14159, "8.2f")
		{"left precision float", "{:<8.2f}", []interface{}{3.14159}, "3.14    "}, // Python: format(3.14159, "<8.2f")
		{"decimal width", "{:5d}", []interface{}{42}, "   42"},                   // Python: format(42, "5d")
		{"decimal left", "{:<5d}", []interface{}{42}, "42   "},                   // Python: format(42, "<5d")
		{"decimal zero pad", "{:05d}", []interface{}{-42}, "-0042"},              // Python: format(-42, "05d")
		{"sign aware align", "{:=5d}", []interface{}{-42}, "-  42"},              // Python: format(-42, "=5d"); DIFF: '=' ignores sign.
		{"sign aware fill", "{:*=6d}", []interface{}{-42}, "-***42"},             // Python: format(-42, "*=6d"); DIFF: '=' ignores sign.
		{"zero string", "{:05}", []interface{}{"ab"}, "ab000"},                   // Python: format("ab", "05"); DIFF: zero flag right-aligns strings.
		{"zero float", "{:08.2f}", []interface{}{-3.14159}, "-0003.14"},          // Python: format(-3.14159, "08.2f")
		{"hex", "{:x}", []interface{}{255}, "ff"},                                // Python: format(255, "x")
		{"hex upper", "{:X}", []interface{}{255}, "FF"},                          // Python: format(255, "X")
		{"octal", "{:o}", []interface{}{8}, "10"},                                // Python: format(8, "o")
		{"binary", "{:b}", []interface{}{5}, "101"},                              // Python: format(5, "b")
		{"negative hex", "{:x}", []interface{}{-255}, "-ff"},                     // Python: format(-255, "x")
		{"bool decimal", "{:d}", []interface{}{true}, "1"},                       // Python: format(True, "d")
		{"float default", "{:f}", []interface{}{1.0}, "1.000000"},                // Python: format(1.0, "f")
		{"float three", "{:.3f}", []interface{}{1.0}, "1.000"},                   // Python: format(1.0, ".3f")
		{"exponent", "{:e}", []interface{}{1234.5678}, "1.234568e+03"},           // Python: format(1234.5678, "e")
		{"exponent precision", "{:.2e}", []interface{}{0.000123}, "1.23e-04"},    // Python: format(0.000123, ".2e")
		{"general", "{:g}", []interface{}{1234.5678}, "1234.57"},                 // Python: format(1234.5678, "g")
		{"general precision", "{:.3g}", []interface{}{1234.5678}, "1.23e+03"},    // Python: format(1234.5678, ".3g")
		{"exponent upper", "{:E}", []interface{}{1234.5678}, "1.234568E+03"},     // Python: format(1234.5678, "E")
		{"general upper", "{:G}", []interface{}{0.000012345}, "1.2345E-05"},      // Python: format(0.000012345, "G")
		{"percent", "{:%}", []interface{}{0.125}, "12.500000%"},                  // Python: format(0.125, "%")
		{"percent precision", "{:.1%}", []interface{}{0.125}, "12.5%"},           // Python: format(0.125, ".1%")
		{"width float", "{:>10.2f}", []interface{}{3.14159}, "      3.14"},       // Python: format(3.14159, ">10.2f")
		{"center float", "{:^9.2f}", []interface{}{3.14159}, "  3.14   "},        // Python: format(3.14159, "^9.2f")
		{"center empty", "{:^5}", []interface{}{""}, "     "},                    // Python: format("", "^5")
		{"invalid type", "{:q}", []interface{}{1}, "1"},                          // Python: ValueError; DIFF: fallback to str(value).
		{"invalid spec", "{0:abc}", []interface{}{"x"}, "x"},                     // Python: ValueError; DIFF: fallback to str(value).
		{"type mismatch decimal", "{:d}", []interface{}{"x"}, "x"},               // Python: ValueError; DIFF: fallback to str(value).
		{"type mismatch float", "{:.2f}", []interface{}{"x"}, "x"},               // Python: ValueError; DIFF: fallback to str(value).
		{"precision integer", "{:.2}", []interface{}{123}, "123"},                // Python: ValueError; DIFF: precision ignored for ints.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.template, tt.args...); got != tt.want {
				t.Fatalf("Format(%q, %#v) = %q, want %q", tt.template, tt.args, got, tt.want)
			}
		})
	}
}

func TestFormatStruct(t *testing.T) {
	type person struct {
		Name string
		Age  int
		City string `json:"city"`
	}
	p := person{Name: "Tom", Age: 18, City: "Beijing"}

	tests := []struct {
		name     string
		template string
		arg      interface{}
		want     string
	}{
		{"field and tag", "{Name}/{city}/{Age:03d}", p, "Tom/Beijing/018"}, // Python: custom object's __format__.
		{"pointer field", "{Name}", &p, "Tom"},
		{"age field", "{Age}", p, "18"},
		{"missing field", "{Missing}", p, "{Missing}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.template, tt.arg); got != tt.want {
				t.Fatalf("Format(%q, %#v) = %q, want %q", tt.template, tt.arg, got, tt.want)
			}
		})
	}
}

func TestFormatMap(t *testing.T) {
	tests := []struct {
		name     string
		template string
		values   map[string]interface{}
		want     string
	}{
		{"name and age", "{name} is {age}", map[string]interface{}{"name": "Tom", "age": 18}, "Tom is 18"}, // Python: "{name} is {age}".format(**values)
		{"name only", "{name}", map[string]interface{}{"name": "Tom"}, "Tom"},
		{"missing", "{missing}", map[string]interface{}{}, "{missing}"}, // Python: KeyError; DIFF: kept verbatim.
		{"missing nil map", "{missing}", nil, "{missing}"},              // Python: KeyError; DIFF.
		{"braces", "{{{name}}}", map[string]interface{}{"name": "Tom"}, "{Tom}"},
		{"unicode template", "你好，{name}！", map[string]interface{}{"name": "Tom"}, "你好，Tom！"},
		{"integer spec", "{age:03d}", map[string]interface{}{"age": 18}, "018"},
		{"float spec", "{score:.1f}", map[string]interface{}{"score": 95.5}, "95.5"},
		{"align", "{name:>5}", map[string]interface{}{"name": "Tom"}, "  Tom"},
		{"unicode value", "{name}", map[string]interface{}{"name": "世界"}, "世界"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatMap(tt.template, tt.values); got != tt.want {
				t.Fatalf("FormatMap(%q, %#v) = %q, want %q", tt.template, tt.values, got, tt.want)
			}
		})
	}
}
