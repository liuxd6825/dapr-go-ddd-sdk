package service

import (
	"context"
)

func (s *GraphService) FindById(ctx context.Context, drawId string) {
	s.graphDao.FindGraphByDrawId(ctx, drawId)
}
