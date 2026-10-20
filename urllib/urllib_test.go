package urllib

import (
	"slices"
	"testing"
)

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		in   string
		safe []string
		want string
	}{
		{"a b/c", nil, "a%20b/c"},
		{"a b/c", []string{""}, "a%20b%2Fc"},
		{"你好", nil, "%E4%BD%A0%E5%A5%BD"},
		{"a+b", nil, "a%2Bb"},
	}
	for _, tt := range tests {
		if got := Quote(tt.in, tt.safe...); got != tt.want {
			t.Fatalf("Quote(%q,%v) = %q, want %q", tt.in, tt.safe, got, tt.want)
		}
	}
	got, err := Unquote("a%20b%2Fc")
	if err != nil || got != "a b/c" {
		t.Fatalf("Unquote = (%q,%v)", got, err)
	}
	if _, err := Unquote("%zz"); err == nil {
		t.Fatal("invalid escape should error")
	}
}

func TestQuotePlus(t *testing.T) {
	if got := QuotePlus("a b/c"); got != "a+b%2Fc" {
		t.Fatalf("QuotePlus = %q", got)
	}
	got, err := UnquotePlus("a+b%2Fc")
	if err != nil || got != "a b/c" {
		t.Fatalf("UnquotePlus = (%q,%v)", got, err)
	}
}

func TestQuery(t *testing.T) {
	encoded := Encode(map[string][]string{"b": {"2"}, "a": {"1", "x y"}})
	if encoded != "a=1&a=x+y&b=2" {
		t.Fatalf("Encode = %q", encoded)
	}
	values, err := ParseQuery(encoded)
	if err != nil || !slices.Equal(values["a"], []string{"1", "x y"}) {
		t.Fatalf("ParseQuery = (%v,%v)", values, err)
	}
	if _, err := ParseQuery("%zz"); err == nil {
		t.Fatal("invalid query should error")
	}
}

func TestParseJoin(t *testing.T) {
	parsed, err := Parse("https://example.com/a/b?x=1#frag")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scheme != "https" || parsed.Host != "example.com" || parsed.Path != "/a/b" || parsed.Fragment != "frag" {
		t.Fatalf("Parse = %#v", parsed)
	}
	joined, err := Join("https://example.com/a/", "b", "c")
	if err != nil || joined != "https://example.com/a/b/c" {
		t.Fatalf("Join = (%q,%v)", joined, err)
	}
}
