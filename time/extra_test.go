package time

import (
	"testing"
	stdtime "time"
)

// extraCST 是新增测试使用的固定时区，避免依赖本机时区与 Now。
var extraCST = stdtime.FixedZone("CST", 8*3600)

func extraDate(year, month, day, hour, minute, second, nsec int) stdtime.Time {
	return stdtime.Date(year, stdtime.Month(month), day, hour, minute, second, nsec, extraCST)
}

func extraDay(year, month, day int) stdtime.Time {
	return extraDate(year, month, day, 10, 30, 45, 123456789)
}

func TestAddWeeks(t *testing.T) {
	tests := []struct {
		name  string
		t     stdtime.Time
		weeks int
		want  stdtime.Time
	}{
		{"zero", extraDay(2024, 1, 15), 0, extraDay(2024, 1, 15)},
		{"one week across month", extraDay(2024, 1, 25), 1, extraDay(2024, 2, 1)},
		{"two weeks across month", extraDay(2024, 2, 20), 2, extraDay(2024, 3, 5)},
		{"one week across leap day", extraDay(2024, 2, 20), 1, extraDay(2024, 2, 27)},
		{"leap february to march", extraDay(2024, 2, 25), 1, extraDay(2024, 3, 3)},
		{"non leap february to march", extraDay(2023, 2, 25), 1, extraDay(2023, 3, 4)},
		{"from leap day", extraDay(2024, 2, 29), 1, extraDay(2024, 3, 7)},
		{"one week across year", extraDay(2024, 12, 26), 1, extraDay(2025, 1, 2)},
		{"minus one week across year", extraDay(2025, 1, 2), -1, extraDay(2024, 12, 26)},
		{"fifty two weeks", extraDay(2024, 1, 1), 52, extraDay(2024, 12, 30)},
		{"minus fifty two weeks", extraDay(2024, 12, 30), -52, extraDay(2024, 1, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddWeeks(tt.t, tt.weeks)
			if !got.Equal(tt.want) {
				t.Errorf("AddWeeks(%v, %d) = %v, want %v", tt.t, tt.weeks, got, tt.want)
			}
			if got.Location() != extraCST {
				t.Errorf("AddWeeks location = %v, want %v", got.Location(), extraCST)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("AddWeeks offset = %d, want %d", offset, 8*3600)
			}
			if got.Hour() != tt.t.Hour() || got.Minute() != tt.t.Minute() ||
				got.Second() != tt.t.Second() || got.Nanosecond() != tt.t.Nanosecond() {
				t.Errorf("AddWeeks changed clock: got %v, input %v", got, tt.t)
			}
		})
	}
}

func TestQuarter(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want int
	}{
		{"january start", extraDay(2024, 1, 1), 1},
		{"february leap day", extraDay(2024, 2, 29), 1},
		{"march end", extraDay(2024, 3, 31), 1},
		{"april start", extraDay(2024, 4, 1), 2},
		{"june end", extraDay(2024, 6, 30), 2},
		{"july start", extraDay(2024, 7, 1), 3},
		{"september end", extraDay(2024, 9, 30), 3},
		{"october start", extraDay(2024, 10, 1), 4},
		{"december end", extraDay(2024, 12, 31), 4},
		{"next year january", extraDay(2025, 1, 1), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quarter(tt.t); got != tt.want {
				t.Errorf("Quarter(%v) = %d, want %d", tt.t, got, tt.want)
			}
		})
	}

	t.Run("every month", func(t *testing.T) {
		wantByMonth := []int{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4}
		for month := 1; month <= 12; month++ {
			got := Quarter(extraDate(2024, month, 15, 8, 0, 0, 0))
			if want := wantByMonth[month-1]; got != want {
				t.Errorf("Quarter(month %d) = %d, want %d", month, got, want)
			}
		}
	})
}

