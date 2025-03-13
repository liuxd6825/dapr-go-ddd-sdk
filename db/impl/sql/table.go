package sql

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/jsonschema/v6"
	"gorm.io/gorm"
	dbschema "gorm.io/gorm/schema"
	"reflect"
	"time"
)

type Table struct {
	tableName string
	schema    *jsonschema.Schema
	dbSchema  *dbschema.Schema
	db        *gorm.DB
}

func NewTable(db *gorm.DB, schema *jsonschema.Schema) db.Table {
	dbSchema, err := NewDBSchema(schema)
	if err != nil {
		panic(err)
	}
	return newTable(db, schema, dbSchema)
}

func newTable(db *gorm.DB, schema *jsonschema.Schema, dbSchema *dbschema.Schema) *Table {
	tableName := dbSchema.Table
	return &Table{db: db, tableName: tableName, schema: schema, dbSchema: dbSchema}
}

func (t *Table) GetTableName() string {
	return t.tableName
}

func (t *Table) GetSchema() *jsonschema.Schema {
	return t.schema
}

func (t *Table) AutoMigrate(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := t.db.AutoMigrate(t.dbSchema)
	if err != nil {
		panic(err)
	}
}

func (t *Table) Exist(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	var result interface{}
	err := t.db.Table(t.dbSchema.Table).First(&result).Error
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
	err := migrator.DropTable(t.dbSchema.Table)
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
