package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/src/human/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/src/human/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/src/human/model"
)

type HumanService struct {
	dao *dao.HumanDao
}

func NewHumanService() *HumanService {
	return &HumanService{
		dao: dao.NewHumanDao(),
	}
}

func (s *HumanService) Create(ctx context.Context, cmd *command.HumanCreateCommand) error {
	human := &model.Human{
		Id: cmd.Data.Id,
	}
	return s.dao.Create(ctx, human)
}