func TestStartOfQuarter(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want stdtime.Time
	}{
		{"Q1 leap february", extraDay(2024, 2, 29), extraDate(2024, 1, 1, 0, 0, 0, 0)},
		{"Q1 non leap march", extraDay(2023, 3, 31), extraDate(2023, 1, 1, 0, 0, 0, 0)},
		{"Q2 may", extraDay(2024, 5, 15), extraDate(2024, 4, 1, 0, 0, 0, 0)},
		{"Q3 august", extraDay(2024, 8, 31), extraDate(2024, 7, 1, 0, 0, 0, 0)},
		{"Q4 november", extraDay(2024, 11, 30), extraDate(2024, 10, 1, 0, 0, 0, 0)},
		{"Q4 year end", extraDate(2024, 12, 31, 23, 59, 59, 999999999), extraDate(2024, 10, 1, 0, 0, 0, 0)},
		{"already quarter start", extraDate(2024, 7, 1, 0, 0, 0, 0), extraDate(2024, 7, 1, 0, 0, 0, 0)},
		{"Q1 next year", extraDay(2025, 1, 1), extraDate(2025, 1, 1, 0, 0, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StartOfQuarter(tt.t)
			if !got.Equal(tt.want) {
				t.Errorf("StartOfQuarter(%v) = %v, want %v", tt.t, got, tt.want)
			}
			if got.Location() != extraCST {
				t.Errorf("StartOfQuarter location = %v, want %v", got.Location(), extraCST)
			}
			if h, m, s, n := got.Hour(), got.Minute(), got.Second(), got.Nanosecond(); h != 0 || m != 0 || s != 0 || n != 0 {
				t.Errorf("StartOfQuarter clock = %02d:%02d:%02d.%09d, want 00:00:00.000000000", h, m, s, n)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("StartOfQuarter offset = %d, want %d", offset, 8*3600)
			}
		})
	}
}

func TestEndOfQuarter(t *testing.T) {
	end := func(year, month, day int) stdtime.Time {
		return extraDate(year, month, day, 23, 59, 59, 999999999)
	}

	tests := []struct {
		name string
		t    stdtime.Time
		want stdtime.Time
	}{
		{"Q1 leap february", extraDay(2024, 2, 10), end(2024, 3, 31)},
		{"Q1 leap day", extraDay(2024, 2, 29), end(2024, 3, 31)},
		{"Q1 non leap", extraDay(2023, 1, 15), end(2023, 3, 31)},
		{"Q2 april", extraDay(2024, 4, 1), end(2024, 6, 30)},
		{"Q3 july", extraDay(2024, 7, 15), end(2024, 9, 30)},
		{"Q4 october", extraDay(2024, 10, 1), end(2024, 12, 31)},
		{"Q4 year end", extraDate(2024, 12, 31, 0, 0, 0, 0), end(2024, 12, 31)},
		{"already quarter end", end(2024, 6, 30), end(2024, 6, 30)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EndOfQuarter(tt.t)
			if !got.Equal(tt.want) {
				t.Errorf("EndOfQuarter(%v) = %v, want %v", tt.t, got, tt.want)
			}
			if got.Location() != extraCST {
				t.Errorf("EndOfQuarter location = %v, want %v", got.Location(), extraCST)
			}
			if h, m, s, n := got.Hour(), got.Minute(), got.Second(), got.Nanosecond(); h != 23 || m != 59 || s != 59 || n != 999999999 {
				t.Errorf("EndOfQuarter clock = %02d:%02d:%02d.%09d, want 23:59:59.999999999", h, m, s, n)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("EndOfQuarter offset = %d, want %d", offset, 8*3600)
			}
			if next := got.Add(stdtime.Nanosecond); !next.Equal(StartOfQuarter(tt.t).AddDate(0, 3, 0)) {
				t.Errorf("EndOfQuarter(%v)+1ns = %v, want next quarter start %v", tt.t, next, StartOfQuarter(tt.t).AddDate(0, 3, 0))
			}
		})
	}
}

