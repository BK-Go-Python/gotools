package time

import (
	"time"
)

// "2006-01-02 15:04:05"格式日期转化为时间戳
func Strtimestamp(date string) int64{

	// 日期转化为时间戳
	timeLayout := "2006-01-02 15:04:05"  	//转化所需模板 这个时间为固定时间
	loc, _ := time.LoadLocation("Local")    //获取时区

	//调用转化方法，传入上面准备好的参数
	tmp, _ := time.ParseInLocation(timeLayout, date, loc)
	timestamp := tmp.Unix()    				//转化为时间戳（秒级） 类型是int64

	return timestamp

}

// 时间戳转化为日期格式"2006-01-02 15:04:05"字符串
func Intimestamp(timestamp int64) string{

	ti := time.Unix(timestamp, 0)
	timestr := ti.Format("2006-01-02 15:04:05")
	return timestr
}
