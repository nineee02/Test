package constant

import "time"

func GetCurrentDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func GetCurrentDay(curDate, startDate time.Time) int32 {
	days := curDate.Sub(startDate).Hours() / 24
	return int32(days) + 1
}
