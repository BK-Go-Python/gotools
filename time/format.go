package time

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	stdtime "time"
)

// Strftime 等价于 Python 的 datetime.strftime(format)，支持常用的
// %-指令；不认识的指令会原样输出。
//
// 支持：%Y %y %m %d %e %H %I %M %S %f %p %a %A %b %B %z %Z %j
//
//	%U %W %w %u %G %V %c %x %X %s %n %t %%
func Strftime(t stdtime.Time, format string) string {
	var b strings.Builder
	b.Grow(len(format) + 16)
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			b.WriteByte(format[i])
			continue
		}
		if i+1 >= len(format) {
			b.WriteByte('%')
			break
		}
		i++
		switch format[i] {
		case 'Y':
			fmt.Fprintf(&b, "%04d", t.Year())
		case 'y':
			fmt.Fprintf(&b, "%02d", t.Year()%100)
		case 'm':
			fmt.Fprintf(&b, "%02d", int(t.Month()))
		case 'd':
			fmt.Fprintf(&b, "%02d", t.Day())
		case 'e':
			fmt.Fprintf(&b, "%2d", t.Day())
		case 'H':
			fmt.Fprintf(&b, "%02d", t.Hour())
		case 'I':
			hour := t.Hour() % 12
			if hour == 0 {
				hour = 12
			}
			fmt.Fprintf(&b, "%02d", hour)
		case 'M':
			fmt.Fprintf(&b, "%02d", t.Minute())
		case 'S':
			fmt.Fprintf(&b, "%02d", t.Second())
		case 'f':
			fmt.Fprintf(&b, "%06d", t.Nanosecond()/1000)
		case 'p':
			if t.Hour() < 12 {
				b.WriteString("AM")
			} else {
				b.WriteString("PM")
			}
		case 'a':
			b.WriteString(t.Format("Mon"))
		case 'A':
			b.WriteString(t.Format("Monday"))
		case 'b':
			b.WriteString(t.Format("Jan"))
		case 'B':
			b.WriteString(t.Format("January"))
		case 'z':
			b.WriteString(t.Format("-0700"))
		case 'Z':
			b.WriteString(t.Format("MST"))
		case 'j':
			fmt.Fprintf(&b, "%03d", t.YearDay())
		case 'U':
			yday := t.YearDay() - 1
			fmt.Fprintf(&b, "%02d", (yday+7-int(t.Weekday()))/7)
		case 'W':
			yday := t.YearDay() - 1
			mondayBased := (int(t.Weekday()) + 6) % 7
			fmt.Fprintf(&b, "%02d", (yday+7-mondayBased)/7)
		case 'w':
			fmt.Fprintf(&b, "%d", int(t.Weekday()))
		case 'u':
			weekday := int(t.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			fmt.Fprintf(&b, "%d", weekday)
		case 'G':
			isoYear, _ := t.ISOWeek()
			fmt.Fprintf(&b, "%04d", isoYear)
		case 'V':
			_, isoWeek := t.ISOWeek()
			fmt.Fprintf(&b, "%02d", isoWeek)
		case 'c':
			b.WriteString(t.Format("Mon Jan _2 15:04:05 2006"))
		case 'x':
			b.WriteString(t.Format("01/02/06"))
		case 'X':
			b.WriteString(t.Format("15:04:05"))
		case 's':
			b.WriteString(strconv.FormatInt(t.Unix(), 10))
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case '%':
			b.WriteByte('%')
		default:
			b.WriteByte('%')
			b.WriteByte(format[i])
		}
	}
	return b.String()
}

// Strptime 等价于 Python 的 datetime.strptime(s, format)。
//
// 支持与 Strftime 相同的常用指令。%f 位于秒之后并带有 "." 或 "," 分隔符时
// 接受 1~6 位小数；字面百分号前带数字（如 "100%%"）时会先剥离该字面片段。
// %j、%U、%W、%V、%G、%w、%u、%s 暂不支持解析。
func Strptime(s, format string) (stdtime.Time, error) {
	if strings.Contains(format, "%f") && fractionalDigitsTooLong(s) {
		return stdtime.Time{}, fmt.Errorf("time: fractional seconds must have at most 6 digits")
	}
	if format == "%s" {
		sec, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			return stdtime.Time{}, fmt.Errorf("time: invalid unix timestamp %q: %v", s, err)
		}
		return ParseTimestamp(sec), nil
	}
	s, format = stripLiteralPercentDigits(s, format)
	layout, err := ToGoLayout(format)
	if err != nil {
		return stdtime.Time{}, err
	}
	return stdtime.ParseInLocation(layout, s, stdtime.Local)
}

// ToGoLayout 把 Python strftime 风格的格式串转换为 Go 的布局字符串。
func ToGoLayout(format string) (string, error) {
	var b strings.Builder
	b.Grow(len(format) + 8)
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			b.WriteByte(format[i])
			continue
		}
		if i+1 >= len(format) {
			b.WriteByte('%')
			break
		}
		i++
		switch format[i] {
		case 'Y':
			b.WriteString("2006")
		case 'y':
			b.WriteString("06")
		case 'm':
			b.WriteString("01")
		case 'd':
			b.WriteString("02")
		case 'e':
			b.WriteString("_2")
		case 'H':
			b.WriteString("15")
		case 'I':
			b.WriteString("03")
		case 'M':
			b.WriteString("04")
		case 'S':
			b.WriteString("05")
		case 'f':
			// Go 的 "999999" 布局在解析时接受 1~9 位小数，比固定 6 位更接近
			// Python strptime 的 %f。
			b.WriteString("999999")
		case 'p':
			b.WriteString("PM")
		case 'a':
			b.WriteString("Mon")
		case 'A':
			b.WriteString("Monday")
		case 'b':
			b.WriteString("Jan")
		case 'B':
			b.WriteString("January")
		case 'z':
			b.WriteString("-0700")
		case 'Z':
			b.WriteString("MST")
		case 'c':
			b.WriteString("Mon Jan _2 15:04:05 2006")
		case 'x':
			b.WriteString("01/02/06")
		case 'X':
			b.WriteString("15:04:05")
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case '%':
			b.WriteByte('%')
		default:
			return "", fmt.Errorf("time: unsupported directive %%%c in layout %q", format[i], format)
		}
	}
	return b.String(), nil
}

// stripLiteralPercentDigits 处理 "100%%" 这类字面百分号前带数字的情况。
// Go 的布局解析器会把 "100%" 中的数字当成日期字段，导致无法解析；
// 这里把该字面片段从格式与输入中同时删除，再交给 time.Parse。
func stripLiteralPercentDigits(s, format string) (string, string) {
	for {
		idx := strings.Index(format, "%%")
		if idx < 0 {
			return s, format
		}
		prev := strings.LastIndex(format[:idx], "%")
		litStart := 0
		if prev >= 0 {
			litStart = prev + 2
		}
		literal := format[litStart:idx]
		if literal == "" || !strings.ContainsAny(literal, "0123456789") {
			return s, format
		}
		marker := literal + "%"
		pos := strings.LastIndex(s, marker)
		if pos < 0 {
			return s, format
		}
		s = s[:pos] + s[pos+len(marker):]
		format = format[:litStart] + format[idx+2:]
	}
}

var fractionalPattern = regexp.MustCompile(`[.,](\d+)`)

func fractionalDigitsTooLong(s string) bool {
	for _, match := range fractionalPattern.FindAllStringSubmatch(s, -1) {
		if len(match[1]) > 6 {
			return true
		}
	}
	return false
}
