package service_test

import (
	"context"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
)

const CASE_ID = "PnRAAeb4liYYMyLfpsA7e9bu"

func Test_FindByCaseId(t *testing.T) {
	query, ctx := newQueryService()
	graphView := query.FindByCaseId(ctx, CASE_ID)
	println(graphView)
}

func Test_FindById(t *testing.T) {
	query, ctx := newQueryService()
	graphView := query.FindById(ctx, CASE_ID, "PnRAAeb4liYYMyLfpsA7e9bu_任萌")
	println(graphView)
}

func Test_FindByName(t *testing.T) {
	query, ctx := newQueryService()
	graphView := query.FindByName(ctx, CASE_ID, "任萌")
	println(graphView)
}

func Test_FindByStartName(t *testing.T) {
	query, ctx := newQueryService()
	graphView := query.FindByStartName(ctx, CASE_ID, "任")
	println(graphView)
}

func Test_FindByContainName(t *testing.T) {
	query, ctx := newQueryService()
	graphView := query.FindByContainName(ctx, CASE_ID, "任")
	println(graphView)
}

func newQueryService() (*service.MasterQueryService, context.Context) {
	ctx := appctx.NewTenantContext(context.Background(), "dev")
	query := service.NewMasterQueryService()
	return query, ctx
}
