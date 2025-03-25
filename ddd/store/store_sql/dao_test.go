package store_sql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type Entity = map[string]any

type User struct {
	Id       string
	Name     string
	Age      int
	Score    int
	Email    string
	TenantId string
}

const test = "test"

func Test_Dao(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
		return
	}

	err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id TEXT  PRIMARY KEY,
            tenant_id TEXT,
            name TEXT,
            age INTEGER,
            score INTEGER,
            email TEXT,
			created_time TEXT
			creator_id  TEXT,
			creator_name TEXT,
			updated_time TEXT,
			updater_id TEXT,
			updater_name TEXT,
			deleted_time TEXT,
			deleter_id TEXT,
			deleter_name TEXT,
			is_deleted INTEGER
        )
	`).Error

	if err != nil {
		t.Error(err)
		return
	}

	eb := store.NewAnyEntityBuilder[map[string]any]()
	dao := NewDao[map[string]any](db, "dbKey", eb, "users")
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	id := randomutils.NewId()
	name := randomutils.NameCN()

	t.Run("Insert", func(t *testing.T) {
		user := map[string]any{}
		user["id"] = id
		user["tenant_id"] = "test"
		user["name"] = name

		res := dao.Insert(ctx, user)
		if res.Error != nil {
			t.Error(res.Error)
		}
	})

	t.Run("Update", func(t *testing.T) {
		user := map[string]any{}
		user["id"] = id
		user["tenant_id"] = "test"
		user["name"] = name + "2"

		res := dao.Update(ctx, user)
		if res.Error != nil {
			t.Error(res.Error)
		}
	})

	return

	t.Run("DeleteById", func(t *testing.T) {
		res := dao.DeleteById(ctx, "test", id)
		if res.Error != nil {
			t.Error(res.Error)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		listRes := dao.FindAll(ctx, test)
		if listRes.Error != nil {
			t.Error(listRes.Error)
			return
		}
		t.Log(marshal(listRes.Data))
	})

	t.Run("FindById", func(t *testing.T) {
		idRes := dao.FindById(ctx, test, id)
		if idRes.Error != nil {
			t.Error(idRes.Error)
			return
		}
		t.Log(marshal(idRes.Data))
	})

	t.Run("FindPaging", func(t *testing.T) {
		qry := ddd_query.NewFindPagingQuery()
		qry.SetTenantId(test)
		qry.SetFilter(fmt.Sprintf("id==\"%s\"", id))
		pagingRes := dao.FindPaging(ctx, qry)
		if pagingRes.Error != nil {
			t.Error(pagingRes.Error)
			return
		}
		if pagingRes.Data == nil {
			t.Error("pagingRes.Data is nil")
			return
		}
		t.Log(marshal(pagingRes.Data))
	})

	t.Run("Count", func(t *testing.T) {
		count, err := dao.CountByRSQL(ctx, test, "")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("count = ", count)
	})

	t.Run("FindByIds", func(t *testing.T) {
		idsRes := dao.FindByIds(ctx, test, []string{id})
		if idsRes.Error != nil {
			t.Error(idsRes.Error)
			return
		}
		t.Log(marshal(idsRes.Data))
	})

	t.Run("Sum", func(t *testing.T) {
		qry := ddd_query.NewFindPagingQuery()
		qry.SetTenantId(test)
		valueCols := []*store.ValueCol{}
		valueCols = append(valueCols, &store.ValueCol{AggFunc: store.AggFuncSum, Field: "age"})
		valueCols = append(valueCols, &store.ValueCol{AggFunc: store.AggFuncSum, Field: "score"})
		qry.SetValueCols(valueCols)
		sumData := map[string]any{}
		sumRes, _, sumErr := dao.Sum(ctx, qry, &sumData)
		if sumErr != nil {
			t.Error(sumErr)
			return
		}
		t.Log(marshal(sumRes))
	})

	t.Run("FindAutoComplete", func(t *testing.T) {
		qry := ddd_query.NewFindAutoCompleteQuery()
		qry.SetTenantId(test)
		qry.SetFilter(fmt.Sprintf("name ==~ \"%s\"", "K"))
		res := dao.FindAutoComplete(ctx, qry)
		if res.Error != nil {
			t.Error(res.Error)
		}
		t.Log(marshal(res.Data))
	})

}

func marshal(data any) string {
	dataJson, _ := json.Marshal(data)
	return string(dataJson)
}
