package conv

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestFloatHex(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want string
	}{
		{"one", 1, "0x1.0000000000000p+0"},
		{"half", 0.5, "0x1.0000000000000p-1"},
		{"tenth", 0.1, "0x1.999999999999ap-4"},
		{"minSubnormal", 5e-324, "0x0.0000000000001p-1022"},
		{"zero", 0, "0x0.0p+0"},
		{"negativeZero", math.Copysign(0, -1), "-0x0.0p+0"},
		{"positiveInf", math.Inf(1), "inf"},
		{"negativeInf", math.Inf(-1), "-inf"},
		{"nan", math.NaN(), "nan"},
		{"negativeNan", math.Copysign(math.NaN(), -1), "nan"},
		{"oneAndHalf", 1.5, "0x1.8000000000000p+0"},
		{"negativeOneAndHalf", -1.5, "-0x1.8000000000000p+0"},
		{"two", 2, "0x1.0000000000000p+1"},
		{"third", 1.0 / 3.0, "0x1.5555555555555p-2"},
		{"pi", math.Pi, "0x1.921fb54442d18p+1"},
		{"decimal", 123.456, "0x1.edd2f1a9fbe77p+6"},
		{"large", 1e100, "0x1.249ad2594c37dp+332"},
		{"small", 1e-100, "0x1.bff2ee48e0530p-333"},
		{"maxFloat", math.MaxFloat64, "0x1.fffffffffffffp+1023"},
		{"minNormal", 2.2250738585072014e-308, "0x1.0000000000000p-1022"},
		{"negativeMinSubnormal", -5e-324, "-0x0.0000000000001p-1022"},
		{"twoMinSubnormals", 1e-323, "0x0.0000000000002p-1022"},
		{"fortyTwo", 42, "0x1.5000000000000p+5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FloatHex(tt.f); got != tt.want {
				t.Errorf("FloatHex(%v) = %q, want %q", tt.f, got, tt.want)
			}
		})
	}
}

func TestFromFloatHex(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    float64
		wantErr bool
	}{
		{"pythonOne", "0x1.0000000000000p+0", 1, false},
		{"pythonHalf", "0x1.0000000000000p-1", 0.5, false},
		{"pythonTenth", "0x1.999999999999ap-4", 0.1, false},
		{"pythonMinSubnormal", "0x0.0000000000001p-1022", 5e-324, false},
		{"pythonZero", "0x0.0p+0", 0, false},
		{"pythonNegativeZero", "-0x0.0p+0", math.Copysign(0, -1), false},
		{"pythonMaxFloat", "0x1.fffffffffffffp+1023", math.MaxFloat64, false},

		{"goWithoutFraction", "0x1p-2", 0.25, false},
		{"goUpper", "0X1.8P+1", 3, false},
		{"goLeadingPlus", "+0x1p0", 1, false},
		{"goNegative", "-0x1.8p0", -1.5, false},
		{"goLeadingZeroExponent", "0x1p-0002", 0.25, false},
		{"goMinSubnormal", "0x1p-1074", 5e-324, false},

		{"pythonNoPrefix", "1.8p+1", 3, false},
		{"pythonNoPrefixLeadingDot", ".8p0", 0.5, false},
		{"pythonNoPrefixTrailingDot", "1.p0", 1, false},
		{"pythonNoExponent", "0x1.8", 1.5, false},
		{"pythonNoExponentLeadingDot", "0x.8", 0.5, false},
		{"surroundingSpace", "  -0x1.8p+2\t", -6, false},
		{"decimal", "2.5", 2.5, false},

		{"inf", "inf", math.Inf(1), false},
		{"plusInf", "+Inf", math.Inf(1), false},
		{"negativeInf", "-inf", math.Inf(-1), false},
		{"infinity", "Infinity", math.Inf(1), false},
		{"negativeInfinity", "-Infinity", math.Inf(-1), false},
		{"nan", "nan", math.NaN(), false},
		{"upperNaN", "NaN", math.NaN(), false},

		{"underflow", "0x1p-1075", 0, false},
		{"empty", "", 0, true},
		{"blank", "   ", 0, true},
		{"barePrefix", "0x", 0, true},
		{"missingExponentDigits", "0x1.8p", 0, true},
		{"doubleDot", "0x1.2.3p0", 0, true},
		{"garbage", "hello", 0, true},
		{"doubleExponent", "0x1p1p1", 0, true},
		{"signedNaN", "-nan", 0, true},
		{"overflow", "0x1p+1024", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromFloatHex(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("FromFloatHex(%q) = %v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("FromFloatHex(%q): %v", tt.in, err)
			}
			if math.IsNaN(tt.want) {
				if !math.IsNaN(got) {
					t.Fatalf("FromFloatHex(%q) = %v, want NaN", tt.in, got)
				}
				return
			}
			if math.Float64bits(got) != math.Float64bits(tt.want) {
				t.Errorf("FromFloatHex(%q) = %v (bits %#016x), want %v (bits %#016x)",
					tt.in, got, math.Float64bits(got), tt.want, math.Float64bits(tt.want))
			}
		})
	}
}

func TestFloatHexRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		f    float64
	}{
		{"zero", 0},
		{"negativeZero", math.Copysign(0, -1)},
		{"one", 1},
		{"tenth", 0.1},
		{"third", 1.0 / 3.0},
		{"minSubnormal", 5e-324},
		{"minNormal", 2.2250738585072014e-308},
		{"maxFloat", math.MaxFloat64},
		{"negative", -123.456},
		{"positiveInf", math.Inf(1)},
		{"negativeInf", math.Inf(-1)},
		{"nan", math.NaN()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FloatHex(tt.f)
			got, err := FromFloatHex(s)
			if err != nil {
				t.Fatalf("FromFloatHex(%q): %v", s, err)
			}
			if math.IsNaN(tt.f) {
				if !math.IsNaN(got) {
					t.Fatalf("FromFloatHex(%q) = %v, want NaN", s, got)
				}
				return
			}
			if math.Float64bits(got) != math.Float64bits(tt.f) {
				t.Fatalf("FromFloatHex(%q) = %v (bits %#016x), want %v (bits %#016x)",
					s, got, math.Float64bits(got), tt.f, math.Float64bits(tt.f))
			}
		})
	}
}

func TestFloatHexAllExponents(t *testing.T) {
	fractions := []uint64{0, 1, 0x8000000000000, 0x8000000000001, 0xfffffffffffff, 0x5555555555555}
	for exponent := uint64(0); exponent <= 0x7fe; exponent++ {
		for _, fraction := range fractions {
			for _, sign := range []uint64{0, 1 << 63} {
				bits := sign | exponent<<52 | fraction
				f := math.Float64frombits(bits)
				s := FloatHex(f)
				got, err := FromFloatHex(s)
				if err != nil {
					t.Fatalf("FromFloatHex(FloatHex(%#016x) = %q): %v", bits, s, err)
				}
				if math.Float64bits(got) != bits {
					t.Fatalf("round trip %#016x via %q = %#016x", bits, s, math.Float64bits(got))
				}
			}
		}
	}
}

