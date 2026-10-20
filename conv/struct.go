package conv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	stdtime "time"

	"github.com/BK-Go-Python/gotools/internal/pyval"
)

var timeType = reflect.TypeOf(stdtime.Time{})

// StructToMap 把结构体转换为 map，字段名遵循 encoding/json 的规则：
// 优先使用 json tag，忽略 json:"-"，支持 omitempty 和匿名嵌入字段展开。
//
// 与 JSON 序列化不同，数字保持原本的 int/float 类型，不会全部变成 float64。
func StructToMap(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("conv.StructToMap: cannot convert NoneType to map")
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil, fmt.Errorf("conv.StructToMap: cannot convert nil pointer to map")
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("conv.StructToMap: expected struct, got %s", pyval.TypeName(v))
	}
	return structToMap(rv), nil
}

// FromStruct 是 StructToMap 的别名，更贴近 Python 的 vars(obj)。
func FromStruct(v interface{}) (map[string]interface{}, error) {
	return StructToMap(v)
}

// MapToStruct 把 map 填充到结构体指针，遵循 encoding/json 的字段匹配规则。
//
// 由于底层使用 JSON 编解码，map 中的值必须可 JSON 序列化；time.Time、
// 嵌套结构体、切片等常见类型都可以正常工作。
func MapToStruct(m map[string]interface{}, out interface{}) error {
	if out == nil {
		return fmt.Errorf("conv.MapToStruct: out must be a non-nil pointer to struct")
	}
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("conv.MapToStruct: out must be a non-nil pointer to struct")
	}
	if rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("conv.MapToStruct: out must point to struct, got %s", rv.Elem().Kind())
	}
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("conv.MapToStruct: marshal map: %w", err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("conv.MapToStruct: unmarshal into struct: %w", err)
	}
	return nil
}

func structToMap(rv reflect.Value) map[string]interface{} {
	t := rv.Type()
	out := make(map[string]interface{}, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := rv.Field(i)
		// 未导出类型的匿名嵌入字段本身不可 Interface，但它的导出字段
		// 会被 encoding/json 提升；这里按同样规则展开。
		if field.PkgPath != "" {
			if field.Anonymous {
				if inner, ok := indirectStruct(fv); ok {
					for k, val := range structToMap(inner) {
						if _, exists := out[k]; !exists {
							out[k] = val
						}
					}
				}
			}
			continue
		}
		name := field.Name
		omitempty := false
		tag, hasTag := field.Tag.Lookup("json")
		if hasTag {
			parts := strings.Split(tag, ",")
			if parts[0] == "-" {
				continue
			}
			if parts[0] != "" {
				name = parts[0]
			}
			for _, option := range parts[1:] {
				if option == "omitempty" {
					omitempty = true
				}
			}
		}
		if omitempty && fv.IsZero() {
			continue
		}
		if field.Anonymous && (!hasTag || name == field.Name) {
			inner := fv
			for inner.Kind() == reflect.Ptr {
				if inner.IsNil() {
					inner = reflect.Value{}
					break
				}
				inner = inner.Elem()
			}
			if inner.IsValid() && inner.Kind() == reflect.Struct {
				for k, val := range structToMap(inner) {
					if _, exists := out[k]; !exists {
						out[k] = val
					}
				}
				continue
			}
		}
		out[name] = mapValue(fv)
	}
	return out
}

func mapValue(fv reflect.Value) interface{} {
	if !fv.IsValid() {
		return nil
	}
	for fv.Kind() == reflect.Interface {
		if fv.IsNil() {
			return nil
		}
		fv = fv.Elem()
	}
	switch fv.Kind() {
	case reflect.Ptr:
		if fv.IsNil() {
			return nil
		}
		return mapValue(fv.Elem())
	case reflect.Struct:
		if fv.Type() == timeType {
			return fv.Interface()
		}
		return structToMap(fv)
	case reflect.Slice:
		if fv.Type().Elem().Kind() == reflect.Uint8 {
			return fv.Bytes()
		}
		return sliceToInterfaces(fv)
	case reflect.Array:
		return sliceToInterfaces(fv)
	case reflect.Map:
		if fv.Type().Key().Kind() == reflect.String {
			out := make(map[string]interface{}, fv.Len())
			for _, key := range fv.MapKeys() {
				out[key.String()] = mapValue(fv.MapIndex(key))
			}
			return out
		}
		return fv.Interface()
	default:
		return fv.Interface()
	}
}

func sliceToInterfaces(fv reflect.Value) []interface{} {
	out := make([]interface{}, fv.Len())
	for i := 0; i < fv.Len(); i++ {
		out[i] = mapValue(fv.Index(i))
	}
	return out
}

func indirectStruct(fv reflect.Value) (reflect.Value, bool) {
	for fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return reflect.Value{}, false
		}
		fv = fv.Elem()
	}
	if fv.Kind() == reflect.Struct {
		return fv, true
	}
	return reflect.Value{}, false
}
