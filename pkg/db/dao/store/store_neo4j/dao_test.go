package store_neo4j

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/randomutils"
	"github.com/stretchr/testify/assert"
)

const test = "test"

func TestDao(t *testing.T) {

	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}

	dao := NewDao[map[string]any](driver, NewNodeCypher(&Config[map[string]any]{}), nil)
	id := randomutils.NewId()
	name := randomutils.NameCN()

	t.Run("Insert", func(t *testing.T) {
		user := map[string]any{}
		user["id"] = id
		user["tenantId"] = "test"
		user["name"] = name
		res := dao.Insert(ctx, user)
		assert.NoError(t, res.Error)
		assert.Equal(t, int64(1), res.RowsAffected)
	})

	t.Run("Update", func(t *testing.T) {
		user := map[string]any{}
		user["id"] = id
		user["tenantId"] = "test"
		user["name"] = name + "2"

		res := dao.Update(ctx, user)
		assert.NoError(t, res.Error)
		assert.Equal(t, int64(1), res.RowsAffected)
	})

	t.Run("FindAll", func(t *testing.T) {
		res := dao.FindAll(ctx, test)
		assert.NoError(t, res.Error)
		t.Log(marshal(res.Data))
	})

	t.Run("FindById", func(t *testing.T) {
		res := dao.FindById(ctx, test, id)
		assert.NoError(t, res.Error)
		if res.Data != nil {
			t.Log(marshal(res.Data))
		}
	})

	t.Run("FindPaging", func(t *testing.T) {
		qry := ddd_query.NewFindPagingQuery()
		qry.SetTenantId(test)
		qry.SetFilter(rsql.NewBuilder().Eq("id", id).Build())

		pagingRes := dao.FindPaging(ctx, qry)
		assert.NoError(t, pagingRes.Error)
		assert.NotNil(t, pagingRes.Data)
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
		res := dao.FindByIds(ctx, test, []string{id})
		assert.NoError(t, res.Error)
		t.Log(marshal(res.Data))
	})

	/*
		t.Run("Sum", func(t *testing.T) {
			qry := ddd_query.NewFindPagingQuery()
			qry.SetTenantId(test)
			valueCols := []*store.ValueCol{}
			valueCols = append(valueCols, &store.ValueCol{AggFunc: store.AggFuncSum, Field: "age"})
			valueCols = append(valueCols, &store.ValueCol{AggFunc: store.AggFuncSum, Field: "score"})
			qry.SetValueCols(valueCols)
			sumData := map[string]any{}
			sumRes, _, sumErr := dao.SumEntity(ctx, qry, &sumData)
			assert.NoError(t, sumErr)
			t.Log(marshal(sumRes))
		})
	*/

	t.Run("FindAutoComplete", func(t *testing.T) {
		qry := ddd_query.NewFindAutoCompleteQuery()
		qry.SetTenantId(test)
		qry.SetFilter(fmt.Sprintf("name ==~ \"%s\"", "K"))
		res := dao.FindAutoComplete(ctx, qry)
		assert.NoError(t, res.Error)
		t.Log(marshal(res.Data))
	})

	t.Run("DeleteById", func(t *testing.T) {
		res := dao.DeleteById(ctx, "test", id)
		assert.NoError(t, res.Error)
		assert.Equal(t, int64(1), res.RowsAffected)
	})
}

func marshal(data any) string {
	dataJson, _ := json.Marshal(data)
	return string(dataJson)
}
