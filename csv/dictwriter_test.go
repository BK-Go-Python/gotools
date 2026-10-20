package csv

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestDictWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := NewDictWriter(&buf, []string{"name", "age"})
	if err := writer.WriteHeader(); err != nil {
		t.Fatal(err)
	}
	if err := writer.Write(map[string]string{"name": "Tom", "age": "18"}); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteAll([]map[string]string{{"name": "Jerry"}}); err != nil {
		t.Fatal(err)
	}
	want := "name,age\r\nTom,18\r\nJerry,\r\n"
	if buf.String() != want {
		t.Fatalf("DictWriter = %q, want %q", buf.String(), want)
	}

	buf.Reset()
	if err := WriteDicts(&buf, []string{"a"}, []map[string]string{{"a": "1"}}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "a\r\n1\r\n" {
		t.Fatalf("WriteDicts = %q", buf.String())
	}

	if err := NewDictWriter(&buf, nil).WriteHeader(); err == nil {
		t.Fatal("empty fieldnames should error")
	}
}

func TestDictWriterFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.csv")
	rows := []map[string]string{
		{"name": "Tom", "age": "18"},
		{"name": "Jerry", "age": "20", "extra": "ignored"},
	}
	if err := WriteDictsFile(path, []string{"name", "age"}, rows); err != nil {
		t.Fatal(err)
	}
	dicts, err := ReadDictsFile(path)
	if err != nil || len(dicts) != 2 || dicts[1]["name"] != "Jerry" {
		t.Fatalf("ReadDictsFile = (%v,%v)", dicts, err)
	}
}
