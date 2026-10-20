package set_test

import (
	"testing"

	"github.com/BK-Go-Python/gotools/set"
)

var benchmarkSetSink set.Set[int]

// benchmarkSet 返回从 offset 开始的 size 个连续整数构成的集合。
func benchmarkSet(size, offset int) set.Set[int] {
	items := make([]int, size)
	for i := range items {
		items[i] = offset + i
	}
	return set.New(items...)
}

func BenchmarkSetUnion(b *testing.B) {
	a := benchmarkSet(256, 0)
	c := benchmarkSet(256, 128)
	b.ReportAllocs()
	b.ResetTimer()

	var union set.Set[int]
	for i := 0; i < b.N; i++ {
		union = a.Union(c)
	}
	benchmarkSetSink = union
}

func BenchmarkSetIntersection(b *testing.B) {
	a := benchmarkSet(256, 0)
	c := benchmarkSet(256, 128)
	b.ReportAllocs()
	b.ResetTimer()

	var intersection set.Set[int]
	for i := 0; i < b.N; i++ {
		intersection = a.Intersection(c)
	}
	benchmarkSetSink = intersection
}
