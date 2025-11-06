package times

import (
	"encoding/json"
	"testing"
)

type Entity struct {
	Date  *Date `json:"date"`
	Time  *Time `json:"time"`
	Time2 *Time `json:"time2"`
}

func Test_Unmarshal(t *testing.T) {
	jsonText := `{
		"date": "2022-01-01",
		"time": "2012-09-10 19:20:00",
		"time2": "2012-09-10T19:20:00Z"
	}`
	var entity Entity
	if err := json.Unmarshal([]byte(jsonText), &entity); err != nil {
		t.Error(err)
	} else {
		t.Log(entity.Date)
		t.Log(entity.Time)
		t.Log(entity.Time2)
	}
}
