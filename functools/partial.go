package functools

import (
	"fmt"
	"reflect"
)

// Partial 是反射版 functools.partial：绑定部分参数后调用。
//
// 由于 Go 没有关键字参数和可变参数函数值，Partial.Call 返回 []any。
// 对类型安全场景，优先使用 Partial1/Partial2。
type Partial struct {
	fn    reflect.Value
	bound []any
}

// NewPartial 创建 Partial，并预绑定 bound 参数。
func NewPartial(fn any, bound ...any) (*Partial, error) {
	value := reflect.ValueOf(fn)
	if !value.IsValid() || value.Kind() != reflect.Func {
		return nil, fmt.Errorf("functools.NewPartial: not a function")
	}
	return &Partial{fn: value, bound: append([]any(nil), bound...)}, nil
}

// Call 追加 args 并调用原函数；返回值按原顺序放入 []any。
func (p *Partial) Call(args ...any) (result []any, err error) {
	if p == nil {
		return nil, fmt.Errorf("functools.Partial: nil receiver")
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("functools.Partial: panic: %v", recovered)
			result = nil
		}
	}()

	all := append(append([]any(nil), p.bound...), args...)
	inputs, err := buildInputs(p.fn.Type(), all)
	if err != nil {
		return nil, err
	}
	var outputs []reflect.Value
	if p.fn.Type().IsVariadic() {
		outputs = p.fn.CallSlice(inputs)
	} else {
		outputs = p.fn.Call(inputs)
	}
	result = make([]any, len(outputs))
	for i, output := range outputs {
		result[i] = output.Interface()
	}
	return result, nil
}

// Partial1 类型安全地绑定一元函数的第一个参数。
func Partial1[A, B any](fn func(A) B, a A) func() B {
	return func() B { return fn(a) }
}

// Partial2 类型安全地绑定二元函数的第一个参数。
func Partial2[A, B, C any](fn func(A, B) C, a A) func(B) C {
	return func(b B) C { return fn(a, b) }
}

func buildInputs(fnType reflect.Type, args []any) ([]reflect.Value, error) {
	fixed := fnType.NumIn()
	variadic := fnType.IsVariadic()
	if variadic {
		fixed--
	}
	if len(args) < fixed {
		return nil, fmt.Errorf("functools.Partial: too few arguments: got %d, need %d", len(args), fixed)
	}
	if !variadic && len(args) != fixed {
		return nil, fmt.Errorf("functools.Partial: wrong argument count: got %d, want %d", len(args), fixed)
	}

	inputs := make([]reflect.Value, 0, fnType.NumIn())
	for i := 0; i < fixed; i++ {
		value, err := convertArgument(args[i], fnType.In(i))
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, value)
	}
	if !variadic {
		return inputs, nil
	}

	elemType := fnType.In(fnType.NumIn() - 1).Elem()
	variadicArgs := reflect.MakeSlice(fnType.In(fnType.NumIn()-1), 0, len(args)-fixed)
	for _, arg := range args[fixed:] {
		value, err := convertArgument(arg, elemType)
		if err != nil {
			return nil, err
		}
		variadicArgs = reflect.Append(variadicArgs, value)
	}
	inputs = append(inputs, variadicArgs)
	return inputs, nil
}

func convertArgument(value any, target reflect.Type) (reflect.Value, error) {
	if value == nil {
		switch target.Kind() {
		case reflect.Interface, reflect.Pointer, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
			return reflect.Zero(target), nil
		default:
			return reflect.Value{}, fmt.Errorf("functools.Partial: nil for %s", target)
		}
	}
	rv := reflect.ValueOf(value)
	if rv.Type().AssignableTo(target) {
		return rv, nil
	}
	if rv.Type().ConvertibleTo(target) {
		return rv.Convert(target), nil
	}
	return reflect.Value{}, fmt.Errorf("functools.Partial: cannot use %s as %s", rv.Type(), target)
}
