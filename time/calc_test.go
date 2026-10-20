package time

import (
	"testing"
	stdtime "time"
)

func TestAddDays(t *testing.T) {
	at := func(year, month, day int) stdtime.Time {
		return cstDate(year, month, day, 10, 20, 30, 123456789)
	}

	tests := []struct {
		name string
		t    stdtime.Time
		days int
		want stdtime.Time
	}{
		{"zero", at(2024, 1, 31), 0, at(2024, 1, 31)},
		{"next day", at(2024, 1, 31), 1, at(2024, 2, 1)},
		{"to leap day", at(2024, 1, 31), 29, at(2024, 2, 29)},
		{"to march", at(2024, 1, 31), 30, at(2024, 3, 1)},
		{"plus 31", at(2024, 1, 31), 31, at(2024, 3, 2)},
		{"leap feb 28 plus one", at(2024, 2, 28), 1, at(2024, 2, 29)},
		{"non leap feb 28 plus one", at(2023, 2, 28), 1, at(2023, 3, 1)},
		{"year rollover", at(2024, 12, 31), 1, at(2025, 1, 1)},
		{"minus one", at(2024, 1, 1), -1, at(2023, 12, 31)},
		{"minus one from march", at(2024, 3, 1), -1, at(2024, 2, 29)},
		{"leap year 366", at(2024, 1, 1), 366, at(2025, 1, 1)},
		{"minus 366", at(2024, 1, 1), -366, at(2022, 12, 31)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddDays(tt.t, tt.days)
			if !got.Equal(tt.want) {
				t.Errorf("AddDays(%v, %d) = %v, want %v", tt.t, tt.days, got, tt.want)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("AddDays offset = %d, want %d", offset, 8*3600)
			}
		})
	}
}

func TestAddMonths(t *testing.T) {
	at := func(year, month, day int) stdtime.Time {
		return cstDate(year, month, day, 10, 20, 30, 123456789)
	}

	tests := []struct {
		name   string
		t      stdtime.Time
		months int
		want   stdtime.Time
	}{
		{"jan 31 plus 1 leap", at(2024, 1, 31), 1, at(2024, 2, 29)},
		{"jan 31 plus 1 non leap", at(2023, 1, 31), 1, at(2023, 2, 28)},
		{"mar 31 minus 1 leap", at(2024, 3, 31), -1, at(2024, 2, 29)},
		{"mar 31 minus 1 non leap", at(2023, 3, 31), -1, at(2023, 2, 28)},
		{"plus two", at(2024, 1, 31), 2, at(2024, 3, 31)},
		{"plus three clamps april", at(2024, 1, 31), 3, at(2024, 4, 30)},
		{"dec plus two leap", at(2023, 12, 31), 2, at(2024, 2, 29)},
		{"dec plus two non leap", at(2024, 12, 31), 2, at(2025, 2, 28)},
		{"cross year negative", at(2024, 1, 15), -1, at(2023, 12, 15)},
		{"negative clamp", at(2024, 1, 31), -2, at(2023, 11, 30)},
		{"large negative", at(2024, 1, 31), -14, at(2022, 11, 30)},
		{"zero", at(2024, 1, 31), 0, at(2024, 1, 31)},
		{"aug plus six", at(2023, 8, 31), 6, at(2024, 2, 29)},
		{"leap day minus twelve", at(2024, 2, 29), -12, at(2023, 2, 28)},
		{"century non leap", at(1900, 1, 31), 1, at(1900, 2, 28)},
		{"year 2000 leap", at(2000, 1, 31), 1, at(2000, 2, 29)},
		{"plus 25", at(2024, 1, 31), 25, at(2026, 2, 28)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddMonths(tt.t, tt.months)
			if !got.Equal(tt.want) {
				t.Errorf("AddMonths(%v, %d) = %v, want %v", tt.t, tt.months, got, tt.want)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("AddMonths offset = %d, want %d", offset, 8*3600)
			}
		})
	}
}

