package graph

import "context"

type QueryService struct {
}

func NewQueryService() *QueryService {
	return &QueryService{}
}

func (s *QueryService) FindByName(ctx context.Context, caseId, nodeLabel, name string) (any, error) {
	return "query", nil
}

func (s *QueryService) FindById(ctx context.Context, caseId, nodeLabel, id string) (any, error) {
	return "query", nil
}
