package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"
)

type QueryService struct {
	graphDao *dao.MasterGraphDao
}

func NewQueryService() *QueryService {
	return &QueryService{
		graphDao: dao.NewMasterGraphDao(),
	}
}

func (s *QueryService) FindByCaseId(ctx context.Context, caseId string) *graph.GraphView {
	return s.graphDao.FindByCaseId(ctx, caseId)
}

func (s *QueryService) FindById(ctx context.Context, caseId, id string) *graph.GraphView {
	return s.graphDao.FindById(ctx, caseId, id)
}

func (s *QueryService) FindByName(ctx context.Context, caseId, name string) *graph.GraphView {
	return s.graphDao.FindByName(ctx, caseId, name)
}
