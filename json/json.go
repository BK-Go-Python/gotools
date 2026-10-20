// Package json 提供 Python json 模块风格的高频 API：loads/dumps/load/dump/indent。
//
// 底层使用 encoding/json；与 Python 的主要差异写在函数注释的 DIFF 中。
package json

import (
	stdjson "encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/BK-Go-Python/gotools/conv"
	"github.com/BK-Go-Python/gotools/file"
)

// Loads 等价于 Python 的 json.loads(s)：把 JSON 字符串解析为 Go 值。
//
// 整数返回 int64，浮点返回 float64，对象返回 map[string]any，数组返回 []any。
// 超出 int64 的整数保留为 json.Number（Python 是任意精度 int）。
//
// Python: json.loads(s)
// DIFF: 大整数不保证精确为 int；不支持 object_hook/parse_float 等参数。
func Loads(data string) (any, error) {
	decoder := stdjson.NewDecoder(strings.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("json.Loads: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("json.Loads: trailing JSON value")
		}
		return nil, fmt.Errorf("json.Loads: trailing data: %w", err)
	}
	return normalize(value), nil
}

// LoadsInto 把 JSON 字符串解析到目标值，字段规则遵循 encoding/json。
func LoadsInto(data string, out any) error {
	if err := stdjson.Unmarshal([]byte(data), out); err != nil {
		return fmt.Errorf("json.LoadsInto: %w", err)
	}
	return nil
}

// Dumps 等价于 Python 的 json.dumps(v)：返回紧凑 JSON，不转义 HTML 字符。
//
// Python: json.dumps(v)
// DIFF: 默认不开启 ensure_ascii；非 ASCII 字符原样输出。
func Dumps(v any) (string, error) {
	return conv.ToJSON(v)
}

// DumpsIndent 等价于 Python 的 json.dumps(v, indent=indent)。
func DumpsIndent(v any, indent string) (string, error) {
	return conv.ToJSONIndent(v, indent)
}

// Load 等价于 Python 的 json.load(fp)：从文件读取并解析 JSON。
func Load(path string) (any, error) {
	data, err := file.Read(path)
	if err != nil {
		return nil, fmt.Errorf("json.Load: %w", err)
	}
	return Loads(data)
}

// Dump 等价于 Python 的 json.dump(v, fp)：把值写入文件。
func Dump(path string, v any) error {
	data, err := Dumps(v)
	if err != nil {
		return fmt.Errorf("json.Dump: %w", err)
	}
	return file.Write(path, data)
}

// Indent 等价于 Python 的 json.dumps(json.loads(s), indent=indent)。
func Indent(data string, indent string) (string, error) {
	value, err := Loads(data)
	if err != nil {
		return "", err
	}
	return DumpsIndent(value, indent)
}

func normalize(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for key, item := range v {
			v[key] = normalize(item)
		}
		return v
	case []any:
		for i, item := range v {
			v[i] = normalize(item)
		}
		return v
	case stdjson.Number:
		text := v.String()
		if !strings.ContainsAny(text, ".eE") {
			if integer, err := v.Int64(); err == nil {
				return integer
			}
			return v
		}
		if number, err := v.Float64(); err == nil {
			return number
		}
		return v
	}
	return value
}
