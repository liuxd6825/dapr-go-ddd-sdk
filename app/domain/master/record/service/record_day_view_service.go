package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
)

type RecordDayViewService struct {
	*dao.RecordDayViewDao
}

func NewRecordDayViewService() *RecordDayViewService {
	viewDao := dao.NewRecordDayViewDao(config.DBKey)
	return &RecordDayViewService{
		RecordDayViewDao: viewDao,
	}
}
