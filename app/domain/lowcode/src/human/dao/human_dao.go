package dao

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/src/human/model"
)

type HumanDao struct {
}

func NewHumanDao() *HumanDao {
	return &HumanDao{}
}

func (d *HumanDao) Create(ctx context.Context, human *model.Human) error {
	return nil
}

func (d *HumanDao) Update(ctx context.Context, human *model.Human) {

}

func (d *HumanDao) FindById(ctx context.Context, id string) (*model.Human, error) {
	return nil, nil
}

func (d *HumanDao) FindPaging(ctx context.Context, human *model.Human) {

}
