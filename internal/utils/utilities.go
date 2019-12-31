package utils

import (
	"strconv"
	"time"
)

func FormatDate(t time.Time) int {
	timeString := t.Format("060102")
	time, _ := strconv.Atoi(timeString)
	return time
}

