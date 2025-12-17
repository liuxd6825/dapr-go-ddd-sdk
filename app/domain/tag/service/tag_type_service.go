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
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"sync"
)

type TagTypeService struct {
	dao *dao.TagTypeDao
	xbase.Service
}

var (
	_tagTypeOnce          sync.Once
	_tagTypeDomainService *TagTypeService
)

func NewTagTypeService() *TagTypeService {
	_tagTypeOnce.Do(func() {
		_tagTypeDomainService = &TagTypeService{
			dao: dao.NewTagTypeDao(config.DBKey),
		}
	})
	return _tagTypeDomainService
}

func (t *TagTypeService) Create(ctx context.Context, cmd *command.TagTypeCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *TagTypeService) CreateData(ctx context.Context, data *model.TagType, opts ...idao.CallOptions) error {
	return t.dao.Create(ctx, data, opts...).GetError()
}

func (t *TagTypeService) Delete(ctx context.Context, cmd *command.TagTypeDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *TagTypeService) Update(ctx context.Context, data *model.TagType, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, opts...).GetError()
}

func (t *TagTypeService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.TagType], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *TagTypeService) InitTagType(ctx context.Context, tenantId string) (*model.TagType, error) {
	tt := &model.TagType{}
	tt.Id = idutils.NewId()
	tt.TenantId = tenantId
	tt.Name = "风险级别"
	tt.IsETag = true
	tt.ParentId = ""
	err := t.CreateData(ctx, tt, idao.NewCallOptions().SetTenantId(tenantId))
	if err != nil {
		return nil, err
	}
	return tt, nil
}
