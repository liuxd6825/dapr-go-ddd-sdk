package ddd_sql

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl/tests"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	dbschema "gorm.io/gorm/schema"
	"testing"
)

func Test_CreateTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
		return
	}
	err = db.Migrator().CreateTable(&Human{})
	if err != nil {
		t.Error(err)
	}
}

func Test_Schema(t *testing.T) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
		return
	}

	dest, err := schema.NewSchemaWithJson("human.json", tests.HumanSchema)
	if err != nil {
		t.Error(err)
		return
	}

	var dbSchema *dbschema.Schema
	t.Run("NewDBSchema", func(t *testing.T) {
		dschema := dbschema.NewSchema()
		dschema.Name = "table"
		v, err := NewGormSchema(dest.GetSchema())
		if err != nil {
			t.Error(err)
		}
		dbSchema = v
	})

	t.Run("AutoMigrate", func(t *testing.T) {
		err = db.Migrator().AutoMigrate(dbSchema)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("DropTable", func(t *testing.T) {
		err = db.Migrator().DropTable(dbSchema)
		if err != nil {
			t.Error(err)
		}
	})

}
