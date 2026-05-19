package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
)

type MasterQueryService struct {
	graphDao *dao.MasterGraphDao
}

func NewMasterQueryService() *MasterQueryService {
	return &MasterQueryService{
		graphDao: dao.NewMasterGraphDao(),
	}
}

func (s *MasterQueryService) FindByCaseId(ctx context.Context, caseId string) *graph.GraphView {
	return s.graphDao.FindByCaseId(ctx, caseId)
}

func (s *MasterQueryService) FindById(ctx context.Context, caseId, id string) *graph.GraphView {
	return s.graphDao.FindById(ctx, caseId, id)
}

func (s *MasterQueryService) FindByName(ctx context.Context, caseId, name string) *graph.GraphView {
	return s.graphDao.FindByName(ctx, caseId, name)
}
