package xbase

import (
	"context"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
	"github.com/stretchr/testify/assert"
)

type UnitTestEvent = events.Event[*UnitTestData]
type UnitTestData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

const UnitTestEventType = "test.unit.test-event"

func NewUnitTestEvent(ctx context.Context, data *UnitTestData) *UnitTestEvent {
	event := &UnitTestEvent{}
	event.SetData(ctx, "", data, &events.EventOptions{EventType: UnitTestEventType})
	return event
}
func Test_PublishEvent(t *testing.T) {
	xtest2.InitEnv_MongoLocal(xtest2.NewMongoOptions().SetDBName("master"))
	ctx := xtest2.NewContext()
	event := NewUnitTestEvent(ctx, &UnitTestData{
		Id:   "1",
		Name: "张三",
		Age:  18,
	})
	event.Meta = map[string]any{
		"key": "value",
	}
	err := PublishEvent(ctx, event)
	assert.Nil(t, err)
}
