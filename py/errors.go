package py

import "errors"

// 本文件定义 py 包共享的哨兵错误。新增 API 一律通过 error 返回，
// 不使用 panic；可用 errors.Is 做判断。
var (
	// ErrEmpty 表示对空序列求最小值/最大值。
	ErrEmpty = errors.New("py: empty sequence")

	// ErrZeroDivision 表示整数除零或模零。
	ErrZeroDivision = errors.New("py: division by zero")

	// ErrNotFound 表示属性、键或元素不存在。
	ErrNotFound = errors.New("py: not found")

	// ErrType 表示类型不匹配或值不可哈希/不可转换。
	ErrType = errors.New("py: invalid type")

	// ErrOutOfRange 表示下标、码点或参数越界。
	ErrOutOfRange = errors.New("py: out of range")

	// ErrNotSupported 预留给部分 API 的未支持模式；不为不可模拟函数提供空壳。
	ErrNotSupported = errors.New("py: operation not supported")
)
