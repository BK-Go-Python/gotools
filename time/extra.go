package time

import (
	stdtime "time"
)

// AddWeeks 增加周数，等价于 Python 的 t + timedelta(weeks=n)，
// 即 AddDate(0, 0, 7*n)。
func AddWeeks(t stdtime.Time, weeks int) stdtime.Time {
	return t.AddDate(0, 0, 7*weeks)
}

// Quarter 返回 t 所在的季度，取值为 1..4。
func Quarter(t stdtime.Time) int {
	return (int(t.Month())-1)/3 + 1
}

// StartOfQuarter 返回 t 所在季度第一天 00:00:00，保留 t 的 Location。
func StartOfQuarter(t stdtime.Time) stdtime.Time {
	month := stdtime.Month((Quarter(t)-1)*3 + 1)
	return stdtime.Date(t.Year(), month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfQuarter 返回 t 所在季度最后一天 23:59:59.999999999，
// 保留 t 的 Location；与 Python 的 date/datetime 边界语义一致。
func EndOfQuarter(t stdtime.Time) stdtime.Time {
	month := stdtime.Month(Quarter(t) * 3)
	day := DaysInMonth(t.Year(), month)
	return stdtime.Date(t.Year(), month, day, 23, 59, 59, 999999999, t.Location())
}

// IsWeekend 判断 t 是否为周六或周日。
func IsWeekend(t stdtime.Time) bool {
	weekday := t.Weekday()
	return weekday == stdtime.Saturday || weekday == stdtime.Sunday
}

// WeekdayISO 返回 ISO 8601 星期编号：周一为 1，周日为 7，
// 与 Python datetime.isoweekday() 一致。
func WeekdayISO(t stdtime.Time) int {
	weekday := int(t.Weekday())
	if weekday == 0 {
		return 7
	}
	return weekday
}

// DaysInYear 返回指定年份的天数：闰年 366，平年 365。
func DaysInYear(year int) int {
	if IsLeapYear(year) {
		return 366
	}
	return 365
}
