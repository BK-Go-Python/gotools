package py

import (
	"math"
	"testing"
)

func TestBitLength(t *testing.T) {
	tests := []struct {
		name string
		got  int
		want int
	}{
		{"zero", BitLength(0), 0},
		{"one", BitLength(1), 1},
		{"eight", BitLength(8), 4},
		{"negative eight", BitLength(-8), 4},
		{"negative one", BitLength(-1), 1},
		{"min int64", BitLength(math.MinInt64), 64},
		{"unsigned", BitLength[uint](255), 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %d, want %d", tt.got, tt.want)
			}
		})
	}
}
