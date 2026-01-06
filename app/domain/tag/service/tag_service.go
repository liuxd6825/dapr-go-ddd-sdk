package service

import (
	"context"
	"fmt"
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

type TagService struct {
	dao *dao.TagDao
	xbase.Service
	tagTypeService *TagTypeService
}

var (
	_tagOnce          sync.Once
	_tagDomainService *TagService
)

func NewTagService() *TagService {
	_tagOnce.Do(func() {
		_tagDomainService = &TagService{
			dao:            dao.NewTagDao(config.DBKey),
			tagTypeService: NewTagTypeService(),
		}
	})
	return _tagDomainService
}

func (t *TagService) Create(ctx context.Context, cmd *command.TagCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *TagService) CreateMany(ctx context.Context, entities []*model.Tag, opts ...idao.CallOptions) error {
	return t.dao.CreateMany(ctx, entities, opts...).GetError()
}

func (t *TagService) Delete(ctx context.Context, cmd *command.TagDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *TagService) Update(ctx context.Context, data *model.Tag, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, opts...).GetError()
}

func (t *TagService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Tag], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *TagService) FindByTenantId(ctx context.Context, tenantId string, opts ...idao.CallOptions) (int64, error) {
	count, err := t.dao.CountByRSQL(ctx, fmt.Sprintf("tenant_id=='%s'", tenantId), opts...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (t *TagService) InitTag(ctx context.Context, tenantId string) error {
	count, err := t.FindByTenantId(ctx, tenantId, idao.NewCallOptions().SetTenantId(tenantId))
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	tt, err := t.tagTypeService.InitTagType(ctx, tenantId)
	if err != nil {
		return err
	}
	tag1 := &model.Tag{}
	tag1.Id = idutils.NewId()
	tag1.TenantId = tenantId
	tag1.Name = "高风险"
	tag1.TagTypeId = tt.Id
	tag1.IsETag = tt.IsETag
	tag1.Color = "#E62412"
	tag2 := &model.Tag{}
	tag2.Id = idutils.NewId()
	tag2.TenantId = tenantId
	tag2.Name = "中风险"
	tag2.TagTypeId = tt.Id
	tag2.IsETag = tt.IsETag
	tag2.Color = "#FA8C15"
	tag3 := &model.Tag{}
	tag3.Id = idutils.NewId()
	tag3.TenantId = tenantId
	tag3.Name = "低风险"
	tag3.TagTypeId = tt.Id
	tag3.IsETag = tt.IsETag
	tag3.Color = "#797EC9"
	return t.CreateMany(ctx, []*model.Tag{tag1, tag2, tag3}, idao.NewCallOptions().SetTenantId(tenantId))
}
