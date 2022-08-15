package ddd

import (
	"encoding/json"
	"github.com/google/uuid"
	"testing"
	"time"
)

func Test_GetRelationTags(t *testing.T) {
	data := fields{
		Id:     "1",
		UserId: "userId-001",
		SysId:  "sysId-001",
	}
	// agg := testAgg{Id: "aggId", AggregateType: "TestAggType", TenantId: "001"}
	event := one{Data: data}
	rel, ok, err := GetRelation(event.Data)
	if err != nil {
		t.Error(err)
	}
	if ok {
		if bs, err := json.Marshal(rel); err != nil {
			t.Error(err)
		} else {
			print(string(bs))
		}
	}
}

func Test_GetRelations_Items(t *testing.T) {
	f := fields{
		Id:     "1",
		UserId: "userId-001",
		SysId:  "sysId-001",
	}
	data := []*fields{&f}
	items := items{
		Data: data,
	}
	rel, ok, err := GetRelation(items.Data)
	if err != nil {
		t.Error(err)
	}
	if ok {
		if bs, err := json.Marshal(rel); err != nil {
			t.Error(err)
		} else {
			print(string(bs))
		}
	}
}

type fields struct {
	Id     string `json:"id" ddd-rel:""`
	UserId string `json:"userId" ddd-rel:"userId"`
	SysId  string `json:"sysId" ddd-rel:"sysId"`
	NilId  string `json:"nilId" `
}

type items struct {
	Data []*fields
}

type one struct {
	Data fields
}

func (t *items) GetTenantId() string {
	return "001"
}

func (t *items) GetCommandId() string {
	return "commandId001"
}

func (t *items) GetEventId() string {
	return uuid.New().String()
}

func (t *items) GetEventType() string {
	return "one"
}

func (t *items) GetEventVersion() string {
	return "v1.0"
}

func (t *items) GetAggregateId() string {
	return ""
}

func (t *items) GetCreatedTime() time.Time {
	return time.Now()
}

func (t *items) GetData() interface{} {
	return t.Data
}

func (t *one) GetTenantId() string {
	return "001"
}

func (t *one) GetCommandId() string {
	return "commandId001"
}

func (t *one) GetEventId() string {
	return uuid.New().String()
}

func (t *one) GetEventType() string {
	return "one"
}

func (t *one) GetEventVersion() string {
	return "v1.0"
}

func (t *one) GetAggregateId() string {
	return t.Data.Id
}

func (t *one) GetCreatedTime() time.Time {
	return time.Now()
}

func (t *one) GetData() interface{} {
	return t.Data
}
