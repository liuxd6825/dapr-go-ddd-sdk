package micro_states

import (
	"context"
)

// IMicroState
// @Description: 幂等接口
type IMicroState interface {
	GetState(ctx context.Context, key string, meta ...StateMeta) (bool, any, error)
	SetState(ctx context.Context, key string, data any, meta ...StateMeta) error
}
