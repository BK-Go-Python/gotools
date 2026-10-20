package conv

import "testing"

func TestBytesMaketrans(t *testing.T) {
	table, err := Bytes_maketrans([]byte("abc"), []byte("xyz"))
	if err != nil {
		t.Fatal(err)
	}
	if got := Bytes_translate([]byte("abc"), table); string(got) != "xyz" {
		t.Fatalf("translate = %q", got)
	}
	if got := Bytes_translate([]byte("a-b-c"), table, []byte("-")); string(got) != "xyz" {
		t.Fatalf("translate delete = %q", got)
	}
	if _, err := Bytes_maketrans([]byte("ab"), []byte("x")); err == nil {
		t.Fatal("length mismatch should error")
	}
}

func TestBytesTranslateIdentity(t *testing.T) {
	table, err := Bytes_maketrans(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte{0, 1, 127, 128, 255}
	if got := Bytes_translate(data, table); string(got) != string(data) {
		t.Fatalf("identity = %v", got)
	}
	if got := Bytes_translate(nil, table); len(got) != 0 {
		t.Fatalf("nil translate = %v", got)
	}
}
