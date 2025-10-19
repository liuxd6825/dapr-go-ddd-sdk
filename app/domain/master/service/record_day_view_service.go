package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
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
