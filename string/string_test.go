package string

import (
	"strings"
	"testing"
)

func TestCapwords(t *testing.T) {
	if got := Capwords("hello   world"); got != "Hello World" {
		t.Fatalf("Capwords = %q", got)
	}
	if got := Capwords("HELLO"); got != "Hello" {
		t.Fatalf("Capwords upper = %q", got)
	}
	if got := Capwords(""); got != "" {
		t.Fatalf("Capwords empty = %q", got)
	}
}

func TestTemplateSubstitute(t *testing.T) {
	tmpl := NewTemplate("$name is ${age}, $$5")
	got, err := tmpl.Substitute(map[string]string{"name": "Tom", "age": "18"})
	if err != nil || got != "Tom is 18, $5" {
		t.Fatalf("Substitute = (%q,%v)", got, err)
	}
	if _, err := tmpl.Substitute(map[string]string{"name": "Tom"}); err == nil {
		t.Fatal("missing key should error")
	}
	if _, err := NewTemplate("$").Substitute(nil); err == nil {
		t.Fatal("trailing $ should error")
	}
	if _, err := NewTemplate("${}").Substitute(nil); err == nil {
		t.Fatal("empty placeholder should error")
	}
	if got := tmpl.Safe_substitute(map[string]string{"name": "Tom"}); got != "Tom is ${age}, $5" {
		t.Fatalf("Safe_substitute = %q", got)
	}
	if got := NewTemplate("${bad-name}").Safe_substitute(nil); got != "${bad-name}" {
		t.Fatalf("Safe invalid = %q", got)
	}
}

func TestConstants(t *testing.T) {
	if !strings.Contains(Ascii_letters, "a") || !strings.Contains(Ascii_letters, "Z") ||
		len(Digits) != 10 || len(Octdigits) != 8 {
		t.Fatal("constants wrong")
	}
}

func TestTemplateInvalidDollar(t *testing.T) {
	if got := NewTemplate("$!x").Safe_substitute(nil); got != "$!x" {
		t.Fatalf("Safe invalid dollar = %q", got)
	}
	if _, err := NewTemplate("$!x").Substitute(nil); err == nil {
		t.Fatal("invalid dollar should error")
	}
}
