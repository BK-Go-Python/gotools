package py

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/BK-Go-Python/gotools/conv"
)

// Callable 等价于 Python 的 callable(v)，Go 中只有函数类型算可调用。
//
// Python: callable(v)
// DIFF: 不支持实现了 __call__ 的对象。
func Callable(v any) bool {
	return v != nil && reflect.TypeOf(v).Kind() == reflect.Func
}

// GetAttr 等价于 Python 的 getattr(v, name[, default])。
//
// Python: getattr(v, name[, default])
// DIFF: 仅支持 map 键、导出字段（含 json tag）和导出方法；不支持 __getattr__。
func GetAttr(v any, name string, def ...any) (any, error) {
	if v != nil {
		original := reflect.ValueOf(v)
		rv := original
		for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				break
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Map:
			if rv.Type().Key().Kind() == reflect.String {
				key := reflect.ValueOf(name).Convert(rv.Type().Key())
				if value := rv.MapIndex(key); value.IsValid() {
					return value.Interface(), nil
				}
			}
		case reflect.Struct:
			if field := structField(rv, name); field.IsValid() && field.CanInterface() {
				return field.Interface(), nil
			}
			if method := original.MethodByName(name); method.IsValid() {
				return method.Interface(), nil
			}
			if method := rv.MethodByName(name); method.IsValid() {
				return method.Interface(), nil
			}
		}
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return nil, fmt.Errorf("py.GetAttr: attribute %q not found: %w", name, ErrNotFound)
}

// HasAttr 等价于 Python 的 hasattr(v, name)。
//
// Python: hasattr(v, name)
func HasAttr(v any, name string) bool {
	_, err := GetAttr(v, name)
	return err == nil
}

// SetAttr 等价于 Python 的 setattr(v, name, value)。
//
// Python: setattr(v, name, value)
// DIFF: 仅支持非 nil map 的 string 键和可设置的导出 struct 字段。
func SetAttr(v any, name string, value any) error {
	if v == nil {
		return fmt.Errorf("py.SetAttr: cannot set attribute on NoneType: %w", ErrType)
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return fmt.Errorf("py.SetAttr: cannot set attribute on nil: %w", ErrType)
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("py.SetAttr: map key must be string: %w", ErrType)
		}
		if rv.IsNil() {
			return fmt.Errorf("py.SetAttr: cannot assign into nil map: %w", ErrType)
		}
		key := reflect.ValueOf(name).Convert(rv.Type().Key())
		assigned, err := assignValue(value, rv.Type().Elem())
		if err != nil {
			return fmt.Errorf("py.SetAttr: %w", err)
		}
		rv.SetMapIndex(key, assigned)
		return nil
	case reflect.Struct:
		field := structField(rv, name)
		if !field.IsValid() {
			return fmt.Errorf("py.SetAttr: attribute %q not found: %w", name, ErrNotFound)
		}
		if !field.CanSet() {
			return fmt.Errorf("py.SetAttr: attribute %q is not settable: %w", name, ErrType)
		}
		assigned, err := assignValue(value, field.Type())
		if err != nil {
			return fmt.Errorf("py.SetAttr: %w", err)
		}
		field.Set(assigned)
		return nil
	default:
		return fmt.Errorf("py.SetAttr: unsupported type %s: %w", rv.Kind(), ErrType)
	}
}

// DelAttr 等价于 Python 的 delattr(v, name)。
//
// Python: delattr(v, name)
// DIFF: 只支持 map 键；Go 的 struct 字段不能删除。
func DelAttr(v any, name string) error {
	if v == nil {
		return fmt.Errorf("py.DelAttr: cannot delete attribute on NoneType: %w", ErrType)
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return fmt.Errorf("py.DelAttr: object is nil: %w", ErrType)
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("py.DelAttr: only map[string]T supports deletion: %w", ErrType)
	}
	key := reflect.ValueOf(name).Convert(rv.Type().Key())
	if !rv.MapIndex(key).IsValid() {
		return fmt.Errorf("py.DelAttr: attribute %q not found: %w", name, ErrNotFound)
	}
	rv.SetMapIndex(key, reflect.Value{})
	return nil
}

