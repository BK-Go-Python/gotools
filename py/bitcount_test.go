package py

import (
	"math"
	"testing"
)

func TestBitCount(t *testing.T) {
	t.Run("signed", func(t *testing.T) {
		tests := []struct {
			name string
			x    int64
			want int
		}{
			{"zero", 0, 0},
			{"one", 1, 1},
			{"255", 255, 8},
			{"256", 256, 1},
			{"minus one", -1, 1},
			{"minus eight", -8, 1},
			{"MinInt64", math.MinInt64, 1},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := BitCount(tt.x); got != tt.want {
					t.Fatalf("BitCount(%d) = %d, want %d", tt.x, got, tt.want)
				}
			})
		}
	})

	t.Run("unsigned", func(t *testing.T) {
		tests := []struct {
			name string
			x    uint64
			want int
		}{
			{"zero", 0, 0},
			{"255", 255, 8},
			{"256", 256, 1},
			{"high bit", 1 << 63, 1},
			{"all bits", math.MaxUint64, 64},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := BitCount(tt.x); got != tt.want {
					t.Fatalf("BitCount(%d) = %d, want %d", tt.x, got, tt.want)
				}
			})
		}
	})

	t.Run("narrow types", func(t *testing.T) {
		if got := BitCount(uint(255)); got != 8 {
			t.Fatalf("BitCount(uint(255)) = %d, want 8", got)
		}
		if got := BitCount(uint8(255)); got != 8 {
			t.Fatalf("BitCount(uint8(255)) = %d, want 8", got)
		}
		if got := BitCount(int8(-128)); got != 1 {
			t.Fatalf("BitCount(int8(-128)) = %d, want 1", got)
		}
	})
}
