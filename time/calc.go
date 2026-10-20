package time

import (
	stdtime "time"
)

// AddDays 等价于 Python 的 t + timedelta(days=n)。
func AddDays(t stdtime.Time, days int) stdtime.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths 增加月份，日期超出目标月份时自动截断到月末
// （例如 1 月 31 日加一个月得到 2 月 28/29 日）。
func AddMonths(t stdtime.Time, months int) stdtime.Time {
	year, month, day := t.Date()
	total := int(month) - 1 + months
	year += total / 12
	total %= 12
	if total < 0 {
		total += 12
		year--
	}
	target := stdtime.Month(total + 1)
	if max := DaysInMonth(year, target); day > max {
		day = max
	}
	return stdtime.Date(year, target, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// AddYears 增加年份，2 月 29 日遇到平年会截断到 2 月 28 日。
func AddYears(t stdtime.Time, years int) stdtime.Time {
	return AddMonths(t, years*12)
}

// DaysBetween 等价于 Python 的 (b - a).days，向负无穷取整。
func DaysBetween(a, b stdtime.Time) int {
	d := b.Sub(a)
	const day = 24 * stdtime.Hour
	days := int(d / day)
	if d%day != 0 && d < 0 {
		days--
	}
	return days
}

// StartOfDay 返回当天 00:00:00。
func StartOfDay(t stdtime.Time) stdtime.Time {
	year, month, day := t.Date()
	return stdtime.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 返回当天 23:59:59.999999999。
func EndOfDay(t stdtime.Time) stdtime.Time {
	return StartOfDay(t).AddDate(0, 0, 1).Add(-stdtime.Nanosecond)
}

// StartOfMonth 返回当月 1 日 00:00:00。
func StartOfMonth(t stdtime.Time) stdtime.Time {
	year, month, _ := t.Date()
	return stdtime.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 返回当月最后一天的 23:59:59.999999999。
func EndOfMonth(t stdtime.Time) stdtime.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-stdtime.Nanosecond)
}

// StartOfWeek 返回本周一 00:00:00（与 Python 的 isoweekday 一致）。
func StartOfWeek(t stdtime.Time) stdtime.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return StartOfDay(t).AddDate(0, 0, 1-weekday)
}

// IsLeapYear 判断是否闰年。
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysInMonth 返回指定年月的天数。
func DaysInMonth(year int, month stdtime.Month) int {
	return stdtime.Date(year, month+1, 0, 0, 0, 0, 0, stdtime.UTC).Day()
}

// IsSameDay 判断两个时间是否是同一天（忽略时区差异，按各自本地日期比较）。
func IsSameDay(a, b stdtime.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// Age 计算周岁年龄；now 省略时使用当前时间。
func Age(birth stdtime.Time, now ...stdtime.Time) int {
	current := Now()
	if len(now) > 0 {
		current = now[0]
	}
	years := current.Year() - birth.Year()
	if current.Month() < birth.Month() ||
		(current.Month() == birth.Month() && current.Day() < birth.Day()) {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}
