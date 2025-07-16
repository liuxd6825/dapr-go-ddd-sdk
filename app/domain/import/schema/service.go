package schema

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type SchemaService struct {
	dao *SchemaDao
	xbase.Service
}

const (
	CacheKey_Create     = "import.SchemaService.Create()"
	CacheKey_Update     = "import.SchemaService.Update()"
	CacheKey_DeleteById = "import.SchemaService.DeleteById()"
	CacheKey_FindPaging = "import.SchemaService.FindPaging()"
	CacheKey_FindById   = "import.SchemaService.FindById()"
)

func NewSchemaService() *SchemaService {
	return &SchemaService{
		dao: NewSchemaDao(config.DBKey),
	}
}

func (s *SchemaService) Create(ctx context.Context, cmd *SchemaCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (s *SchemaService) Update(ctx context.Context, cmd *SchemaUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (s *SchemaService) DeleteById(ctx context.Context, cmd *SchemaDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (s *SchemaService) FindPaging(ctx context.Context, qry idao.FindPagingByCaseIdQuery) (idao.FindPagingResult[*SchemaModel], error) {
	return xbase.DoQuery[idao.FindPagingResult[*SchemaModel]](ctx, qry, func(ctx context.Context) (idao.FindPagingResult[*SchemaModel], error) {
		res := s.dao.FindPaging(ctx, qry)
		return res, res.GetError()
	}, xbase.DoQueryOptions{CacheKey: CacheKey_FindPaging})
}

func (s *SchemaService) FindById(ctx context.Context, id string) (*SchemaModel, error) {
	return xbase.DoQuery[*SchemaModel](ctx, id, func(ctx context.Context) (*SchemaModel, error) {
		return s.dao.FindById(ctx, id)
	}, xbase.DoQueryOptions{CacheKey: CacheKey_FindById})
}
