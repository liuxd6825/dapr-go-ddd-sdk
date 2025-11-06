package sql

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
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

	dest := schema.NewJsonSchemaWithJson("human.json", xtest.HumanSchema)

	var dbSchema *gormschema.Schema
	t.Run("NewDBSchema", func(t *testing.T) {
		dschema := dbschema.NewDBSchemaWithJsonSchema(dest)
		dschema.Name = "table"
		v, err := NewGormSchema(dschema)
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
