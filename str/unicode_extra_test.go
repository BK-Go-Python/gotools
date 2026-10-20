package str

import "testing"

func TestUnicodeSpecialMappings(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"upper sharp s", Upper("ß"), "SS"},
		{"upper ligature", Upper("ﬀﬁﬂﬃﬄﬅﬆ"), "FFFIFLFFIFFLSTST"},
		{"casefold sharp s", Casefold("ß"), "ss"},
		{"casefold ligature", Casefold("ﬀﬁﬂﬃﬄﬅﬆ"), "fffiflffifflstst"},
		{"title ligature", Title("ﬀoo"), "Ffoo"},
		{"capitalize ligature", Capitalize("ﬃx"), "Ffix"},
		{"swapcase sharp s", Swapcase("ßﬀ"), "SSFF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestIsIdentifierUnicode(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{in: "a\u0301", want: true},
		{in: "\u0301a", want: false},
		{in: "a\u203fb", want: true},
		{in: "_a", want: true},
		{in: "a-b", want: false},
	}
	for _, tt := range tests {
		if got := IsIdentifier(tt.in); got != tt.want {
			t.Fatalf("IsIdentifier(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
