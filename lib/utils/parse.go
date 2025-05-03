package utils

import (
	"fmt"
	"time"
)

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04:05.000", timeStr)
}

func ParseDuration(timeStr string, layout string) (time.Duration, error) {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid time format", err)
	}
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond())*time.Nanosecond, nil
}

func FormatDurationToTime(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
