package mapper

import (
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"testing"
	"time"
)

type Person struct {
	Id       int64          `json:"id"`
	Name     string         `json:"name"`
	Birthday types.JSONTime `json:"birthday"`
}

func TestTimeJson(t *testing.T) {
	now := types.JSONTime(time.Now())
	t.Log(now)
	src := `{"id":5,"name":"xiaoming","birthday":"2016-06-30 16:09:51"}`
	p := new(Person)
	err := json.Unmarshal([]byte(src), p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
	t.Log(time.Time(p.Birthday))
	js, _ := json.Marshal(p)
	t.Log(string(js))
}

func TestTimeJson_NewFormat(t *testing.T) {
	now := types.JSONTime(time.Now())
	t.Log(now)
	types.SetTimeJSONFormat("2006-01-02T15:04:05")
	src := `{"id":5,"name":"xiaoming","birthday":"2016-06-30T16:09:51"}`
	p := new(Person)
	err := json.Unmarshal([]byte(src), p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
	t.Log(time.Time(p.Birthday))
	js, _ := json.Marshal(p)
	t.Log(string(js))
}
