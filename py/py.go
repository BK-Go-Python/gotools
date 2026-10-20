// Package py 提供 Python 内置函数风格的辅助函数：range、zip、enumerate、
// map、filter、reduce、sorted、sum、min、max，以及真值判断、repr、type 等。
//
// 设计原则：
//   - 参数顺序尽量与 Python 保持一致，例如 py.Map(fn, iterable)；
//   - 旧的动态 API（Len/Range/Reduce 等）为兼容历史可能 panic；
//   - 新增的泛型、反射与 IO API 一律返回 error，不 panic；
//   - 结果确定：map 迭代按 key 排序，排序是稳定的。
package py

import (
	"fmt"
	"math"
	"math/big"

	"github.com/BK-Go-Python/gotools/internal/pyval"
	"github.com/BK-Go-Python/gotools/list"
)

// Pair 是 Enumerate 返回的 (index, value) 对。
type Pair struct {
	Index int
	Value interface{}
}

// Range 等价于 Python 的 range()，返回切片。
//
// 支持 1~3 个参数：
//
//	Range(5)        // [0 1 2 3 4]
//	Range(2, 5)     // [2 3 4]
//	Range(0, 10, 2) // [0 2 4 6 8]
//	Range(5, 0, -1) // [5 4 3 2 1]
func Range(args ...int) []int {
	var start, stop, step int
	switch len(args) {
	case 1:
		stop, step = args[0], 1
	case 2:
		start, stop, step = args[0], args[1], 1
	case 3:
		start, stop, step = args[0], args[1], args[2]
	default:
		panic("py.Range() expects 1 to 3 arguments")
	}
	if step == 0 {
		panic("py.Range() arg 3 must not be zero")
	}

	if step > 0 {
		if stop <= start {
			return []int{}
		}
		out := make([]int, 0, (stop-start+step-1)/step)
		for i := start; i < stop; i += step {
			out = append(out, i)
		}
		return out
	}
	if stop >= start {
		return []int{}
	}
	out := make([]int, 0, (start-stop-step-1)/(-step))
	for i := start; i > stop; i += step {
		out = append(out, i)
	}
	return out
}

// Len 等价于 Python 的 len(v)，非法类型会 panic。
func Len(v interface{}) int {
	n, err := pyval.Len(v)
	if err != nil {
		panic(err)
	}
	return n
}

// LenE 与 Len 相同，但返回 error。
func LenE(v interface{}) (int, error) {
	return pyval.Len(v)
}

// Truthy 等价于 Python 的 bool(v) 真值判断。
func Truthy(v interface{}) bool {
	return pyval.Truthy(v)
}

// Bool 等价于 Python 的 bool(v)。
func Bool(v interface{}) bool {
	return pyval.Truthy(v)
}

// Str 等价于 Python 的 str(v)。
func Str(v interface{}) string {
	return pyval.Str(v)
}

// Repr 等价于 Python 的 repr(v)，字符串会带引号，None/True/False 使用
// Python 写法。
func Repr(v interface{}) string {
	return pyval.Repr(v)
}

// Type 等价于 Python 的 type(v).__name__。
func Type(v interface{}) string {
	return pyval.TypeName(v)
}

// IsNone 判断 v 是否为 nil（包括值为 nil 的指针、接口、函数和 channel）。
func IsNone(v interface{}) bool {
	return pyval.IsNone(v)
}

// AbsInt 返回整数的绝对值。
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// AbsFloat 返回浮点数的绝对值。
func AbsFloat(x float64) float64 {
	return math.Abs(x)
}

// Round 等价于 Python 3 的 round()：使用银行家舍入（四舍六入五成双）。
//
//	Round(2.5)      // 2
//	Round(3.5)      // 4
//	Round(2.675, 2) // 2.67
func Round(x float64, ndigits ...int) float64 {
	n := 0
	if len(ndigits) > 0 {
		n = ndigits[0]
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}
	if n == 0 {
		return math.RoundToEven(x)
	}
	if n > 308 {
		return x
	}
	if n < -308 {
		return math.Copysign(0, x)
	}

	absN := n
	if absN < 0 {
		absN = -absN
	}
	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(absN)), nil)
	rat := new(big.Rat).SetFloat64(x)
	if rat == nil {
		return x
	}
	if n > 0 {
		rat.Mul(rat, new(big.Rat).SetInt(pow))
	} else {
		rat.Quo(rat, new(big.Rat).SetInt(pow))
	}

	quotient := roundHalfEven(rat)
	result := new(big.Rat).SetInt(quotient)
	if n > 0 {
		result.Quo(result, new(big.Rat).SetInt(pow))
	} else {
		result.Mul(result, new(big.Rat).SetInt(pow))
	}
	f, _ := result.Float64()
	return f
}

