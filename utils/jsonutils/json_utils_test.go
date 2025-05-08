package jsonutils

import "testing"

func TestMarshals(t *testing.T) {

	t.Run("MarshalNoKeyMarks", func(t *testing.T) {
		data := make(map[string]interface{})
		data["name"] = "name"
		data["max"] = 10
		data["create"] = nil
		json, err := MarshalNoKeyMarks(data)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(json)
		}
	})

}
