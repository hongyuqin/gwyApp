package util

import "time"

//获取当前日期的下一天零点
func GetNextDayBegin() int64 {
	t := time.Now().AddDate(0, 0, 1)
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return tm1.Unix()
}
