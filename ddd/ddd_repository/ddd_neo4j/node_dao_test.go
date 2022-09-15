package ddd_neo4j

import (
	"context"
	"github.com/google/uuid"
	"log"
	"testing"
)

type CompanyNode struct {
	BaseNode
	CaseId    string `json:"caseId"`
	GraphId   string `json:"graphId"`
	IsVisible bool   `json:"isVisible"`
	Key       string `json:"key"`
	Name      string `json:"name"`
}

type CompanyNodeDao struct {
	NodeDao[*CompanyNode]
}

func TestNodeDao(t *testing.T) {
	ctx := context.Background()
	tenantId := "test"

	dao := NewNodeDao[*CompanyNode](driver, NewNodeCypher("CompanyNode"))
	company := &CompanyNode{}
	company.Id = uuid.New().String()
	company.Name = "company"
	company.TenantId = tenantId
	company.CaseId = "caseId001"
	company.GraphId = "graphId001"

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
		if findV, ok, err := dao.FindById(ctx, tenantId, company.Id); err != nil {
			t.Error(err)
			return
		} else if ok {
			t.Logf("findV.id = %v", findV.GetId())
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

	t.Run("DeleteById", func(t *testing.T) {
		if err := dao.DeleteById(ctx, tenantId, company.Id); err != nil {
			t.Error(err)
		}
	})

}

// Neo4j-测试获取结果集
func TestGetList(t *testing.T) {
	repos := NewNodeDao[*CompanyNode](driver, NewNodeCypher("TestNode"))
	cypher := "MATCH (n:graph_T1_N3eb0982799464cf199f2182d130e4a32_company)-[r*0..]->(m) RETURN n, r "
	if result, err := repos.Query(context.Background(), cypher, nil); err != nil {
		log.Println("error connecting to neo4j:", err)
	} else {
		var comps []CompanyNode
		var rels []CompanyRelation
		if err := result.GetLists([]string{"n", "r"}, &comps, &rels); err != nil {
			t.Error(err)
			return
		}
		log.Printf("rels.length = %d ; company.length=%d", len(comps), len(rels))
	}
}

func NewCompanyNodeDao() *CompanyNodeDao {
	dao := &CompanyNodeDao{}
	dao.init(driver, NewNodeCypher("CompanyNode"))
	return dao
}
