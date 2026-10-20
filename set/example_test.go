package set_test

import (
	"fmt"
	"slices"

	"github.com/BK-Go-Python/gotools/set"
)

func ExampleNew() {
	a := set.New(1, 2, 3)
	b := set.New(3, 4)

	union := a.Union(b).Items()
	slices.Sort(union)
	fmt.Println(union)

	intersection := a.Intersection(b).Items()
	slices.Sort(intersection)
	fmt.Println(intersection)

	fmt.Println(a.Has(2), b.Has(2))
	// Output:
	// [1 2 3 4]
	// [3]
	// true false
}
