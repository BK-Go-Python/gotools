package conv

import "fmt"

// Complex 等价于 Python 的 complex(real[, imag])。
//
// Python: complex(real=0, imag=0)
// DIFF: 不支持 complex("1+2j") 字符串形式，也不支持对象的 __complex__ 协议。
func Complex(args ...float64) (complex128, error) {
	switch len(args) {
	case 0:
		return 0, nil
	case 1:
		return complex(args[0], 0), nil
	case 2:
		return complex(args[0], args[1]), nil
	default:
		return 0, fmt.Errorf("conv.Complex: expected at most 2 arguments, got %d", len(args))
	}
}
