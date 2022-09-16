package ddd_neo4j

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

type CompanyRelationDao struct {
	RelationDao[*CompanyRelation]
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

	dao := NewRelationDao[*CompanyRelation](driver, NewRelationCypher("C"))

	rel := &CompanyRelation{}
	rel.Id = uuid.New().String()
	rel.Display = true
	rel.StartId = "111abd9f-5392-4928-bab5-27fd688f4824"
	rel.EndId = "90a49b8e-953c-4135-9690-f3f4daa54dc6"
	// rel.Type = "A"
	rel.TenantId = tenantId
	rel.Name = "Name"

	t.Run("Insert", func(t *testing.T) {
		if res, err := dao.Insert(ctx, rel).Result(); err != nil {
			t.Error(err)
		} else {
			t.Log(res)
		}
	})

	t.Run("FindById", func(t *testing.T) {
		if v, ok, err := dao.FindById(ctx, rel.TenantId, rel.Id); err != nil {
			t.Error(err)
		} else if !ok {
			t.Error("Not Found ")
		} else {
			t.Log(v)
		}
	})

	t.Run("FindByFilter", func(t *testing.T) {
		filter := "name=='Name'"
		if vList, ok, err := dao.FindByFilter(ctx, tenantId, filter).Result(); err != nil {
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
