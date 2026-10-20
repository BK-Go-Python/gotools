package conv

import (
	"bytes"
	"testing"
)

func TestBytes_contains(t *testing.T) {
	tests := []struct {
		name string
		data string
		sub  string
		want bool
	}{
		{"found", "hello world", "lo w", true},
		{"not found", "hello", "xyz", false},
		{"empty sub", "abc", "", true},
		{"empty both", "", "", true},
		{"empty data", "", "a", false},
		{"whole data", "abc", "abc", true},
		{"sub longer", "abc", "abcd", false},
		{"chinese", "你好，世界", "世界", true},
		{"nul", "\x00a\x00", "\x00", true},
		{"binary", "\xff\xfe", "\xfe", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes_contains([]byte(tt.data), []byte(tt.sub)); got != tt.want {
				t.Errorf("Bytes_contains(%q, %q) = %v, want %v", tt.data, tt.sub, got, tt.want)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		if !Bytes_contains(nil, nil) {
			t.Error("Bytes_contains(nil, nil) = false, want true")
		}
		if Bytes_contains(nil, []byte("a")) {
			t.Error(`Bytes_contains(nil, "a") = true, want false`)
		}
	})
}

func TestBytes_index(t *testing.T) {
	tests := []struct {
		name string
		data string
		sub  string
		want int
	}{
		{"empty sub", "abc", "", 0},
		{"empty both", "", "", 0},
		{"empty data", "", "a", -1},
		{"not found", "abc", "z", -1},
		{"at start", "abcabc", "abc", 0},
		{"middle", "abcabc", "bca", 1},
		{"repeated", "aaaa", "aa", 0},
		{"sub longer", "ab", "abc", -1},
		{"chinese", "你好世界", "世界", 6},
		{"nul", "\x00a", "\x00", 0},
		{"binary", "\xff\x00\xff", "\x00", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes_index([]byte(tt.data), []byte(tt.sub)); got != tt.want {
				t.Errorf("Bytes_index(%q, %q) = %d, want %d", tt.data, tt.sub, got, tt.want)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		if got := Bytes_index(nil, nil); got != 0 {
			t.Errorf("Bytes_index(nil, nil) = %d, want 0", got)
		}
		if got := Bytes_index(nil, []byte("a")); got != -1 {
			t.Errorf(`Bytes_index(nil, "a") = %d, want -1`, got)
		}
	})
}

func TestBytes_count(t *testing.T) {
	tests := []struct {
		name string
		data string
		sub  string
		want int
	}{
		{"empty sub", "abc", "", 4},                     // Python: b"abc".count(b"") == 4
		{"empty sub chinese", "你好", "", 7},              // Python: len(bytes) + 1
		{"empty sub empty data", "", "", 1},             // Python: b"".count(b"") == 1
		{"single byte", "banana", "a", 3},               // Python: b"banana".count(b"a") == 3
		{"non overlapping", "aaaa", "aa", 2},            // Python: b"aaaa".count(b"aa") == 2
		{"overlapping candidates", "abababa", "aba", 2}, // Python: 2
		{"not found", "abc", "z", 0},
		{"whole data", "abc", "abc", 1},
		{"sub longer", "ab", "abc", 0},
		{"chinese", "你好你好", "你好", 2},
		{"nul", "\x00\x00\x00", "\x00", 3},
		{"binary", "\xff\xff", "\xff", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes_count([]byte(tt.data), []byte(tt.sub)); got != tt.want {
				t.Errorf("Bytes_count(%q, %q) = %d, want %d", tt.data, tt.sub, got, tt.want)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		if got := Bytes_count(nil, nil); got != 1 {
			t.Errorf("Bytes_count(nil, nil) = %d, want 1", got)
		}
		if got := Bytes_count(nil, []byte("a")); got != 0 {
			t.Errorf(`Bytes_count(nil, "a") = %d, want 0`, got)
		}
	})
}

func TestBytes_replace(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		old    string
		new    string
		counts []int
		want   string
	}{
		{"default all", "banana", "a", "A", nil, "bAnAnA"},
		{"count one", "banana", "a", "A", []int{1}, "bAnana"},
		{"count two", "banana", "a", "A", []int{2}, "bAnAna"},
		{"count zero", "banana", "a", "A", []int{0}, "banana"},
		{"negative count", "banana", "a", "A", []int{-1}, "bAnAnA"},
		{"large negative count", "banana", "a", "A", []int{-100}, "bAnAnA"},
		{"count exceeds matches", "banana", "a", "A", []int{99}, "bAnAnA"},
		{"not found", "abc", "z", "Y", nil, "abc"},
		{"delete all", "aaa", "a", "", nil, ""},
		{"non overlapping", "aaa", "aa", "b", nil, "ba"},
		{"non overlapping limit", "aaaa", "aa", "b", []int{1}, "baa"},
		{"whole data", "abc", "abc", "X", nil, "X"},
		{"whole data count zero", "abc", "abc", "X", []int{0}, "abc"},
		{"grow", "abc", "b", "bb", nil, "abbc"},
		{"binary old", "\xe4\xbd\xa0\xe5\xa5\xbd", "\xbd\xa0", "X", nil, "\xe4X\xe5\xa5\xbd"},
		{"empty old default", "abc", "", "-", nil, "-a-b-c-"},
		{"empty old count zero", "abc", "", "-", []int{0}, "abc"},
		{"empty old count one", "abc", "", "-", []int{1}, "-abc"},
		{"empty old count two", "abc", "", "-", []int{2}, "-a-bc"},
		{"empty old count three", "abc", "", "-", []int{3}, "-a-b-c"},
		{"empty old count exceeds", "abc", "", "-", []int{9}, "-a-b-c-"},
		{"empty old and new", "abc", "", "", nil, "abc"},
		// Python 按字节边界插入（包括 UTF-8 多字节序列内部）：
		// b"\xe4\xbd\xa0\xe5\xa5\xbd".replace(b"", b"-")
		{"empty old non ascii", "你好", "", "-", nil, "-\xe4-\xbd-\xa0-\xe5-\xa5-\xbd-"},
		{"empty old non ascii count two", "你好", "", "-", []int{2}, "-\xe4-\xbd\xa0\xe5\xa5\xbd"},
		{"empty old nil data", "", "", "x", nil, "x"},
		{"empty old multiple counts", "abc", "", "-", []int{1, 2}, "-abc"},
		{"empty data non empty old", "", "a", "b", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []byte
			if tt.counts == nil {
				got = Bytes_replace([]byte(tt.data), []byte(tt.old), []byte(tt.new))
			} else {
				got = Bytes_replace([]byte(tt.data), []byte(tt.old), []byte(tt.new), tt.counts...)
			}
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("Bytes_replace(%q, %q, %q, %v) = %q, want %q",
					tt.data, tt.old, tt.new, tt.counts, got, tt.want)
			}
		})
	}

	t.Run("does not modify input", func(t *testing.T) {
		src := []byte("banana")
		got := Bytes_replace(src, []byte("a"), []byte("A"))
		if !bytes.Equal(src, []byte("banana")) {
			t.Fatalf("input modified: %q", src)
		}
		got[0] = 'X'
		if !bytes.Equal(src, []byte("banana")) {
			t.Fatalf("result shares input backing array: %q", src)
		}

		zero := Bytes_replace(src, []byte("a"), []byte("A"), 0)
		zero[0] = 'X'
		if !bytes.Equal(src, []byte("banana")) {
			t.Fatalf("count=0 result shares input backing array: %q", src)
		}
	})

	t.Run("nil", func(t *testing.T) {
		if got := Bytes_replace(nil, nil, []byte("x")); !bytes.Equal(got, []byte("x")) {
			t.Errorf("Bytes_replace(nil, nil, \"x\") = %q, want %q", got, "x")
		}
		if got := Bytes_replace(nil, []byte("a"), []byte("b")); len(got) != 0 {
			t.Errorf("Bytes_replace(nil, \"a\", \"b\") = %q, want empty", got)
		}
		if got := Bytes_replace(nil, nil, nil, 0); len(got) != 0 {
			t.Errorf("Bytes_replace(nil, nil, nil, 0) = %q, want empty", got)
		}
	})
}

