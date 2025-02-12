package ddd_sql

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type Entity = map[string]any

type User struct {
	ID       string
	Name     string
	Age      int
	Email    string
	TenantID string
}

func Test_Insert(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
		return
	}

	entityBuilder := NewMapEntityBuilder()

	err = db.AutoMigrate(&User{})
	if err != nil {
		t.Error(err)
		return
	}

	err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id TEXT  PRIMARY KEY AUTOINCREMENT,
            tenant_id TEXT,
            name TEXT,
            age INTEGER,
            email TEXT
        )
	`).Error
	if err != nil {
		t.Error(err)
		return
	}

	dao := NewDao[MapEntity](db, entityBuilder, "users")
	ctx := context.Background()

	user := NewMapEntity()
	user["tenant_id"] = "test"
	user["name"] = "name"
	user["id"] = idutils.NewId()
	res := dao.Insert(ctx, user)
	if res.Error != nil {
		t.Error(res.Error)
		return
	}

	listRes := dao.FindAll(ctx, "test")
	if listRes.Error != nil {
		t.Error(listRes.Error)
		return
	}
	print(listRes.Data)

	idRes := dao.FindById(ctx, "test", "7MC10GTPH63KJJN19YYAREZDS5CJMF")
	if idRes.Error != nil {
		t.Error(idRes.Error)
		return
	}
	print(idRes.Data)
}
