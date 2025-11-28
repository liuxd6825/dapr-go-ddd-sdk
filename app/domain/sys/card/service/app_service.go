package service

import (
	"context"
	"sort"
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

type AppService struct {
	dao *dao.AppDao
	xbase.Service
}

var (
	_appOnce          sync.Once
	_appDomainService *AppService
)

func NewAppService() *AppService {
	_appOnce.Do(func() {
		_appDomainService = &AppService{
			dao: dao.NewAppDao(config.DBKey),
		}
	})
	return _appDomainService
}

func (t *AppService) Create(ctx context.Context, cmd *command.AppCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *AppService) Delete(ctx context.Context, cmd *command.AppDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *AppService) Update(ctx context.Context, cmd *command.AppUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *AppService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.App, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *AppService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.App], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *AppService) FindAll(ctx context.Context) ([]*model.App, error) {
	return t.dao.FindByRSQL(ctx, "")
}

func (t *AppService) TransformTree(ctx context.Context, apps []*model.App, funs []*model.Fun) []*model.AppTree {
	arr := make([]*model.AppTree, 0)
	for _, app := range apps {
		m := &model.AppTree{
			Id:       app.Id,
			Name:     app.Name,
			Children: []*model.AppTree{},
		}
		t.transformFuns(true, app.Id, funs, &m.Children)
		arr = append(arr, m)
	}
	return arr
}

func (t *AppService) transformFuns(isApp bool, parentId string, funs []*model.Fun, arr *[]*model.AppTree) {
	arrFuns := make([]*model.Fun, 0)
	for _, fun := range funs {
		if (isApp && parentId == fun.AppId && len(fun.ParentId) == 0) || (!isApp && parentId == fun.ParentId) {
			arrFuns = append(arrFuns, fun)
		}
	}
	if len(arrFuns) > 0 {
		sort.Slice(arrFuns, func(i, j int) bool { return arrFuns[i].OrderNum < arrFuns[j].OrderNum })
		for _, fun := range arrFuns {
			m := &model.AppTree{
				Id:       fun.Id,
				Name:     fun.Name,
				Children: []*model.AppTree{},
			}
			t.transformFuns(false, fun.Id, funs, &m.Children)
			*arr = append(*arr, m)
		}
	}
}
