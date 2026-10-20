package conv

import (
	"bytes"
	"testing"
)

func TestBytesExtra(t *testing.T) {
	if !Bytes_startswith([]byte("hello"), []byte("he")) || Bytes_startswith([]byte("hello"), []byte("x")) {
		t.Fatal("Bytes_startswith mismatch")
	}
	if !Bytes_endswith([]byte("hello"), []byte("lo")) || Bytes_endswith([]byte("hello"), []byte("x")) {
		t.Fatal("Bytes_endswith mismatch")
	}
	if Bytes_find([]byte("banana"), []byte("na")) != 2 || Bytes_rfind([]byte("banana"), []byte("na")) != 4 {
		t.Fatal("find/rfind mismatch")
	}
	if Bytes_find([]byte("abc"), []byte("x")) != -1 {
		t.Fatal("find missing should be -1")
	}

	before, middle, after := Bytes_partition([]byte("a=b=c"), []byte("="))
	if string(before) != "a" || string(middle) != "=" || string(after) != "b=c" {
		t.Fatalf("partition = (%q,%q,%q)", before, middle, after)
	}
	before, middle, after = Bytes_partition([]byte("abc"), []byte("="))
	if string(before) != "abc" || middle != nil || after != nil {
		t.Fatalf("partition missing = (%q,%q,%q)", before, middle, after)
	}
	before, middle, after = Bytes_rpartition([]byte("a=b=c"), []byte("="))
	if string(before) != "a=b" || string(middle) != "=" || string(after) != "c" {
		t.Fatalf("rpartition = (%q,%q,%q)", before, middle, after)
	}

	lines := Bytes_splitlines([]byte("a\r\nb\rc\nd\n"))
	want := [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}
	if len(lines) != len(want) {
		t.Fatalf("splitlines = %q", lines)
	}
	for i := range want {
		if !bytes.Equal(lines[i], want[i]) {
			t.Fatalf("splitlines[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
	if len(Bytes_splitlines(nil)) != 0 {
		t.Fatal("splitlines nil should be empty")
	}

	if got := Bytes_capitalize([]byte("hELLO")); string(got) != "Hello" {
		t.Fatalf("capitalize = %q", got)
	}
	if got := Bytes_title([]byte("hello world")); string(got) != "Hello World" {
		t.Fatalf("title = %q", got)
	}
	if got := Bytes_swapcase([]byte("aBc")); string(got) != "AbC" {
		t.Fatalf("swapcase = %q", got)
	}
	if got := Bytes_zfill([]byte("-42"), 5); string(got) != "-0042" {
		t.Fatalf("zfill = %q", got)
	}
	if got := Bytes_zfill([]byte("42"), 2); string(got) != "42" {
		t.Fatalf("zfill no-op = %q", got)
	}
}
