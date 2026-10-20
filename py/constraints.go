package py

// Number 是 Python 数值内置函数（abs、sum）接受的实数类型集合。
//
// 不包含 complex：Python 的 abs(complex) 返回 float64，与泛型返回 T 冲突；
// 复数请使用 math/cmplx.Abs。
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Integer 是 bin/hex/oct/divmod 接受的整数类型集合。
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
