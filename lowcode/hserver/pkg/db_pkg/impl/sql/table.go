package sql

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type Table struct {
	name   string
	schema *schema.Schema
	db     *gorm.DB
}

func NewTable(db *gorm.DB, name string, schema *schema.Schema) *Table {
	return &Table{db: db, name: name, schema: schema}
}

func (t *Table) GetName() string {
	return t.name
}

func (t *Table) GetSchema() *schema.Schema {
	return t.schema
}

func (t *Table) AutoMigrate(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	model := t.newModel(t.schema)
	err := t.db.AutoMigrate(model)
	if err != nil {
		panic(err)
	}
}

func (t *Table) Exist(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	var result interface{}
	err := t.db.Table(t.name).First(&result).Error
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
	err := migrator.DropTable(t.name)
	if err != nil {
		panic(err)
	}
}

func (t *Table) newModel(schema *schema.Schema) interface{} {
	type DynamicModel struct {
		ID string `gorm:"primaryKey"`
	}

	for fieldName, field := range schema.Properties {
		fieldType := reflect.TypeOf("")
		if field.Type == "string" {
			fieldType = reflect.TypeOf("")
		} else if field.Type == "integer" {
			fieldType = reflect.TypeOf(0)
		} else if field.Type == "int" {
			fieldType = reflect.TypeOf(0)
		} else if field.Type == "number" {
			fieldType = reflect.TypeOf(0)
		} else if field.Type == "date" {
			fieldType = reflect.TypeOf(time.Time{})
		} else if field.Type == "dateTime" {
			fieldType = reflect.TypeOf(time.Time{})
		} else if field.Type == "bool" {
			fieldType = reflect.TypeOf(false)
		}
		fieldStruct := reflect.StructField{
			Name: fieldName,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, fieldName)),
		}
		reflect.ValueOf(&DynamicModel{}).Elem().FieldByName(fieldName).Set(reflect.New(fieldStruct.Type).Elem())
	}
	return DynamicModel{}
}
