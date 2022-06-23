package ddd

import (
	"encoding/json"
	"github.com/google/uuid"
	"testing"
	"time"
)

type testAgg struct {
	Id            string
	AggregateType string
	TenantId      string
}

func (t *testAgg) GetTenantId() string {
	return t.TenantId
}

func (t *testAgg) GetAggregateId() string {
	return t.Id
}

func (t *testAgg) GetAggregateType() string {
	return t.AggregateType
}

func (t *testAgg) GetAggregateVersion() string {
	return "v1.0"
}

type testEvent struct {
	Data testFields
}

func (t *testEvent) GetTenantId() string {
	return "001"
}

func (t *testEvent) GetCommandId() string {
	return "commandId001"
}

func (t *testEvent) GetEventId() string {
	return uuid.New().String()
}

func (t *testEvent) GetEventType() string {
	return "testEvent"
}

func (t *testEvent) GetEventVersion() string {
	return "v1.0"
}

func (t *testEvent) GetAggregateId() string {
	return t.Data.Id
}

func (t *testEvent) GetCreatedTime() time.Time {
	return time.Now()
}

func (t *testEvent) GetData() interface{} {
	return t.Data
}

type testFields struct {
	Id     string `json:"id" ddd-rel:""`
	UserId string `json:"userId" ddd-rel:""`
	SysId  string `json:"sysId" ddd-rel:""`
	NilId  string `json:"nilId" `
}

func Test_GetRelationTags(t *testing.T) {
	data := testFields{
		Id:     "1",
		UserId: "userId-001",
		SysId:  "sysId-001",
	}
	// agg := testAgg{Id: "aggId", AggregateType: "TestAggType", TenantId: "001"}
	event := testEvent{Data: data}
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
