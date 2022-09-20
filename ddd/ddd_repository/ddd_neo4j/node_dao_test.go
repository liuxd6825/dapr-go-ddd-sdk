package ddd_neo4j

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"testing"
)

type CompanyNode struct {
	BaseNode
	CaseId    string `json:"caseId"`
	GraphId   string `json:"graphId"`
	IsVisible bool   `json:"isVisible"`
	Name      string `json:"name"`
}

func TestNodeDao(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	tenantId := "test"

	dao := NewCompanyNodeDao()
	company := &CompanyNode{}
	company.Id = id
	company.Name = "company"
	company.TenantId = tenantId
	company.CaseId = "caseId001"
	company.GraphId = "graphId001"

	println("id:" + company.Id)

	t.Run("graph", func(t *testing.T) {
		match := "match  ()-[r:`分公司`|`母公司`]->() return r"
		if res, err := dao.Query(ctx, match, nil); err != nil {
			t.Error(err)
		} else {
			for k, v := range res.data {
				t.Logf("key:%v; count:%v", k, len(v))
			}
		}

		match = "match (n:`公司`) return n"
		if res, err := dao.Query(ctx, match, nil); err != nil {
			t.Error(err)
		} else {
			for k, v := range res.data {
				t.Logf("key:%v; count:%v", k, len(v))
			}
		}
	})

	t.Run("Insert", func(t *testing.T) {
		if err := dao.Insert(ctx, company).GetError(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		company.Name = "company_update"
		if err := dao.Update(ctx, company).GetError(); err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("FindById", func(t *testing.T) {
		if findV, ok, err := dao.FindById(ctx, tenantId, id); err != nil {
			t.Error(err)
			return
		} else if ok {
			t.Logf("findV.id = %v", findV.GetId())
		} else {
			t.Error(errors.New("not found "))
		}
	})

	t.Run("FindByGraphId", func(t *testing.T) {
		if res := dao.FindByGraphId(ctx, tenantId, "graphId001"); res.GetError() != nil {
			t.Error(res.GetError())
			return
		} else {
			t.Logf("graph list = %v ", res.GetData())
		}
	})

	/*	t.Run("DeleteById", func(t *testing.T) {
			if err := dao.DeleteById(ctx, tenantId, company.Id); err != nil {
				t.Error(err)
			}
		})
	*/
}

// Neo4j-测试获取结果集
func TestGetList(t *testing.T) {
	repos := NewCompanyNodeDao()
	cypher := "MATCH (n:graph_T1_N3eb0982799464cf199f2182d130e4a32_company)-[r*0..]->(m) RETURN n, r "
	if result, err := repos.Query(context.Background(), cypher, nil); err != nil {
		log.Println("error connecting to neo4j:", err)
	} else {
		var comps []CompanyNode
		var rels []CompanyRelation
		if err := result.GetList("n", &comps); err != nil {
			t.Error(err)
			return
		}
		log.Printf("rels.length = %d ; company.length=%d", len(comps), len(rels))
	}
}

func NewCompanyNodeDao() *Dao[*CompanyNode] {
	return NewNodeDao[*CompanyNode](driver, NewNodeCypher("CompanyNode"))
}
