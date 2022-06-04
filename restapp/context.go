package restapp

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"strconv"
)

type serverContext struct {
	ctx iris.Context
}

func NewContext(irisCtx iris.Context) context.Context {
	metadata := make(map[string]string, 0)
	header := irisCtx.Request().Header
	for k, v := range header {
		metadata[k] = v[0]
	}
	serverCtx := newServerContext(irisCtx)
	return ddd_context.NewContext(context.Background(), metadata, serverCtx)
}

func newServerContext(ctx iris.Context) ddd_context.ServerContext {
	return &serverContext{
		ctx: ctx,
	}
}

func (s *serverContext) SetResponseHeader(name string, value string) {
	s.ctx.Header(name, value)
}

func (s *serverContext) URLParamDefault(name, def string) string {
	return s.ctx.URLParamDefault(name, def)
}

func NewListQuery(ctx context.Context, tenantId string) (ddd_repository.FindPagingQuery, error) {
	svrCtx := ddd_context.GetServerContext(ctx)
	fields := svrCtx.URLParamDefault("fields", "")
	filter := svrCtx.URLParamDefault("filter", "")
	sort := svrCtx.URLParamDefault("sort", "")
	pageStr := svrCtx.URLParamDefault("page", "0")
	sizeStr := svrCtx.URLParamDefault("size", "20")

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

	return ddd_repository.NewFindPagingQuery(
		tenantId,
		fields,
		filter,
		sort,
		page,
		size), nil
}
