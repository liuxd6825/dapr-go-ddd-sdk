package server

import "time"

type Time struct {
}

func NewTime() *Time {
	return &Time{}
}

func (t *Time) Now() time.Time {
	return time.Now()
}
