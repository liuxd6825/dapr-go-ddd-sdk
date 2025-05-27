package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service/model"
)

type DrawService struct {
	drawDao *dao.DrawDao
}

func NewDrawService() *DrawService {
	return &DrawService{
		drawDao: dao.NewDrawDao(),
	}
}

func (s *DrawService) FindById(ctx context.Context, id string) *model.Draw {
	return s.drawDao.FindById(ctx, id)
}
