package events

import "testing"

func Test_StructToMap(t *testing.T) {
	data := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "张三",
		Age:  18,
	}

	mapData, err := StructToMap(data)
	if err != nil {
		t.Errorf("StructToMap() error = %v", err)
		return
	}
	t.Logf("mapData = %v", mapData)
}

func Test_LoadEvent(t *testing.T) {

}
