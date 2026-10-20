package csv

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	records, err := Read(strings.NewReader("a,b\n1,2\n"))
	if err != nil || len(records) != 2 || records[1][1] != "2" {
		t.Fatalf("Read = (%v,%v)", records, err)
	}
	records, err = Read(strings.NewReader(`"a,b",c` + "\r\n"))
	if err != nil || records[0][0] != "a,b" {
		t.Fatalf("quoted = (%v,%v)", records, err)
	}
	reader := NewReader(strings.NewReader("a;b\n"))
	reader.SetComma(';')
	record, err := reader.Read()
	if err != nil || record[1] != "b" {
		t.Fatalf("custom comma = (%v,%v)", record, err)
	}
	if _, err := Read(strings.NewReader("a,b\n1\n")); err == nil {
		t.Fatal("wrong field count should error")
	}
}

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, [][]string{{"a", "b"}, {"1", "2"}}); err != nil {
		t.Fatal(err)
	}
	want := "a,b\r\n1,2\r\n"
	if buf.String() != want {
		t.Fatalf("Write = %q, want %q", buf.String(), want)
	}

	buf.Reset()
	writer := NewWriter(&buf)
	writer.SetUseCRLF(false)
	writer.SetComma(';')
	if err := writer.WriteAll([][]string{{"a", "b"}}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "a;b\n" {
		t.Fatalf("custom writer = %q", buf.String())
	}
}

func TestFileAndDicts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.csv")
	records := [][]string{{"name", "age"}, {"Tom", "18"}, {"Jerry", "20"}}
	if err := WriteFile(path, records); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil || len(got) != 3 || got[2][0] != "Jerry" {
		t.Fatalf("ReadFile = (%v,%v)", got, err)
	}
	dicts, err := ReadDictsFile(path)
	if err != nil || len(dicts) != 2 || dicts[0]["name"] != "Tom" || dicts[0]["age"] != "18" {
		t.Fatalf("ReadDictsFile = (%v,%v)", dicts, err)
	}
	if _, err := ReadFile(filepath.Join(dir, "missing.csv")); err == nil {
		t.Fatal("missing file should error")
	}
}
