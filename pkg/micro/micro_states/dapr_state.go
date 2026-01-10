package micro_states

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/jsonutils"
)

type DaprState struct {
	storeName  string
	daprClient dapr.Client
}

func NewDaprState() *DaprState {
	return &DaprState{
		storeName:  "state",
		daprClient: dapr.GetDaprClient(),
	}
}

func (c *DaprState) GetState(ctx context.Context, key string, meta ...StateMeta) (bool, any, error) {
	metaMap := NewStateMeta(meta...).ToMap()
	item, err := c.daprClient.GetState(ctx, c.storeName, key, metaMap)
	if err != nil {
		return false, nil, err
	}
	if item == nil {
		return false, nil, nil
	}
	if item.Value == nil {
		return false, nil, nil
	}
	var data map[string]any
	if err := jsonutils.Unmarshal(item.Value, data); err != nil {
		return false, nil, err
	}
	return true, data, nil
}

func (c *DaprState) SetState(ctx context.Context, key string, data any, meta ...StateMeta) error {
	metaMap := NewStateMeta(meta...).ToMap()
	bytes, err := jsonutils.MarshalBytes(data)
	if err != nil {
		return err
	}
	err = c.daprClient.SaveState(ctx, c.storeName, key, bytes, metaMap)
	return err
}
