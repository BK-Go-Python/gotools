package py

import (
	"fmt"

	"github.com/BK-Go-Python/gotools/internal/pyval"
	"github.com/BK-Go-Python/gotools/list"
)

// Catch 执行 fn 并把 panic 转换为 error；没有 panic 时返回 nil。
//
// Go 没有 Python 的 try/except；这是给旧动态 API 或第三方回调兜底的
// 通用入口。新 API 应直接返回 error。
func Catch(fn func()) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = panicError(recovered)
		}
	}()
	if fn != nil {
		fn()
	}
	return nil
}

// Try 执行 fn，把 panic 转换为 error，并返回 fn 的结果。
func Try[T any](fn func() T) (value T, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			var zero T
			value = zero
			err = panicError(recovered)
		}
	}()
	if fn == nil {
		return value, fmt.Errorf("py.Try: nil function")
	}
	return fn(), nil
}

// Try2 是返回两个值的 Try。
func Try2[A, B any](fn func() (A, B)) (a A, b B, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			var zeroA A
			var zeroB B
			a, b = zeroA, zeroB
			err = panicError(recovered)
		}
	}()
	if fn == nil {
		return a, b, fmt.Errorf("py.Try2: nil function")
	}
	a, b = fn()
	return a, b, nil
}

// Try3 是返回三个值的 Try。
func Try3[A, B, C any](fn func() (A, B, C)) (a A, b B, c C, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			var zeroA A
			var zeroB B
			var zeroC C
			a, b, c = zeroA, zeroB, zeroC
			err = panicError(recovered)
		}
	}()
	if fn == nil {
		return a, b, c, fmt.Errorf("py.Try3: nil function")
	}
	a, b, c = fn()
	return a, b, c, nil
}

// RangeE 是 Range 的安全版本：参数非法/step 为 0 时返回 error，不 panic。
func RangeE(args ...int) (out []int, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			out = nil
			err = panicError(recovered)
		}
	}()
	return Range(args...), nil
}

// FloorDivE 是 FloorDiv 的安全版本，除数为 0 时返回 ErrZeroDivision。
func FloorDivE(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("py.FloorDivE: %w", ErrZeroDivision)
	}
	return FloorDiv(a, b), nil
}

// ModE 是 Mod 的安全版本，除数为 0 时返回 ErrZeroDivision。
func ModE(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("py.ModE: %w", ErrZeroDivision)
	}
	return Mod(a, b), nil
}

// ReduceE 是 Reduce 的安全版本；空序列无初始值返回 ErrEmpty。
func ReduceE(fn func(acc, item interface{}) interface{}, iterable interface{}, initial ...interface{}) (result interface{}, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = nil
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return nil, fmt.Errorf("py.ReduceE: %w: %v", ErrType, err)
	}
	start := 0
	if len(initial) > 0 {
		result = initial[0]
	} else {
		if len(items) == 0 {
			return nil, fmt.Errorf("py.ReduceE: %w", ErrEmpty)
		}
		result = items[0]
		start = 1
	}
	for _, item := range items[start:] {
		result = fn(result, item)
	}
	return result, nil
}

// AnyE 是 Any 的安全版本；非法 iterable 返回 error。
func AnyE(iterable interface{}) (result bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = false
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return false, fmt.Errorf("py.AnyE: %w: %v", ErrType, err)
	}
	for _, item := range items {
		if pyval.Truthy(item) {
			return true, nil
		}
	}
	return false, nil
}

// AllE 是 All 的安全版本；非法 iterable 返回 error。
func AllE(iterable interface{}) (result bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = false
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return false, fmt.Errorf("py.AllE: %w: %v", ErrType, err)
	}
	for _, item := range items {
		if !pyval.Truthy(item) {
			return false, nil
		}
	}
	return true, nil
}

// EnumerateAnyE 是 EnumerateAny 的安全版本。
func EnumerateAnyE(iterable interface{}, start ...int) (out []Pair, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			out = nil
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return nil, fmt.Errorf("py.EnumerateAnyE: %w: %v", ErrType, err)
	}
	offset := 0
	if len(start) > 0 {
		offset = start[0]
	}
	out = make([]Pair, len(items))
	for i, item := range items {
		out[i] = Pair{Index: offset + i, Value: item}
	}
	return out, nil
}

// ZipE 是 Zip 的安全版本。
func ZipE(iterables ...interface{}) (out [][]interface{}, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			out = nil
			err = panicError(recovered)
		}
	}()
	if len(iterables) == 0 {
		return [][]interface{}{}, nil
	}
	lists := make([][]interface{}, len(iterables))
	length := -1
	for i, iterable := range iterables {
		items, err := pyval.Iter(iterable)
		if err != nil {
			return nil, fmt.Errorf("py.ZipE: %w: %v", ErrType, err)
		}
		lists[i] = items
		if length < 0 || len(items) < length {
			length = len(items)
		}
	}
	out = make([][]interface{}, length)
	for i := 0; i < length; i++ {
		row := make([]interface{}, len(lists))
		for j, items := range lists {
			row[j] = items[i]
		}
		out[i] = row
	}
	return out, nil
}

// SortedAnyE 是 SortedAny 的安全版本。
func SortedAnyE(iterable interface{}, less ...func(a, b interface{}) bool) (out list.List, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			out = nil
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return nil, fmt.Errorf("py.SortedAnyE: %w: %v", ErrType, err)
	}
	out = list.New(items...)
	if len(less) > 0 && less[0] != nil {
		out.Sort(less[0])
	} else {
		out.Sort()
	}
	return out, nil
}

// SortedByE 是 SortedBy 的安全版本。
func SortedByE(iterable interface{}, key func(interface{}) interface{}, reverse ...bool) (out list.List, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			out = nil
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return nil, fmt.Errorf("py.SortedByE: %w: %v", ErrType, err)
	}
	out = list.New(items...)
	out.SortBy(key, reverse...)
	return out, nil
}

// ReversedAnyE 是 ReversedAny 的安全版本。
func ReversedAnyE(iterable interface{}) (out list.List, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			out = nil
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return nil, fmt.Errorf("py.ReversedAnyE: %w: %v", ErrType, err)
	}
	out = make(list.List, len(items))
	for i, item := range items {
		out[len(items)-1-i] = item
	}
	return out, nil
}

// MapAnyE 是 MapAny 的安全版本。
func MapAnyE(fn func(interface{}) interface{}, iterable interface{}) (out list.List, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			out = nil
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return nil, fmt.Errorf("py.MapAnyE: %w: %v", ErrType, err)
	}
	out = make(list.List, 0, len(items))
	for _, item := range items {
		out = append(out, fn(item))
	}
	return out, nil
}

// FilterAnyE 是 FilterAny 的安全版本。
func FilterAnyE(fn func(interface{}) bool, iterable interface{}) (out list.List, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			out = nil
			err = panicError(recovered)
		}
	}()
	items, err := pyval.Iter(iterable)
	if err != nil {
		return nil, fmt.Errorf("py.FilterAnyE: %w: %v", ErrType, err)
	}
	out = make(list.List, 0, len(items))
	for _, item := range items {
		if fn(item) {
			out = append(out, item)
		}
	}
	return out, nil
}

func panicError(recovered any) error {
	if err, ok := recovered.(error); ok {
		return fmt.Errorf("py: recovered panic: %w", err)
	}
	return fmt.Errorf("py: recovered panic: %v", recovered)
}
