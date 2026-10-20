package py_test

import (
	"strconv"
	"testing"

	"github.com/BK-Go-Python/gotools/py"
)

var (
	benchmarkSumSink  int
	benchmarkIntsSink []int
	benchmarkStrsSink []string
)

// benchmarkInts 返回 n 个降序整数，保证排序类基准每次都要真正重排。
func benchmarkInts(n int) []int {
	seq := make([]int, n)
	for i := range seq {
		seq[i] = n - i
	}
	return seq
}

func BenchmarkSum(b *testing.B) {
	seq := benchmarkInts(1024)
	b.ReportAllocs()
	b.ResetTimer()

	var total int
	for i := 0; i < b.N; i++ {
		total = py.Sum(seq)
	}
	benchmarkSumSink = total
}

func BenchmarkSortedOf(b *testing.B) {
	seq := benchmarkInts(1024)
	b.ReportAllocs()
	b.ResetTimer()

	var sorted []int
	for i := 0; i < b.N; i++ {
		sorted = py.SortedOf(seq)
	}
	benchmarkIntsSink = sorted
}

func BenchmarkMapOf(b *testing.B) {
	seq := benchmarkInts(1024)
	b.ReportAllocs()
	b.ResetTimer()

	var mapped []string
	for i := 0; i < b.N; i++ {
		mapped = py.MapOf(seq, strconv.Itoa)
	}
	benchmarkStrsSink = mapped
}
