package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/mapperutils"
)

// RagRecordService
// @Description: 知识库流水关系
type RagRecordService struct {
	recordDao *dao.RagRecordDao
}

func NewRecordService() *RagRecordService {
	return &RagRecordService{
		recordDao: dao.NewRagRecordDao(config.Neo4jDBKey),
	}
}

func (s *RagRecordService) DataChange(ctx context.Context, cdcRecord *dbevent.CDCRecord) error {
	record, err := s.NewRecord(ctx, cdcRecord)
	if err != nil {
		return err
	}
	switch cdcRecord.OpType {
	case dbevent.OpTypeRead:
		break
	case dbevent.OpTypeCreate:
		return s.Create(ctx, record)
	case dbevent.OpTypeUpdate:
		return s.Update(ctx, record)
	case dbevent.OpTypeDelete:
		return s.Delete(ctx, record)
	}
	return nil
}

func (s *RagRecordService) Create(ctx context.Context, record *model.Record) error {
	return s.recordDao.Create(ctx, record)
}

func (s *RagRecordService) Update(ctx context.Context, record *model.Record) error {
	return s.recordDao.Update(ctx, record)
}

func (s *RagRecordService) Delete(ctx context.Context, record *model.Record) error {
	return s.recordDao.Delete(ctx, record)
}
func (s *RagRecordService) NewRecord(ctx context.Context, cdcRecord *dbevent.CDCRecord) (record *model.Record, err error) {
	switch cdcRecord.OpType {
	case dbevent.OpTypeRead:
		record, err = s.newRecord(ctx, cdcRecord.AfterMap())
	case dbevent.OpTypeCreate:
		record, err = s.newRecord(ctx, cdcRecord.AfterMap())
	case dbevent.OpTypeUpdate:
		record, err = s.newRecord(ctx, cdcRecord.AfterMap())
	case dbevent.OpTypeDelete:
		record, err = s.newRecord(ctx, cdcRecord.AfterMap())
	}
	return record, err
}

func (s *RagRecordService) newRecord(ctx context.Context, data map[string]any) (*model.Record, error) {
	record := &model.Record{}
	err := mapperutils.MapToStructWithTag(data, record, "bson")
	return record, err
}
