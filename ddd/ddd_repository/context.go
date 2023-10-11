package ddd_repository

import "context"

type ctxKey struct {
}

type IGetTableName interface {
	GetTableName(ctx context.Context) string
}

func NewContext(parent context.Context, data any) context.Context {
	ctx, _ := context.WithCancel(parent)
	return context.WithValue(ctx, ctxKey{}, data)
}

//
// GetContextData
//  @Description:  获取仓储上下文数据对象
//  @param ctx
//  @return any
//
func GetContextData(ctx context.Context) any {
	return ctx.Value(ctxKey{})
}

//
// GetTableName
//  @Description: 从仓储上下文中获取表名称
//  @param ctx
//  @return string
//  @return bool
//
func GetTableName(ctx context.Context) (string, bool) {
	value := ctx.Value(ctxKey{})
	if getTable, ok := value.(IGetTableName); ok {
		table := getTable.GetTableName(ctx)
		return table, len(table) > 0
	}
	return "", false
}
