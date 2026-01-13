package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
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
	record.Id, _ = maputils.GetString(data, "id", "")
	record.CaseId, _ = maputils.GetString(data, "case_id", "")
	record.TenantId, _ = maputils.GetString(data, "tenant_id", "")
	record.TaskId, _ = maputils.GetString(data, "task_id", "")

	record.DocId, _ = maputils.GetString(data, "doc_id", "")
	record.FileId, _ = maputils.GetString(data, "file_id", "")
	record.FileName, _ = maputils.GetString(data, "file_name", "")
	record.RowNum, _ = maputils.GetInt64(data, "row_num", 0)
	record.Remark, _ = maputils.GetString(data, "remark", "")

	record.Iden, _ = maputils.GetString(data, "iden", "")
	record.Name, _ = maputils.GetString(data, "name", "")
	record.Type, _ = maputils.GetString(data, "type", "")
	record.Acct, _ = maputils.GetString(data, "acct", "")
	record.AcctType, _ = maputils.GetString(data, "acct_type", "")
	record.BankName, _ = maputils.GetString(data, "bank_name", "")
	record.Balance, _ = maputils.GetPFloat64(data, "balance", nil)
	record.Category, _ = maputils.GetString(data, "category", "")

	record.OppIden, _ = maputils.GetString(data, "opp_iden", "")
	record.OppName, _ = maputils.GetString(data, "opp_name", "")
	record.OppAcctType, _ = maputils.GetString(data, "opp_acct_type", "")
	record.OppAcct, _ = maputils.GetString(data, "opp_acct", "")
	record.OppBankName, _ = maputils.GetString(data, "opp_bank_name", "")

	record.Payout, _ = maputils.GetPFloat64(data, "payout", nil)
	record.Income, _ = maputils.GetPFloat64(data, "income", nil)
	record.Date, _ = maputils.GetDate(data, "date", nil)
	record.Amount, _ = maputils.GetPFloat64(data, "amount", nil)
	record.Ccy, _ = maputils.GetString(data, "ccy", "")

	record.CreatorId, _ = maputils.GetString(data, "creator_id", "")
	record.CreatorName, _ = maputils.GetString(data, "creator_name", "")
	record.CreatedTime, _ = maputils.GetDate(data, "created_time", nil)
	record.UpdaterId, _ = maputils.GetString(data, "updater_id", "")
	record.UpdaterName, _ = maputils.GetString(data, "updater_name", "")
	record.UpdatedTime, _ = maputils.GetDate(data, "updated_time", nil)

	record.Year, _ = maputils.GetInt(data, "year", 0)
	record.Month, _ = maputils.GetInt(data, "month", 0)
	record.Day, _ = maputils.GetInt(data, "day", 0)

	return record, nil
}
