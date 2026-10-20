package str

import (
	"slices"
	"testing"
)

// Expected values in this file were checked against CPython 3.13.15.

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sep  string
		want []string
	}{
		{"comma", "a,b,c", ",", []string{"a", "b", "c"}},                         // Python: "a,b,c".split(",")
		{"leading and trailing", ",a,", ",", []string{"", "a", ""}},              // Python: ",a,".split(",")
		{"not found", "abc", ",", []string{"abc"}},                               // Python: "abc".split(",")
		{"empty input", "", ",", []string{""}},                                   // Python: "".split(",")
		{"empty sep empty input", "", "", []string{}},                            // Python: "".split() == [] (Python raises on split("")).
		{"whitespace", "  a  b  ", "", []string{"a", "b"}},                       // Python: "  a  b  ".split()
		{"mixed whitespace", "a\tb\nc\u00a0d", "", []string{"a", "b", "c", "d"}}, // Python: .split()
		{"ideographic space", "\u3000a\u3000b\u3000", "", []string{"a", "b"}},    // Python: .split()
		{"chinese", "你好 世界", "", []string{"你好", "世界"}},                           // Python: "你好 世界".split()
		{"empty fields", "a,,b", ",", []string{"a", "", "b"}},                    // Python: "a,,b".split(",")
		{"literal space sep", "a  b", " ", []string{"a", "", "b"}},               // Python: "a  b".split(" ")
		{"control separator", "a\x1cb", "", []string{"a", "b"}},                  // Python: "a\x1cb".split() == ["a", "b"]; DIFF: unicode.IsSpace.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.s, tt.sep); !slices.Equal(got, tt.want) {
				t.Fatalf("Split(%q, %q) = %#v, want %#v", tt.s, tt.sep, got, tt.want)
			}
		})
	}
}

func TestSplitN(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      string
		maxsplit int
		want     []string
	}{
		{"one", "a,b,c", ",", 1, []string{"a", "b,c"}},          // Python: "a,b,c".split(",", 1)
		{"zero", "a,b,c", ",", 0, []string{"a,b,c"}},            // Python: "a,b,c".split(",", 0)
		{"two", "a,b,c", ",", 2, []string{"a", "b", "c"}},       // Python: "a,b,c".split(",", 2)
		{"beyond", "a,b,c", ",", 5, []string{"a", "b", "c"}},    // Python: "a,b,c".split(",", 5)
		{"negative", "a,b,c", ",", -1, []string{"a", "b", "c"}}, // Python: "a,b,c".split(",", -1)
		{"negative two", "a,b,c", ",", -2, []string{"a", "b", "c"}},
		{"separator not found", "abc", ",", 1, []string{"abc"}},           // Python: "abc".split(",", 1)
		{"empty input with sep", "", ",", 0, []string{""}},                // Python: "".split(",", 0)
		{"whitespace one", "  a b c  ", "", 1, []string{"a", "b c  "}},    // Python: .split(None, 1)
		{"whitespace zero", "  a b c  ", "", 0, []string{"a b c  "}},      // Python: .split(None, 0)
		{"whitespace two", "  a b c  ", "", 2, []string{"a", "b", "c  "}}, // Python: .split(None, 2)
		{"whitespace only", "   ", "", 0, []string{}},                     // Python: "   ".split(None, 0) == []
		{"whitespace empty input", "", "", 0, []string{}},                 // Python: "".split(None, 0) == []
		{"whitespace empty input one", "", "", 1, []string{}},             // Python: "".split(None, 1) == []
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitN(tt.s, tt.sep, tt.maxsplit); !slices.Equal(got, tt.want) {
				t.Fatalf("SplitN(%q, %q, %d) = %#v, want %#v", tt.s, tt.sep, tt.maxsplit, got, tt.want)
			}
		})
	}
}

