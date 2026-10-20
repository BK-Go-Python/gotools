package conv

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestIntToBytes(t *testing.T) {
	minI64Big := []byte{0x80, 0, 0, 0, 0, 0, 0, 0}
	minI64Little := []byte{0, 0, 0, 0, 0, 0, 0, 0x80}
	maxI64Big := []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	maxI64Little := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}

	tests := []struct {
		name    string
		x       int64
		length  int
		order   string
		signed  []bool
		want    []byte
		wantErr bool
	}{
		// 期望值来自 Python 3.13: x.to_bytes(length, order, signed=signed)
		{name: "1024 big", x: 1024, length: 2, order: "big", want: []byte{0x04, 0x00}},
		{name: "1024 little", x: 1024, length: 2, order: "little", want: []byte{0x00, 0x04}},
		{name: "1024 big uppercase", x: 1024, length: 2, order: "BIG", want: []byte{0x04, 0x00}},
		{name: "1024 mixed case", x: 1024, length: 2, order: "LiTtLe", want: []byte{0x00, 0x04}},
		{name: "-1 signed big", x: -1, length: 2, order: "big", signed: []bool{true}, want: []byte{0xff, 0xff}},
		{name: "-1 signed little", x: -1, length: 2, order: "little", signed: []bool{true}, want: []byte{0xff, 0xff}},
		{name: "-1 unsigned default", x: -1, length: 2, order: "big", wantErr: true},
		{name: "-1 unsigned explicit", x: -1, length: 2, order: "big", signed: []bool{false}, wantErr: true},
		{name: "255 one byte", x: 255, length: 1, order: "big", want: []byte{0xff}},
		{name: "256 overflows one byte", x: 256, length: 1, order: "big", wantErr: true},
		{name: "signed 127 one byte", x: 127, length: 1, order: "big", signed: []bool{true}, want: []byte{0x7f}},
		{name: "signed 128 overflows one byte", x: 128, length: 1, order: "big", signed: []bool{true}, wantErr: true},
		{name: "signed -128 one byte", x: -128, length: 1, order: "big", signed: []bool{true}, want: []byte{0x80}},
		{name: "signed -129 overflows one byte", x: -129, length: 1, order: "big", signed: []bool{true}, wantErr: true},
		{name: "-2 signed big", x: -2, length: 2, order: "big", signed: []bool{true}, want: []byte{0xff, 0xfe}},
		{name: "-257 signed little", x: -257, length: 2, order: "little", signed: []bool{true}, want: []byte{0xff, 0xfe}},
		{name: "min int64 big signed", x: math.MinInt64, length: 8, order: "big", signed: []bool{true}, want: minI64Big},
		{name: "min int64 little signed", x: math.MinInt64, length: 8, order: "little", signed: []bool{true}, want: minI64Little},
		{name: "max int64 big signed", x: math.MaxInt64, length: 8, order: "big", signed: []bool{true}, want: maxI64Big},
		{name: "max int64 big unsigned", x: math.MaxInt64, length: 8, order: "big", want: maxI64Big},
		{name: "max int64 little unsigned", x: math.MaxInt64, length: 8, order: "little", want: maxI64Little},
		{name: "zero one byte", x: 0, length: 1, order: "big", want: []byte{0x00}},
		{name: "10 bytes big zero padded", x: 1, length: 10, order: "big", want: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}},
		{name: "10 bytes little zero padded", x: 1, length: 10, order: "little", want: []byte{0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "10 bytes signed negative padded", x: -1, length: 10, order: "little", signed: []bool{true}, want: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{name: "zero length", x: 0, length: 0, order: "big", wantErr: true},
		{name: "negative length", x: 0, length: -1, order: "big", wantErr: true},
		{name: "invalid order", x: 1, length: 1, order: "middle", wantErr: true},
		{name: "empty order", x: 1, length: 1, order: "", wantErr: true},
		{name: "order prefix not accepted", x: 1, length: 1, order: "bigger", wantErr: true},
		{name: "length exceeds limit", x: 1, length: maxIntToBytesLen + 1, order: "big", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntToBytes(tt.x, tt.length, tt.order, tt.signed...)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("IntToBytes(%d, %d, %q, %v) = %v, want error", tt.x, tt.length, tt.order, tt.signed, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("IntToBytes(%d, %d, %q, %v) unexpected error: %v", tt.x, tt.length, tt.order, tt.signed, err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Fatalf("IntToBytes(%d, %d, %q, %v) = %#v, want %#v", tt.x, tt.length, tt.order, tt.signed, got, tt.want)
			}
		})
	}
}

