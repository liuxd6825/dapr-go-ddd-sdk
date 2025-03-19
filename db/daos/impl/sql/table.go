package sql

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
	"reflect"
	"time"
)

type Table struct {
	tableName string
	dbSch     *dbschema.DBSchema
	gormSch   *gormschema.Schema
	db        *gorm.DB
	entity    any
}

func NewTable(db *gorm.DB, dbSch *dbschema.DBSchema, entity any) idao.Table {
	gormSchema, err := NewGormSchema(dbSch)
	if err != nil {
		panic(err)
	}
	return newTable(db, dbSch.TableName, entity, dbSch, gormSchema)
}

func newTable(db *gorm.DB, tableName string, entity any, schema *dbschema.DBSchema, gormSchema *gormschema.Schema) *Table {
	return &Table{db: db, tableName: tableName, entity: entity, dbSch: schema, gormSch: gormSchema}
}

func (t *Table) GetTableName() string {
	return t.tableName
}

func (t *Table) GetSchema() *dbschema.DBSchema {
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

func (t *Table) newModel(schema *schema.Schema) interface{} {
	type DynamicModel struct {
		ID string `gorm:"primaryKey"`
	}
	obj := &DynamicModel{}
	val := reflect.ValueOf(obj).Elem()

	for fieldName, field := range schema.Properties {
		fieldType := reflect.TypeOf("")
		switch field.Type {
		case nil, "string":
			fieldType = reflect.TypeOf("")
			break
		case "int", "number", "integer":
			fieldType = reflect.TypeOf(0)
			break
		case "bool":
			fieldType = reflect.TypeOf(false)
			break
		case "float", "double":
			fieldType = reflect.TypeOf(float64(0))
			break
		case "datetime", "time", "date":
			fieldType = reflect.TypeOf(time.Time{})
			break
		default:
			err := errors.New(fmt.Sprintf("unsupported schema type %v ", fieldType))
			panic(err)
		}
		fieldStruct := reflect.StructField{
			Name: fieldName,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, fieldName)),
		}
		val.Set(reflect.Append(val, reflect.New(fieldStruct.Type)))
	}
	return DynamicModel{}
}
