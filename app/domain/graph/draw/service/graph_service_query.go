package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
)

func (s *GraphService) FindById(ctx context.Context, caseId string, drawId string) *graph.GraphView {
	return s.graphDao.FindGraphByDrawId(ctx, caseId, drawId)
}

func (s *GraphService) FindInCaseByNames(ctx context.Context, caseId string, names []string) *graph.GraphView {
	return s.graphDao.FindInCaseByNames(ctx, caseId, names)
}
