package py

import "fmt"

// NewSlice 等价于 Python 的 slice(start, stop, step)。
//
// 参数可以是 nil 或 int：0 个参数表示全为 None；1 个是 stop；
// 2 个是 start、stop；3 个是 start、stop、step。nil 表示 Python 的 None。
//
// Python: slice(start, stop, step)
func NewSlice(args ...any) (Slice, error) {
	if len(args) > 3 {
		return Slice{}, fmt.Errorf("py.NewSlice: expected at most 3 arguments, got %d: %w", len(args), ErrOutOfRange)
	}
	var values [3]int
	var present [3]bool
	for i, arg := range args {
		if arg == nil {
			continue
		}
		n, ok := arg.(int)
		if !ok {
			return Slice{}, fmt.Errorf("py.NewSlice: argument %d must be int or nil, got %T: %w", i, arg, ErrType)
		}
		values[i] = n
		present[i] = true
	}

	switch len(args) {
	case 0:
		return Slice{}, nil
	case 1:
		return Slice{Stop: values[0], HasStop: present[0]}, nil
	case 2:
		return Slice{
			Start: values[0], HasStart: present[0],
			Stop: values[1], HasStop: present[1],
		}, nil
	default:
		if present[2] && values[2] == 0 {
			return Slice{}, fmt.Errorf("py.NewSlice: step must not be zero: %w", ErrOutOfRange)
		}
		return Slice{
			Start: values[0], HasStart: present[0],
			Stop: values[1], HasStop: present[1],
			Step: values[2], HasStep: present[2],
		}, nil
	}
}

// Indices 等价于 Python 的 slice.indices(n)，返回规范化后的 start、stop、step。
//
// 直接构造 Slice{Step: 0} 时按 step=1 处理；NewSlice 会拒绝 step=0。
func (s Slice) Indices(n int) (start, stop, step int) {
	step = 1
	if s.HasStep {
		step = s.Step
	}
	if step == 0 {
		step = 1
	}

	if step > 0 {
		start = 0
		if s.HasStart {
			start = s.Start
		}
		stop = n
		if s.HasStop {
			stop = s.Stop
		}
		if start < 0 {
			start += n
		}
		if stop < 0 {
			stop += n
		}
		if start < 0 {
			start = 0
		}
		if start > n {
			start = n
		}
		if stop < 0 {
			stop = 0
		}
		if stop > n {
			stop = n
		}
		return start, stop, step
	}

	if s.HasStart {
		start = s.Start
		if start < 0 {
			start += n
		}
	} else {
		start = n - 1
	}
	if s.HasStop {
		stop = s.Stop
		if stop < 0 {
			stop += n
		}
	} else {
		stop = -1
	}
	if start < -1 {
		start = -1
	}
	if start > n-1 {
		start = n - 1
	}
	if stop < -1 {
		stop = -1
	}
	if stop > n-1 {
		stop = n - 1
	}
	return start, stop, step
}