// Dir 等价于 Python 的 dir(v)，返回排序后的属性名。
//
// Python: dir(v)
// DIFF: 只返回 map 键、导出字段和导出方法；没有 __dir__ 协议。
func Dir(v any) []string {
	names := map[string]struct{}{}
	if v != nil {
		rv := reflect.ValueOf(v)
		for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				rv = reflect.Value{}
				break
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Map:
			if rv.Type().Key().Kind() == reflect.String {
				for _, key := range rv.MapKeys() {
					names[key.String()] = struct{}{}
				}
			}
		case reflect.Struct:
			t := rv.Type()
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if field.PkgPath == "" {
					names[field.Name] = struct{}{}
				}
			}
			for i := 0; i < t.NumMethod(); i++ {
				method := t.Method(i)
				if method.PkgPath == "" {
					names[method.Name] = struct{}{}
				}
			}
		}
	}

	out := make([]string, 0, len(names))
	for name := range names {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// IsInstance 等价于 Python 的 isinstance(v, types)。
//
// Python: isinstance(v, (T1, T2, ...))
// DIFF: types 元素可以是 reflect.Type 或示例值；只支持 Go 的类型/接口，
// 不支持 ABC、虚拟子类和 __instancecheck__。
func IsInstance(v any, types ...any) bool {
	for _, target := range types {
		typ := typeOfTarget(target)
		if typ == nil {
			continue
		}
		if v == nil {
			if typ.Kind() == reflect.Interface {
				return true
			}
			continue
		}
		valueType := reflect.TypeOf(v)
		if valueType.AssignableTo(typ) || (typ.Kind() == reflect.Interface && valueType.Implements(typ)) {
			return true
		}
	}
	return false
}

// IsSubclass 近似等价于 Python 的 issubclass(child, parent)。
//
// Python: issubclass(child, parent)
// DIFF: Go 没有继承；这里只支持类型相同或 child 实现 parent 接口。
func IsSubclass(child, parent any) bool {
	childType := typeOfTarget(child)
	parentType := typeOfTarget(parent)
	if childType == nil || parentType == nil {
		return false
	}
	if childType == parentType {
		return true
	}
	if parentType.Kind() == reflect.Interface {
		return childType.Implements(parentType)
	}
	return false
}

// TypeOf 等价于 Python 的 type(v)，返回 Go 的 reflect.Type。
//
// Python: type(v)
// DIFF: 现有的 py.Type 返回类型名字符串；TypeOf 返回类型对象。
func TypeOf(v any) reflect.Type {
	return reflect.TypeOf(v)
}

// Vars 等价于 Python 的 vars(v)。
//
// Python: vars(v)
// DIFF: struct 使用 json tag 映射为 map；map[string]T 返回浅拷贝；
// Go 没有 __dict__，无法返回对象的动态属性表。
func Vars(v any) (map[string]any, error) {
	if v == nil {
		return nil, fmt.Errorf("py.Vars: cannot get vars of NoneType: %w", ErrType)
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil, fmt.Errorf("py.Vars: cannot get vars of nil: %w", ErrType)
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Struct:
		return conv.StructToMap(v)
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("py.Vars: map key must be string: %w", ErrType)
		}
		out := make(map[string]any, rv.Len())
		for _, key := range rv.MapKeys() {
			out[key.String()] = rv.MapIndex(key).Interface()
		}
		return out, nil
	}
	return nil, fmt.Errorf("py.Vars: unsupported type %s: %w", rv.Kind(), ErrType)
}

func structField(rv reflect.Value, name string) reflect.Value {
	if field := rv.FieldByName(name); field.IsValid() && field.CanInterface() {
		return field
	}
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tagName := field.Name
		if tag, ok := field.Tag.Lookup("json"); ok {
			parts := strings.Split(tag, ",")
			if parts[0] == "-" {
				continue
			}
			if parts[0] != "" {
				tagName = parts[0]
			}
		}
		if tagName == name && rv.Field(i).CanInterface() {
			return rv.Field(i)
		}
	}
	return reflect.Value{}
}

func assignValue(value any, target reflect.Type) (reflect.Value, error) {
	if value == nil {
		switch target.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
			return reflect.Zero(target), nil
		}
		return reflect.Value{}, fmt.Errorf("cannot assign None to %s: %w", target, ErrType)
	}
	rv := reflect.ValueOf(value)
	if rv.Type().AssignableTo(target) {
		return rv, nil
	}
	if rv.Type().ConvertibleTo(target) {
		return rv.Convert(target), nil
	}
	return reflect.Value{}, fmt.Errorf("cannot assign %s to %s: %w", rv.Type(), target, ErrType)
}

func typeOfTarget(target any) reflect.Type {
	if target == nil {
		return nil
	}
	if typ, ok := target.(reflect.Type); ok {
		return typ
	}
	return reflect.TypeOf(target)
}