func TestRsplit(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      string
		maxsplit int
		want     []string
	}{
		{"one", "a,b,c", ",", 1, []string{"a,b", "c"}},                    // Python: "a,b,c".rsplit(",", 1)
		{"zero", "a,b,c", ",", 0, []string{"a,b,c"}},                      // Python: "a,b,c".rsplit(",", 0)
		{"two", "a,b,c", ",", 2, []string{"a", "b", "c"}},                 // Python: "a,b,c".rsplit(",", 2)
		{"beyond", "a,b,c", ",", 5, []string{"a", "b", "c"}},              // Python: "a,b,c".rsplit(",", 5)
		{"negative", "a,b,c", ",", -1, []string{"a", "b", "c"}},           // Python: "a,b,c".rsplit(",", -1)
		{"separator not found", "abc", ",", 1, []string{"abc"}},           // Python: "abc".rsplit(",", 1)
		{"empty input with sep", "", ",", 1, []string{""}},                // Python: "".rsplit(",", 1)
		{"trailing empty", ",a,", ",", 1, []string{",a", ""}},             // Python: ",a,".rsplit(",", 1)
		{"overlap", "aaa", "aa", 1, []string{"a", ""}},                    // Python: "aaa".rsplit("aa", 1)
		{"whitespace zero", "  a b c  ", "", 0, []string{"  a b c"}},      // Python: .rsplit(None, 0)
		{"whitespace one", "  a b c  ", "", 1, []string{"  a b", "c"}},    // Python: .rsplit(None, 1)
		{"whitespace two", "  a b c  ", "", 2, []string{"  a", "b", "c"}}, // Python: .rsplit(None, 2)
		{"whitespace only", "   ", "", 1, []string{}},                     // Python: "   ".rsplit(None, 1) == []
		{"empty input whitespace", "", "", 1, []string{}},                 // Python: "".rsplit(None, 1) == []
		{"empty input whitespace zero", "", "", 0, []string{}},            // Python: "".rsplit(None, 0) == []
		{"whitespace only zero", "   ", "", 0, []string{}},                // Python: "   ".rsplit(None, 0) == []
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rsplit(tt.s, tt.sep, tt.maxsplit); !slices.Equal(got, tt.want) {
				t.Fatalf("Rsplit(%q, %q, %d) = %#v, want %#v", tt.s, tt.sep, tt.maxsplit, got, tt.want)
			}
		})
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{"empty", "", []string{}}, // Python: "".splitlines() == []
		{"lf", "a\nb", []string{"a", "b"}},
		{"trailing lf", "a\n", []string{"a"}},        // Python: "a\n".splitlines() == ["a"]
		{"lone lf", "\n", []string{""}},              // Python: "\n".splitlines() == [""]
		{"double lf", "a\n\n", []string{"a", ""}},    // Python: "a\n\n".splitlines() == ["a", ""]
		{"crlf", "\r\n", []string{""}},               // Python: "\r\n".splitlines() == [""]
		{"cr run", "a\r\rb", []string{"a", "", "b"}}, // Python: "a\r\rb".splitlines()
		{"all boundaries", "a\r\nb\rc\vd\fe\x1cf\u0085g\u2028h\u2029i", []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}},
		{"emoji", "😀\n你", []string{"😀", "你"}},
		{"lone file separator", "\x1c", []string{""}}, // Python: "\x1c".splitlines() == [""]
		{"trailing file separator", "a\x1c", []string{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitLines(tt.s); !slices.Equal(got, tt.want) {
				t.Fatalf("SplitLines(%q) = %#v, want %#v", tt.s, got, tt.want)
			}
		})
	}
}

