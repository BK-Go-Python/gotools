package time

import (
	"strings"
	"testing"
	stdtime "time"
)

func TestStrftimeDirectives(t *testing.T) {
	ts := cstDate(2024, 1, 2, 3, 4, 5, 123456789)
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"literal", "hello", "hello"},
		{"empty", "", ""},
		{"%Y", "%Y", "2024"},
		{"%y", "%y", "24"},
		{"%m", "%m", "01"},
		{"%d", "%d", "02"},
		{"%e", "%e", " 2"},
		{"%H", "%H", "03"},
		{"%I", "%I", "03"},
		{"%M", "%M", "04"},
		{"%S", "%S", "05"},
		{"%f", "%f", "123456"},
		{"%p", "%p", "AM"},
		{"%a", "%a", "Tue"},
		{"%A", "%A", "Tuesday"},
		{"%b", "%b", "Jan"},
		{"%B", "%B", "January"},
		{"%z", "%z", "+0800"},
		{"%Z", "%Z", "CST"},
		{"%j", "%j", "002"},
		{"%U", "%U", "00"},
		{"%W", "%W", "01"},
		{"%w", "%w", "2"},
		{"%u", "%u", "2"},
		{"%G", "%G", "2024"},
		{"%V", "%V", "01"},
		{"%c", "%c", "Tue Jan  2 03:04:05 2024"},
		{"%x", "%x", "01/02/24"},
		{"%X", "%X", "03:04:05"},
		{"%s", "%s", "1704135845"},
		{"%n", "%n", "\n"},
		{"%t", "%t", "\t"},
		{"%%", "%%", "%"},
		{"trailing percent", "%Y%", "2024%"},
		{"only percent", "%", "%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Strftime(ts, tt.format); got != tt.want {
				t.Errorf("Strftime(%v, %q) = %q, want %q", ts, tt.format, got, tt.want)
			}
		})
	}
}

// TestStrftimeCombined 使用 Python 3.13 在同一 CST 时刻的实际输出
// 校验所有指令按顺序组合的结果。
// Python: datetime(2024,1,2,3,4,5,123456,tzinfo=timezone(timedelta(hours=8), "CST")).strftime(fmt)
func TestStrftimeCombined(t *testing.T) {
	ts := cstDate(2024, 1, 2, 3, 4, 5, 123456000)
	format := "%Y|%y|%m|%d|%e|%H|%I|%M|%S|%f|%p|%a|%A|%b|%B|%z|%Z|%j|%U|%W|%w|%u|%G|%V|%c|%x|%X|%s|%n|%t|%%"
	want := "2024|24|01|02| 2|03|03|04|05|123456|AM|Tue|Tuesday|Jan|January|+0800|CST|" +
		"002|00|01|2|2|2024|01|Tue Jan  2 03:04:05 2024|01/02/24|03:04:05|1704135845|\n|\t|%"

	if got := Strftime(ts, format); got != want {
		t.Errorf("Strftime combined:\n got %q\nwant %q", got, want)
	}
}

func TestStrftimeBoundaries(t *testing.T) {
	ist := stdtime.FixedZone("IST", 5*3600+30*60)
	tests := []struct {
		name   string
		t      stdtime.Time
		format string
		want   string
	}{
		{
			name:   "midnight 12-hour",
			t:      cstDate(2024, 1, 2, 0, 5, 9, 0),
			format: "%I:%M:%S %p",
			want:   "12:05:09 AM",
		},
		{
			name:   "noon 12-hour",
			t:      cstDate(2024, 1, 2, 12, 0, 0, 0),
			format: "%I:%M:%S %p",
			want:   "12:00:00 PM",
		},
		{
			name:   "late 12-hour",
			t:      cstDate(2024, 1, 2, 23, 0, 0, 0),
			format: "%I:%M:%S %p",
			want:   "11:00:00 PM",
		},
		{
			name:   "iso year previous",
			t:      cstDate(2021, 1, 1, 0, 0, 0, 0),
			format: "%G-%V",
			want:   "2020-53",
		},
		{
			name:   "iso year next",
			t:      cstDate(2024, 12, 31, 0, 0, 0, 0),
			format: "%G-%V",
			want:   "2025-01",
		},
		{
			name:   "sunday U",
			t:      cstDate(2023, 12, 31, 0, 0, 0, 0),
			format: "%U",
			want:   "53",
		},
		{
			name:   "monday U",
			t:      cstDate(2024, 1, 1, 0, 0, 0, 0),
			format: "%U",
			want:   "00",
		},
		{
			name:   "monday W",
			t:      cstDate(2024, 1, 1, 0, 0, 0, 0),
			format: "%W",
			want:   "01",
		},
		{
			name:   "utc zone",
			t:      stdtime.Date(2024, 1, 2, 3, 4, 5, 0, stdtime.UTC),
			format: "%z %Z",
			want:   "+0000 UTC",
		},
		{
			name:   "half hour zone",
			t:      stdtime.Date(2024, 1, 2, 3, 4, 5, 0, ist),
			format: "%z %Z",
			want:   "+0530 IST",
		},
		{
			name:   "microsecond padding",
			t:      cstDate(2024, 1, 2, 3, 4, 5, 7*1000), // 7 微秒 = 7000 纳秒
			format: "%f",
			want:   "000007",
		},
		{
			name:   "zero microsecond",
			t:      cstDate(2024, 1, 2, 3, 4, 5, 0),
			format: "%f",
			want:   "000000",
		},
		{
			name:   "epoch",
			t:      cstDate(1970, 1, 1, 8, 0, 0, 0),
			format: "%s",
			want:   "0",
		},
		{
			name:   "before epoch",
			t:      cstDate(1970, 1, 1, 7, 59, 59, 0),
			format: "%s",
			want:   "-1",
		},
		{
			name:   "space padded day",
			t:      cstDate(2024, 1, 5, 0, 0, 0, 0),
			format: "%e|%d",
			want:   " 5|05",
		},
		{
			name:   "year padding",
			t:      cstDate(5, 1, 1, 0, 0, 0, 0),
			format: "%Y",
			want:   "0005",
		},
		{
			name:   "leap day year day",
			t:      cstDate(2024, 2, 29, 0, 0, 0, 0),
			format: "%j",
			want:   "060",
		},
		{
			name:   "century year",
			t:      cstDate(2000, 1, 1, 0, 0, 0, 0),
			format: "%y",
			want:   "00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Strftime(tt.t, tt.format); got != tt.want {
				t.Errorf("Strftime(%v, %q) = %q, want %q", tt.t, tt.format, got, tt.want)
			}
		})
	}
}

func TestStrptimeRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		t      stdtime.Time
		format string
	}{
		{"datetime", cstDate(2024, 1, 2, 3, 4, 5, 0), "%Y-%m-%d %H:%M:%S"},
		{"date", cstDate(2024, 1, 2, 0, 0, 0, 0), "%Y/%m/%d"},
		{"fraction", cstDate(2024, 1, 2, 3, 4, 5, 123456000), "%Y-%m-%d %H:%M:%S.%f"},
		{"12 hour", cstDate(2024, 1, 2, 15, 4, 5, 0), "%Y-%m-%d %I:%M:%S %p"},
		{"names", cstDate(2024, 1, 2, 3, 4, 5, 0), "%a %b %e %Y %H:%M:%S"},
		{"c", cstDate(2024, 1, 2, 3, 4, 5, 0), "%c"},
		{"x and X", cstDate(2024, 1, 2, 3, 4, 5, 0), "%x %X"},
		{"zone offset", cstDate(2024, 1, 2, 3, 4, 5, 0), "%Y-%m-%d %H:%M:%S %z"},
		{"zone name", cstDate(2024, 1, 2, 3, 4, 5, 0), "%Y-%m-%d %H:%M:%S %Z"},
		{"safe literal percent", cstDate(2024, 1, 2, 0, 0, 0, 0), "%% %Y-%m-%d"},
		{"newline and tab", cstDate(2024, 1, 2, 0, 0, 0, 0), "%Y%n%m%t%d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Strftime(tt.t, tt.format)
			got, err := Strptime(s, tt.format)
			if err != nil {
				t.Fatalf("Strptime(%q, %q) error: %v", s, tt.format, err)
			}
			if !got.Equal(tt.t) {
				t.Errorf("Strptime(Strftime(%v, %q)) = %v, want %v", tt.t, tt.format, got, tt.t)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("round trip offset = %d, want %d", offset, 8*3600)
			}
		})
	}
}

