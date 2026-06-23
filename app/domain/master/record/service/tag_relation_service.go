package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
)

type TagRelationService struct {
	*dao.TagRelationDao
}

func NewTagRelationService() *TagRelationService {
	return &TagRelationService{
		TagRelationDao: dao.NewTagRelationDao(config.DBKey),
	}
}