func TestFields(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{"empty", "", []string{}}, // Python: "".split() == []
		{"spaces", "  a  b  ", []string{"a", "b"}},
		{"only spaces", "   ", []string{}},
		{"chinese", "你好 世界", []string{"你好", "世界"}},
		{"nbsp", "\u00a0a\u00a0", []string{"a"}},
		{"ideographic space", "\u3000a\u3000b", []string{"a", "b"}},
		{"control separator", "a\x1cb", []string{"a", "b"}}, // Python: "a\x1cb".split() == ["a", "b"]; DIFF: unicode.IsSpace.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fields(tt.s); !slices.Equal(got, tt.want) {
				t.Fatalf("Fields(%q) = %#v, want %#v", tt.s, got, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name  string
		sep   string
		elems []string
		want  string
	}{
		{"comma", "-", []string{"a", "b", "c"}, "a-b-c"},      // Python: "-".join(["a", "b", "c"])
		{"empty sep", "", []string{"a", "b"}, "ab"},           // Python: "".join(["a", "b"])
		{"empty list", ",", []string{}, ""},                   // Python: ",".join([]) == ""
		{"nil list", ",", nil, ""},                            // Python: ",".join([]) == ""
		{"chinese", "你", []string{"a", "b"}, "a你b"},           // Python: "你".join(["a", "b"])
		{"unicode elems", ",", []string{"你好", "世界"}, "你好,世界"}, // Python: ",".join(...)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.sep, tt.elems); got != tt.want {
				t.Fatalf("Join(%q, %#v) = %q, want %q", tt.sep, tt.elems, got, tt.want)
			}
		})
	}
}

func TestJoinAny(t *testing.T) {
	tests := []struct {
		name  string
		sep   string
		elems []interface{}
		want  string
	}{
		{"strings", "-", []interface{}{"a", "b"}, "a-b"},
		{"mixed", "-", []interface{}{"a", 1, true, nil}, "a-1-True-None"}, // Python: "-".join(["a", 1, ...]) raises TypeError; DIFF: elements are converted.
		{"floats and bool", ",", []interface{}{3.5, false}, "3.5,False"},  // Python: TypeError; DIFF.
		{"empty", ",", nil, ""},                           // Python: ",".join([]) == ""
		{"empty sep ints", "", []interface{}{1, 2}, "12"}, // Python: TypeError; DIFF.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinAny(tt.sep, tt.elems); got != tt.want {
				t.Fatalf("JoinAny(%q, %#v) = %q, want %q", tt.sep, tt.elems, got, tt.want)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		name, s, sep    string
		head, mid, tail string
	}{
		{"basic", "a,b,c", ",", "a", ",", "b,c"},                // Python: "a,b,c".partition(",")
		{"not found", "abc", ",", "abc", "", ""},                // Python: "abc".partition(",")
		{"leading", ",a", ",", "", ",", "a"},                    // Python: ",a".partition(",")
		{"trailing", "a,", ",", "a", ",", ""},                   // Python: "a,".partition(",")
		{"multi separator", "a--b--c", "--", "a", "--", "b--c"}, // Python: "a--b--c".partition("--")
		{"chinese", "你好 世界", " ", "你好", " ", "世界"},              // Python: "你好 世界".partition(" ")
		{"empty sep", "abc", "", "abc", "", ""},                 // Python: "abc".partition("") raises ValueError; DIFF documented.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, mid, tail := Partition(tt.s, tt.sep)
			if head != tt.head || mid != tt.mid || tail != tt.tail {
				t.Fatalf("Partition(%q, %q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.s, tt.sep, head, mid, tail, tt.head, tt.mid, tt.tail)
			}
		})
	}
}

