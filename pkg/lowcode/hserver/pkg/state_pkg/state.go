package state_pkg

import (
	"context"

	"github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
)

type StatePkg struct {
	server element.Server
	client dapr.Client
}

func New(server element.Server) *StatePkg {
	cli, err := dapr.GetClient()
	if err != nil {
		panic(err)
	}
	return &StatePkg{
		server: server,
		client: cli,
	}
}

func (s *StatePkg) TryLock(ctx context.Context, storeName string, request *client.LockRequest) *client.LockResponse {
	if request == nil {
		panic("request is nil")
	}
	resp, err := s.client.TryLockAlpha1(ctx, storeName, request)
	if err != nil {
		panic(err)
	}
	return resp
}

func (s *StatePkg) Unlock(ctx context.Context, storeName string, request *client.UnlockRequest) *client.UnlockResponse {
	if request == nil {
		panic("request is nil")
	}
	resp, err := s.client.UnlockAlpha1(ctx, storeName, request)
	if err != nil {
		panic(err)
	}
	return resp
}

func (s *StatePkg) SaveState(ctx context.Context, storeName, key string, data []byte, meta map[string]string, so ...client.StateOption) {
	err := s.client.SaveState(ctx, storeName, key, data, meta, so...)
	if err != nil {
		panic(err)
	}
}

func (s *StatePkg) DeleteState(ctx context.Context, storeName, key string, meta map[string]string) {
	err := s.client.DeleteState(ctx, storeName, key, meta)
	if err != nil {
		panic(err)
	}
}

func (s *StatePkg) SaveBulkState(ctx context.Context, storeName string, items ...*client.SetStateItem) {
	err := s.client.SaveBulkState(ctx, storeName, items...)
	if err != nil {
		panic(err)
	}
}

func (s *StatePkg) DeleteBulkState(ctx context.Context, storeName string, keys []string, meta map[string]string) {
	err := s.client.DeleteBulkState(ctx, storeName, keys, meta)
	if err != nil {
		panic(err)
	}
}