func TestAddYears(t *testing.T) {
	at := func(year, month, day int) stdtime.Time {
		return cstDate(year, month, day, 1, 2, 3, 456)
	}

	tests := []struct {
		name  string
		t     stdtime.Time
		years int
		want  stdtime.Time
	}{
		{"leap day plus one", at(2024, 2, 29), 1, at(2025, 2, 28)},
		{"leap day plus four", at(2024, 2, 29), 4, at(2028, 2, 29)},
		{"leap day minus one", at(2024, 2, 29), -1, at(2023, 2, 28)},
		{"leap day minus four", at(2024, 2, 29), -4, at(2020, 2, 29)},
		{"zero", at(2024, 2, 29), 0, at(2024, 2, 29)},
		{"century non leap", at(2000, 2, 29), 100, at(2100, 2, 28)},
		{"four century leap", at(2000, 2, 29), 400, at(2400, 2, 29)},
		{"march first", at(2023, 3, 1), 1, at(2024, 3, 1)},
		{"leap day minus five", at(2024, 2, 29), -5, at(2019, 2, 28)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddYears(tt.t, tt.years)
			if !got.Equal(tt.want) {
				t.Errorf("AddYears(%v, %d) = %v, want %v", tt.t, tt.years, got, tt.want)
			}
		})
	}
}

