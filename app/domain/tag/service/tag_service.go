package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/dao"
)

type TagService struct {
	*dao.TagDao
}

func NewTagService() *TagService {
	return &TagService{
		TagDao: dao.NewTagDao(DBKey),
	}
}
