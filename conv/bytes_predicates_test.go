package conv

import "testing"

func TestBytesPredicates(t *testing.T) {
	tests := []struct {
		in    string
		alpha bool
		alnum bool
		digit bool
		space bool
		lower bool
		upper bool
		title bool
		ascii bool
	}{
		{"", false, false, false, false, false, false, false, true},
		{"abc", true, true, false, false, true, false, false, true},
		{"abc1", false, true, false, false, true, false, false, true},
		{"123", false, true, true, false, false, false, false, true},
		{"ABC", true, true, false, false, false, true, false, true},
		{"Abc", true, true, false, false, false, false, true, true},
		{"Hello World", false, false, false, false, false, false, true, true},
		{"Hello world", false, false, false, false, false, false, false, true},
		{"\t\n", false, false, false, true, false, false, false, true},
		{"a\xff", false, false, false, false, true, false, false, false},
		{"\xffa", false, false, false, false, true, false, false, false},
		{"A\xff", false, false, false, false, false, true, true, false},
	}
	for _, tt := range tests {
		in := []byte(tt.in)
		if got := Bytes_isalpha(in); got != tt.alpha {
			t.Fatalf("isalpha(%q) = %v", tt.in, got)
		}
		if got := Bytes_isalnum(in); got != tt.alnum {
			t.Fatalf("isalnum(%q) = %v", tt.in, got)
		}
		if got := Bytes_isdigit(in); got != tt.digit {
			t.Fatalf("isdigit(%q) = %v", tt.in, got)
		}
		if got := Bytes_isspace(in); got != tt.space {
			t.Fatalf("isspace(%q) = %v", tt.in, got)
		}
		if got := Bytes_islower(in); got != tt.lower {
			t.Fatalf("islower(%q) = %v", tt.in, got)
		}
		if got := Bytes_isupper(in); got != tt.upper {
			t.Fatalf("isupper(%q) = %v", tt.in, got)
		}
		if got := Bytes_istitle(in); got != tt.title {
			t.Fatalf("istitle(%q) = %v", tt.in, got)
		}
		if got := Bytes_isascii(in); got != tt.ascii {
			t.Fatalf("isascii(%q) = %v", tt.in, got)
		}
	}
}

func TestBytesMultiAffixes(t *testing.T) {
	data := []byte("hello.go")
	if !Bytes_startswith(data, []byte("x"), []byte("hello")) {
		t.Fatal("startswith multi failed")
	}
	if Bytes_startswith(data) {
		t.Fatal("startswith no prefixes should be false")
	}
	if !Bytes_endswith(data, []byte("x"), []byte(".go")) {
		t.Fatal("endswith multi failed")
	}
	if Bytes_endswith(data) {
		t.Fatal("endswith no suffixes should be false")
	}
}
