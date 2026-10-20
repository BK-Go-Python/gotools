package csv

import (
	"fmt"
	"io"
	"os"
)

// DictWriter 等价于 Python 的 csv.DictWriter：按 fieldnames 顺序把 map 写成一行。
type DictWriter struct {
	w          *Writer
	fieldnames []string
}

// NewDictWriter 创建 DictWriter。
func NewDictWriter(w io.Writer, fieldnames []string) *DictWriter {
	names := append([]string(nil), fieldnames...)
	return &DictWriter{w: NewWriter(w), fieldnames: names}
}

// WriteHeader 写入表头。
func (d *DictWriter) WriteHeader() error {
	if len(d.fieldnames) == 0 {
		return fmt.Errorf("csv.DictWriter: fieldnames must not be empty")
	}
	return d.w.Write(d.fieldnames)
}

// Write 按 fieldnames 顺序写入一行；缺失字段写空字符串，多余字段忽略。
//
// Python 默认对缺失字段抛 ValueError（可用 restval 指定默认值）；这里
// 统一写空字符串，写 DIFF。
func (d *DictWriter) Write(record map[string]string) error {
	row := make([]string, len(d.fieldnames))
	for i, name := range d.fieldnames {
		row[i] = record[name]
	}
	return d.w.Write(row)
}

// WriteAll 写入所有行（不写表头）。
func (d *DictWriter) WriteAll(records []map[string]string) error {
	for _, record := range records {
		if err := d.Write(record); err != nil {
			return err
		}
	}
	d.w.Flush()
	return d.w.Error()
}

// Flush 刷新缓冲。
func (d *DictWriter) Flush() { d.w.Flush() }

// Error 返回写入错误。
func (d *DictWriter) Error() error { return d.w.Error() }

// WriteDicts 写出带表头的 DictWriter 数据。
func WriteDicts(w io.Writer, fieldnames []string, records []map[string]string) error {
	writer := NewDictWriter(w, fieldnames)
	if err := writer.WriteHeader(); err != nil {
		return err
	}
	return writer.WriteAll(records)
}

// WriteDictsFile 把 DictWriter 数据写入文件。
func WriteDictsFile(path string, fieldnames []string, records []map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("csv.WriteDictsFile: %w", err)
	}
	if err := WriteDicts(file, fieldnames, records); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("csv.WriteDictsFile: %w", err)
	}
	return nil
}
