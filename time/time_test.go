package time

import (
	"os"
	"strconv"
	"testing"
	stdtime "time"
)

// cst 是所有测试使用的固定时区，避免依赖本机时区。
var cst = stdtime.FixedZone("CST", 8*3600)

// TestMain 将标准库 Local 固定为 CST，使 Parse/ParseIn/Strptime 等
// 依赖 Local 的 API 行为确定。
func TestMain(m *testing.M) {
	stdtime.Local = cst
	os.Exit(m.Run())
}

func cstDate(year, month, day, hour, minute, second, nsec int) stdtime.Time {
	return stdtime.Date(year, stdtime.Month(month), day, hour, minute, second, nsec, cst)
}

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       stdtime.Time
		wantOffset int
		wantErr    bool
	}{
		{
			name:       "datetime",
			input:      "2024-01-02 03:04:05",
			want:       cstDate(2024, 1, 2, 3, 4, 5, 0),
			wantOffset: 8 * 3600,
		},
		{
			name:       "date only",
			input:      "2024-01-02",
			want:       cstDate(2024, 1, 2, 0, 0, 0, 0),
			wantOffset: 8 * 3600,
		},
		{
			name:       "minute precision",
			input:      "2024-01-02 03:04",
			want:       cstDate(2024, 1, 2, 3, 4, 0, 0),
			wantOffset: 8 * 3600,
		},
		{
			name:       "slash datetime",
			input:      "2024/01/02 03:04:05",
			want:       cstDate(2024, 1, 2, 3, 4, 5, 0),
			wantOffset: 8 * 3600,
		},
		{
			name:       "slash minute",
			input:      "2024/01/02 03:04",
			want:       cstDate(2024, 1, 2, 3, 4, 0, 0),
			wantOffset: 8 * 3600,
		},
		{
			name:       "slash date",
			input:      "2024/01/02",
			want:       cstDate(2024, 1, 2, 0, 0, 0, 0),
			wantOffset: 8 * 3600,
		},
		{
			name:       "rfc3339nano offset",
			input:      "2024-01-02T03:04:05.123456789+08:00",
			want:       cstDate(2024, 1, 2, 3, 4, 5, 123456789),
			wantOffset: 8 * 3600,
		},
		{
			name:       "rfc3339 utc",
			input:      "2024-01-02T03:04:05Z",
			want:       stdtime.Date(2024, 1, 2, 3, 4, 5, 0, stdtime.UTC),
			wantOffset: 0,
		},
		{
			name:       "rfc3339 negative offset",
			input:      "2024-01-02T03:04:05-05:00",
			want:       stdtime.Date(2024, 1, 2, 8, 4, 5, 0, stdtime.UTC),
			wantOffset: -5 * 3600,
		},
		{name: "empty", input: "", wantErr: true},
		{name: "garbage", input: "not-a-time", wantErr: true},
		{name: "invalid day", input: "2024-02-30 00:00:00", wantErr: true},
		{name: "non leap feb 29", input: "2023-02-29", wantErr: true},
		{name: "invalid hour", input: "2024-01-02 25:00:00", wantErr: true},
		{name: "invalid offset", input: "2024-01-02T03:04:05+25:00", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q) = %v, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
			_, offset := got.Zone()
			if offset != tt.wantOffset {
				t.Errorf("Parse(%q) offset = %d, want %d", tt.input, offset, tt.wantOffset)
			}
		})
	}
}

func TestParseIn(t *testing.T) {
	tests := []struct {
		name       string
		layout     string
		input      string
		want       stdtime.Time
		wantOffset int
		wantErr    bool
	}{
		{
			name:       "standard layout",
			layout:     "2006-01-02 15:04:05",
			input:      "2024-01-02 03:04:05",
			want:       cstDate(2024, 1, 2, 3, 4, 5, 0),
			wantOffset: 8 * 3600,
		},
		{
			name:       "custom layout",
			layout:     "02 Jan 2006 15:04",
			input:      "02 Jan 2024 03:04",
			want:       cstDate(2024, 1, 2, 3, 4, 0, 0),
			wantOffset: 8 * 3600,
		},
		{
			name:       "fractional seconds",
			layout:     "2006-01-02 15:04:05.000",
			input:      "2024-01-02 03:04:05.123",
			want:       cstDate(2024, 1, 2, 3, 4, 5, 123000000),
			wantOffset: 8 * 3600,
		},
		{name: "bad input", layout: "2006-01-02", input: "2024/01/02", wantErr: true},
		{name: "bad layout", layout: "not-a-layout", input: "2024-01-02", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIn(tt.layout, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseIn(%q, %q) = %v, want error", tt.layout, tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseIn(%q, %q) error: %v", tt.layout, tt.input, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseIn(%q, %q) = %v, want %v", tt.layout, tt.input, got, tt.want)
			}
			_, offset := got.Zone()
			if offset != tt.wantOffset {
				t.Errorf("ParseIn(%q, %q) offset = %d, want %d", tt.layout, tt.input, offset, tt.wantOffset)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want string
	}{
		{"cst", cstDate(2024, 1, 2, 3, 4, 5, 0), "2024-01-02 03:04:05"},
		{"utc", stdtime.Date(2024, 1, 2, 3, 4, 5, 0, stdtime.UTC), "2024-01-02 03:04:05"},
		{"midnight", cstDate(2024, 2, 29, 0, 0, 0, 0), "2024-02-29 00:00:00"},
		{"year five", cstDate(5, 1, 1, 0, 0, 0, 0), "0005-01-01 00:00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.t); got != tt.want {
				t.Errorf("Format(%v) = %q, want %q", tt.t, got, tt.want)
			}
		})
	}
}

func TestTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		t       stdtime.Time
		wantSec int64
		wantMs  int64
	}{
		{
			name:    "epoch in CST",
			t:       cstDate(1970, 1, 1, 8, 0, 0, 0),
			wantSec: 0,
			wantMs:  0,
		},
		{
			name:    "2024 with nanoseconds",
			t:       cstDate(2024, 1, 2, 3, 4, 5, 123456789),
			wantSec: 1704135845,
			wantMs:  1704135845123,
		},
		{
			name:    "one second before epoch",
			t:       cstDate(1970, 1, 1, 7, 59, 59, 0),
			wantSec: -1,
			wantMs:  -1000,
		},
		{
			name:    "leap day end",
			t:       cstDate(2024, 2, 29, 23, 59, 59, 0),
			wantSec: 1709222399,
			wantMs:  1709222399000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Timestamp(tt.t); got != tt.wantSec {
				t.Errorf("Timestamp(%v) = %d, want %d", tt.t, got, tt.wantSec)
			}
			if got := TimestampMs(tt.t); got != tt.wantMs {
				t.Errorf("TimestampMs(%v) = %d, want %d", tt.t, got, tt.wantMs)
			}
		})
	}
}

func TestTimestampRoundTrip(t *testing.T) {
	secs := []int64{0, 1704135845, -1, 946656000, 1709222399}
	for _, sec := range secs {
		sec := sec
		t.Run(strconv.FormatInt(sec, 10), func(t *testing.T) {
			got := ParseTimestamp(sec)
			if back := Timestamp(got); back != sec {
				t.Errorf("Timestamp(ParseTimestamp(%d)) = %d", sec, back)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("ParseTimestamp(%d) offset = %d, want %d", sec, offset, 8*3600)
			}
		})
	}
}

func TestFromMillis(t *testing.T) {
	tests := []struct {
		name string
		ms   int64
		want stdtime.Time
	}{
		{"epoch", 0, cstDate(1970, 1, 1, 8, 0, 0, 0)},
		{"2024", 1704135845123, cstDate(2024, 1, 2, 3, 4, 5, 123000000)},
		{"negative", -1500, cstDate(1970, 1, 1, 7, 59, 58, 500000000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromMillis(tt.ms)
			if !got.Equal(tt.want) {
				t.Errorf("FromMillis(%d) = %v, want %v", tt.ms, got, tt.want)
			}
			if back := TimestampMs(got); back != tt.ms {
				t.Errorf("TimestampMs(FromMillis(%d)) = %d", tt.ms, back)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("FromMillis(%d) offset = %d, want %d", tt.ms, offset, 8*3600)
			}
		})
	}
}

func TestStrtimestamp(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int64
	}{
		{"normal", "2024-01-02 03:04:05", 1704135845},
		{"epoch in CST", "1970-01-01 08:00:00", 0},
		{"epoch UTC is negative in CST", "1970-01-01 00:00:00", -28800},
		{"leap day end", "2024-02-29 23:59:59", 1709222399},
		{"invalid text", "not-a-time", 0},
		{"invalid date", "2023-02-29 00:00:00", 0},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Strtimestamp(tt.input); got != tt.want {
				t.Errorf("Strtimestamp(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestIntimestamp(t *testing.T) {
	tests := []struct {
		name string
		sec  int64
		want string
	}{
		{"normal", 1704135845, "2024-01-02 03:04:05"},
		{"epoch in CST", 0, "1970-01-01 08:00:00"},
		{"before epoch", -1, "1970-01-01 07:59:59"},
		{"leap day end", 1709222399, "2024-02-29 23:59:59"},
		{"epoch UTC", -28800, "1970-01-01 00:00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intimestamp(tt.sec); got != tt.want {
				t.Errorf("Intimestamp(%d) = %q, want %q", tt.sec, got, tt.want)
			}
		})
	}
}
