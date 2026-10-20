// Package csv 提供 Python csv 模块风格的高频读写 API。
//
// 底层使用 encoding/csv；Writer 默认使用 CRLF 行尾，与 Python csv 默认
// lineterminator="\r\n" 一致。
package csv

import (
	stdcsv "encoding/csv"
	"fmt"
	"io"
	"os"
)

// Reader 是 CSV 读取器。
type Reader struct {
	r *stdcsv.Reader
}

// NewReader 创建 CSV 读取器。
func NewReader(rd io.Reader) *Reader {
	return &Reader{r: stdcsv.NewReader(rd)}
}

// SetComma 设置分隔符，默认是逗号。
func (r *Reader) SetComma(comma rune) {
	r.r.Comma = comma
}

// SetFieldsPerRecord 设置每行字段数；0 表示以第一行为准。
func (r *Reader) SetFieldsPerRecord(n int) {
	r.r.FieldsPerRecord = n
}

// SetLazyQuotes 允许非标准引号。
func (r *Reader) SetLazyQuotes(lazy bool) {
	r.r.LazyQuotes = lazy
}

// Read 读取一行；EOF 时返回 io.EOF。
func (r *Reader) Read() ([]string, error) {
	return r.r.Read()
}

// ReadAll 读取全部记录。
func (r *Reader) ReadAll() ([][]string, error) {
	return r.r.ReadAll()
}

// Writer 是 CSV 写入器。
type Writer struct {
	w *stdcsv.Writer
}

// NewWriter 创建 CSV 写入器；默认使用 CRLF 行尾。
func NewWriter(w io.Writer) *Writer {
	writer := stdcsv.NewWriter(w)
	writer.UseCRLF = true
	return &Writer{w: writer}
}

// SetComma 设置分隔符。
func (w *Writer) SetComma(comma rune) {
	w.w.Comma = comma
}

// SetUseCRLF 控制是否使用 CRLF 行尾。
func (w *Writer) SetUseCRLF(use bool) {
	w.w.UseCRLF = use
}

// Write 写入一条记录。
func (w *Writer) Write(record []string) error {
	return w.w.Write(record)
}

// WriteAll 写入全部记录并 flush。
func (w *Writer) WriteAll(records [][]string) error {
	w.w.WriteAll(records)
	return w.w.Error()
}

// Flush 刷新缓冲。
func (w *Writer) Flush() {
	w.w.Flush()
}

// Error 返回写入过程中遇到的错误。
func (w *Writer) Error() error {
	return w.w.Error()
}

// Read 读取 io.Reader 中的全部 CSV 记录。
func Read(rd io.Reader) ([][]string, error) {
	return NewReader(rd).ReadAll()
}

// ReadFile 读取 CSV 文件。
func ReadFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("csv.ReadFile: %w", err)
	}
	defer file.Close()
	return Read(file)
}

// Write 把记录写入 io.Writer。
func Write(w io.Writer, records [][]string) error {
	writer := NewWriter(w)
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("csv.Write: %w", err)
	}
	return nil
}

// WriteFile 把记录写入 CSV 文件。
func WriteFile(path string, records [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("csv.WriteFile: %w", err)
	}
	if err := Write(file, records); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("csv.WriteFile: %w", err)
	}
	return nil
}

// ReadDicts 等价于 Python 的 csv.DictReader：首行作为表头，其余行转为 map。
//
// 行字段少于表头时缺失字段为空字符串；多出的字段会被忽略（Python 默认
// 把它们放到 None 键下）。
func ReadDicts(rd io.Reader) ([]map[string]string, error) {
	records, err := Read(rd)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []map[string]string{}, nil
	}
	header := records[0]
	out := make([]map[string]string, 0, len(records)-1)
	for _, record := range records[1:] {
		row := make(map[string]string, len(header))
		for i, name := range header {
			if i < len(record) {
				row[name] = record[i]
			} else {
				row[name] = ""
			}
		}
		out = append(out, row)
	}
	return out, nil
}

// ReadDictsFile 从文件读取 DictReader 风格数据。
func ReadDictsFile(path string) ([]map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("csv.ReadDictsFile: %w", err)
	}
	defer file.Close()
	return ReadDicts(file)
}
