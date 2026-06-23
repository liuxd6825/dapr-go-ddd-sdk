package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

// CompanyRecordService 公司-流水关联业务服务（只读）
type CompanyRecordService struct {
	dao *dao.CompanyRecordDao
}

var _companyRecordService *CompanyRecordService
var _companyRecordServiceOnce sync.Once

func NewCompanyRecordService() *CompanyRecordService {
	_companyRecordServiceOnce.Do(func() {
		_companyRecordService = &CompanyRecordService{
			dao: dao.NewCompanyRecordDao(config.DBKey),
		}
	})
	return _companyRecordService
}

func (s *CompanyRecordService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.CompanyRecord, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *CompanyRecordService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyRecord] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *CompanyRecordService) FindByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyRecord] {
	return s.dao.FindPagingByCompanyId(ctx, qry, companyId, opts...)
}