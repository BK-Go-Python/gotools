package time

import (
	"testing"
	stdtime "time"
)

// TestHumanizeAt 使用固定基准时间，避免依赖 time.Now。
// 期望值依据 HumanizeAt 的分档规则：小于 1 分钟、整除法向下取整。
func TestHumanizeAt(t *testing.T) {
	now := cstDate(2024, 1, 15, 12, 0, 0, 0)
	tests := []struct {
		name string
		t    stdtime.Time
		want string
	}{
		{"same instant", now, "刚刚"},
		{"past one nanosecond", now.Add(-stdtime.Nanosecond), "刚刚"},
		{"past 59 seconds", now.Add(-59 * stdtime.Second), "刚刚"},
		{"past one minute", now.Add(-stdtime.Minute), "1分钟前"},
		{"past 59m59s", now.Add(-(59*stdtime.Minute + 59*stdtime.Second)), "59分钟前"},
		{"past one hour", now.Add(-stdtime.Hour), "1小时前"},
		{"past 90 minutes", now.Add(-90 * stdtime.Minute), "1小时前"},
		{"past 23h59m59s", now.Add(-(23*stdtime.Hour + 59*stdtime.Minute + 59*stdtime.Second)), "23小时前"},
		{"past 24 hours", now.Add(-24 * stdtime.Hour), "1天前"},
		{"past 6 days", now.Add(-6 * 24 * stdtime.Hour), "6天前"},
		{"past 6d23h59m", now.Add(-(6*24*stdtime.Hour + 23*stdtime.Hour + 59*stdtime.Minute)), "6天前"},
		{"past 7 days", now.Add(-7 * 24 * stdtime.Hour), "1周前"},
		{"past 29 days", now.Add(-29 * 24 * stdtime.Hour), "4周前"},
		{"past 29d23h59m", now.Add(-(29*24*stdtime.Hour + 23*stdtime.Hour + 59*stdtime.Minute)), "4周前"},
		{"past 30 days", now.Add(-30 * 24 * stdtime.Hour), "1个月前"},
		{"past 31 days", now.Add(-31 * 24 * stdtime.Hour), "1个月前"},
		{"past 364 days", now.Add(-364 * 24 * stdtime.Hour), "12个月前"},
		{"past 364d23h59m", now.Add(-(364*24*stdtime.Hour + 23*stdtime.Hour + 59*stdtime.Minute)), "12个月前"},
		{"past 365 days", now.Add(-365 * 24 * stdtime.Hour), "1年前"},
		{"past 400 days", now.Add(-400 * 24 * stdtime.Hour), "1年前"},
		{"past 730 days", now.Add(-730 * 24 * stdtime.Hour), "2年前"},

		{"future one nanosecond", now.Add(stdtime.Nanosecond), "即将"},
		{"future 59 seconds", now.Add(59 * stdtime.Second), "即将"},
		{"future one minute", now.Add(stdtime.Minute), "1分钟后"},
		{"future two hours", now.Add(2 * stdtime.Hour), "2小时后"},
		{"future 3 days", now.Add(3 * 24 * stdtime.Hour), "3天后"},
		{"future 14 days", now.Add(14 * 24 * stdtime.Hour), "2周后"},
		{"future 29 days", now.Add(29 * 24 * stdtime.Hour), "4周后"},
		{"future 30 days", now.Add(30 * 24 * stdtime.Hour), "1个月后"},
		{"future 60 days", now.Add(60 * 24 * stdtime.Hour), "2个月后"},
		{"future 364 days", now.Add(364 * 24 * stdtime.Hour), "12个月后"},
		{"future 365 days", now.Add(365 * 24 * stdtime.Hour), "1年后"},
		{"future 366 days", now.Add(366 * 24 * stdtime.Hour), "1年后"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HumanizeAt(tt.t, now); got != tt.want {
				t.Errorf("HumanizeAt(%v, %v) = %q, want %q", tt.t, now, got, tt.want)
			}
		})
	}
}

func TestDurationHuman(t *testing.T) {
	tests := []struct {
		name string
		d    stdtime.Duration
		want string
	}{
		{"zero", 0, "0秒"},
		{"one millisecond", stdtime.Millisecond, "1毫秒"},
		{"999 milliseconds", 999 * stdtime.Millisecond, "999毫秒"},
		{"one second", stdtime.Second, "1秒"},
		{"45 seconds", 45 * stdtime.Second, "45秒"},
		{"one minute", stdtime.Minute, "1分钟"},
		{"90 seconds", 90 * stdtime.Second, "1分钟30秒"},
		{"one hour", stdtime.Hour, "1小时"},
		{"hour minute second", stdtime.Hour + 2*stdtime.Minute + 3*stdtime.Second, "1小时2分钟3秒"},
		{"one day", 24 * stdtime.Hour, "1天"},
		{"one day two hours three minutes", 24*stdtime.Hour + 2*stdtime.Hour + 3*stdtime.Minute, "1天2小时3分钟"},
		{"day hour minute second", 25*stdtime.Hour + 2*stdtime.Minute + 3*stdtime.Second, "1天1小时2分钟3秒"},
		{"two days", 48 * stdtime.Hour, "2天"},
		{"seven days", 7 * 24 * stdtime.Hour, "7天"},
		{"one second less one nanosecond", stdtime.Second - stdtime.Nanosecond, "999毫秒"},
		{"one hour less one nanosecond", stdtime.Hour - stdtime.Nanosecond, "59分钟59秒"},
		{"1500 milliseconds truncates", 1500 * stdtime.Millisecond, "1秒"},

		{"negative 500 milliseconds", -500 * stdtime.Millisecond, "-500毫秒"},
		{"negative 45 seconds", -45 * stdtime.Second, "-45秒"},
		{"negative composite", -(stdtime.Hour + 2*stdtime.Minute + 3*stdtime.Second), "-1小时2分钟3秒"},
		{"negative day", -25 * stdtime.Hour, "-1天1小时"},
		{"negative seven days", -7 * 24 * stdtime.Hour, "-7天"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationHuman(tt.d); got != tt.want {
				t.Errorf("DurationHuman(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}
