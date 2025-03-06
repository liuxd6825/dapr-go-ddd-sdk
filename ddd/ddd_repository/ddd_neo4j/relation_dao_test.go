package ddd_neo4j

import (
	"context"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"testing"
)

type CompanyRelationDao struct {
	Dao[*CompanyRelation]
}

type CompanyRelation struct {
	BaseRelation
	Name    string `json:"name"`
	Display bool   `json:"display"`
}

func newRelationItem() Relation {
	return &CompanyRelation{}
}

func TestRelationDao(t *testing.T) {
	ctx := context.Background()
	tenantId := "test"

	dao := NewCompanyRelationDao()

	rel := &CompanyRelation{}
	rel.Id = uuid.New().String()
	rel.Display = true
	rel.StartId = "111abd9f-5392-4928-bab5-27fd688f4824"
	rel.EndId = "90a49b8e-953c-4135-9690-f3f4daa54dc6"
	// rel.Type = "A"
	rel.TenantId = tenantId
	rel.Name = "TableName"

	t.Run("Insert", func(t *testing.T) {
		if res, err := dao.Insert(ctx, rel).Result(); err != nil {
			t.Error(err)
		} else {
			t.Log(res)
		}
	})

	t.Run("FindById", func(t *testing.T) {
		if v := dao.FindById(ctx, rel.TenantId, rel.Id); v.Error != nil {
			t.Error(v.Error)
		} else if !v.IsFound {
			t.Error("Not Found ")
		} else {
			t.Log(v)
		}
	})

	t.Run("FindByFilter", func(t *testing.T) {
		filter := "name=='TableName'"
		if vList, ok, err := dao.FindByRSQL(ctx, tenantId, filter).Result(); err != nil {
			t.Error(err)
		} else if !ok {
			t.Log("Not Found ")
		} else {
			t.Log(vList)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		if vList, ok, err := dao.FindAll(ctx, tenantId).Result(); err != nil {
			t.Error(err)
		} else if !ok {
			t.Log("Not Found ")
		} else {
			t.Log(vList)
		}
	})
}

func NewCompanyRelationDao() ddd_repository.Dao[*CompanyRelation] {
	eb := NewRelationEntityBuilder[*CompanyRelation](nil)
	return NewRelationDao[*CompanyRelation](driver, "C", eb)
}
