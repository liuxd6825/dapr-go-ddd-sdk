package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/dao"
)

type RecordDayViewService struct {
	*dao.RecordDayViewDao
}

func NewRecordDayViewService(dbKey string) *RecordDayViewService {
	viewDao := dao.NewRecordDayViewDao(dbKey)
	return &RecordDayViewService{
		RecordDayViewDao: viewDao,
	}
}
