// Package datetime 提供 Python datetime 风格的 Date、DateTime、Timedelta 类型。
//
// 与 gotools/time 的函数式 API 互补：这里用类型方法表达对象式操作。
package datetime

import (
	"fmt"
	stdtime "time"

	gotime "github.com/BK-Go-Python/gotools/time"
)

// Timedelta 表示时间间隔，对应 Python 的 datetime.timedelta。
type Timedelta struct {
	duration stdtime.Duration
}

// NewTimedelta 从 time.Duration 创建 Timedelta。
func NewTimedelta(d stdtime.Duration) Timedelta {
	return Timedelta{duration: d}
}

// Duration 返回底层 time.Duration。
func (t Timedelta) Duration() stdtime.Duration { return t.duration }

// Days 返回天数（可为小数）。
func (t Timedelta) Days() float64 { return t.duration.Hours() / 24 }

// TotalSeconds 返回总秒数。
func (t Timedelta) TotalSeconds() float64 { return t.duration.Seconds() }

// TotalMilliseconds 返回总毫秒数。
func (t Timedelta) TotalMilliseconds() float64 {
	return float64(t.duration) / float64(stdtime.Millisecond)
}

// Add 返回 t + other。
func (t Timedelta) Add(other Timedelta) Timedelta { return Timedelta{t.duration + other.duration} }

// Sub 返回 t - other。
func (t Timedelta) Sub(other Timedelta) Timedelta { return Timedelta{t.duration - other.duration} }

// Neg 返回 -t。
func (t Timedelta) Neg() Timedelta { return Timedelta{-t.duration} }

// Abs 返回绝对值。
func (t Timedelta) Abs() Timedelta {
	if t.duration < 0 {
		return Timedelta{-t.duration}
	}
	return t
}

// String 返回 Go duration 字符串。
func (t Timedelta) String() string { return t.duration.String() }

// Date 表示日期，内部为当天 00:00:00 的 time.Time。
type Date struct {
	t stdtime.Time
}

// NewDate 创建 Date；非法日期（如 2 月 30 日）返回 error。
func NewDate(year int, month stdtime.Month, day int) (Date, error) {
	t := stdtime.Date(year, month, day, 0, 0, 0, 0, stdtime.Local)
	if t.Year() != year || t.Month() != month || t.Day() != day {
		return Date{}, fmt.Errorf("datetime.NewDate: invalid date %04d-%02d-%02d", year, int(month), day)
	}
	return Date{t: t}, nil
}

// DateFromTime 取 t 的日期部分。
func DateFromTime(t stdtime.Time) Date {
	y, m, d := t.Date()
	return Date{t: stdtime.Date(y, m, d, 0, 0, 0, 0, t.Location())}
}

// Time 返回当天 00:00:00。
func (d Date) Time() stdtime.Time { return d.t }

// Year 返回年份。
func (d Date) Year() int { return d.t.Year() }

// Month 返回月份。
func (d Date) Month() stdtime.Month { return d.t.Month() }

// Day 返回日。
func (d Date) Day() int { return d.t.Day() }

// Weekday 返回 ISO 星期：周一 1，周日 7。
func (d Date) Weekday() int { return gotime.WeekdayISO(d.t) }

// ISOFormat 返回 2006-01-02。
func (d Date) ISOFormat() string { return d.t.Format("2006-01-02") }

// Format 按 Go 布局格式化。
func (d Date) Format(layout string) string { return d.t.Format(layout) }

// AddDays 返回 d + n 天。
func (d Date) AddDays(n int) Date { return Date{t: d.t.AddDate(0, 0, n)} }

// Sub 返回两个日期相差的天数（b - a 的方向与 Python date - date 一致）。
func (d Date) Sub(other Date) int { return gotime.DaysBetween(other.t, d.t) }

// Before 判断 d 是否早于 other。
func (d Date) Before(other Date) bool { return d.t.Before(other.t) }

// After 判断 d 是否晚于 other。
func (d Date) After(other Date) bool { return d.t.After(other.t) }

// Equal 判断日期是否相同。
func (d Date) Equal(other Date) bool { return d.t.Equal(other.t) }

