package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"
)

func (s *GraphService) FindById(ctx context.Context, drawId string) *graph.GraphView {
	return s.graphDao.FindGraphByDrawId(ctx, drawId)
}
