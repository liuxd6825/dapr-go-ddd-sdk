package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/model"
)

type DrawService struct {
	drawDao *dao.DrawDao
}

const DBKey string = "$drawDb"
const Neo4jDBKey string = "$neo4jDb"

func NewDrawService() *DrawService {
	return &DrawService{
		drawDao: dao.NewDrawDao(DBKey),
	}
}

func (s *DrawService) FindById(ctx context.Context, id string) *model.Draw {
	return s.drawDao.FindById(ctx, id)
}
