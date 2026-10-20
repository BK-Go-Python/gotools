package py_test

import (
	"fmt"
	"slices"

	"github.com/BK-Go-Python/gotools/py"
)

func ExampleRange() {
	fmt.Println(py.Range(5))
	fmt.Println(py.Range(2, 5))
	fmt.Println(py.Range(0, 10, 2))
	fmt.Println(py.Range(5, 0, -1))
	// Output:
	// [0 1 2 3 4]
	// [2 3 4]
	// [0 2 4 6 8]
	// [5 4 3 2 1]
}

func ExampleSum() {
	fmt.Println(py.Sum([]int{1, 2, 3, 4}))
	fmt.Println(py.Sum([]float64{1.5, 2.5}))
	// Output:
	// 10
	// 4
}

func ExampleSortedOf() {
	nums := []int{3, 1, 2}
	sorted := py.SortedOf(nums)
	fmt.Println(sorted)
	fmt.Println(nums)
	// Output:
	// [1 2 3]
	// [3 1 2]
}

func ExampleIter() {
	it := py.Iter(slices.Values([]string{"a", "b", "c"}))
	defer it.Stop()

	for {
		value, ok := it.Next()
		if !ok {
			break
		}
		fmt.Println(value)
	}
	// Output:
	// a
	// b
	// c
}
