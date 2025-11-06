package jsonutils

import (
	"testing"
	"time"

	times2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
)

func Test_DecodeWithTimeConversion(t *testing.T) {
	text := `{
		"time":"2006-01-02 15:04:05",
		"data":{
			"created":"2006-01-02 15:04:05"
		}, 
		"list":[
			{"time":"2006-01-02 15:04:05"}, 
			{"time":"2006-01-02 15:04:05"}
		]
	}`
	// 定义多级 TimeFields 结构
	timeFields := map[string]any{
		"time": false,
		"data": map[string]any{
			"created": true,
		},
		"list": map[string]any{
			"time": true,
		},
	}
	const dateFormat = "2006-01-02 15:04:05"
	data, err := UnmarshalTime([]byte(text), &UnmarshalTimeOptions{
		TimeFields: timeFields,
		ParseTime: func(val string, key any) (timeVal any, err error) {
			if val == "" || val == "null" {
				return nil, nil
			}
			tm, er := time.Parse(dateFormat, val)
			if er != nil {
				return nil, er
			}
			b, ok := key.(bool)
			if b || !ok {
				timeVal = times2.GetTime(&tm)
			} else {
				timeVal = times2.GetDate(&tm)
			}
			return timeVal, err
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}
