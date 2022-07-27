package ddd_neo4j

import (
	"context"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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

func newCompanyNode() ElementEntity {
	return &CompanyNode{}
}

type Rel struct {
	BaseRelationship
	Id string
}

var (
	neo4jURL = "bolt://192.168.64.4:7687"
	username = "neo4j"
	password = "123456"
)

func CreateDriver(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
}

func CloseDriver(driver neo4j.Driver) error {
	return driver.Close()
}

func TestWriteNode(t *testing.T) {

}

// Neo4j-测试获取结果集
func TestGetList(t *testing.T) {
	driver, err := CreateDriver(neo4jURL, username, password)
	defer func(driver neo4j.Driver) {
		err = CloseDriver(driver)
		if err != nil {
			log.Println("neo4j close error:", err)
		}
	}(driver)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}

	repos := NewNeo4jDao[*CompanyNode](driver, NewReflectBuilder("CompanyNode"))
	cypher := "MATCH (n:graph_T1_N3eb0982799464cf199f2182d130e4a32_company)-[r*0..]->(m) RETURN n, r "
	if result, err := repos.Query(context.Background(), cypher); err != nil {
		log.Println("error connecting to neo4j:", err)
	} else {
		var comps []CompanyNode
		var rels []Rel
		if err := result.GetLists([]string{"n", "r"}, &comps, &rels); err != nil {
			t.Error(err)
			return
		}
		log.Printf("rels.length = %d ; company.length=%d", len(comps), len(rels))
	}
}

type CompanyRepository[T interface{ *CompanyNode }] struct {
	Neo4jDao[*CompanyNode]
}

func NewCompanyRepository(driver neo4j.Driver) *CompanyRepository[*CompanyNode] {
	resp := &CompanyRepository[*CompanyNode]{}
	build := NewReflectBuilder("TestCompanyNode")
	resp.Init(driver, build)
	return resp
}

func TestNeo4JDao_Insert(t *testing.T) {
	driver, err := CreateDriver(neo4jURL, username, password)
	defer func(driver neo4j.Driver) {
		if e := recover(); e != nil {
			if err := ddd_errors.GetRecoverError(e); err != nil {
				t.Error(err)
			}
		}
		err = CloseDriver(driver)
		if err != nil {
			log.Println("neo4j close error:", err)
		}
	}(driver)

	if err != nil {
		t.Error(err)
		return
	}

	repos := NewCompanyRepository(driver)
	/*	if com, err := repos.NewEntity(); err != nil {
			t.Error(err)
			return
		} else {
			t.Logf("new %v", com)
		}*/

	company := &CompanyNode{}
	company.Id = uuid.New().String()
	company.Labels = []string{"TestCompanyNode"}
	company.Name = "company"
	company.TenantId = "001"
	company.CaseId = "caseId001"
	company.GraphId = "graphId001"

	ctx := context.Background()

	/*	err = repos.Insert(ctx, company).GetError()
		if err != nil {
			t.Error(err)
		}

		company.Key = "keys"
		err = repos.Update(ctx, company).GetError()
		if err != nil {
			t.Error(err)
			return
		}

		if company, err = repos.FindById(ctx, "001", company.Id); err != nil {
			t.Error(err)
			return
		} else {
			t.Logf("company.id = %v", company.Id)
		}
	*/
	if res := repos.FindByGraphId(ctx, "001", "graphId001"); res.GetError() != nil {
		t.Error(res.GetError())
		return
	} else {
		t.Logf("graph list = %v ", res.GetData())
	}

	err = repos.DeleteById(ctx, company).GetError()
	if err != nil {
		t.Error(err)
	}

}
