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
	Name    string
	Display bool
}

func TestRelationDao(t *testing.T) {
	ctx := context.Background()
	tenantId := "test"

	dao := NewRelationDao[*CompanyRelation](driver, NewRelationCypher("CompanyNode"))

	t.Run("Insert", func(t *testing.T) {
		rel := &CompanyRelation{}
		rel.Id = uuid.New().String()
		rel.Display = true
		rel.StartId = "001"
		rel.EndId = "002"
		rel.Type = "type"
		rel.TenantId = tenantId
		rel.Name = "Name"

		if res, err := dao.Insert(ctx, rel).Result(); err != nil {
			t.Error(err)
		} else {
			t.Log(res)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		dao.
	})
}
