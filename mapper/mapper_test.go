package mapper

import (
	"testing"
	"time"
)

type Data struct {
	Name  string
	Time  time.Time
	PTime *time.Time
}

func Test_newMap(t *testing.T) {
	dateTime := time.Now()
	data := Data{
		Name:  "Data",
		Time:  dateTime,
		PTime: &dateTime,
	}

	if dataMap, err := NewMap(data); err != nil {
		println(err)
	} else {
		v := dataMap["Time"]
		println(v)
		println(dataMap)
	}
}
