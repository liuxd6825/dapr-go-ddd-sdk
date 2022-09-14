package timeutils

import "time"

func Now() time.Time {
	t := time.Now()
	v := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return v
}
