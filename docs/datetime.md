# datetime：对象式日期时间

`gotools/datetime` 提供 `Date`、`DateTime`、`Timedelta` 三个类型，和 `gotools/time` 的函数式 API 互补。

```go
d, err := datetime.NewDate(2024, time.February, 29)
fmt.Println(d.ISOFormat())        // 2024-02-29
next := d.AddDays(1)              // 2024-03-01
fmt.Println(next.Sub(d))          // 1

now := datetime.Now()
fmt.Println(now.Strftime("%Y-%m-%d %H:%M"))
delta := datetime.NewTimedelta(36 * time.Hour)
fmt.Println(delta.Days())         // 1.5
fmt.Println(now.Add(delta))
```

## API

- `Date`：`NewDate`、`DateFromTime`、`Year/Month/Day`、`Weekday`、`ISOFormat`、`Format`、`AddDays`、`Sub`、`Before/After/Equal`、`String`
- `DateTime`：`Now`、`Today`、`FromTimestamp`、`Parse`、`DateTimeFromTime`、`Timestamp`、`ISOFormat`、`Strftime`、`Format`、`Date`、`Add`、`Sub`、`AddDays/AddMonths/AddYears`、`Before/After/Equal`
- `Timedelta`：`NewTimedelta`、`Duration`、`Days`、`TotalSeconds`、`TotalMilliseconds`、`Add/Sub/Neg/Abs`、`String`

## DIFF

- 没有 Python 的 `datetime` 时区对象（`timezone`/`tzinfo`）；用 Go `time.Location` 表达。
- `DateTime` 的 `Sub` 返回 `Timedelta`；`Date` 的 `Sub` 返回整数天数。
- `ISOFormat` 不自动附带时区偏移；需要时用 `Format("2006-01-02T15:04:05-07:00")`。
- 未提供 `combine`、`replace`、`astimezone`、`isoformat(timespec=...)` 等长尾方法。