// TestDaysBetween 期望值与 Python (b - a).days 的向负无穷取整一致。
func TestDaysBetween(t *testing.T) {
	base := cstDate(2024, 1, 1, 12, 0, 0, 0)
	tests := []struct {
		name string
		a    stdtime.Time
		b    stdtime.Time
		want int
	}{
		{"same instant", base, base, 0},
		{"plus one nanosecond", base, base.Add(stdtime.Nanosecond), 0},
		{"minus one nanosecond", base, base.Add(-stdtime.Nanosecond), -1},
		{"23h59m", base, base.Add(23*stdtime.Hour + 59*stdtime.Minute), 0},
		{"24h", base, base.Add(24 * stdtime.Hour), 1},
		{"24h plus ns", base, base.Add(24*stdtime.Hour + stdtime.Nanosecond), 1},
		{"36h", base, base.Add(36 * stdtime.Hour), 1},
		{"minus 1h", base, base.Add(-stdtime.Hour), -1},
		{"minus 24h", base, base.Add(-24 * stdtime.Hour), -1},
		{"minus 24h minus ns", base, base.Add(-24*stdtime.Hour - stdtime.Nanosecond), -2},
		{"minus 36h", base, base.Add(-36 * stdtime.Hour), -2},
		{"7 days", base, base.Add(7 * 24 * stdtime.Hour), 7},
		{"minus 7 days", base, base.Add(-7 * 24 * stdtime.Hour), -7},
		{
			name: "leap feb 28 to mar 1",
			a:    cstDate(2024, 2, 28, 0, 0, 0, 0),
			b:    cstDate(2024, 3, 1, 0, 0, 0, 0),
			want: 2,
		},
		{
			name: "non leap feb 28 to mar 1",
			a:    cstDate(2023, 2, 28, 0, 0, 0, 0),
			b:    cstDate(2023, 3, 1, 0, 0, 0, 0),
			want: 1,
		},
		{
			name: "year boundary",
			a:    cstDate(2024, 12, 31, 0, 0, 0, 0),
			b:    cstDate(2025, 1, 1, 0, 0, 0, 0),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DaysBetween(tt.a, tt.b); got != tt.want {
				t.Errorf("DaysBetween(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestStartOfDay(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want stdtime.Time
	}{
		{"normal", cstDate(2024, 1, 2, 3, 4, 5, 123456789), cstDate(2024, 1, 2, 0, 0, 0, 0)},
		{"already midnight", cstDate(2024, 1, 2, 0, 0, 0, 0), cstDate(2024, 1, 2, 0, 0, 0, 0)},
		{"leap day", cstDate(2024, 2, 29, 23, 59, 59, 999999999), cstDate(2024, 2, 29, 0, 0, 0, 0)},
		{"year end", cstDate(2024, 12, 31, 12, 0, 0, 0), cstDate(2024, 12, 31, 0, 0, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StartOfDay(tt.t)
			if !got.Equal(tt.want) {
				t.Errorf("StartOfDay(%v) = %v, want %v", tt.t, got, tt.want)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("StartOfDay offset = %d, want %d", offset, 8*3600)
			}
		})
	}
}

func TestEndOfDay(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want stdtime.Time
	}{
		{"normal", cstDate(2024, 1, 2, 3, 4, 5, 0), cstDate(2024, 1, 2, 23, 59, 59, 999999999)},
		{"leap day", cstDate(2024, 2, 29, 12, 0, 0, 0), cstDate(2024, 2, 29, 23, 59, 59, 999999999)},
		{"year end", cstDate(2024, 12, 31, 0, 0, 0, 0), cstDate(2024, 12, 31, 23, 59, 59, 999999999)},
		{"single digit day", cstDate(2024, 1, 5, 6, 0, 0, 0), cstDate(2024, 1, 5, 23, 59, 59, 999999999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfDay(tt.t); !got.Equal(tt.want) {
				t.Errorf("EndOfDay(%v) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

func TestStartOfMonth(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want stdtime.Time
	}{
		{"normal", cstDate(2024, 1, 2, 3, 4, 5, 123), cstDate(2024, 1, 1, 0, 0, 0, 0)},
		{"leap february", cstDate(2024, 2, 29, 23, 59, 59, 0), cstDate(2024, 2, 1, 0, 0, 0, 0)},
		{"december", cstDate(2024, 12, 31, 12, 0, 0, 0), cstDate(2024, 12, 1, 0, 0, 0, 0)},
		{"already first", cstDate(2023, 5, 1, 0, 0, 0, 0), cstDate(2023, 5, 1, 0, 0, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StartOfMonth(tt.t); !got.Equal(tt.want) {
				t.Errorf("StartOfMonth(%v) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

func TestEndOfMonth(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want stdtime.Time
	}{
		{"january", cstDate(2024, 1, 15, 12, 0, 0, 0), cstDate(2024, 1, 31, 23, 59, 59, 999999999)},
		{"leap february", cstDate(2024, 2, 15, 12, 0, 0, 0), cstDate(2024, 2, 29, 23, 59, 59, 999999999)},
		{"non leap february", cstDate(2023, 2, 15, 12, 0, 0, 0), cstDate(2023, 2, 28, 23, 59, 59, 999999999)},
		{"april", cstDate(2024, 4, 1, 0, 0, 0, 0), cstDate(2024, 4, 30, 23, 59, 59, 999999999)},
		{"december", cstDate(2024, 12, 15, 12, 0, 0, 0), cstDate(2024, 12, 31, 23, 59, 59, 999999999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfMonth(tt.t); !got.Equal(tt.want) {
				t.Errorf("EndOfMonth(%v) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

func TestStartOfWeek(t *testing.T) {
	tests := []struct {
		name string
		t    stdtime.Time
		want stdtime.Time
	}{
		{"monday", cstDate(2024, 1, 1, 10, 0, 0, 0), cstDate(2024, 1, 1, 0, 0, 0, 0)},
		{"tuesday", cstDate(2024, 1, 2, 10, 0, 0, 0), cstDate(2024, 1, 1, 0, 0, 0, 0)},
		{"thursday", cstDate(2024, 1, 4, 23, 59, 59, 0), cstDate(2024, 1, 1, 0, 0, 0, 0)},
		{"sunday", cstDate(2024, 1, 7, 10, 0, 0, 0), cstDate(2024, 1, 1, 0, 0, 0, 0)},
		{"sunday across year", cstDate(2023, 1, 1, 10, 0, 0, 0), cstDate(2022, 12, 26, 0, 0, 0, 0)},
		{"wednesday across year", cstDate(2025, 1, 1, 10, 0, 0, 0), cstDate(2024, 12, 30, 0, 0, 0, 0)},
		{"monday midnight", cstDate(2024, 1, 1, 0, 0, 0, 0), cstDate(2024, 1, 1, 0, 0, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StartOfWeek(tt.t)
			if !got.Equal(tt.want) {
				t.Errorf("StartOfWeek(%v) = %v, want %v", tt.t, got, tt.want)
			}
			if _, offset := got.Zone(); offset != 8*3600 {
				t.Errorf("StartOfWeek offset = %d, want %d", offset, 8*3600)
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{2024, true},
		{2023, false},
		{2000, true},
		{1900, false},
		{2100, false},
		{2400, true},
		{4, true},
		{1, false},
		{0, true},
	}

	for _, tt := range tests {
		t.Run(stdtime.Date(tt.year, 1, 1, 0, 0, 0, 0, stdtime.UTC).Format("2006"), func(t *testing.T) {
			if got := IsLeapYear(tt.year); got != tt.want {
				t.Errorf("IsLeapYear(%d) = %v, want %v", tt.year, got, tt.want)
			}
		})
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month stdtime.Month
		want  int
	}{
		{"january", 2024, stdtime.January, 31},
		{"leap february", 2024, stdtime.February, 29},
		{"non leap february", 2023, stdtime.February, 28},
		{"march", 2024, stdtime.March, 31},
		{"april", 2024, stdtime.April, 30},
		{"may", 2024, stdtime.May, 31},
		{"june", 2024, stdtime.June, 30},
		{"july", 2024, stdtime.July, 31},
		{"august", 2024, stdtime.August, 31},
		{"september", 2024, stdtime.September, 30},
		{"october", 2024, stdtime.October, 31},
		{"november", 2024, stdtime.November, 30},
		{"december", 2024, stdtime.December, 31},
		{"year 2000 february", 2000, stdtime.February, 29},
		{"year 1900 february", 1900, stdtime.February, 28},
		{"year 2100 february", 2100, stdtime.February, 28},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DaysInMonth(tt.year, tt.month); got != tt.want {
				t.Errorf("DaysInMonth(%d, %v) = %d, want %d", tt.year, tt.month, got, tt.want)
			}
		})
	}
}

func TestIsSameDay(t *testing.T) {
	otherZone := stdtime.FixedZone("X", -3600)
	tests := []struct {
		name string
		a    stdtime.Time
		b    stdtime.Time
		want bool
	}{
		{
			name: "same day different time",
			a:    cstDate(2024, 1, 2, 0, 0, 0, 0),
			b:    cstDate(2024, 1, 2, 23, 59, 59, 0),
			want: true,
		},
		{
			name: "different day",
			a:    cstDate(2024, 1, 2, 23, 59, 59, 0),
			b:    cstDate(2024, 1, 3, 0, 0, 0, 0),
			want: false,
		},
		{
			name: "different month",
			a:    cstDate(2024, 1, 2, 0, 0, 0, 0),
			b:    cstDate(2024, 2, 2, 0, 0, 0, 0),
			want: false,
		},
		{
			name: "different year",
			a:    cstDate(2023, 1, 2, 0, 0, 0, 0),
			b:    cstDate(2024, 1, 2, 0, 0, 0, 0),
			want: false,
		},
		{
			name: "same components different zones",
			a:    cstDate(2024, 1, 2, 0, 0, 0, 0),
			b:    stdtime.Date(2024, 1, 2, 12, 0, 0, 0, otherZone),
			want: true,
		},
		{
			name: "same instant different local date",
			a:    cstDate(2024, 1, 1, 0, 30, 0, 0),
			b:    cstDate(2024, 1, 1, 0, 30, 0, 0).In(stdtime.UTC),
			want: false,
		},
		{
			name: "leap day",
			a:    cstDate(2024, 2, 29, 1, 0, 0, 0),
			b:    cstDate(2024, 2, 29, 20, 0, 0, 0),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameDay(tt.a, tt.b); got != tt.want {
				t.Errorf("IsSameDay(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAge(t *testing.T) {
	tests := []struct {
		name  string
		birth stdtime.Time
		now   stdtime.Time
		want  int
	}{
		{
			name:  "birthday today",
			birth: cstDate(2000, 6, 15, 0, 0, 0, 0),
			now:   cstDate(2024, 6, 15, 0, 0, 0, 0),
			want:  24,
		},
		{
			name:  "day before birthday",
			birth: cstDate(2000, 6, 15, 0, 0, 0, 0),
			now:   cstDate(2024, 6, 14, 23, 59, 59, 0),
			want:  23,
		},
		{
			name:  "day after birthday",
			birth: cstDate(2000, 6, 15, 0, 0, 0, 0),
			now:   cstDate(2024, 6, 16, 0, 0, 0, 0),
			want:  24,
		},
		{
			name:  "same year before birthday",
			birth: cstDate(2024, 6, 15, 0, 0, 0, 0),
			now:   cstDate(2024, 6, 14, 0, 0, 0, 0),
			want:  0,
		},
		{
			name:  "same day",
			birth: cstDate(2024, 1, 1, 12, 0, 0, 0),
			now:   cstDate(2024, 1, 1, 13, 0, 0, 0),
			want:  0,
		},
		{
			name:  "newborn next day",
			birth: cstDate(2024, 1, 1, 23, 59, 0, 0),
			now:   cstDate(2024, 1, 2, 0, 0, 0, 0),
			want:  0,
		},
		{
			name:  "future birth",
			birth: cstDate(2024, 1, 1, 0, 0, 0, 0),
			now:   cstDate(2023, 1, 1, 0, 0, 0, 0),
			want:  0,
		},
		{
			name:  "leap day before birthday in non leap year",
			birth: cstDate(2000, 2, 29, 0, 0, 0, 0),
			now:   cstDate(2023, 2, 28, 0, 0, 0, 0),
			want:  22,
		},
		{
			name:  "leap day after birthday in non leap year",
			birth: cstDate(2000, 2, 29, 0, 0, 0, 0),
			now:   cstDate(2023, 3, 1, 0, 0, 0, 0),
			want:  23,
		},
		{
			name:  "leap day in leap year",
			birth: cstDate(2000, 2, 29, 0, 0, 0, 0),
			now:   cstDate(2024, 2, 29, 0, 0, 0, 0),
			want:  24,
		},
		{
			name:  "leap day 2028",
			birth: cstDate(2000, 2, 29, 0, 0, 0, 0),
			now:   cstDate(2028, 2, 29, 0, 0, 0, 0),
			want:  28,
		},
		{
			name:  "century",
			birth: cstDate(2000, 1, 1, 0, 0, 0, 0),
			now:   cstDate(2100, 1, 1, 0, 0, 0, 0),
			want:  100,
		},
		{
			name:  "year end birth",
			birth: cstDate(2000, 12, 31, 0, 0, 0, 0),
			now:   cstDate(2024, 1, 1, 0, 0, 0, 0),
			want:  23,
		},
		{
			name:  "time of day ignored",
			birth: cstDate(2000, 1, 1, 23, 0, 0, 0),
			now:   cstDate(2024, 1, 1, 0, 1, 0, 0),
			want:  24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Age(tt.birth, tt.now); got != tt.want {
				t.Errorf("Age(%v, %v) = %d, want %d", tt.birth, tt.now, got, tt.want)
			}
		})
	}
}
