// Package random 提供 Python random 模块风格的高频随机函数。
//
// 使用 math/rand，并由互斥锁保护全局生成器；Seed 可重置全局序列。
package random

import (
	"fmt"
	mathrand "math/rand"
	"sync"
	"time"
)

var (
	mu  sync.Mutex
	rng = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
)

// Seed 等价于 Python 的 random.seed(seed)，重置全局随机序列。
func Seed(seed int64) {
	mu.Lock()
	defer mu.Unlock()
	rng = mathrand.New(mathrand.NewSource(seed))
}

// Float 等价于 Python 的 random.random()，返回 [0, 1) 的浮点数。
func Float() float64 {
	mu.Lock()
	defer mu.Unlock()
	return rng.Float64()
}

// Uniform 等价于 Python 的 random.uniform(a, b)。
func Uniform(a, b float64) float64 {
	if a > b {
		a, b = b, a
	}
	return a + (b-a)*Float()
}

// Randrange 等价于 Python 的 random.randrange(start, stop[, step])。
func Randrange(start, stop int, step ...int) (int, error) {
	delta := 1
	if len(step) > 0 {
		delta = step[0]
	}
	if delta == 0 {
		return 0, fmt.Errorf("random.Randrange: step must not be zero")
	}

	var count int
	if delta > 0 {
		if start >= stop {
			return 0, fmt.Errorf("random.Randrange: empty range")
		}
		count = (stop - start + delta - 1) / delta
	} else {
		if start <= stop {
			return 0, fmt.Errorf("random.Randrange: empty range")
		}
		count = (start - stop - delta - 1) / -delta
	}
	if count <= 0 {
		return 0, fmt.Errorf("random.Randrange: empty range")
	}

	mu.Lock()
	index := rng.Intn(count)
	mu.Unlock()
	return start + index*delta, nil
}

// Randint 等价于 Python 的 random.randint(a, b)，包含两端。
func Randint(a, b int) int {
	if a > b {
		a, b = b, a
	}
	value, err := Randrange(a, b+1)
	if err != nil {
		return a
	}
	return value
}

// Choice 等价于 Python 的 random.choice(seq)；空切片返回零值和 false。
func Choice[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	mu.Lock()
	index := rng.Intn(len(items))
	mu.Unlock()
	return items[index], true
}

// Shuffle 等价于 Python 的 random.shuffle(seq)，原地打乱。
func Shuffle[T any](items []T) {
	mu.Lock()
	defer mu.Unlock()
	rng.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
}

// Sample 等价于 Python 的 random.sample(seq, k)：无放回抽样。
func Sample[T any](items []T, k int) ([]T, error) {
	if k < 0 || k > len(items) {
		return nil, fmt.Errorf("random.Sample: sample larger than population or negative: %d", k)
	}
	out := append([]T(nil), items...)
	mu.Lock()
	for i := 0; i < k; i++ {
		j := i + rng.Intn(len(out)-i)
		out[i], out[j] = out[j], out[i]
	}
	mu.Unlock()
	return out[:k], nil
}

// Choices 等价于 Python 的 random.choices(seq, k=k)：有放回抽样。
func Choices[T any](items []T, k int) ([]T, error) {
	if len(items) == 0 && k > 0 {
		return nil, fmt.Errorf("random.Choices: empty population")
	}
	if k < 0 {
		return nil, fmt.Errorf("random.Choices: negative count: %d", k)
	}
	out := make([]T, k)
	mu.Lock()
	for i := 0; i < k; i++ {
		out[i] = items[rng.Intn(len(items))]
	}
	mu.Unlock()
	return out, nil
}
