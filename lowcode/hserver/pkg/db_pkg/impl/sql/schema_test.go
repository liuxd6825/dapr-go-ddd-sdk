package sql

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	dbschema "gorm.io/gorm/schema"
	"testing"
)

const HumanSchema = `
	{
	  "name": "humans",
	  "type": "object",
	  "required": ["id","name"],
	  "description": "人员基本信息",
	  "allOf": [
		{"$ref": "/definition/db/base.json"}
	  ],
	  "properties": {
		"age": {
		  "type": ["integer", "null"],
		  "title": "年龄",
		  "order": 100
		},
		"analyse": {
		  "type": ["string", "null"],
		  "title": "分析状态",
		  "order": 101
		},
		"birthday": {
		  "type": ["date", "null"],
		  "title": "出生日期",
		  "order": 102
		},
		"gender": {
		  "type": ["string", "null"],
		  "title": "性别",
		  "order": 103
		},
		"name": {
		  "type": ["string", "null"],
		  "title": "姓名",
		  "order": 104
		},
		"peopleType": {
		  "type": "array",
		  "items": { "type": "string" },
		  "title": "人员类型",
		  "order": 105
		},
		"tags": {
		  "title": "标签",
		  "type": "array",
		  "items": { "type": "string" },
		  "order": 106
		}
	  }
	}
	`

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

	dest, err := schema.NewSchemaWithJson("human.json", HumanSchema)
	if err != nil {
		t.Error(err)
		return
	}

	var dbSchema *dbschema.Schema
	t.Run("NewDBSchema", func(t *testing.T) {
		v, err := NewDBSchema(dest)
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