func TestStrptimeExplicit(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		format string
		want   stdtime.Time
	}{
		{
			name:   "fraction",
			s:      "2024-01-02 03:04:05.123456",
			format: "%Y-%m-%d %H:%M:%S.%f",
			want:   cstDate(2024, 1, 2, 3, 4, 5, 123456000),
		},
		{
			// 依据：Strptime 文档说明 %f 前可为 "." 或 ","。
			name:   "comma fraction",
			s:      "2024-01-02 03:04:05,123456",
			format: "%Y-%m-%d %H:%M:%S.%f",
			want:   cstDate(2024, 1, 2, 3, 4, 5, 123456000),
		},
		{
			name:   "12 hour pm",
			s:      "2024-01-02 03:04:05 PM",
			format: "%Y-%m-%d %I:%M:%S %p",
			want:   cstDate(2024, 1, 2, 15, 4, 5, 0),
		},
		{
			name:   "space padded day",
			s:      "2024-01-02  5",
			format: "%Y-%m-%d %e",
			want:   cstDate(2024, 1, 5, 0, 0, 0, 0),
		},
		{
			name:   "c layout",
			s:      "Tue Jan  2 03:04:05 2024",
			format: "%c",
			want:   cstDate(2024, 1, 2, 3, 4, 5, 0),
		},
		{
			name:   "x and X",
			s:      "01/02/24 03:04:05",
			format: "%x %X",
			want:   cstDate(2024, 1, 2, 3, 4, 5, 0),
		},
		{
			name:   "positive offset",
			s:      "2024-01-02 03:04:05 +0800",
			format: "%Y-%m-%d %H:%M:%S %z",
			want:   cstDate(2024, 1, 2, 3, 4, 5, 0),
		},
		{
			name:   "negative offset",
			s:      "2024-01-02 03:04:05 -0500",
			format: "%Y-%m-%d %H:%M:%S %z",
			want:   stdtime.Date(2024, 1, 2, 8, 4, 5, 0, stdtime.UTC),
		},
		{
			name:   "zone name",
			s:      "2024-01-02 03:04:05 CST",
			format: "%Y-%m-%d %H:%M:%S %Z",
			want:   cstDate(2024, 1, 2, 3, 4, 5, 0),
		},
		{
			name:   "literal percent",
			s:      "2024-01-02 %",
			format: "%Y-%m-%d %%",
			want:   cstDate(2024, 1, 2, 0, 0, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Strptime(tt.s, tt.format)
			if err != nil {
				t.Fatalf("Strptime(%q, %q) error: %v", tt.s, tt.format, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("Strptime(%q, %q) = %v, want %v", tt.s, tt.format, got, tt.want)
			}
		})
	}
}

// TestStrptimeFractionalSecondsPythonCompatible 记录与 Python 的兼容性差异：
// Python 3.13 的 %f 接受 1-6 位小数，当前实现只接受恰好 6 位。
func TestStrptimeFractionalSecondsPythonCompatible(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantNS int
	}{
		{"one digit", "2024-01-02 03:04:05.1", 100000000},
		{"three digits", "2024-01-02 03:04:05.123", 123000000},
		{"six digits", "2024-01-02 03:04:05.123456", 123456000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Python: datetime.strptime(input, "%Y-%m-%d %H:%M:%S.%f").microsecond
			got, err := Strptime(tt.input, "%Y-%m-%d %H:%M:%S.%f")
			if err != nil {
				t.Fatalf("Strptime(%q, %%f) error: %v", tt.input, err)
			}
			if got.Nanosecond() != tt.wantNS {
				t.Errorf("Strptime(%q, %%f) nanosecond = %d, want %d", tt.input, got.Nanosecond(), tt.wantNS)
			}
		})
	}
}

// TestStrptimeLiteralPercentAfterDigits 记录 %% 紧跟数字时的兼容性问题。
// Python 3.13: datetime.strptime("2024-01-02 100%", "%Y-%m-%d 100%%") 可正常解析。
func TestStrptimeLiteralPercentAfterDigits(t *testing.T) {
	got, err := Strptime("2024-01-02 100%", "%Y-%m-%d 100%%")
	if err != nil {
		t.Fatalf("Strptime with literal %% after digits error: %v", err)
	}
	want := cstDate(2024, 1, 2, 0, 0, 0, 0)
	if !got.Equal(want) {
		t.Errorf("Strptime with literal %% after digits = %v, want %v", got, want)
	}
}

func TestStrptimeErrors(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		format string
	}{
		{"invalid day", "2024-02-30", "%Y-%m-%d"},
		{"invalid month", "2024-13-01", "%Y-%m-%d"},
		{"invalid hour", "2024-01-02 25:00:00", "%Y-%m-%d %H:%M:%S"},
		{"non leap feb 29", "2023-02-29", "%Y-%m-%d"},
		{"trailing junk", "2024-01-02 extra", "%Y-%m-%d"},
		{"wrong literal", "2024/01/02", "%Y-%m-%d"},
		{"too many fraction digits", "2024-01-02 03:04:05.1234567", "%Y-%m-%d %H:%M:%S.%f"},
		{"empty fraction", "2024-01-02 03:04:05.", "%Y-%m-%d %H:%M:%S.%f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Strptime(tt.s, tt.format); err == nil {
				t.Fatalf("Strptime(%q, %q) = %v, want error", tt.s, tt.format, got)
			}
		})
	}
}

func TestStrptimeUnsupportedDirectives(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"%j", "%j"},
		{"%U", "%U"},
		{"%W", "%W"},
		{"%w", "%w"},
		{"%u", "%u"},
		{"%G", "%G"},
		{"%V", "%V"},
		{"embedded %s", "%Y %s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Strptime("2024", tt.format)
			if err == nil {
				t.Fatalf("Strptime with unsupported directive %q = nil error, want error", tt.format)
			}
			if !strings.Contains(err.Error(), "unsupported directive") {
				t.Errorf("error = %v, want unsupported directive message", err)
			}
		})
	}
}
