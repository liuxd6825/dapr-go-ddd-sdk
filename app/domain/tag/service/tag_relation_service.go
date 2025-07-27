package service

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/dao"

type TagRelationService struct {
	*dao.TagRelationDao
}

func NewTagRelationService() *TagRelationService {
	return &TagRelationService{
		TagRelationDao: dao.NewTagRelationDao(DBKey),
	}
}
