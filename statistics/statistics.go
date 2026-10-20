// Package statistics 提供 Python statistics 模块的常用统计函数。
//
// 数值函数接受任意数值切片（包括命名切片），返回 float64；空输入返回 error。
package statistics

import (
	"fmt"
	"math"
	"sort"
)

// Number 是统计函数接受的数值类型集合。
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Mean 等价于 Python 的 statistics.mean(data)。
func Mean[S ~[]E, E Number](data S) (float64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("statistics.Mean: empty data")
	}
	total := 0.0
	for _, value := range data {
		total += float64(value)
	}
	return total / float64(len(data)), nil
}

// Median 等价于 Python 的 statistics.median(data)。
func Median[S ~[]E, E Number](data S) (float64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("statistics.Median: empty data")
	}
	values := toFloats(data)
	sort.Float64s(values)
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle], nil
	}
	return (values[middle-1] + values[middle]) / 2, nil
}

// Mode 等价于 Python 的 statistics.mode(data)：返回出现次数最多的一个值。
//
// 并列时按 fmt.Sprint 排序取第一个，以保证确定性（Python 3.8+ 返回首个
// 遇到的值，但 Go map 没有插入顺序）。
func Mode[T comparable](data []T) (T, error) {
	modes, err := Multimode(data)
	if err != nil {
		var zero T
		return zero, err
	}
	return modes[0], nil
}

// Multimode 等价于 Python 的 statistics.multimode(data)：返回所有众数。
func Multimode[T comparable](data []T) ([]T, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("statistics.Multimode: empty data")
	}
	counts := make(map[T]int, len(data))
	max := 0
	for _, value := range data {
		counts[value]++
		if counts[value] > max {
			max = counts[value]
		}
	}
	modes := make([]T, 0)
	for value, count := range counts {
		if count == max {
			modes = append(modes, value)
		}
	}
	sort.Slice(modes, func(i, j int) bool {
		return fmt.Sprint(modes[i]) < fmt.Sprint(modes[j])
	})
	return modes, nil
}

// Variance 等价于 Python 的 statistics.variance(data)：样本方差（n-1）。
func Variance[S ~[]E, E Number](data S) (float64, error) {
	mean, err := Mean(data)
	if err != nil {
		return 0, fmt.Errorf("statistics.Variance: %w", err)
	}
	if len(data) < 2 {
		return 0, fmt.Errorf("statistics.Variance: need at least two data points")
	}
	sum := 0.0
	for _, value := range data {
		diff := float64(value) - mean
		sum += diff * diff
	}
	return sum / float64(len(data)-1), nil
}

// Pvariance 等价于 Python 的 statistics.pvariance(data)：总体方差（n）。
func Pvariance[S ~[]E, E Number](data S) (float64, error) {
	mean, err := Mean(data)
	if err != nil {
		return 0, fmt.Errorf("statistics.Pvariance: %w", err)
	}
	sum := 0.0
	for _, value := range data {
		diff := float64(value) - mean
		sum += diff * diff
	}
	return sum / float64(len(data)), nil
}

// Stdev 等价于 Python 的 statistics.stdev(data)：样本标准差。
func Stdev[S ~[]E, E Number](data S) (float64, error) {
	variance, err := Variance(data)
	if err != nil {
		return 0, fmt.Errorf("statistics.Stdev: %w", err)
	}
	return math.Sqrt(variance), nil
}

// Pstdev 等价于 Python 的 statistics.pstdev(data)：总体标准差。
func Pstdev[S ~[]E, E Number](data S) (float64, error) {
	variance, err := Pvariance(data)
	if err != nil {
		return 0, fmt.Errorf("statistics.Pstdev: %w", err)
	}
	return math.Sqrt(variance), nil
}

func toFloats[S ~[]E, E Number](data S) []float64 {
	values := make([]float64, len(data))
	for i, value := range data {
		values[i] = float64(value)
	}
	return values
}
