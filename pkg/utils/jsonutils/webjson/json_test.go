package webjson

import (
	"testing"
	"time"

	times2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
)

type Example struct {
	Time1 time.Time  `json:"dateTime"`
	Time2 *time.Time `json:"pTime"`
}

func Test_JsonMarshal(t *testing.T) {
	dateTime := time.Now()
	example := &Example{
		Time1: dateTime,
		Time2: &dateTime,
	}
	data, e := Marshal(example)
	if e != nil {
		t.Error(e)
	} else {
		t.Log(data)
		var e Example
		err := Unmarshal([]byte(data), &e)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("time1:", e.Time1)
			t.Log("time2:", e.Time2)
		}
	}
}

type Example2 struct {
	Time1 times2.Time  `json:"time1"`
	Time2 *times2.Time `json:"time2"`
}

func Test_JsonMarshal2(t *testing.T) {
	dateTime := times2.Now()
	example := &Example2{
		Time1: dateTime,
		Time2: &dateTime,
	}
	data, e := Marshal(example)
	if e != nil {
		t.Error(e)
	} else {
		t.Log(data)
		var e Example2
		err := Unmarshal([]byte(data), &e)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("time1:", e.Time1)
			t.Log("time2:", e.Time2)
		}
	}
}

type DateExample struct {
	Date1 times2.Date  `json:"date1"`
	Date2 *times2.Date `json:"data2"`
	Date3 *times2.Date `json:"data3"`
}

func Test_JsonMarshal3(t *testing.T) {
	dateTime := times2.NewDate()
	example := &DateExample{
		Date1: *dateTime,
		Date2: dateTime,
	}
	data, e := Marshal(example)
	if e != nil {
		t.Error(e)
	} else {
		t.Log(data)
		var e DateExample
		err := Unmarshal([]byte(data), &e)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Date1:", e.Date1)
			t.Log("Date2:", e.Date2)
			t.Log("Date3:", e.Date3)
		}
	}
}

func Test_JsonMarshalMap(t *testing.T) {
	dateTime := times2.NewDate()
	timeVal := time.Now()
	example := map[string]any{
		"date1": *dateTime,
		"date2": dateTime,
		"time1": timeVal,
	}
	data, e := Marshal(example)
	if e != nil {
		t.Error(e)
	} else {
		t.Log(data)
		var e DateExample
		err := Unmarshal([]byte(data), &e)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Date1:", e.Date1)
			t.Log("Date2:", e.Date2)
			t.Log("Date3:", e.Date3)
		}
	}
}