func roundHalfEven(r *big.Rat) *big.Int {
	num, den := r.Num(), r.Denom()
	quotient, remainder := new(big.Int).QuoRem(num, den, new(big.Int))
	twiceRemainder := new(big.Int).Abs(remainder)
	twiceRemainder.Lsh(twiceRemainder, 1)
	cmp := twiceRemainder.Cmp(den)
	if cmp > 0 || (cmp == 0 && new(big.Int).Abs(quotient).Bit(0) == 1) {
		if num.Sign() >= 0 {
			quotient.Add(quotient, big.NewInt(1))
		} else {
			quotient.Sub(quotient, big.NewInt(1))
		}
	}
	return quotient
}

// FloorDiv 等价于 Python 的 a // b，向负无穷取整。
func FloorDiv(a, b int) int {
	if b == 0 {
		panic("py.FloorDiv(): integer division or modulo by zero")
	}
	q := a / b
	if r := a % b; r != 0 && (r < 0) != (b < 0) {
		q--
	}
	return q
}

// Mod 等价于 Python 的 a % b，结果符号与除数一致。
func Mod(a, b int) int {
	if b == 0 {
		panic("py.Mod(): integer division or modulo by zero")
	}
	r := a % b
	if r != 0 && (r < 0) != (b < 0) {
		r += b
	}
	return r
}

// DivMod 等价于 Python 的 divmod(a, b)。
func DivMod(a, b int) (int, int) {
	return FloorDiv(a, b), Mod(a, b)
}

// Pow 等价于 Python 的 pow(x, y)，底层使用 math.Pow。
func Pow(x, y float64) float64 {
	return math.Pow(x, y)
}

// SumAny 是动态版 sum(iterable)：接受任意可迭代对象，以 float64 返回。
//
// 新代码优先使用泛型 Sum，后者能保持元素类型；只有在处理 []any、map
// 等异构/非切片可迭代对象时才需要 SumAny。
//
// Python: sum(iterable)
func SumAny(iterable interface{}) (float64, error) {
	items := mustIter(iterable)
	total := 0.0
	for _, item := range items {
		f, ok := numericFloat(item)
		if !ok {
			return 0, fmt.Errorf("py.SumAny(): unsupported operand type %s", pyval.TypeName(item))
		}
		total += f
	}
	return total, nil
}

// Min 等价于 Python 的 min(iterable[, key])。
func Min(iterable interface{}, key ...func(interface{}) interface{}) (interface{}, error) {
	return extreme(iterable, key, true)
}

// Max 等价于 Python 的 max(iterable[, key])。
func Max(iterable interface{}, key ...func(interface{}) interface{}) (interface{}, error) {
	return extreme(iterable, key, false)
}

func extreme(iterable interface{}, key []func(interface{}) interface{}, wantMin bool) (interface{}, error) {
	items := mustIter(iterable)
	if len(items) == 0 {
		if wantMin {
			return nil, fmt.Errorf("py.Min(): arg is an empty sequence")
		}
		return nil, fmt.Errorf("py.Max(): arg is an empty sequence")
	}
	best := items[0]
	bestKey := best
	if len(key) > 0 && key[0] != nil {
		bestKey = key[0](best)
	}
	for _, item := range items[1:] {
		itemKey := item
		if len(key) > 0 && key[0] != nil {
			itemKey = key[0](item)
		}
		c, err := pyval.Compare(itemKey, bestKey)
		if err != nil {
			return nil, err
		}
		if (wantMin && c < 0) || (!wantMin && c > 0) {
			best, bestKey = item, itemKey
		}
	}
	return best, nil
}

