package conv

import (
	"encoding/hex"
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string // UTF-8 字节的十六进制表示
	}{
		{"empty", "", ""},
		{"ascii", "hello", "68656c6c6f"},
		{"nul", "\x00", "00"},
		{"chinese", "你好，世界", "e4bda0e5a5bdefbc8ce4b896e7958c"},
		{"mixedEmoji", "Hello, 世界! 😀🚀", "48656c6c6f2c20e4b896e7958c2120f09f9880f09f9a80"},
		{"bom", "\ufeffBOM", "efbbbf424f4d"},
		{"maxRune", "\U0010ffff", "f48fbfbf"},
		{"invalidGoString", string([]byte{0xff, 0xfe}), "fffe"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.in)
			if gotHex := hex.EncodeToString(got); gotHex != tt.want {
				t.Errorf("Encode(%q) = %s, want %s", tt.in, gotHex, tt.want)
			}
			if utf8.ValidString(tt.in) && !utf8.Valid(got) {
				t.Errorf("Encode(%q) returned invalid UTF-8", tt.in)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name  string
		inHex string
		want  string
	}{
		{"empty", "", ""},
		{"ascii", "68656c6c6f", "hello"},
		{"nul", "00", "\x00"},
		{"chinese", "e4bda0e5a5bdefbc8ce4b896e7958c", "你好，世界"},
		{"emoji", "f09f9880f09f9a80", "😀🚀"},
		{"bom", "efbbbf424f4d", "\ufeffBOM"},
		{"maxRune", "f48fbfbf", "\U0010ffff"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.inHex)
			if err != nil {
				t.Fatalf("bad test input %q: %v", tt.inHex, err)
			}
			got, err := Decode(data)
			if err != nil {
				t.Fatalf("Decode(%x): %v", data, err)
			}
			if got != tt.want {
				t.Errorf("Decode(%x) = %q, want %q", data, got, tt.want)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		got, err := Decode(nil)
		if err != nil || got != "" {
			t.Fatalf("Decode(nil) = (%q, %v), want (\"\", nil)", got, err)
		}
	})
}

func TestDecodeInvalid(t *testing.T) {
	tests := []struct {
		name  string
		inHex string
	}{
		{"strayContinuation", "80"},
		{"invalidByte", "ff"},
		{"badContinuation", "c328"},
		{"truncatedTwoByte", "c3"},
		{"truncatedThreeByte", "e4b8"},
		{"truncatedFourByte", "f09f98"},
		{"overlongTwoByte", "c0af"},
		{"overlongThreeByte", "e080af"},
		{"overlongFourByte", "f0808080"},
		{"surrogate", "eda080"},
		{"outOfRange", "f4908080"},
		{"validPrefixThenInvalid", "e4bda0ff"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.inHex)
			if err != nil {
				t.Fatalf("bad test input %q: %v", tt.inHex, err)
			}
			got, err := Decode(data)
			if err == nil {
				t.Fatalf("Decode(%x) = %q, want error", data, got)
			}
			if got != "" {
				t.Fatalf("Decode(%x) returned %q with error", data, got)
			}
		})
	}
}

func TestEncodingRoundTrip(t *testing.T) {
	values := []string{
		"",
		"hello",
		"\x00",
		"你好，世界",
		"Hello, 世界! 😀🚀",
		"\ufeffBOM",
		"\U0010ffff",
		"line1\nline2\tend",
	}
	for i, s := range values {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, err := Decode(Encode(s))
			if err != nil {
				t.Fatalf("Decode(Encode(%q)): %v", s, err)
			}
			if got != s {
				t.Fatalf("Decode(Encode(%q)) = %q", s, got)
			}
		})
	}
}
