package restapi

import (
	"context"
	"fmt"

	icontext "github.com/kataras/iris/v12/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/micro/micro_states"
)

func (c *ApiController) getState(ctx context.Context, ictx *icontext.Context, params any) (bool, any, error) {
	var stateKey = c.getStateKey(ctx, ictx, params)
	return c.stateProvider.GetState(ctx, stateKey)
}

func (c *ApiController) getStateKey(ctx context.Context, ictx *icontext.Context, params any) string {
	var stateKey string = "unknown"
	if cmd, ok := params.(ICommand); ok {
		stateKey = fmt.Sprintf("command-%s", cmd.GetCommandId())
	} else if ev, ok := params.(IEvent); ok {
		stateKey = fmt.Sprintf("event-%s", ev.GetEventId())
	} else {
		stateKey = ictx.Request().RequestURI
	}
	return stateKey
}

func (c *ApiController) setState(ctx context.Context, ictx *icontext.Context, params any, data any) error {
	var stateKey = c.getStateKey(ctx, ictx, params)
	return c.stateProvider.SetState(ctx, stateKey, data, micro_states.StateMeta{
		TtlInSeconds: 3600,
	})
}