func TestQuarterKeepsLocation(t *testing.T) {
	zones := []*stdtime.Location{
		extraCST,
		stdtime.FixedZone("NPT", 5*3600+45*60),
		stdtime.FixedZone("WEST", -5*3600),
	}

	for _, loc := range zones {
		t.Run(loc.String(), func(t *testing.T) {
			in := stdtime.Date(2024, stdtime.November, 15, 12, 0, 0, 0, loc)
			start := StartOfQuarter(in)
			stop := EndOfQuarter(in)

			if start.Location() != loc || stop.Location() != loc {
				t.Errorf("location not preserved: start=%v end=%v, want %v", start.Location(), stop.Location(), loc)
			}
			_, wantOffset := in.Zone()
			if _, got := start.Zone(); got != wantOffset {
				t.Errorf("StartOfQuarter offset = %d, want %d", got, wantOffset)
			}
			if _, got := stop.Zone(); got != wantOffset {
				t.Errorf("EndOfQuarter offset = %d, want %d", got, wantOffset)
			}
			if start.Year() != 2024 || start.Month() != stdtime.October || start.Day() != 1 {
				t.Errorf("StartOfQuarter = %v, want 2024-10-01", start)
			}
			if stop.Year() != 2024 || stop.Month() != stdtime.December || stop.Day() != 31 {
				t.Errorf("EndOfQuarter = %v, want 2024-12-31", stop)
			}
		})
	}
}

func TestIsWeekend(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want bool
	}{
		{"monday", extraDay(2024, 1, 1), false},
		{"tuesday", extraDay(2024, 1, 2), false},
		{"wednesday", extraDay(2024, 1, 3), false},
		{"thursday", extraDay(2024, 1, 4), false},
		{"friday", extraDay(2024, 1, 5), false},
		{"saturday", extraDay(2024, 1, 6), true},
		{"sunday", extraDay(2024, 1, 7), true},
		{"saturday year start", extraDay(2022, 1, 1), true},
		{"sunday year start", extraDay(2022, 1, 2), true},
		{"friday year end", extraDay(2024, 12, 27), false},
		{"saturday year end", extraDay(2024, 12, 28), true},
		{"sunday year end", extraDay(2024, 12, 29), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWeekend(tt.t); got != tt.want {
				t.Errorf("IsWeekend(%v) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

func TestWeekdayISO(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want int
	}{
		{"monday", extraDay(2024, 1, 1), 1},
		{"tuesday", extraDay(2024, 1, 2), 2},
		{"wednesday", extraDay(2024, 1, 3), 3},
		{"thursday", extraDay(2024, 1, 4), 4},
		{"friday", extraDay(2024, 1, 5), 5},
		{"saturday", extraDay(2024, 1, 6), 6},
		{"sunday", extraDay(2024, 1, 7), 7},
		{"tuesday year end", extraDay(2024, 12, 31), 2},
		{"wednesday new year", extraDay(2025, 1, 1), 3},
		{"saturday year end", extraDay(2024, 12, 28), 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeekdayISO(tt.t)
			if got != tt.want {
				t.Errorf("WeekdayISO(%v) = %d, want %d", tt.t, got, tt.want)
			}
			if got < 1 || got > 7 {
				t.Errorf("WeekdayISO(%v) = %d, want value in 1..7", tt.t, got)
			}
		})
	}
}

func TestDaysInYear(t *testing.T) {
	tests := []struct {
		name string
		year int
		want int
	}{
		{"leap 2024", 2024, 366},
		{"non leap 2025", 2025, 365},
		{"century leap 2000", 2000, 366},
		{"century non leap 1900", 1900, 365},
		{"century non leap 2100", 2100, 365},
		{"non leap 2023", 2023, 365},
		{"leap 1600", 1600, 366},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DaysInYear(tt.year); got != tt.want {
				t.Errorf("DaysInYear(%d) = %d, want %d", tt.year, got, tt.want)
			}
		})
	}
}

func TestExtraDoesNotMutateInput(t *testing.T) {
	base := extraDate(2024, 2, 29, 13, 14, 15, 123456789)
	before := base

	AddWeeks(base, 3)
	Quarter(base)
	StartOfQuarter(base)
	EndOfQuarter(base)
	IsWeekend(base)
	WeekdayISO(base)
	DaysInYear(base.Year())

	if !base.Equal(before) {
		t.Errorf("input changed: got %v, want %v", base, before)
	}
	if base.Location() != before.Location() {
		t.Errorf("input location changed: got %v, want %v", base.Location(), before.Location())
	}
}
