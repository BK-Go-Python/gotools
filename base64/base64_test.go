package base64

import "testing"

func TestBase64(t *testing.T) {
	encoded := B64encode([]byte("hello"))
	if encoded != "aGVsbG8=" {
		t.Fatalf("B64encode = %q", encoded)
	}
	decoded, err := B64decode(encoded)
	if err != nil || string(decoded) != "hello" {
		t.Fatalf("B64decode = (%q,%v)", decoded, err)
	}
	if _, err := B64decode("!!!"); err == nil {
		t.Fatal("invalid base64 should error")
	}
	data := []byte{0xfb, 0xff}
	urlsafe := Urlsafe_b64encode(data)
	if urlsafe != "-_8=" {
		t.Fatalf("Urlsafe_b64encode = %q", urlsafe)
	}
	decoded, err = Urlsafe_b64decode(urlsafe)
	if err != nil || len(decoded) != 2 || decoded[0] != 0xfb || decoded[1] != 0xff {
		t.Fatalf("Urlsafe_b64decode = (%v,%v)", decoded, err)
	}
}
