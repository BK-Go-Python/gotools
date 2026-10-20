package py

import (
	"errors"
	"math"
	"testing"
)

func TestAbs(t *testing.T) {
	if got := Abs(-5); got != 5 {
		t.Fatalf("Abs(-5) = %d, want 5", got)
	}
	if got := Abs(0); got != 0 {
		t.Fatalf("Abs(0) = %d, want 0", got)
	}
	if got := Abs(-1.5); got != 1.5 {
		t.Fatalf("Abs(-1.5) = %v, want 1.5", got)
	}
}

func TestRadix(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"bin positive", Bin(5), "0b101"},
		{"bin negative", Bin(-5), "-0b101"},
		{"bin zero", Bin(0), "0b0"},
		{"hex positive", Hex(255), "0xff"},
		{"hex negative", Hex(-255), "-0xff"},
		{"oct positive", Oct(8), "0o10"},
		{"oct negative", Oct(-8), "-0o10"},
		{"unsigned", Bin[uint8](5), "0b101"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestDivModOf(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		wantQ   int
		wantR   int
		wantErr error
	}{
		{"positive", 7, 3, 2, 1, nil},
		{"negative dividend", -7, 3, -3, 2, nil},
		{"negative divisor", 7, -3, -3, -2, nil},
		{"both negative", -7, -3, 2, -1, nil},
		{"zero divisor", 7, 0, 0, 0, ErrZeroDivision},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, r, err := DivModOf(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && (q != tt.wantQ || r != tt.wantR) {
				t.Fatalf("DivModOf(%d,%d) = (%d,%d), want (%d,%d)", tt.a, tt.b, q, r, tt.wantQ, tt.wantR)
			}
		})
	}
}

func TestPowMod(t *testing.T) {
	tests := []struct {
		name          string
		base, exp, md int64
		want          int64
		wantErr       error
	}{
		{"basic", 2, 10, 1000, 24, nil},
		{"negative exponent inverse", 3, -1, 11, 4, nil},
		{"negative base", -2, 3, 5, 2, nil},
		{"zero exponent", 2, 0, 5, 1, nil},
		{"zero to zero", 0, 0, 5, 1, nil},
		{"non invertible", 2, -1, 4, 0, ErrType},
		{"zero modulus", 2, 3, 0, 0, ErrOutOfRange},
		{"min exponent", 2, math.MinInt64, 5, 0, ErrOutOfRange},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PowMod(tt.base, tt.exp, tt.md)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("PowMod(%d,%d,%d) = %d, want %d", tt.base, tt.exp, tt.md, got, tt.want)
			}
		})
	}
}

func TestRoundExisting(t *testing.T) {
	tests := []struct {
		x    float64
		n    []int
		want float64
	}{
		{2.5, nil, 2},
		{3.5, nil, 4},
		{-2.5, nil, -2},
		{2.675, []int{2}, 2.67},
	}
	for _, tt := range tests {
		if got := Round(tt.x, tt.n...); got != tt.want {
			t.Fatalf("Round(%v, %v) = %v, want %v", tt.x, tt.n, got, tt.want)
		}
	}
}
