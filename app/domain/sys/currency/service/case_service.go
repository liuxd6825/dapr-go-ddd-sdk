package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"strconv"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
)

// CurrencyService
// @Description: 币种
type CurrencyService struct {
	dao *dao.CurrencyDao
}

var (
	_currencyOnce    sync.Once
	_currencyService *CurrencyService
)

func NewCurrencyService() *CurrencyService {
	_currencyOnce.Do(func() {
		_currencyService = &CurrencyService{
			dao: dao.NewCurrencyDao(config.DBKey),
		}
	})
	return _currencyService
}

func (s *CurrencyService) GetConfig() *idao.DaoConfig {
	return s.dao.GetConfig()
}

func (s *CurrencyService) Create(ctx context.Context, cmd *command.CurrencyCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (s *CurrencyService) CreateMany(ctx context.Context, cmd *command.CurrencyCreateManyCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.CreateManyData(ctx, cmd.Data)
	})
}

func (s *CurrencyService) CreateManyData(ctx context.Context, entities []*model.Currency, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, entities, opts...).GetError()
}

func (s *CurrencyService) Delete(ctx context.Context, cmd *command.CurrencyDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (s *CurrencyService) DeleteAll(ctx context.Context, cmd *command.CurrencyDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.DeleteAll(ctx).GetError()
	})
}

func (s *CurrencyService) DeleteBatch(ctx context.Context, cmd *command.CurrencyDeleteBatchCommand) error {
	return s.dao.DeleteByIds(ctx, cmd.Data.Ids).GetError()
}

func (s *CurrencyService) Update(ctx context.Context, cmd *command.CurrencyUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (s *CurrencyService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Currency, error) {
	return s.dao.FindById(ctx, qry.Id)
}

func (s *CurrencyService) FindByAll(ctx context.Context) ([]*model.Currency, error) {
	data := s.dao.FindAll(ctx)
	return data.Data, data.Error
}

func (s *CurrencyService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Currency], error) {
	res := s.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

// UpdateAllCnyRate
// @Description:  更新以人民币计价的汇率
// @receiver s
// @param ctx
// @param cnyRate 人民币对美元
// @return error
func (s *CurrencyService) UpdateAllCnyRate(ctx context.Context, cnyRate float64) error {
	dbKey := config.DBKey
	return tx.StartTx(ctx, []string{dbKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		s.dao.UpdateByRSQL(ctx, "name='人民币'", &model.Currency{Rate: cnyRate, CnRate: 1}, idao.NewCallOptions().SetUpdateFields([]string{"Rate"}))
		list, err := s.dao.FindByRSQL(ctx, "name!='人民币'")
		if err != nil {
			return err
		}
		for _, c := range list {
			cnRate := c.Rate / cnyRate
			valStr := strconv.FormatFloat(cnRate, 'f', 4, 64)
			if val, err := strconv.ParseFloat(valStr, 64); err == nil {
				c.CnRate = val
			}
		}
		return s.dao.UpdateMany(ctx, list, idao.NewCallOptions().SetUpdateFields([]string{"Rate"})).GetError()
	})
}

func (s *CurrencyService) InitCurrency(ctx context.Context, tenantId string) error {
	arr, err := s.dao.FindByRSQL(ctx, fmt.Sprintf("tenant_id=='%s'", "test"), idao.NewCallOptions().SetTenantId("test"))
	if err != nil {
		return err
	}
	if arr == nil || len(arr) == 0 {
		return errors.New("初始化币种时没有找到币种数据")
	}

	for _, v := range arr {
		v.Id = idutils.NewId()
		v.TenantId = tenantId
	}

	return s.CreateManyData(ctx, arr, idao.NewCallOptions().SetTenantId(tenantId))
}