func TestIntFromBytes(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		order   string
		signed  []bool
		want    int64
		wantErr bool
	}{
		// 期望值来自 Python 3.13: int.from_bytes(data, order, signed=signed)
		{name: "1024 big", data: []byte{0x04, 0x00}, order: "big", want: 1024},
		{name: "1024 little", data: []byte{0x00, 0x04}, order: "little", want: 1024},
		{name: "1024 big uppercase", data: []byte{0x04, 0x00}, order: "BIG", want: 1024},
		{name: "1024 mixed case", data: []byte{0x00, 0x04}, order: "LiTtLe", want: 1024},
		{name: "ff ff unsigned big", data: []byte{0xff, 0xff}, order: "big", want: 65535},
		{name: "ff ff signed big", data: []byte{0xff, 0xff}, order: "big", signed: []bool{true}, want: -1},
		{name: "80 unsigned big", data: []byte{0x80}, order: "big", want: 128},
		{name: "80 signed big", data: []byte{0x80}, order: "big", signed: []bool{true}, want: -128},
		{name: "80 signed little", data: []byte{0x80}, order: "little", signed: []bool{true}, want: -128},
		{name: "7f signed", data: []byte{0x7f}, order: "big", signed: []bool{true}, want: 127},
		{name: "ff fe signed little", data: []byte{0xff, 0xfe}, order: "little", signed: []bool{true}, want: -257},
		{name: "ff fe signed big", data: []byte{0xff, 0xfe}, order: "big", signed: []bool{true}, want: -2},
		{name: "min int64 big signed", data: []byte{0x80, 0, 0, 0, 0, 0, 0, 0}, order: "big", signed: []bool{true}, want: math.MinInt64},
		{name: "min int64 little signed", data: []byte{0, 0, 0, 0, 0, 0, 0, 0x80}, order: "little", signed: []bool{true}, want: math.MinInt64},
		{name: "max int64 big unsigned", data: []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, order: "big", want: math.MaxInt64},
		{name: "8 zero bytes signed", data: []byte{0, 0, 0, 0, 0, 0, 0, 0}, order: "big", signed: []bool{true}, want: 0},
		{name: "01 00 little signed", data: []byte{0x01, 0x00}, order: "little", signed: []bool{true}, want: 1},
		{name: "00 80 big unsigned", data: []byte{0x00, 0x80}, order: "big", want: 128},
		{name: "high bit overflows int64", data: []byte{0x80, 0, 0, 0, 0, 0, 0, 0}, order: "big", wantErr: true},
		{name: "high bit little overflows int64", data: []byte{0, 0, 0, 0, 0, 0, 0, 0x80}, order: "little", wantErr: true},
		{name: "all ff unsigned overflows int64", data: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, order: "big", wantErr: true},
		{name: "all ff signed is -1", data: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, order: "big", signed: []bool{true}, want: -1},
		{name: "nil data", data: nil, order: "big", wantErr: true},
		{name: "empty data signed", data: []byte{}, order: "little", signed: []bool{true}, wantErr: true},
		{name: "nine bytes", data: make([]byte, 9), order: "big", wantErr: true},
		{name: "invalid order", data: []byte{0}, order: "middle", wantErr: true},
		{name: "empty order", data: []byte{0}, order: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntFromBytes(tt.data, tt.order, tt.signed...)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("IntFromBytes(%#v, %q, %v) = (%d, nil), want error", tt.data, tt.order, tt.signed, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("IntFromBytes(%#v, %q, %v) unexpected error: %v", tt.data, tt.order, tt.signed, err)
			}
			if got != tt.want {
				t.Fatalf("IntFromBytes(%#v, %q, %v) = %d, want %d", tt.data, tt.order, tt.signed, got, tt.want)
			}
		})
	}
}

func TestIntBytesRoundTrip(t *testing.T) {
	for _, order := range []string{"big", "little"} {
		for _, length := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
			for _, x := range []int64{0, 1, -1} {
				name := fmt.Sprintf("%s/%d/%d", order, length, x)
				t.Run(name, func(t *testing.T) {
					data, err := IntToBytes(x, length, order, true)
					if err != nil {
						t.Fatalf("IntToBytes(%d, %d, %q, true): %v", x, length, order, err)
					}
					got, err := IntFromBytes(data, order, true)
					if err != nil {
						t.Fatalf("IntFromBytes(%#v, %q, true): %v", data, order, err)
					}
					if got != x {
						t.Fatalf("round trip = %d, want %d", got, x)
					}
				})
			}
		}
	}
}
