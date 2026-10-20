package py

import (
	"fmt"
	"iter"
)

// Iter 把 Go 的推式迭代器包装成可按 next 拉取的 Iterator。
//
// Python: iter(iterable)
// DIFF: Go 没有统一的 __iter__ 协议，需先转换为 iter.Seq。
func Iter[T any](seq iter.Seq[T]) *Iterator[T] {
	if seq == nil {
		return &Iterator[T]{}
	}
	next, stop := iter.Pull(seq)
	return &Iterator[T]{next: next, stop: stop}
}

// IterOf 把切片包装成 Iterator，不复制底层切片。
//
// Python: iter(list)
func IterOf[T any](items []T) *Iterator[T] {
	index := 0
	return &Iterator[T]{
		next: func() (T, bool) {
			if index >= len(items) {
				var zero T
				return zero, false
			}
			value := items[index]
			index++
			return value, true
		},
	}
}

// Next 等价于 Python 的 next(iterator)，耗尽时返回零值和 false。
//
// Python: next(iterator)
// DIFF: Python 抛 StopIteration，Go 用 ok=false 表达。
func Next[T any](it *Iterator[T]) (T, bool) {
	if it == nil || it.next == nil {
		var zero T
		return zero, false
	}
	return it.next()
}

// Next 返回下一个元素；迭代器为 nil 或已耗尽时返回零值和 false。
func (it *Iterator[T]) Next() (T, bool) {
	return Next(it)
}

// Stop 释放迭代器资源，可重复调用。
func (it *Iterator[T]) Stop() {
	if it == nil {
		return
	}
	if it.stop != nil {
		it.stop()
		it.stop = nil
	}
	it.next = nil
}

// RangeSeq 是 Python range() 的惰性版本，返回 iter.Seq[int]。
//
// Python: range(start, stop, step)
// DIFF: Python 的 range 支持任意精度整数；这里使用 int。
func RangeSeq(start, stop, step int) (iter.Seq[int], error) {
	if step == 0 {
		return nil, fmt.Errorf("py.RangeSeq: step must not be zero: %w", ErrZeroDivision)
	}
	return func(yield func(int) bool) {
		limit := int(^uint(0) >> 1)
		if step > 0 {
			for i := start; i < stop; {
				if !yield(i) {
					return
				}
				// 防止 i+step 溢出后回绕成负数，导致无限循环。
				if i > limit-step {
					return
				}
				i += step
			}
			return
		}
		minimum := -limit - 1
		for i := start; i > stop; {
			if !yield(i) {
				return
			}
			if i < minimum-step {
				return
			}
			i += step
		}
	}, nil
}

// EnumerateOf 等价于 Python 的 enumerate(seq, start)。
//
// Python: enumerate(seq, start=0)
// DIFF: 返回切片而不是惰性迭代器。
func EnumerateOf[T any](seq []T, start ...int) []PairOf[T] {
	offset := 0
	if len(start) > 0 {
		offset = start[0]
	}
	out := make([]PairOf[T], len(seq))
	for i, value := range seq {
		out[i] = PairOf[T]{Index: offset + i, Value: value}
	}
	return out
}

// MapOf 等价于 Python 的 map(fn, seq)，返回新切片。
//
// Python: map(fn, seq)
// DIFF: 返回切片而不是惰性迭代器。
func MapOf[T, R any](seq []T, fn func(T) R) []R {
	out := make([]R, len(seq))
	for i, value := range seq {
		out[i] = fn(value)
	}
	return out
}

// FilterOf 等价于 Python 的 filter(fn, seq)，返回新切片。
//
// Python: filter(fn, seq)
// DIFF: 返回切片而不是惰性迭代器。
func FilterOf[T any](seq []T, keep func(T) bool) []T {
	out := make([]T, 0, len(seq))
	for _, value := range seq {
		if keep(value) {
			out = append(out, value)
		}
	}
	return out
}

// ReversedOf 等价于 Python 的 reversed(seq)，返回新切片且不修改入参。
//
// Python: reversed(seq)
// DIFF: 返回切片而不是反向迭代器。
func ReversedOf[T any](seq []T) []T {
	out := make([]T, len(seq))
	for i, value := range seq {
		out[len(seq)-1-i] = value
	}
	return out
}

// Zip2 等价于 Python 的 zip(a, b)，长度取两者较短者。
//
// Python: zip(a, b)
func Zip2[A, B any](a []A, b []B) []Pair2[A, B] {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	out := make([]Pair2[A, B], n)
	for i := 0; i < n; i++ {
		out[i] = Pair2[A, B]{First: a[i], Second: b[i]}
	}
	return out
}

// Zip3 等价于 Python 的 zip(a, b, c)，长度取三者最短者。
//
// Python: zip(a, b, c)
func Zip3[A, B, C any](a []A, b []B, c []C) []Pair3[A, B, C] {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if len(c) < n {
		n = len(c)
	}
	out := make([]Pair3[A, B, C], n)
	for i := 0; i < n; i++ {
		out[i] = Pair3[A, B, C]{First: a[i], Second: b[i], Third: c[i]}
	}
	return out
}
