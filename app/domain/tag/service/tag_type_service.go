package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/dao"
)

const DBKey string = "$db"

type TagTypeService struct {
	*dao.TagTypeDao
}

func NewTagTypeService() *TagTypeService {
	return &TagTypeService{
		TagTypeDao: dao.NewTagTypeDao(DBKey),
	}
}
