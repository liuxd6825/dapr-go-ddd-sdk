package sql

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_sql"
	"gorm.io/gorm"
	dbschema "gorm.io/gorm/schema"
)

func init() {
	// 安装插件，用于取得数据表
	ddd_sql.RegisterPlugin(NewTablePlugin())
}

type TablePlugin struct {
}

func NewTablePlugin() *TablePlugin {
	return &TablePlugin{}
}

func (p *TablePlugin) PluginName() string {
	return "PluginName"
}

type metadata interface {
	GetMetadata() map[string]any
}

// GetTable 通过Schema取得数据表
func (p *TablePlugin) GetTable(ctx context.Context, dao any, db *gorm.DB, tableName string, opts ...ddd_repository.Options) *gorm.DB {
	if d, ok := dao.(metadata); ok {
		if val, ok := d.GetMetadata()["dbSchema"]; ok {
			if dbSchema, ok := val.(*dbschema.Schema); ok {
				return db.Table(dbSchema.Table).MapSchema(dbSchema)
			}
		}
	}
	return db.Table(tableName)
}
