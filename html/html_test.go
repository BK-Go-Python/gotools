package html

import "testing"

func TestEscapeUnescape(t *testing.T) {
	input := `<a href="x">'&`
	want := `&lt;a href=&quot;x&quot;&gt;&#x27;&amp;`
	if got := Escape(input); got != want {
		t.Fatalf("Escape = %q, want %q", got, want)
	}
	wantNoQuote := `&lt;a href="x"&gt;'&amp;`
	if got := Escape(input, false); got != wantNoQuote {
		t.Fatalf("Escape quote=false = %q", got)
	}
	if got := Unescape(want); got != input {
		t.Fatalf("Unescape = %q", got)
	}
}
