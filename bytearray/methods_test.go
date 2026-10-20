package bytearray

import (
	"bytes"
	"testing"
)

func TestBytesMethodWrappers(t *testing.T) {
	b := New(' ', 'a', 'b', ' ')
	if got := b.Strip(); string(got) != "ab" {
		t.Fatalf("Strip = %q", got)
	}
	if got := b.Lstrip(); string(got) != "ab " {
		t.Fatalf("Lstrip = %q", got)
	}
	if got := b.Rstrip(); string(got) != " ab" {
		t.Fatalf("Rstrip = %q", got)
	}
	if got := b.StripChars([]byte(" ")); string(got) != "ab" {
		t.Fatalf("StripChars = %q", got)
	}
	if !b.Startswith([]byte("x"), []byte(" ")) || !b.Endswith([]byte(" ")) {
		t.Fatal("Startswith/Endswith failed")
	}
	if b.Find([]byte("ab")) != 1 || b.Rfind([]byte("b")) != 2 {
		t.Fatal("Find/Rfind failed")
	}
	if _, err := b.Rindex([]byte("z")); err == nil {
		t.Fatal("Rindex missing should error")
	}
}

func TestByteArraySplitJoin(t *testing.T) {
	b := New('a', ',', 'b', ',', 'c')
	parts, err := b.Split([]byte(","))
	if err != nil || len(parts) != 3 || string(parts[2]) != "c" {
		t.Fatalf("Split = (%v,%v)", parts, err)
	}
	parts, err = b.SplitN([]byte(","), 1)
	if err != nil || len(parts) != 2 || string(parts[1]) != "b,c" {
		t.Fatalf("SplitN = (%v,%v)", parts, err)
	}
	parts, err = b.RSplit([]byte(","), 1)
	if err != nil || len(parts) != 2 || string(parts[0]) != "a,b" {
		t.Fatalf("RSplit = (%v,%v)", parts, err)
	}
	if got := New(',').Join([][]byte{[]byte("a"), []byte("b")}); string(got) != "a,b" {
		t.Fatalf("Join = %q", got)
	}
	lines := New('a', '\n', 'b', '\n').Splitlines()
	if len(lines) != 2 || string(lines[1]) != "b" {
		t.Fatalf("Splitlines = %v", lines)
	}
	lines = New('a', '\n').Splitlines(true)
	if len(lines) != 1 || string(lines[0]) != "a\n" {
		t.Fatalf("Splitlines keepends = %v", lines)
	}
}

func TestByteArrayPadAndCase(t *testing.T) {
	if got := New('a', 'b').Center(5); string(got) != "  ab " {
		t.Fatalf("Center = %q", got)
	}
	if got := New('a', 'b').Ljust(4, '.'); string(got) != "ab.." {
		t.Fatalf("Ljust = %q", got)
	}
	if got := New('a', 'b').Rjust(4, '.'); string(got) != "..ab" {
		t.Fatalf("Rjust = %q", got)
	}
	if got := New('-', '4', '2').Zfill(5); string(got) != "-0042" {
		t.Fatalf("Zfill = %q", got)
	}
	if got := New('a', '\t', 'b').Expandtabs(4); string(got) != "a   b" {
		t.Fatalf("Expandtabs = %q", got)
	}
	if got := New('h', 'E', 'L', 'L', 'O').Capitalize(); string(got) != "Hello" {
		t.Fatalf("Capitalize = %q", got)
	}
	if got := New('h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd').Title(); string(got) != "Hello World" {
		t.Fatalf("Title = %q", got)
	}
	if got := New('a', 'B', 'c').Swapcase(); string(got) != "AbC" {
		t.Fatalf("Swapcase = %q", got)
	}
}

func TestByteArrayAffixTranslatePredicates(t *testing.T) {
	if got := New('p', 'r', 'e', 'f', 'i', 'x', '-', 'x').RemovePrefix([]byte("prefix-")); string(got) != "x" {
		t.Fatalf("RemovePrefix = %q", got)
	}
	if got := New('x', '-', 's', 'u', 'f', 'f', 'i', 'x').Removesuffix([]byte("-suffix")); string(got) != "x" {
		t.Fatalf("Removesuffix = %q", got)
	}
	table, err := Maketrans([]byte("abc"), []byte("xyz"))
	if err != nil {
		t.Fatal(err)
	}
	if got := New('a', 'b', 'c').Translate(table); string(got) != "xyz" {
		t.Fatalf("Translate = %q", got)
	}
	if got, err := New(0xde, 0xad).Hex(":"); err != nil || got != "de:ad" {
		t.Fatalf("Hex = (%q,%v)", got, err)
	}
	if !New('a', 'b', '1').Isalnum() || !New('a', 'b').Isalpha() ||
		!New('1', '2').Isdigit() || !New('a').Isascii() ||
		!New('a', 'b').Islower() || !New('A', 'B').Isupper() ||
		!New('A', 'b').Istitle() || !New(' ').Isspace() {
		t.Fatal("predicate wrapper failed")
	}
	if !bytes.Equal(New('a').Bytes(), []byte("a")) {
		t.Fatal("Bytes failed")
	}
}