func TestRpartition(t *testing.T) {
	tests := []struct {
		name, s, sep    string
		head, mid, tail string
	}{
		{"basic", "a,b,c", ",", "a,b", ",", "c"},                // Python: "a,b,c".rpartition(",")
		{"not found", "abc", ",", "", "", "abc"},                // Python: "abc".rpartition(",")
		{"leading", ",a", ",", "", ",", "a"},                    // Python: ",a".rpartition(",")
		{"trailing", "a,", ",", "a", ",", ""},                   // Python: "a,".rpartition(",")
		{"multi separator", "a--b--c", "--", "a--b", "--", "c"}, // Python: "a--b--c".rpartition("--")
		{"chinese", "你好 世界", " ", "你好", " ", "世界"},              // Python: "你好 世界".rpartition(" ")
		{"whole string", "abc", "abc", "", "abc", ""},           // Python: "abc".rpartition("abc")
		{"empty input", "", ",", "", "", ""},                    // Python: "".rpartition(",")
		{"empty sep", "abc", "", "", "", "abc"},                 // Python: "abc".rpartition("") raises ValueError; DIFF documented.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, mid, tail := Rpartition(tt.s, tt.sep)
			if head != tt.head || mid != tt.mid || tail != tt.tail {
				t.Fatalf("Rpartition(%q, %q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.s, tt.sep, head, mid, tail, tt.head, tt.mid, tt.tail)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name, s, sub string
		want         bool
	}{
		{"found", "hello", "ell", true},             // Python: "ell" in "hello"
		{"empty sub", "hello", "", true},            // Python: "" in "hello"
		{"both empty", "", "", true},                // Python: "" in ""
		{"not found", "hello", "z", false},          // Python: "z" in "hello"
		{"chinese", "你好世界", "世界", true},             // Python: "世界" in "你好世界"
		{"emoji", "😀a", "😀", true},                  // Python: "😀" in "😀a"
		{"case sensitive", "Hello", "hello", false}, // Python: "hello" in "Hello"
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.s, tt.sub); got != tt.want {
				t.Fatalf("Contains(%q, %q) = %v, want %v", tt.s, tt.sub, got, tt.want)
			}
		})
	}
}

func TestFindRfind(t *testing.T) {
	tests := []struct {
		name, s, sub string
		find, rfind  int
	}{
		{"ascii", "hello", "l", 2, 3},            // Python: "hello".find("l"), .rfind("l")
		{"empty sub", "hello", "", 0, 5},         // Python: "hello".find(""), .rfind("")
		{"chinese", "你好世界", "世界", 2, 2},          // Python: rune indices 2
		{"emoji before ascii", "a😀b", "b", 2, 2}, // Python: "a😀b".find("b") == 2
		{"emoji", "a😀b", "😀", 1, 1},              // Python: "a😀b".find("😀") == 1
		{"repeated emoji", "a😀b😀", "😀", 1, 3},    // Python: rune indices 1, 3
		{"not found", "abc", "z", -1, -1},        // Python: "abc".find("z") == -1
		{"both empty", "", "", 0, 0},             // Python: "".find("") == 0
		{"empty haystack", "", "a", -1, -1},      // Python: "".find("a") == -1
		{"repeated chinese", "你好你好", "你好", 0, 2}, // Python: rune index of last occurrence is 2
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.s, tt.sub); got != tt.find {
				t.Errorf("Find(%q, %q) = %d, want %d", tt.s, tt.sub, got, tt.find)
			}
			if got := Rfind(tt.s, tt.sub); got != tt.rfind {
				t.Errorf("Rfind(%q, %q) = %d, want %d", tt.s, tt.sub, got, tt.rfind)
			}
		})
	}
}

