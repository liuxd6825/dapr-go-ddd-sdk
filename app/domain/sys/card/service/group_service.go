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

type GroupService struct {
	dao *dao.GroupDao
	xbase.Service
}

var (
	_groupOnce          sync.Once
	_groupDomainService *GroupService
)

func NewGroupService() *GroupService {
	_groupOnce.Do(func() {
		_groupDomainService = &GroupService{
			dao: dao.NewGroupDao(config.DBKey),
		}
	})
	return _groupDomainService
}

func (t *GroupService) Create(ctx context.Context, cmd *command.GroupCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *GroupService) Delete(ctx context.Context, cmd *command.GroupDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *GroupService) Update(ctx context.Context, cmd *command.GroupUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *GroupService) CreateMany(ctx context.Context, groups []*model.Group) (int64, error) {
	if len(groups) == 0 {
		return 0, nil
	}
	res := t.dao.CreateMany(ctx, groups)
	return res.GetRowsAffected(), res.GetError()
}

func (t *GroupService) UpdateMany(ctx context.Context, groups []*model.Group) (int64, error) {
	if len(groups) == 0 {
		return 0, nil
	}
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"name", "order_num", "is_hide"})
	res := t.dao.UpdateMany(ctx, groups, opts)
	return res.GetRowsAffected(), res.GetError()
}

func (t *GroupService) DeleteMany(ctx context.Context, listId []string) (int64, error) {
	if len(listId) == 0 {
		return 0, nil
	}
	res := t.dao.DeleteByIds(ctx, listId)
	return res.GetRowsAffected(), res.GetError()
}

func (t *GroupService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Group, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *GroupService) FindPaging(ctx context.Context, caseId string, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Group], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *GroupService) FindByHomeId(ctx context.Context, homeId string) ([]*model.Group, error) {
	rsql := fmt.Sprintf("home_id=='%s'", homeId)
	return t.dao.FindByRSQL(ctx, rsql)
}
