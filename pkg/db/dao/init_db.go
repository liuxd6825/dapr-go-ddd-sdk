package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

type initDBOption func(ctx context.Context, table idao.Table, err error)

func CreateTables(ctx context.Context, tables []idao.Table, options ...initDBOption) {
	if len(options) == 0 {
		options = append(options, dbLogError)
	}
	for _, table := range tables {
		err := table.AutoMigrate(ctx)
		for _, f := range options {
			f(ctx, table, err)
		}
	}
}

func dbLogError(ctx context.Context, table idao.Table, err error) {
	if err == nil {
		logs.Infofmt(ctx, "tableName:%s; succeed", table.GetTableName())
	} else {
		logs.Errorfmt(ctx, "tableName:%s; error:%s", table.GetTableName(), err.Error())
	}
}
