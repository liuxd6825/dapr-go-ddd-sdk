package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/mapperutils"
)

type RecordService struct {
	recordDao *dao.RecordDao
}

func NewRecordService() *RecordService {
	return &RecordService{
		recordDao: dao.NewRecordDao(config.Neo4jDBKey),
	}
}

func (s *RecordService) DataChange(ctx context.Context, cdcRecord *model.CDCRecord) error {
	record, err := s.NewRecord(ctx, cdcRecord)
	if err != nil {
		return err
	}
	switch cdcRecord.OpType {
	case model.OpTypeRead:
		break
	case model.OpTypeCreate:
		return s.Create(ctx, record)
	case model.OpTypeUpdate:
		return s.Update(ctx, record)
	case model.OpTypeDelete:
		return s.Delete(ctx, record)
	}
	return nil
}

func (s *RecordService) Create(ctx context.Context, record *model.Record) error {
	return s.recordDao.Create(ctx, record)
}

func (s *RecordService) Update(ctx context.Context, record *model.Record) error {
	return s.recordDao.Update(ctx, record)
}

func (s *RecordService) Delete(ctx context.Context, record *model.Record) error {
	return s.recordDao.Delete(ctx, record)
}
func (s *RecordService) NewRecord(ctx context.Context, cdcRecord *model.CDCRecord) (record *model.Record, err error) {
	switch cdcRecord.OpType {
	case model.OpTypeRead:
		record, err = s.newRecord(ctx, cdcRecord.AfterMap())
	case model.OpTypeCreate:
		record, err = s.newRecord(ctx, cdcRecord.AfterMap())
	case model.OpTypeUpdate:
		record, err = s.newRecord(ctx, cdcRecord.AfterMap())
	case model.OpTypeDelete:
		record, err = s.newRecord(ctx, cdcRecord.AfterMap())
	}
	return record, err
}

func (s *RecordService) newRecord(ctx context.Context, data map[string]any) (*model.Record, error) {
	record := &model.Record{}
	err := mapperutils.MapToStructWithTag(data, record, "bson")
	return record, err
}