func TestIndexRindex(t *testing.T) {
	tests := []struct {
		name, s, sub  string
		index, rindex int
		wantErr       bool
	}{
		{"ascii", "hello", "l", 2, 3, false},    // Python: "hello".index("l"), .rindex("l")
		{"empty sub", "hello", "", 0, 5, false}, // Python: "hello".index(""), .rindex("")
		{"chinese", "你好世界", "世界", 2, 2, false},  // Python: rune indices 2
		{"emoji", "a😀b😀", "😀", 1, 3, false},     // Python: rune indices 1, 3
		{"both empty", "", "", 0, 0, false},     // Python: "".index("") == 0
		{"not found", "abc", "z", 0, 0, true},   // Python: ValueError: substring not found
		{"empty haystack", "", "a", 0, 0, true}, // Python: ValueError
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, err := Index(tt.s, tt.sub)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Index(%q, %q) error = %v, wantErr %v", tt.s, tt.sub, err, tt.wantErr)
			}
			if err == nil && idx != tt.index {
				t.Errorf("Index(%q, %q) = %d, want %d", tt.s, tt.sub, idx, tt.index)
			}

			ridx, err := Rindex(tt.s, tt.sub)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Rindex(%q, %q) error = %v, wantErr %v", tt.s, tt.sub, err, tt.wantErr)
			}
			if err == nil && ridx != tt.rindex {
				t.Errorf("Rindex(%q, %q) = %d, want %d", tt.s, tt.sub, ridx, tt.rindex)
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name, s, sub string
		want         int
	}{
		{"single", "hello", "l", 2},          // Python: "hello".count("l")
		{"empty sub", "hello", "", 6},        // Python: "hello".count("") == 6
		{"repeated", "abcabc", "abc", 2},     // Python: "abcabc".count("abc")
		{"non overlapping", "aaaa", "aa", 2}, // Python: "aaaa".count("aa")
		{"both empty", "", "", 1},            // Python: "".count("") == 1
		{"empty haystack", "", "a", 0},       // Python: "".count("a") == 0
		{"chinese", "你好你好", "你好", 2},         // Python: "你好你好".count("你好")
		{"emoji", "a😀a😀a", "😀", 2},           // Python: "a😀a😀a".count("😀")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Count(tt.s, tt.sub); got != tt.want {
				t.Fatalf("Count(%q, %q) = %d, want %d", tt.s, tt.sub, got, tt.want)
			}
		})
	}
}

func TestStartswith(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		prefixes []string
		want     bool
	}{
		{"single", "hello", []string{"he"}, true},                   // Python: "hello".startswith("he")
		{"multiple match", "hello", []string{"he", "hi"}, true},     // Python: startswith(("he", "hi"))
		{"multiple no match", "hello", []string{"hi", "ho"}, false}, // Python: startswith(("hi", "ho"))
		{"no prefixes", "hello", nil, false},                        // Python: "hello".startswith(()) == False
		{"empty s empty prefix", "", []string{""}, true},            // Python: "".startswith("") == True
		{"empty prefix", "hello", []string{""}, true},               // Python: "hello".startswith("") == True
		{"chinese", "你好世界", []string{"你好"}, true},                   // Python: "你好世界".startswith("你好")
		{"emoji", "a😀", []string{"a😀"}, true},                       // Python: "a😀".startswith("a😀")
		{"case sensitive", "Hello", []string{"he"}, false},          // Python: "Hello".startswith("he")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Startswith(tt.s, tt.prefixes...); got != tt.want {
				t.Fatalf("Startswith(%q, %#v) = %v, want %v", tt.s, tt.prefixes, got, tt.want)
			}
		})
	}
}

func TestEndswith(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		suffixes []string
		want     bool
	}{
		{"single", "hello", []string{"lo"}, true},                   // Python: "hello".endswith("lo")
		{"multiple match", "hello", []string{"lo", "x"}, true},      // Python: endswith(("lo", "x"))
		{"multiple no match", "hello", []string{"hi", "ho"}, false}, // Python: endswith(("hi", "ho"))
		{"no suffixes", "hello", nil, false},                        // Python: "hello".endswith(()) == False
		{"empty s empty suffix", "", []string{""}, true},            // Python: "".endswith("") == True
		{"empty suffix", "hello", []string{""}, true},               // Python: "hello".endswith("") == True
		{"chinese", "你好世界", []string{"世界"}, true},                   // Python: "你好世界".endswith("世界")
		{"emoji", "a😀", []string{"😀"}, true},                        // Python: "a😀".endswith("😀")
		{"case sensitive", "Hello", []string{"lo"}, true},           // Python: "Hello".endswith("lo")
		{"case sensitive miss", "Hello", []string{"LO"}, false},     // Python: "Hello".endswith("LO")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Endswith(tt.s, tt.suffixes...); got != tt.want {
				t.Fatalf("Endswith(%q, %#v) = %v, want %v", tt.s, tt.suffixes, got, tt.want)
			}
		})
	}
}
