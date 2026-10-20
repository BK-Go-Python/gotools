package conv

import "testing"

func TestComplex(t *testing.T) {
	tests := []struct {
		name    string
		real    []float64
		want    complex128
		wantErr bool
	}{
		{"no args", nil, 0, false},
		{"real only", []float64{1}, 1, false},
		{"real and imag", []float64{1, 2}, complex(1, 2), false},
		{"too many", []float64{1, 2, 3}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Complex(tt.real...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("Complex(%v) = %v, want %v", tt.real, got, tt.want)
			}
		})
	}
}
