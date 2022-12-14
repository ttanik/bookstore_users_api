package date_utils

import (
	"time"
)

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDBLayout   = "2006-01-02 15:04:05"
)

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}
func GetNow() time.Time {
	return time.Now().UTC()
}
func GetNowDBFormat() string {
	return GetNow().Format(apiDBLayout)
}
