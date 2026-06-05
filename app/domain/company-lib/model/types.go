package model

import (
	"encoding/json"
	"time"
)

type ESTime struct {
	time.Time
}

const esTimeFormat = "2006-01-02T15:04:05.000000"
const esTimeFormatShort = "2006-01-02T15:04:05"

func (t *ESTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = s[1 : len(s)-1]

	var err error
	t.Time, err = time.Parse(esTimeFormat, s)
	if err != nil {
		t.Time, err = time.Parse(esTimeFormatShort, s)
	}
	return err
}

type RawJSON struct {
	json.RawMessage
}

func (r *RawJSON) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	r.RawMessage = data
	return nil
}

type JSONString struct {
	String string
}

func (j *JSONString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		return json.Unmarshal(data, &j.String)
	}
	j.String = string(data)
	return nil
}