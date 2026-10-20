package textwrap

import (
	"slices"
	"testing"
)

func TestWrapFill(t *testing.T) {
	got := Wrap("The quick brown fox", 10)
	if !slices.Equal(got, []string{"The quick", "brown fox"}) {
		t.Fatalf("Wrap = %v", got)
	}
	if got := Wrap("abcdefghij", 5); !slices.Equal(got, []string{"abcde", "fghij"}) {
		t.Fatalf("Wrap long = %v", got)
	}
	if got := Wrap("", 10); len(got) != 0 {
		t.Fatalf("Wrap empty = %v", got)
	}
	if got := Fill("a b c", 3); got != "a b\nc" {
		t.Fatalf("Fill = %q", got)
	}
}

func TestDedentIndent(t *testing.T) {
	input := "\n    a\n      b\n    c\n"
	want := "\na\n  b\nc\n"
	if got := Dedent(input); got != want {
		t.Fatalf("Dedent = %q, want %q", got, want)
	}
	text := "a\n\nb"
	if got := Indent(text, "> "); got != "> a\n\n> b" {
		t.Fatalf("Indent = %q", got)
	}
}

func TestShorten(t *testing.T) {
	if got := Shorten("The quick brown fox", 10); got != "The qui..." {
		t.Fatalf("Shorten = %q", got)
	}
	if got := Shorten("short", 10); got != "short" {
		t.Fatalf("Shorten short = %q", got)
	}
	if got := Shorten("abcdef", 5, ".."); got != "abc.." {
		t.Fatalf("Shorten custom = %q", got)
	}
}
