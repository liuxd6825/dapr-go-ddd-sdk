package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"sync"
)

type TagRelationService struct {
	dao *dao.TagRelationDao
	xbase.Service
}

var (
	_tagRelationOnce          sync.Once
	_tagRelationDomainService *TagRelationService
)

func NewTagRelationService() *TagRelationService {
	_tagRelationOnce.Do(func() {
		_tagRelationDomainService = &TagRelationService{
			dao: dao.NewTagRelationDao(config.DBKey),
		}
	})
	return _tagRelationDomainService
}

func (t *TagRelationService) Create(ctx context.Context, cmd *command.TagRelationCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *TagRelationService) Delete(ctx context.Context, cmd *command.TagRelationDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *TagRelationService) DeleteByRSQL(ctx context.Context, rSql string) error {
	return t.dao.DeleteByRSQL(ctx, rSql).GetError()
}

func (t *TagRelationService) Update(ctx context.Context, data *model.TagRelation, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, opts...).GetError()
}

func (t *TagRelationService) UpdateMany(ctx context.Context, data []*model.TagRelation, opts ...idao.CallOptions) error {
	return t.dao.UpdateMany(ctx, data, opts...).GetError()
}

func (t *TagRelationService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.TagRelation], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *TagRelationService) FindByRSQL(ctx context.Context, rSql string) ([]*model.TagRelation, error) {
	return t.dao.FindByRSQL(ctx, rSql)
}
