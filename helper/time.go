package helper

import (
	"time"

	"github.com/tanenking/svrframe/constants"
)

const Second_one_day = uint32(24 * 3600)
const Millisecond_one_day = int64(Second_one_day) * 1000
const Second_one_hour = uint32(3600)
const Second_one_week = Second_one_day * 7

func GetNowTime() time.Time {
	local_milli := local_timestamp_milli + (time.Now().UnixMilli() - local_timestamp_milli_begin)
	return time.UnixMilli(local_milli)
}
func GetNowTimestamp() uint32 {
	return uint32(GetNowTime().Unix())
}
func GetNowTimestampMilli() int64 {
	return GetNowTime().UnixMilli()
}

func ModifyTime(tMilli int64) time.Time {
	local_timestamp_milli = tMilli
	local_timestamp_milli_begin = time.Now().UnixMilli()

	return GetNowTime()
}

func StrToUTCTimestamp(datetime string) uint32 {
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(constants.TimeFormatString, datetime, loc)
	if err != nil {
		return 0
	}
	unixTime := uint32(theTime.Unix())
	return unixTime
}

// 获取指定时间的时间戳
func GetAppointDate(year int, month time.Month, day, hour, min, sec int, local *time.Location) time.Time {
	loc := local
	if loc == nil {
		loc, _ = time.LoadLocation("Local")
	}
	t := time.Date(year, month, day, hour, min, sec, 0, loc)
	return t
}

// 获取当天0点时间戳
func GetToday0ClockTimestamp() uint32 {
	currentTime := GetNowTime()
	t := GetAppointDate(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, currentTime.Location()).Unix()
	return uint32(t)
}
