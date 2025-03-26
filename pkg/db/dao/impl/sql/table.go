package sql

import (
	"context"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
)

type Table struct {
	tableName string
	dbSch     *store.DBSchema
	gormSch   *gormschema.Schema
	db        *gorm.DB
	entity    any
}

func NewTable(db *gorm.DB, dbSch *store.DBSchema, entity any) idao.Table {
	gormSchema, err := NewGormSchema(dbSch)
	if err != nil {
		panic(err)
	}
	return newTable(db, dbSch.TableName, entity, dbSch, gormSchema)
}

func newTable(db *gorm.DB, tableName string, entity any, schema *store.DBSchema, gormSchema *gormschema.Schema) *Table {
	return &Table{db: db, tableName: tableName, entity: entity, dbSch: schema, gormSch: gormSchema}
}

func (t *Table) GetTableName() string {
	return t.tableName
}

func (t *Table) GetSchema() *store.DBSchema {
	return t.dbSch
}

func (t *Table) AutoMigrate(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	if _, ok := t.entity.(map[string]any); ok {
		err = t.db.Table(t.tableName).AutoMigrate(t.gormSch)
	} else {
		err = t.db.Table(t.tableName).CustomSchema(t.gormSch).AutoMigrate(t.entity)
	}

	if err != nil {
		panic(err)
	}
}

func (t *Table) Exist(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	var result interface{}
	err := t.db.Table(t.tableName).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 表存在但没有数据
			return true
		} else {
			panic(err)
		}
	}
	return true
}

func (t *Table) Drop(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	migrator := t.db.Migrator()
	err := migrator.DropTable(t.gormSch.Table)
	if err != nil {
		panic(err)
	}
}
