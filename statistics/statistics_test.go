package statistics

import (
	"math"
	"testing"
)

func closeTo(a, b float64) bool {
	return math.Abs(a-b) < 1e-12
}

func TestMeanMedian(t *testing.T) {
	if got, err := Mean([]int{1, 2, 3}); err != nil || got != 2 {
		t.Fatalf("Mean = (%v,%v)", got, err)
	}
	if _, err := Mean([]int{}); err == nil {
		t.Fatal("Mean empty should error")
	}
	if got, err := Median([]int{3, 1, 2}); err != nil || got != 2 {
		t.Fatalf("Median odd = (%v,%v)", got, err)
	}
	if got, err := Median([]float64{4, 1, 3, 2}); err != nil || got != 2.5 {
		t.Fatalf("Median even = (%v,%v)", got, err)
	}
	if _, err := Median([]int{}); err == nil {
		t.Fatal("Median empty should error")
	}
}

func TestMode(t *testing.T) {
	if got, err := Mode([]int{1, 1, 2}); err != nil || got != 1 {
		t.Fatalf("Mode = (%v,%v)", got, err)
	}
	modes, err := Multimode([]int{3, 3, 1, 1, 2})
	if err != nil || len(modes) != 2 || modes[0] != 1 || modes[1] != 3 {
		t.Fatalf("Multimode = (%v,%v)", modes, err)
	}
	if _, err := Mode([]int{}); err == nil {
		t.Fatal("Mode empty should error")
	}
}

func TestVarianceStdev(t *testing.T) {
	if got, err := Variance([]int{1, 2, 3}); err != nil || !closeTo(got, 1) {
		t.Fatalf("Variance = (%v,%v)", got, err)
	}
	if got, err := Pvariance([]int{1, 2, 3}); err != nil || !closeTo(got, 2.0/3.0) {
		t.Fatalf("Pvariance = (%v,%v)", got, err)
	}
	if got, err := Stdev([]int{1, 2, 3}); err != nil || !closeTo(got, 1) {
		t.Fatalf("Stdev = (%v,%v)", got, err)
	}
	if got, err := Pstdev([]int{1, 2, 3}); err != nil || !closeTo(got, math.Sqrt(2.0/3.0)) {
		t.Fatalf("Pstdev = (%v,%v)", got, err)
	}
	if _, err := Variance([]int{1}); err == nil {
		t.Fatal("Variance single value should error")
	}
	if _, err := Pvariance([]int{}); err == nil {
		t.Fatal("Pvariance empty should error")
	}
}
