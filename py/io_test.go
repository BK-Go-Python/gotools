package py

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintTo(t *testing.T) {
	var buf bytes.Buffer
	if err := PrintTo(&buf, "-", "!", 1, "a", true); err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), "1-a-True!"; got != want {
		t.Fatalf("PrintTo = %q, want %q", got, want)
	}

	buf.Reset()
	if err := PrintTo(&buf, "", "\n"); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != "\n" {
		t.Fatalf("empty PrintTo = %q", got)
	}

	if err := PrintTo(nil, " ", "\n", 1); err == nil {
		t.Fatal("PrintTo(nil) returned nil error")
	}
}

func TestInputFrom(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"newline", "hello\n", "hello", false},
		{"crlf", "hello\r\n", "hello", false},
		{"last line without newline", "hello", "hello", false},
		{"empty line", "\n", "", false},
		{"eof", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			got, err := inputFrom(strings.NewReader(tt.input), &out, "> ")
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("inputFrom = %q, want %q", got, tt.want)
			}
			if out.String() != "> " {
				t.Fatalf("prompt = %q", out.String())
			}
		})
	}
}

func TestBoolTruthiness(t *testing.T) {
	if !Bool("false") {
		t.Fatal(`Bool("false") should be true, Python bool("false") is True`)
	}
	if Bool(0) {
		t.Fatal("Bool(0) should be false")
	}
}
