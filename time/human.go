package time

import (
	"fmt"
	stdtime "time"
)

// Humanize 返回中文的相对时间描述，例如 "刚刚"、"3分钟前"、"2天前"、"1年前"。
func Humanize(t stdtime.Time) string {
	return HumanizeAt(t, stdtime.Now())
}

// HumanizeAt 以 now 为基准返回中文相对时间描述，便于测试或指定基准时间。
func HumanizeAt(t, now stdtime.Time) string {
	delta := now.Sub(t)
	future := delta < 0
	if future {
		delta = -delta
	}
	suffix := "前"
	if future {
		suffix = "后"
	}

	switch {
	case delta < stdtime.Minute:
		if future {
			return "即将"
		}
		return "刚刚"
	case delta < stdtime.Hour:
		return fmt.Sprintf("%d分钟%s", int(delta/stdtime.Minute), suffix)
	case delta < 24*stdtime.Hour:
		return fmt.Sprintf("%d小时%s", int(delta/stdtime.Hour), suffix)
	case delta < 7*24*stdtime.Hour:
		return fmt.Sprintf("%d天%s", int(delta/(24*stdtime.Hour)), suffix)
	case delta < 30*24*stdtime.Hour:
		return fmt.Sprintf("%d周%s", int(delta/(7*24*stdtime.Hour)), suffix)
	case delta < 365*24*stdtime.Hour:
		return fmt.Sprintf("%d个月%s", int(delta/(30*24*stdtime.Hour)), suffix)
	default:
		return fmt.Sprintf("%d年%s", int(delta/(365*24*stdtime.Hour)), suffix)
	}
}

// DurationHuman 把时长格式化为中文，例如 "1天2小时3分钟"、"45秒"。
func DurationHuman(d stdtime.Duration) string {
	if d == 0 {
		return "0秒"
	}
	prefix := ""
	if d < 0 {
		prefix = "-"
		d = -d
	}
	if d < stdtime.Second {
		return fmt.Sprintf("%s%d毫秒", prefix, d/stdtime.Millisecond)
	}

	days := int(d / (24 * stdtime.Hour))
	d -= stdtime.Duration(days) * 24 * stdtime.Hour
	hours := int(d / stdtime.Hour)
	d -= stdtime.Duration(hours) * stdtime.Hour
	minutes := int(d / stdtime.Minute)
	d -= stdtime.Duration(minutes) * stdtime.Minute
	seconds := int(d / stdtime.Second)

	out := prefix
	if days > 0 {
		out += fmt.Sprintf("%d天", days)
	}
	if hours > 0 {
		out += fmt.Sprintf("%d小时", hours)
	}
	if minutes > 0 {
		out += fmt.Sprintf("%d分钟", minutes)
	}
	if seconds > 0 || out == prefix {
		out += fmt.Sprintf("%d秒", seconds)
	}
	return out
}