// String 返回 2006-01-02。
func (d Date) String() string { return d.ISOFormat() }

// DateTime 表示日期时间。
type DateTime struct {
	t stdtime.Time
}

// DateTimeFromTime 包装 time.Time。
func DateTimeFromTime(t stdtime.Time) DateTime { return DateTime{t: t} }

// Now 等价于 Python 的 datetime.now()。
func Now() DateTime { return DateTime{t: stdtime.Now()} }

// Today 等价于 Python 的 datetime.today()（当天 00:00:00）。
func Today() DateTime { return DateTime{t: gotime.Today()} }

// FromTimestamp 等价于 Python 的 datetime.fromtimestamp(sec)。
func FromTimestamp(sec int64) DateTime {
	return DateTime{t: stdtime.Unix(sec, 0)}
}

// Parse 按 Go 布局解析本地时间。
func Parse(layout, value string) (DateTime, error) {
	t, err := stdtime.ParseInLocation(layout, value, stdtime.Local)
	if err != nil {
		return DateTime{}, fmt.Errorf("datetime.Parse: %w", err)
	}
	return DateTime{t: t}, nil
}

// Time 返回底层 time.Time。
func (d DateTime) Time() stdtime.Time { return d.t }

// Timestamp 返回 Unix 秒。
func (d DateTime) Timestamp() int64 { return d.t.Unix() }

// Year 返回年份。
func (d DateTime) Year() int { return d.t.Year() }

// Month 返回月份。
func (d DateTime) Month() stdtime.Month { return d.t.Month() }

// Day 返回日。
func (d DateTime) Day() int { return d.t.Day() }

// Hour 返回小时。
func (d DateTime) Hour() int { return d.t.Hour() }

// Minute 返回分钟。
func (d DateTime) Minute() int { return d.t.Minute() }

// Second 返回秒。
func (d DateTime) Second() int { return d.t.Second() }

// Nanosecond 返回纳秒。
func (d DateTime) Nanosecond() int { return d.t.Nanosecond() }

// Weekday 返回 ISO 星期：周一 1，周日 7。
func (d DateTime) Weekday() int { return gotime.WeekdayISO(d.t) }

// ISOFormat 返回 2006-01-02T15:04:05。
func (d DateTime) ISOFormat() string { return d.t.Format("2006-01-02T15:04:05") }

// Strftime 等价于 Python 的 datetime.strftime(format)。
func (d DateTime) Strftime(format string) string { return gotime.Strftime(d.t, format) }

// Format 按 Go 布局格式化。
func (d DateTime) Format(layout string) string { return d.t.Format(layout) }

// Date 返回日期部分。
func (d DateTime) Date() Date { return DateFromTime(d.t) }

// Add 返回 d + delta。
func (d DateTime) Add(delta Timedelta) DateTime { return DateTime{t: d.t.Add(delta.duration)} }

// Sub 返回 d - other 的间隔。
func (d DateTime) Sub(other DateTime) Timedelta {
	return Timedelta{duration: d.t.Sub(other.t)}
}

// AddDays 返回 d + n 天。
func (d DateTime) AddDays(n int) DateTime { return DateTime{t: d.t.AddDate(0, 0, n)} }

// AddMonths 返回 d + n 个月（月末截断）。
func (d DateTime) AddMonths(n int) DateTime { return DateTime{t: gotime.AddMonths(d.t, n)} }

// AddYears 返回 d + n 年（闰日截断）。
func (d DateTime) AddYears(n int) DateTime { return DateTime{t: gotime.AddYears(d.t, n)} }

// Before 判断 d 是否早于 other。
func (d DateTime) Before(other DateTime) bool { return d.t.Before(other.t) }

// After 判断 d 是否晚于 other。
func (d DateTime) After(other DateTime) bool { return d.t.After(other.t) }

// Equal 判断时间是否相同。
func (d DateTime) Equal(other DateTime) bool { return d.t.Equal(other.t) }

// String 返回 2006-01-02 15:04:05。
func (d DateTime) String() string { return d.t.Format("2006-01-02 15:04:05") }
