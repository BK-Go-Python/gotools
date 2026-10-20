package py

// PairOf 是 Enumerate 返回的 (index, value) 对。
type PairOf[T any] struct {
	Index int
	Value T
}

// Pair2 是 Zip2 返回的二元组。
type Pair2[A, B any] struct {
	First  A
	Second B
}

// Pair3 是 Zip3 返回的三元组。
type Pair3[A, B, C any] struct {
	First  A
	Second B
	Third  C
}

// Iterator 把 Go 1.23+ 的推式迭代器 iter.Seq 包装成可 next 的拉式迭代器。
//
// 使用完必须调用 Stop 释放底层资源；Stop 可以重复调用。
type Iterator[T any] struct {
	next func() (T, bool)
	stop func()
}

// Slice 对应 Python 的 slice(start, stop, step) 对象。
//
// HasStart/HasStop/HasStep 用于区分“未提供”和“提供了 0”，
// 因为 Python 的 None 与 0 在切片语义中不同。
type Slice struct {
	Start    int
	Stop     int
	Step     int
	HasStart bool
	HasStop  bool
	HasStep  bool
}
