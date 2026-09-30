package service

import (
	"context"

	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao/impl/hugedao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
)

type MasterQueryService struct {
	graphDao dao2.IMasterGraphDao
}

func NewMasterQueryService() *MasterQueryService {
	return &MasterQueryService{
		graphDao: hugedao.NewMasterGraphDao(),
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

func (s *MasterQueryService) FindByContainName(ctx context.Context, caseId, name string) *graph.GraphView {
	return s.graphDao.FindByContainName(ctx, caseId, name)
}

func (s *MasterQueryService) FindByStartName(ctx context.Context, caseId, name string) *graph.GraphView {
	return s.graphDao.FindByStartWithName(ctx, caseId, name)
}
