package py

import (
	"bufio"
	"bytes"
	"errors"
	"testing"
)

func TestDivmod(t *testing.T) {
	q, r, err := Divmod(-7, 3)
	if err != nil || q != -3 || r != 2 {
		t.Fatalf("Divmod = (%d,%d,%v)", q, r, err)
	}
	if _, _, err := Divmod(1, 0); !errors.Is(err, ErrZeroDivision) {
		t.Fatalf("Divmod zero err = %v", err)
	}
}

func TestPrintWith(t *testing.T) {
	var buf bytes.Buffer
	opts := Default_print_options()
	opts.File = &buf
	if err := Print_with(opts, 1, "a", true); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != "1 a True\n" {
		t.Fatalf("Print_with = %q", got)
	}

	var flushed bytes.Buffer
	writer := bufio.NewWriter(&flushed)
	opts = Default_print_options()
	opts.File = writer
	opts.Sep = "-"
	opts.End = "!"
	opts.Flush = true
	if err := Print_with(opts, "x", "y"); err != nil {
		t.Fatal(err)
	}
	if got := flushed.String(); got != "x-y!" {
		t.Fatalf("flushed = %q", got)
	}
}