func TestFloatAsIntegerRatio(t *testing.T) {
	// Python 3.13: (5e-324).as_integer_ratio() == (1, 2**1074).
	const minSubnormalDen = "202402253307310618352495346718917307049556649764142118356901358027430339567995346891960383701437124495187077864316811911389808737385793476867013399940738509921517424276566361364466907742093216341239767678472745068562007483424692698618103355649159556340810056512358769552333414615230502532186327508646006263307707741093494784"

	tests := []struct {
		name    string
		f       float64
		wantNum string
		wantDen string
	}{
		{"half", 0.5, "1", "2"},
		{"tenth", 0.1, "3602879701896397", "36028797018963968"},
		{"one", 1, "1", "1"},
		{"negativeTwoAndHalf", -2.5, "-5", "2"},
		{"twoAndHalf", 2.5, "5", "2"},
		{"zero", 0, "0", "1"},
		{"negativeZero", math.Copysign(0, -1), "0", "1"},
		{"third", 1.0 / 3.0, "6004799503160661", "18014398509481984"},
		{"pi", math.Pi, "884279719003555", "281474976710656"},
		{"minSubnormal", 5e-324, "1", minSubnormalDen},
		{"powerOfTwo", math.Ldexp(1, 60), "1152921504606846976", "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, den, err := FloatAsIntegerRatio(tt.f)
			if err != nil {
				t.Fatalf("FloatAsIntegerRatio(%v): %v", tt.f, err)
			}
			wantNum, ok := new(big.Int).SetString(tt.wantNum, 10)
			if !ok {
				t.Fatalf("bad test numerator %q", tt.wantNum)
			}
			wantDen, ok := new(big.Int).SetString(tt.wantDen, 10)
			if !ok {
				t.Fatalf("bad test denominator %q", tt.wantDen)
			}
			if num.Cmp(wantNum) != 0 || den.Cmp(wantDen) != 0 {
				t.Errorf("FloatAsIntegerRatio(%v) = (%s, %s), want (%s, %s)",
					tt.f, num, den, tt.wantNum, tt.wantDen)
			}
		})
	}
}

func TestFloatAsIntegerRatioSpecial(t *testing.T) {
	tests := []struct {
		name string
		f    float64
	}{
		{"nan", math.NaN()},
		{"positiveInf", math.Inf(1)},
		{"negativeInf", math.Inf(-1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, den, err := FloatAsIntegerRatio(tt.f)
			if err == nil {
				t.Fatalf("FloatAsIntegerRatio(%v) = (%v, %v), want error", tt.f, num, den)
			}
			if num != nil || den != nil {
				t.Fatalf("FloatAsIntegerRatio(%v) returned (%v, %v) with error", tt.f, num, den)
			}
		})
	}
}

func TestFloatAsIntegerRatioExact(t *testing.T) {
	one := big.NewInt(1)
	values := []float64{
		0, math.Copysign(0, -1), 0.1, 0.5, 1, -1, 1.0 / 3.0, math.Pi,
		5e-324, 1e-323, 2.2250738585072014e-308, math.MaxFloat64,
		-math.MaxFloat64, 1e100, 1e-100, 42, 123.456,
	}
	for i, f := range values {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			num, den, err := FloatAsIntegerRatio(f)
			if err != nil {
				t.Fatalf("FloatAsIntegerRatio(%v): %v", f, err)
			}
			if den.Sign() <= 0 {
				t.Fatalf("FloatAsIntegerRatio(%v) denominator = %s, want positive", f, den)
			}
			gcd := new(big.Int).GCD(nil, nil, new(big.Int).Abs(num), den)
			if gcd.Cmp(one) != 0 {
				t.Fatalf("FloatAsIntegerRatio(%v) = %s/%s, gcd = %s", f, num, den, gcd)
			}
			if f == 0 {
				if num.Sign() != 0 || den.Cmp(one) != 0 {
					t.Fatalf("FloatAsIntegerRatio(%v) = (%s, %s), want (0, 1)", f, num, den)
				}
				return
			}
			if (num.Sign() < 0) != math.Signbit(f) {
				t.Fatalf("FloatAsIntegerRatio(%v) sign = %s", f, num)
			}
			got, exact := new(big.Rat).SetFrac(num, den).Float64()
			if !exact || math.Float64bits(got) != math.Float64bits(f) {
				t.Fatalf("ratio %s/%s -> %v (exact=%v), want %v", num, den, got, exact, f)
			}
		})
	}
}
