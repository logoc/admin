package gf

import (
	"fmt"
	"time"
)

func ParseDate(input string) (time.Time, error) {
	// 定义可能的日期格式
	formats := []string{
		"20060102",
		"2006-01-02",
		"2006/01/02",
		"01/02/06",
		"01-02-06",
	}

	var parsedDate time.Time
	var err error

	// 尝试按顺序解析每种格式
	for _, format := range formats {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		parsedDate, err = time.ParseInLocation(format, input, loc)
		if err == nil {
			return parsedDate, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse date: %s", input)
}

func StringTimestamp2(timeLayout string) int64 {
	timetpl := "2006/1/2"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	times, _ := time.ParseInLocation(timetpl, timeLayout, loc)
	timeUnix := times.Unix()
	return timeUnix
}

// 日期时间转时间戳
// timetype时间格式类型  datetime=日期时间 datesecond=日期时间秒date=日期
func StringTimestamp(timeLayout string, timetype string) int64 {
	timetpl := "2006-01-02 15:04:05"
	if timetype == "date" {
		timetpl = "2006-01-02"
	} else if timetype == "datetime" {
		timetpl = "2006-01-02 15:04"
	}
	times, _ := time.ParseInLocation(timetpl, timeLayout, time.Local)
	timeUnix := times.Unix()
	return timeUnix
}

// 时间戳格式化为日期字符串
// timetype时间格式类型 date=日期 datetime=日期时间 datesecond=日期时间秒
func TimestampString(timedata interface{}, timetype string) string {
	timetpl := "2006-01-02 15:04:05"
	if timetype == "date" {
		timetpl = "2006-01-02"
	} else if timetype == "datetime" {
		timetpl = "2006-01-02 15:04"
	}
	return time.Unix(timedata.(int64), 0).Format(timetpl)
}

// 获取当前时间戳
func NowTimestamp() int64 {
	return time.Now().Unix()
}

// 获取当前时间戳
func NowBeforeTimestamp(day int64) int64 {
	t := time.Now().Add(time.Duration(-day*24) * time.Hour)
	return t.Unix()
}