func TestBytes_split(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		sep     string
		want    []string
		wantErr bool
	}{
		{"normal", "a,b,c", ",", []string{"a", "b", "c"}, false},
		{"consecutive separators", "a,,b", ",", []string{"a", "", "b"}, false},
		{"leading and trailing", ",a,", ",", []string{"", "a", ""}, false},
		{"multi byte sep", "a::b::c", "::", []string{"a", "b", "c"}, false},
		{"repeated sep", "----", "--", []string{"", "", ""}, false},
		{"not found", "abc", ",", []string{"abc"}, false},
		{"empty data", "", ",", []string{""}, false},
		{"data equals sep", ",", ",", []string{"", ""}, false},
		{"sep longer than data", "ab", "abc", []string{"ab"}, false},
		{"overlapping sep", "aaa", "aa", []string{"", "a"}, false},
		{"chinese", "你好，世界", "，", []string{"你好", "世界"}, false},
		{"binary", "a\xffb\xffc", "\xff", []string{"a", "b", "c"}, false},
		{"empty sep", "abc", "", nil, true},
		{"empty data empty sep", "", "", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Bytes_split([]byte(tt.data), []byte(tt.sep))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Bytes_split(%q, %q) = %q, want error", tt.data, tt.sep, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Bytes_split(%q, %q): %v", tt.data, tt.sep, err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Bytes_split(%q, %q) = %q, want %q", tt.data, tt.sep, got, tt.want)
			}
			for i := range tt.want {
				if !bytes.Equal(got[i], []byte(tt.want[i])) {
					t.Errorf("part %d = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}

	t.Run("parts do not share data", func(t *testing.T) {
		data := []byte("a,b,c")
		parts, err := Bytes_split(data, []byte(","))
		if err != nil {
			t.Fatal(err)
		}
		for _, part := range parts {
			for i := range part {
				part[i] = 'X'
			}
		}
		if !bytes.Equal(data, []byte("a,b,c")) {
			t.Fatalf("split parts share data backing array: %q", data)
		}
	})

	t.Run("not found part is a copy", func(t *testing.T) {
		data := []byte("abc")
		parts, err := Bytes_split(data, []byte(","))
		if err != nil {
			t.Fatal(err)
		}
		parts[0][0] = 'X'
		if !bytes.Equal(data, []byte("abc")) {
			t.Fatalf("split result shares data backing array: %q", data)
		}
	})

	t.Run("nil data", func(t *testing.T) {
		parts, err := Bytes_split(nil, []byte(","))
		if err != nil {
			t.Fatalf("Bytes_split(nil, \",\"): %v", err)
		}
		if len(parts) != 1 || len(parts[0]) != 0 {
			t.Fatalf("Bytes_split(nil, \",\") = %q, want one empty part", parts)
		}
	})
}

func TestBytes_join(t *testing.T) {
	tests := []struct {
		name  string
		sep   string
		parts []string
		want  string
	}{
		{"nil parts", ",", nil, ""},
		{"empty parts", ",", []string{}, ""},
		{"single part", ",", []string{"a"}, "a"},
		{"multiple parts", ",", []string{"a", "b", "c"}, "a,b,c"},
		{"empty sep", "", []string{"a", "b", "c"}, "abc"},
		{"empty elements", ",", []string{"", "a", ""}, ",a,"},
		{"all empty elements", ",", []string{"", "", ""}, ",,"},
		{"single empty element", ",", []string{""}, ""},
		{"chinese", "-", []string{"你好", "世界"}, "你好-世界"},
		{"binary sep", "\xff", []string{"a", "b"}, "a\xffb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := make([][]byte, len(tt.parts))
			for i := range tt.parts {
				parts[i] = []byte(tt.parts[i])
			}
			got := Bytes_join([]byte(tt.sep), parts)
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("Bytes_join(%q, %q) = %q, want %q", tt.sep, tt.parts, got, tt.want)
			}
		})
	}

	t.Run("empty result is non nil", func(t *testing.T) {
		if got := Bytes_join([]byte(","), nil); got == nil {
			t.Error("Bytes_join(\",\", nil) = nil, want non-nil empty slice")
		}
		if got := Bytes_join(nil, [][]byte{[]byte("a"), []byte("b")}); !bytes.Equal(got, []byte("ab")) {
			t.Errorf(`Bytes_join(nil, ["a" "b"]) = %q, want "ab"`, got)
		}
	})

	t.Run("does not modify inputs", func(t *testing.T) {
		sep := []byte("--")
		parts := [][]byte{[]byte("a"), []byte("b")}
		got := Bytes_join(sep, parts)
		if !bytes.Equal(got, []byte("a--b")) {
			t.Fatalf("Bytes_join = %q, want %q", got, "a--b")
		}
		got[0] = 'X'
		if !bytes.Equal(sep, []byte("--")) || !bytes.Equal(parts[0], []byte("a")) || !bytes.Equal(parts[1], []byte("b")) {
			t.Fatalf("result shares input backing array: sep=%q parts=%q", sep, parts)
		}
	})
}

func TestBytes_strip(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"no whitespace", "abc", "abc"},
		{"spaces", "  abc  ", "abc"},
		{"tabs", "\t\tabc\t", "abc"},
		{"newlines", "\nabc\n", "abc"},
		{"carriage returns", "\rabc\r", "abc"},
		{"vertical tabs", "\vabc\v", "abc"},
		{"form feeds", "\fabc\f", "abc"},
		{"all ascii whitespace", " \t\n\r\v\f", ""},
		{"mixed ascii whitespace", "\t \nabc\r\v\f", "abc"},
		{"internal whitespace", "  a b\tc  ", "a b\tc"},
		{"chinese", "  你好  ", "你好"},
		{"ideographic space", "\u3000你好\u3000", "\u3000你好\u3000"},
		{"unicode next line", "\u0085abc\u0085", "\u0085abc\u0085"},
		{"non breaking space", "\u00a0abc\u00a0", "\u00a0abc\u00a0"},
		{"nul", "\x00abc\x00", "\x00abc\x00"},
		{"binary", " \xff\x00 ", "\xff\x00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes_strip([]byte(tt.in)); !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("Bytes_strip(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}

	t.Run("does not share input", func(t *testing.T) {
		src := []byte("  abc  ")
		got := Bytes_strip(src)
		got[0] = 'X'
		if !bytes.Equal(src, []byte("  abc  ")) {
			t.Fatalf("result shares input backing array: %q", src)
		}
	})

	t.Run("nil", func(t *testing.T) {
		if got := Bytes_strip(nil); len(got) != 0 {
			t.Errorf("Bytes_strip(nil) = %q, want empty", got)
		}
	})
}

func TestBytes_upper(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"lowercase", "hello", "HELLO"},
		{"mixed", "Hello, World!", "HELLO, WORLD!"},
		{"already uppercase", "ABC", "ABC"},
		{"digits and punctuation", "abc123!@#", "ABC123!@#"},
		{"chinese", "你好abc", "你好ABC"},
		{"accented utf8", "café", "CAFé"}, // b"caf\xc3\xa9".upper() keeps \xc3\xa9
		{"greek utf8", "αβγabc", "αβγABC"},
		{"high bytes", "\xffa\x00z", "\xffA\x00Z"},
		{"invalid utf8", "\xe4a", "\xe4A"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bytes_upper([]byte(tt.in))
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("Bytes_upper(%q) = %q, want %q", tt.in, got, tt.want)
			}
			if len(got) != len(tt.in) {
				t.Fatalf("Bytes_upper(%q) changed length: %d -> %d", tt.in, len(tt.in), len(got))
			}
			for i := 0; i < len(tt.in); i++ {
				if tt.in[i] >= 0x80 && got[i] != tt.in[i] {
					t.Errorf("Bytes_upper(%q) changed non-ASCII byte at %d: %#x -> %#x", tt.in, i, tt.in[i], got[i])
				}
			}
		})
	}

	t.Run("does not share input", func(t *testing.T) {
		src := []byte("aAzZ\xff")
		got := Bytes_upper(src)
		got[0] = 'x'
		if !bytes.Equal(src, []byte("aAzZ\xff")) {
			t.Fatalf("result shares input backing array: %q", src)
		}

		nonASCII := []byte("你好123")
		unchanged := Bytes_upper(nonASCII)
		unchanged[0] = 0
		if !bytes.Equal(nonASCII, []byte("你好123")) {
			t.Fatalf("unchanged result shares input backing array: %q", nonASCII)
		}
	})

	t.Run("nil", func(t *testing.T) {
		if got := Bytes_upper(nil); len(got) != 0 {
			t.Errorf("Bytes_upper(nil) = %q, want empty", got)
		}
	})
}

