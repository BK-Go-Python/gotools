package py

import "testing"

func TestFormat(t *testing.T) {
	got, err := Format(3.14159, ".2f")
	if err != nil || got != "3.14" {
		t.Fatalf("Format = (%q,%v), want (3.14,nil)", got, err)
	}
	if _, err := Format(1, "q"); err == nil {
		t.Fatal("Format with invalid spec returned nil error")
	}
}
