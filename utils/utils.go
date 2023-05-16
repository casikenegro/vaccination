package utils

import "time"

const Layout = "2006-01-02"

func ParseStringToTime(date string) (time.Time, error) {
	return time.Parse(Layout, date)
}