func TestBytes_lower(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"uppercase", "HELLO", "hello"},
		{"mixed", "Hello, World!", "hello, world!"},
		{"already lowercase", "abc", "abc"},
		{"digits and punctuation", "ABC123!@#", "abc123!@#"},
		{"chinese", "你好ABC", "你好abc"},
		{"accented utf8", "CAFÉ", "cafÉ"}, // b"CAF\xc3\x89".lower() keeps \xc3\x89
		{"greek utf8", "ΑΒΓABC", "ΑΒΓabc"},
		{"high bytes", "\xffA\x00Z", "\xffa\x00z"},
		{"invalid utf8", "\xe4A", "\xe4a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bytes_lower([]byte(tt.in))
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("Bytes_lower(%q) = %q, want %q", tt.in, got, tt.want)
			}
			if len(got) != len(tt.in) {
				t.Fatalf("Bytes_lower(%q) changed length: %d -> %d", tt.in, len(tt.in), len(got))
			}
			for i := 0; i < len(tt.in); i++ {
				if tt.in[i] >= 0x80 && got[i] != tt.in[i] {
					t.Errorf("Bytes_lower(%q) changed non-ASCII byte at %d: %#x -> %#x", tt.in, i, tt.in[i], got[i])
				}
			}
		})
	}

	t.Run("does not share input", func(t *testing.T) {
		src := []byte("aAzZ\xff")
		got := Bytes_lower(src)
		got[0] = 'X'
		if !bytes.Equal(src, []byte("aAzZ\xff")) {
			t.Fatalf("result shares input backing array: %q", src)
		}

		nonASCII := []byte("你好123")
		unchanged := Bytes_lower(nonASCII)
		unchanged[0] = 0
		if !bytes.Equal(nonASCII, []byte("你好123")) {
			t.Fatalf("unchanged result shares input backing array: %q", nonASCII)
		}
	})

	t.Run("nil", func(t *testing.T) {
		if got := Bytes_lower(nil); len(got) != 0 {
			t.Errorf("Bytes_lower(nil) = %q, want empty", got)
		}
	})
}
