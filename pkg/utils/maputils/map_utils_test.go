package maputils

import (
	"testing"
	"time"
)

type userView struct {
	TenantId    string
	Remarks     string
	CreatedTime time.Time
	UpdatedTime *time.Time
}

type Date struct {
	Value time.Time
}

type Int64 struct {
	Value int64
}

func TestMapStructure_Decode(t *testing.T) {
	user := &userView{}
	props := map[string]interface{}{
		"graphId":     "graphId",
		"tenantId":    "tenantId",
		"remarks":     "remarks",
		"createdTime": "2022-08-31 13:47:36",
		"updatedTime": "2022-08-31T13:47:36.255551+08:00",
	}

	if err := Decode(props, user); err != nil {
		t.Error(err)
		return
	}
	t.Logf("tenantId = %v", user.TenantId)
	t.Logf("remarks = %v", user.Remarks)
	t.Logf("createdTime = %v", user.CreatedTime)
	t.Logf("updatedTime = %v", user.UpdatedTime)
}

func TestMapStructure_NewMap(t *testing.T) {
	vTime := time.Now()
	date := &Date{
		Value: vTime,
	}
	mapData, err := NewMap(date)
	if err != nil {
		t.Error(err)
		return
	}
	if v, ok := mapData["Value"]; ok {
		if _, ok := v.(time.Time); !ok {
			t.Error("mapData.Value is not time.Time")
		}
	}
	t.Log(mapData)
}

func TestMapStructure_NewMap_Int64(t *testing.T) {
	data := Int64{
		Value: 10,
	}
	mapData, err := NewMap(&data)
	if err != nil {
		t.Error(err)
		return
	}
	if v, ok := mapData["Value"]; ok {
		if _, ok := v.(int64); !ok {
			t.Error("mapData.Value is not time.Time")
		}
	}
	t.Log(mapData)
}
