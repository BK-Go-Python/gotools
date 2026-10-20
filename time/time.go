// Package time 提供 Python datetime/time 风格的时间工具：时间戳与字符串
// 互转、strftime/strptime、日期加减、月初月末、人类可读的时间描述等。
//
// 包名与标准库 time 相同，同时导入时请给标准库起别名，例如：
//
//	import (
//		gotime "github.com/BK-Go-Python/gotools/time"
//		"time"
//	)
package time

import (
	"fmt"
	stdtime "time"
)

// 常用布局。
const (
	// Layout 是项目最早使用的日期时间布局。
	Layout = "2006-01-02 15:04:05"
	// DateLayout 是仅有日期的布局。
	DateLayout = "2006-01-02"
	// TimeLayout 是仅有时间的布局。
	TimeLayout = "15:04:05"
)

// Now 返回当前本地时间。
func Now() stdtime.Time {
	return stdtime.Now()
}

// Today 返回今天 00:00:00 的本地时间。
func Today() stdtime.Time {
	now := stdtime.Now()
	y, m, d := now.Date()
	return stdtime.Date(y, m, d, 0, 0, 0, 0, stdtime.Local)
}

// Parse 尝试用常见布局解析日期时间字符串。
func Parse(s string) (stdtime.Time, error) {
	layouts := []string{
		Layout,
		DateLayout,
		"2006-01-02 15:04",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		stdtime.RFC3339Nano,
		stdtime.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := stdtime.ParseInLocation(layout, s, stdtime.Local); err == nil {
			return t, nil
		}
	}
	return stdtime.Time{}, fmt.Errorf("time: cannot parse %q with any known layout", s)
}

// ParseIn 按指定 Go 布局解析时间字符串，时区为本地时区。
func ParseIn(layout, s string) (stdtime.Time, error) {
	return stdtime.ParseInLocation(layout, s, stdtime.Local)
}

// ParseTimestamp 把 Unix 秒时间戳转成时间。
func ParseTimestamp(sec int64) stdtime.Time {
	return stdtime.Unix(sec, 0)
}

// Timestamp 返回 Unix 秒时间戳。
func Timestamp(t stdtime.Time) int64 {
	return t.Unix()
}

// TimestampMs 返回 Unix 毫秒时间戳。
func TimestampMs(t stdtime.Time) int64 {
	return t.UnixNano() / int64(stdtime.Millisecond)
}

// FromMillis 把 Unix 毫秒时间戳转成时间。
func FromMillis(ms int64) stdtime.Time {
	return stdtime.Unix(0, ms*int64(stdtime.Millisecond))
}

// Format 按 Layout（2006-01-02 15:04:05）格式化时间。
func Format(t stdtime.Time) string {
	return t.Format(Layout)
}

// Strtimestamp 把 "2006-01-02 15:04:05" 格式的日期转成时间戳（秒）。
//
// 解析失败返回 0；需要错误信息请使用 Parse + Timestamp。
func Strtimestamp(date string) int64 {
	t, err := ParseIn(Layout, date)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// Intimestamp 把时间戳（秒）转成 "2006-01-02 15:04:05" 格式字符串。
func Intimestamp(timestamp int64) string {
	return ParseTimestamp(timestamp).Format(Layout)
}
