package server

import "time"

type Time1 struct {
}

func NewTime() *Time1 {
	return &Time1{}
}

func (t *Time1) Now() time.Time {
	return time.Now()
}
