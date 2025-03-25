package store_sql

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"gorm.io/gorm"
)

type Plugin interface {
	PluginName() string
}

var plugins = types.NewCMap[Plugin]()

type GetTablePlugin interface {
	GetTable(ctx context.Context, dao any, db *gorm.DB, tableName string, opts ...store.Options) *gorm.DB
}

func RegisterPlugin(plugin Plugin) {
	plugins.Set(plugin.PluginName(), plugin)
}

func RemovePlugin(plugin Plugin) {
	plugins.Remove(plugin.PluginName())
}
