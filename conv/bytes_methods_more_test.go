package conv

import (
	"bytes"
	"testing"
)

func TestBytesStripVariants(t *testing.T) {
	if got := Bytes_strip_chars([]byte("xxhixx"), []byte("x")); string(got) != "hi" {
		t.Fatalf("strip_chars = %q", got)
	}
	if got := Bytes_strip_chars([]byte("abc"), nil); string(got) != "abc" {
		t.Fatalf("strip_chars empty cutset = %q", got)
	}
	if got := Bytes_lstrip([]byte("  a ")); string(got) != "a " {
		t.Fatalf("lstrip = %q", got)
	}
	if got := Bytes_rstrip([]byte(" a  ")); string(got) != " a" {
		t.Fatalf("rstrip = %q", got)
	}
	if got := Bytes_lstrip_chars([]byte("xxab"), []byte("x")); string(got) != "ab" {
		t.Fatalf("lstrip_chars = %q", got)
	}
	if got := Bytes_rstrip_chars([]byte("abxx"), []byte("x")); string(got) != "ab" {
		t.Fatalf("rstrip_chars = %q", got)
	}
}

func TestBytesAffixes(t *testing.T) {
	if got := Bytes_removeprefix([]byte("prefix-x"), []byte("prefix-")); string(got) != "x" {
		t.Fatalf("removeprefix = %q", got)
	}
	if got := Bytes_removesuffix([]byte("x-suffix"), []byte("-suffix")); string(got) != "x" {
		t.Fatalf("removesuffix = %q", got)
	}
	if got := Bytes_removeprefix([]byte("abc"), []byte("x")); string(got) != "abc" {
		t.Fatalf("removeprefix miss = %q", got)
	}
}

func TestBytesPad(t *testing.T) {
	if got := Bytes_center([]byte("ab"), 5); string(got) != "  ab " {
		t.Fatalf("center = %q", got)
	}
	if got := Bytes_center([]byte("ab"), 6, 'x'); string(got) != "xxabxx" {
		t.Fatalf("center fill = %q", got)
	}
	if got := Bytes_ljust([]byte("ab"), 4, '.'); string(got) != "ab.." {
		t.Fatalf("ljust = %q", got)
	}
	if got := Bytes_rjust([]byte("ab"), 4, '.'); string(got) != "..ab" {
		t.Fatalf("rjust = %q", got)
	}
	if got := Bytes_center([]byte("abcdef"), 3); string(got) != "abcdef" {
		t.Fatalf("center no-op = %q", got)
	}
}

func TestBytesExpandtabs(t *testing.T) {
	if got := Bytes_expandtabs([]byte("a\tb"), 4); string(got) != "a   b" {
		t.Fatalf("expandtabs = %q", got)
	}
	if got := Bytes_expandtabs([]byte("a\n\tb"), 4); string(got) != "a\n    b" {
		t.Fatalf("expandtabs newline = %q", got)
	}
	if got := Bytes_expandtabs([]byte("a\tb"), 0); string(got) != "ab" {
		t.Fatalf("expandtabs zero = %q", got)
	}
}

func TestBytesSplitMaxsplit(t *testing.T) {
	parts, err := Bytes_split_n([]byte("a,b,c"), []byte(","), 1)
	if err != nil || len(parts) != 2 || string(parts[0]) != "a" || string(parts[1]) != "b,c" {
		t.Fatalf("split_n = (%q,%v)", parts, err)
	}
	parts, err = Bytes_split_n([]byte("a,b,c"), []byte(","), 0)
	if err != nil || len(parts) != 1 || string(parts[0]) != "a,b,c" {
		t.Fatalf("split_n 0 = (%q,%v)", parts, err)
	}
	if _, err := Bytes_split_n([]byte("abc"), nil, -1); err == nil {
		t.Fatal("split_n empty sep should error")
	}

	parts, err = Bytes_rsplit([]byte("a,b,c"), []byte(","), 1)
	if err != nil || len(parts) != 2 || string(parts[0]) != "a,b" || string(parts[1]) != "c" {
		t.Fatalf("rsplit = (%q,%v)", parts, err)
	}
	parts, err = Bytes_rsplit([]byte("a,b,c"), []byte(","), 0)
	if err != nil || len(parts) != 1 || string(parts[0]) != "a,b,c" {
		t.Fatalf("rsplit 0 = (%q,%v)", parts, err)
	}
	if _, err := Bytes_rsplit([]byte("abc"), nil, -1); err == nil {
		t.Fatal("rsplit empty sep should error")
	}
}

func TestBytesSplitDoesNotAlias(t *testing.T) {
	data := []byte("a,b")
	parts, err := Bytes_split_n(data, []byte(","), -1)
	if err != nil {
		t.Fatal(err)
	}
	parts[0][0] = 'X'
	if bytes.Equal(data, []byte("X,b")) {
		t.Fatal("split parts alias input")
	}
}
