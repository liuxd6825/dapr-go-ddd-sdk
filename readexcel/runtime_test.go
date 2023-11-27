package readexcel

import "testing"

func TestField_toDateTime(t *testing.T) {
	date := "20110812"
	time := "11.56.23"
	res := toDateTime(date, time)
	println(res)
}
