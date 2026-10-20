package str

import "testing"

func TestPythonSpellingNames(t *testing.T) {
	if !Isdigit("123") || !Isascii("abc") || !Isalpha("abc") || !Isalnum("a1") ||
		!Isspace(" ") || !Isupper("AB") || !Islower("ab") || !Istitle("Ab") ||
		!Isidentifier("_x") || !Isprintable("abc") {
		t.Fatal("predicate alias failed")
	}
	if got := Splitlines("a\nb"); len(got) != 2 {
		t.Fatalf("Splitlines = %v", got)
	}
	if got := Splitlines_keepends("a\n"); len(got) != 1 || got[0] != "a\n" {
		t.Fatalf("Splitlines_keepends = %v", got)
	}
	if got := Removeprefix("prefix-x", "prefix-"); got != "x" {
		t.Fatalf("Removeprefix = %q", got)
	}
	if got := Removesuffix("x-suffix", "-suffix"); got != "x" {
		t.Fatalf("Removesuffix = %q", got)
	}
	if got := Format_map("{a}", map[string]interface{}{"a": 1}); got != "1" {
		t.Fatalf("Format_map = %q", got)
	}
}
