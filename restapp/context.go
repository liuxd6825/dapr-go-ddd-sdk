package restapp

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"strconv"
)

func NewListQuery(ctx iris.Context, tenantId string) (*ddd_repository.FindPagingQuery, error) {
	fields := ctx.URLParamDefault("fields", "")
	filter := ctx.URLParamDefault("filter", "")
	sort := ctx.URLParamDefault("sort", "")
	pageStr := ctx.URLParamDefault("page", "0")
	sizeStr := ctx.URLParamDefault("size", "20")

	if len(tenantId) == 0 {
		return nil, errors.New("tenantId is null")
	}

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &ddd_repository.FindPagingQuery{
		TenantId: tenantId,
		Fields:   fields,
		Filter:   filter,
		Sort:     sort,
		PageNum:  page,
		PageSize: size,
	}, nil
}
