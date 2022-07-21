package ddd_neo4j

import (
	"context"
	"github.com/google/uuid"
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

	repos := NewBaseRepository[*CompanyNode](driver, NewReflectBuilder())
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
	BaseRepository[*CompanyNode]
}

func NewCompanyRepository(driver neo4j.Driver) *CompanyRepository[*CompanyNode] {
	resp := &CompanyRepository[*CompanyNode]{}
	build := NewReflectBuilder()
	resp.Init(driver, build)
	return resp
}

func Test_Insert(t *testing.T) {
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

	company := &CompanyNode{}
	company.Id = uuid.New().String()
	company.Labels = []string{"TestCompanyNode"}
	company.Name = "company"
	company.TenantId = "001"
	company.CaseId = "caseId001"
	company.GraphId = "graphId001"

	repos := NewCompanyRepository(driver)
	err = repos.Insert(context.Background(), company).GetError()
	if err != nil {
		t.Error(err)
	}

	company.Key = "keys"
	err = repos.Update(context.Background(), company).GetError()
	if err != nil {
		t.Error(err)
		return
	}

	if company, err = repos.FindById(context.Background(), company.TenantId, company.Id); err != nil {
		t.Error(err)
		return
	} else {
		t.Logf("company.id = %v", company.Id)
	}

	err = repos.DeleteById(context.Background(), company).GetError()
	if err != nil {
		t.Error(err)
	}

}
