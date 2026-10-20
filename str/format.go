package str

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

// Format 实现 Python str.format / f-string 的常用子集。
//
// 支持：
//   - 自动编号 {}、位置编号 {0}、关键字 {name}（args 中的 map 或 struct 字段）
//   - 转义 {{ 和 }}
//   - 填充与对齐 [fill][< > ^]、宽度、精度和类型 s d f e g x X o b %
//
// 例如：
//
//	str.Format("{} 有 {:.2f} 元", "Tom", 12.345) // "Tom 有 12.35 元"
//	str.Format("{0}-{0}", "a")                   // "a-a"
//	str.Format("{name}", map[string]interface{}{"name": "Go"})
//
// 无法解析的占位符会原样保留，方便排查；不会 panic。
func Format(template string, args ...interface{}) string {
	return format(template, args, nil)
}

// FormatMap 等价于 Python 的 template.format(**values)，只使用关键字参数。
//
//	str.FormatMap("{name} is {age}", map[string]interface{}{"name": "Tom", "age": 18})
func FormatMap(template string, values map[string]interface{}) string {
	return format(template, []interface{}{values}, values)
}

func format(template string, args []interface{}, named map[string]interface{}) string {
	var b strings.Builder
	b.Grow(len(template) + 16)
	auto := 0

	for i := 0; i < len(template); {
		switch template[i] {
		case '{':
			if i+1 < len(template) && template[i+1] == '{' {
				b.WriteByte('{')
				i += 2
				continue
			}
			rel := strings.IndexByte(template[i+1:], '}')
			if rel < 0 {
				b.WriteString(template[i:])
				return b.String()
			}
			field := template[i+1 : i+1+rel]
			name, spec := splitField(field)
			value, found := lookupField(name, args, named, &auto)
			if !found {
				b.WriteString(template[i : i+rel+2])
			} else {
				b.WriteString(formatValue(value, spec))
			}
			i += rel + 2
		case '}':
			if i+1 < len(template) && template[i+1] == '}' {
				b.WriteByte('}')
				i += 2
				continue
			}
			b.WriteByte('}')
			i++
		default:
			b.WriteByte(template[i])
			i++
		}
	}
	return b.String()
}

func splitField(field string) (string, string) {
	if i := strings.IndexByte(field, ':'); i >= 0 {
		return field[:i], field[i+1:]
	}
	return field, ""
}

func lookupField(name string, args []interface{}, named map[string]interface{}, auto *int) (interface{}, bool) {
	if name == "" {
		if *auto >= len(args) {
			return nil, false
		}
		v := args[*auto]
		*auto++
		return v, true
	}
	if idx, err := strconv.Atoi(name); err == nil {
		if idx >= 0 && idx < len(args) {
			return args[idx], true
		}
		return nil, false
	}
	if named != nil {
		if v, ok := named[name]; ok {
			return v, true
		}
	}
	for _, arg := range args {
		if v, ok := fieldOf(arg, name); ok {
			return v, true
		}
	}
	return nil, false
}

func fieldOf(v interface{}, name string) (interface{}, bool) {
	if v == nil {
		return nil, false
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil, false
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return nil, false
		}
		key := reflect.ValueOf(name).Convert(rv.Type().Key())
		if value := rv.MapIndex(key); value.IsValid() {
			return value.Interface(), true
		}
	case reflect.Struct:
		if f := rv.FieldByName(name); f.IsValid() && f.CanInterface() {
			return f.Interface(), true
		}
		t := rv.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" {
				continue
			}
			nameByTag := f.Name
			if tag, ok := f.Tag.Lookup("json"); ok {
				parts := strings.Split(tag, ",")
				if parts[0] == "-" {
					continue
				}
				if parts[0] != "" {
					nameByTag = parts[0]
				}
			}
			if nameByTag == name && rv.Field(i).CanInterface() {
				return rv.Field(i).Interface(), true
			}
		}
	}
	return nil, false
}

type formatSpec struct {
	fill      rune
	align     byte
	width     int
	precision int
	typ       byte
	zero      bool
}

func parseFormatSpec(spec string) (formatSpec, bool) {
	fs := formatSpec{fill: ' ', precision: -1}
	runes := []rune(spec)
	i := 0
	if len(runes) >= 2 && isAlignRune(runes[1]) {
		fs.fill = runes[0]
		fs.align = byte(runes[1])
		i = 2
	} else if len(runes) >= 1 && isAlignRune(runes[0]) {
		fs.align = byte(runes[0])
		i = 1
	}
	if i < len(runes) && (runes[i] == '+' || runes[i] == '-' || runes[i] == ' ') {
		i++
	}
	if i < len(runes) && runes[i] == '0' {
		fs.zero = true
		fs.fill = '0'
		i++
	}
	start := i
	for i < len(runes) && runes[i] >= '0' && runes[i] <= '9' {
		i++
	}
	if i > start {
		fs.width, _ = strconv.Atoi(string(runes[start:i]))
	}
	if i < len(runes) && runes[i] == '.' {
		i++
		start = i
		for i < len(runes) && runes[i] >= '0' && runes[i] <= '9' {
			i++
		}
		if i == start {
			fs.precision = 0
		} else {
			fs.precision, _ = strconv.Atoi(string(runes[start:i]))
		}
	}
	if i < len(runes) {
		if runes[i] > 0x7F {
			return fs, false
		}
		fs.typ = byte(runes[i])
		i++
	}
	if i != len(runes) {
		return fs, false
	}
	return fs, true
}

