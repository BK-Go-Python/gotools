package str

import "testing"

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		spec    string
		want    string
		wantErr bool
	}{
		{"float fixed", 3.14159, ".2f", "3.14", false},
		{"right align", "ab", ">5", "   ab", false},
		{"center", "ab", "^7", "  ab   ", false},
		{"hex", 255, "x", "ff", false},
		{"zero pad negative", -42, "05d", "-0042", false},
		{"percent", 0.125, ".1%", "12.5%", false},
		{"default float", 3.5, "", "3.5", false},
		{"default string", "abc", "", "abc", false},
		{"binary", 5, "b", "101", false},
		{"invalid type", 1, "q", "", true},
		{"unsupported flag", 1, "#x", "", true},
		{"float spec on string", "abc", ".2f", "", true},
		{"int spec on nil", nil, "d", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatValue(tt.value, tt.spec)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("FormatValue(%#v, %q) = %q, want %q", tt.value, tt.spec, got, tt.want)
			}
		})
	}
}
