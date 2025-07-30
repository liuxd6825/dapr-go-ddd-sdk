package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
)

type RecordService struct {
	recordDao *dao.RecordDao
}

func NewRecordService() *RecordService {
	return &RecordService{}
}

func (s *RecordService) Save(ctx context.Context, record *model.CDCRecord) error {
	return x
}
