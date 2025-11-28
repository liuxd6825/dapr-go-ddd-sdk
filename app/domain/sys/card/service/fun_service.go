package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type FunService struct {
	dao *dao.FunDao
	xbase.Service
}

var (
	_funOnce          sync.Once
	_funDomainService *FunService
)

func NewFunService() *FunService {
	_funOnce.Do(func() {
		_funDomainService = &FunService{
			dao: dao.NewFunDao(config.DBKey),
		}
	})
	return _funDomainService
}

func (t *FunService) Create(ctx context.Context, cmd *command.FunCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *FunService) Delete(ctx context.Context, cmd *command.FunDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *FunService) Update(ctx context.Context, cmd *command.FunUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *FunService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Fun, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *FunService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Fun], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *FunService) FindAll(ctx context.Context) ([]*model.Fun, error) {
	return t.dao.FindByRSQL(ctx, "")
}

func (t *FunService) FindByAppId(ctx context.Context, appId string) ([]*model.Fun, error) {
	rsql := fmt.Sprintf("app_id=='%s'", appId)
	return t.dao.FindByRSQL(ctx, rsql)
}

func (t *FunService) FindViewByAppId(ctx context.Context, appId string) ([]*model.FunView, error) {
	arrFun, err := t.FindByAppId(ctx, appId)
	if err != nil {
		return nil, err
	}

	arrLeafFun := make([]*model.Fun, 0)
	for _, v := range arrFun {
		if t.isLeaf(v, arrFun) {
			arrLeafFun = append(arrLeafFun, v)
		}
	}

	arrFunPath := make([]*model.FunView, 0)
	for _, v := range arrLeafFun {
		fp := &model.FunView{}
		path := make([]string, 0)
		multiple := int64(100)
		t.buildFunPath(v, arrFun, &path, &v.OrderNum, &multiple)
		fp.Id = v.Id
		fp.AppId = v.AppId
		fp.Name = v.Name
		fp.Path = path
		fp.OrderNum = float64(v.OrderNum) / float64(multiple)
		fp.Cards = []*model.CardFile{}
		arrFunPath = append(arrFunPath, fp)
	}
	return arrFunPath, nil
}

func (t *FunService) buildFunPath(fun *model.Fun, arr []*model.Fun, path *[]string, orderNum *int64, multiple *int64) {
	*path = append([]string{fun.Name}, *path...)
	if len(fun.ParentId) > 0 {
		for _, v := range arr {
			if v.Id == fun.ParentId {
				*orderNum = (*multiple)*v.OrderNum + (*orderNum)
				*multiple = *multiple * 100
				t.buildFunPath(v, arr, path, orderNum, multiple)
			}
		}
	}
}

func (t *FunService) isLeaf(fun *model.Fun, arr []*model.Fun) bool {
	for _, f := range arr {
		if fun.Id == f.ParentId {
			return false
		}
	}
	return true
}