func formatValue(v interface{}, spec string) string {
	if spec == "" {
		return pyval.Str(v)
	}
	fs, ok := parseFormatSpec(spec)
	if !ok {
		return pyval.Str(v)
	}

	numeric := false
	var out string
	switch fs.typ {
	case 0:
		out = pyval.Str(v)
		numeric = isNumeric(v)
		if fs.precision >= 0 {
			if isFloat(v) {
				if f, ok := toFloat64(v); ok {
					out = strconv.FormatFloat(f, 'g', fs.precision, 64)
				}
			} else if !numeric {
				out = truncateRunes(out, fs.precision)
			}
		}
	case 's':
		out = pyval.Str(v)
		if fs.precision >= 0 {
			out = truncateRunes(out, fs.precision)
		}
	case 'd', 'x', 'X', 'o', 'b':
		n, ok := toInt64(v)
		if !ok {
			return pyval.Str(v)
		}
		numeric = true
		switch fs.typ {
		case 'd':
			out = strconv.FormatInt(n, 10)
		case 'x':
			out = strconv.FormatInt(n, 16)
		case 'X':
			out = strings.ToUpper(strconv.FormatInt(n, 16))
		case 'o':
			out = strconv.FormatInt(n, 8)
		case 'b':
			out = strconv.FormatInt(n, 2)
		}
	case 'f', 'F', 'e', 'E', 'g', 'G':
		f, ok := toFloat64(v)
		if !ok {
			return pyval.Str(v)
		}
		numeric = true
		precision := fs.precision
		if precision < 0 {
			precision = 6
		}
		verb := fs.typ
		if verb == 'F' {
			verb = 'f'
		}
		out = strconv.FormatFloat(f, verb, precision, 64)
	case '%':
		f, ok := toFloat64(v)
		if !ok {
			return pyval.Str(v)
		}
		numeric = true
		precision := fs.precision
		if precision < 0 {
			precision = 6
		}
		out = strconv.FormatFloat(f*100, 'f', precision, 64) + "%"
	default:
		return pyval.Str(v)
	}

	return padFormat(out, fs, numeric)
}

func padFormat(out string, fs formatSpec, numeric bool) string {
	width := fs.width
	if width <= 0 {
		return out
	}
	n := utf8.RuneCountInString(out)
	if n >= width {
		return out
	}
	padding := width - n
	align := fs.align
	if align == 0 {
		if numeric {
			align = '>'
		} else {
			align = '<'
		}
	}
	if fs.zero && (align == '>' || align == '=') {
		sign := ""
		body := out
		if strings.HasPrefix(out, "-") || strings.HasPrefix(out, "+") || strings.HasPrefix(out, " ") {
			sign = out[:1]
			body = out[1:]
		}
		return sign + strings.Repeat("0", padding) + body
	}
	switch align {
	case '<':
		return out + strings.Repeat(string(fs.fill), padding)
	case '^':
		left := padding / 2
		return strings.Repeat(string(fs.fill), left) + out + strings.Repeat(string(fs.fill), padding-left)
	case '=':
		sign := ""
		body := out
		if strings.HasPrefix(out, "-") || strings.HasPrefix(out, "+") || strings.HasPrefix(out, " ") {
			sign = out[:1]
			body = out[1:]
		}
		return sign + strings.Repeat(string(fs.fill), padding) + body
	default:
		return strings.Repeat(string(fs.fill), padding) + out
	}
}

func isAlignRune(r rune) bool {
	return r == '<' || r == '>' || r == '^' || r == '='
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if n >= len(r) {
		return s
	}
	if n < 0 {
		n = 0
	}
	return string(r[:n])
}

func isFloat(v interface{}) bool {
	switch v.(type) {
	case float32, float64:
		return true
	}
	return false
}

func isNumeric(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64:
		return true
	}
	return false
}

func toInt64(v interface{}) (int64, bool) {
	switch x := v.(type) {
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	case int:
		return int64(x), true
	case int8:
		return int64(x), true
	case int16:
		return int64(x), true
	case int32:
		return int64(x), true
	case int64:
		return x, true
	case uint:
		if uint64(x) <= math.MaxInt64 {
			return int64(x), true
		}
	case uint8:
		return int64(x), true
	case uint16:
		return int64(x), true
	case uint32:
		return int64(x), true
	case uint64:
		if x <= math.MaxInt64 {
			return int64(x), true
		}
	case uintptr:
		if uint64(x) <= math.MaxInt64 {
			return int64(x), true
		}
	}
	return 0, false
}

func toFloat64(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	case int:
		return float64(x), true
	case int8:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint8:
		return float64(x), true
	case uint16:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	case uintptr:
		return float64(x), true
	case float32:
		return float64(x), true
	case float64:
		return x, true
	}
	return 0, false
}
