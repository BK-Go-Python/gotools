package str_test

import (
	"fmt"

	"github.com/BK-Go-Python/gotools/str"
)

func ExampleFormat() {
	fmt.Println(str.Format("{} 有 {:.2f} 元", "Tom", 12.345))
	fmt.Println(str.Format("{0}-{0}", "a"))
	fmt.Println(str.Format("{name} 使用 {lang}", map[string]interface{}{"name": "Ann", "lang": "Go"}))
	// Output:
	// Tom 有 12.35 元
	// a-a
	// Ann 使用 Go
}

func ExampleSplit() {
	fmt.Println(str.Split("a,b,c", ","))
	fmt.Printf("%q\n", str.Split("  hello   world  ", ""))
	// Output:
	// [a b c]
	// ["hello" "world"]
}

func ExampleTitle() {
	fmt.Println(str.Title("hello world"))
	fmt.Println(str.Title("they're"))
	fmt.Println(str.Title("go 1.26"))
	// Output:
	// Hello World
	// They'Re
	// Go 1.26
}
