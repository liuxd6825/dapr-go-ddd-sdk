package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/dao"
)

type TagRelationService struct {
	*dao.TagRelationDao
}

func NewTagRelationService(dbKey string) *TagRelationService {
	return &TagRelationService{
		TagRelationDao: dao.NewTagRelationDao(dbKey),
	}
}