// Any 等价于 Python 的 any(iterable)。
func Any(iterable interface{}) bool {
	for _, item := range mustIter(iterable) {
		if pyval.Truthy(item) {
			return true
		}
	}
	return false
}

// All 等价于 Python 的 all(iterable)。
func All(iterable interface{}) bool {
	for _, item := range mustIter(iterable) {
		if !pyval.Truthy(item) {
			return false
		}
	}
	return true
}

// Enumerate 等价于 Python 的 enumerate(iterable, start=0)。
func Enumerate(iterable interface{}, start ...int) []Pair {
	offset := 0
	if len(start) > 0 {
		offset = start[0]
	}
	items := mustIter(iterable)
	pairs := make([]Pair, len(items))
	for i, item := range items {
		pairs[i] = Pair{Index: offset + i, Value: item}
	}
	return pairs
}

// Zip 等价于 Python 的 zip(*iterables)，结果长度取最短的可迭代对象。
func Zip(iterables ...interface{}) [][]interface{} {
	if len(iterables) == 0 {
		return [][]interface{}{}
	}
	lists := make([][]interface{}, len(iterables))
	length := -1
	for i, iterable := range iterables {
		lists[i] = mustIter(iterable)
		if length < 0 || len(lists[i]) < length {
			length = len(lists[i])
		}
	}
	rows := make([][]interface{}, length)
	for i := 0; i < length; i++ {
		row := make([]interface{}, len(lists))
		for j, items := range lists {
			row[j] = items[i]
		}
		rows[i] = row
	}
	return rows
}

// Map 等价于 Python 的 map(fn, iterable)，返回 List。
func Map(fn func(interface{}) interface{}, iterable interface{}) list.List {
	items := mustIter(iterable)
	out := make(list.List, len(items))
	for i, item := range items {
		out[i] = fn(item)
	}
	return out
}

// Filter 等价于 Python 的 filter(fn, iterable)，返回 List。
func Filter(fn func(interface{}) bool, iterable interface{}) list.List {
	items := mustIter(iterable)
	out := make(list.List, 0, len(items))
	for _, item := range items {
		if fn(item) {
			out = append(out, item)
		}
	}
	return out
}

// Reduce 等价于 Python 的 functools.reduce(fn, iterable[, initial])。
//
// 没有 initial 且 iterable 为空时 panic。
func Reduce(fn func(acc, item interface{}) interface{}, iterable interface{}, initial ...interface{}) interface{} {
	items := mustIter(iterable)
	var acc interface{}
	start := 0
	if len(initial) > 0 {
		acc = initial[0]
	} else {
		if len(items) == 0 {
			panic("py.Reduce() of empty sequence with no initial value")
		}
		acc = items[0]
		start = 1
	}
	for _, item := range items[start:] {
		acc = fn(acc, item)
	}
	return acc
}

// Sorted 等价于 Python 的 sorted(iterable)，返回新的 List。
//
// 可选的 less 参数用于自定义比较，例如按字符串长度排序。
func Sorted(iterable interface{}, less ...func(a, b interface{}) bool) list.List {
	out := list.New(mustIter(iterable)...)
	if len(less) > 0 && less[0] != nil {
		out.Sort(less[0])
	} else {
		out.Sort()
	}
	return out
}

// SortedBy 等价于 Python 的 sorted(iterable, key=..., reverse=...)。
func SortedBy(iterable interface{}, key func(interface{}) interface{}, reverse ...bool) list.List {
	out := list.New(mustIter(iterable)...)
	out.SortBy(key, reverse...)
	return out
}

// Reversed 等价于 Python 的 reversed(iterable)，返回新的 List。
func Reversed(iterable interface{}) list.List {
	items := mustIter(iterable)
	out := make(list.List, len(items))
	for i, item := range items {
		out[len(items)-1-i] = item
	}
	return out
}

func mustIter(v interface{}) []interface{} {
	items, err := pyval.Iter(v)
	if err != nil {
		panic(err)
	}
	return items
}

func numericFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	case int:
		return float64(x), true
	case int8:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint8:
		return float64(x), true
	case uint16:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	case uintptr:
		return float64(x), true
	case float32:
		return float64(x), true
	case float64:
		return x, true
	}
	return 0, false
}
