package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"strings"
)

type relationCypher struct {
	labels string
}

func NewRelationCypher(labels ...string) Cypher {
	var s string
	for _, l := range labels {
		s = s + ":" + l
	}
	return &relationCypher{labels: s}
}

func (c *relationCypher) Insert(ctx context.Context, data interface{}) (CypherResult, error) {
	rel := data.(Relation)
	props, dataMap, err := getCreateProperties(ctx, data)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf(`
	MATCH (a{tenantId:'%v'}),(b{tenantId:'%v'})
	WHERE a.id = '%v' AND b.id = '%v'
	CREATE (a)-[r:%v{%v}]->(b)
	RETURN r`, rel.GetTenantId(), rel.GetTenantId(), rel.GetStartId(), rel.GetEndId(), rel.GetType(), props)

	return NewCypherBuilderResult(cypher, dataMap, nil), nil
}

func (c *relationCypher) InsertMany(ctx context.Context, list interface{}) (CypherResult, error) {
	rels := list.([]Relation)
	println(rels)

	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) Update(ctx context.Context, data interface{}, setFields ...string) (CypherResult, error) {
	rel := data.(Relation)
	prosNames, mapData, err := getUpdateProperties(ctx, data, setFields...)
	if err != nil {
		return nil, err
	}
	// `MATCH (a)-[r{id:'r001'}]-(b) set r.name='rName' return a, r, b`
	// cypher := fmt.Sprintf("MATCH (n{id:$id}) SET %s RETURN n ", prosNames)
	cypher := fmt.Sprintf("MATCH (a)-[r{tenantId:'%v',id:'%v'}]-(b) SET %s ", rel.GetTenantId(), rel.GetId(), prosNames)
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (c *relationCypher) UpdateMany(ctx context.Context, list interface{}) (CypherResult, error) {
	rels := list.([]Relation)
	println(rels)

	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) DeleteById(ctx context.Context, tenantId string, id string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v',id:'%v'}]-(b) delete r `, tenantId, id)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) DeleteAll(ctx context.Context, tenantId string) (CypherResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) DeleteByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) GetFilter(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) FindById(ctx context.Context, tenantId, id string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v',id:'%v'}]-(b) RETURN r `, tenantId, id)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	strIds := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (a{tenantId:'%v'})-[r{tenantId:'%v'}]-(b{tenantId:'%v'}) where r.id in [%v] RETURN r", tenantId, tenantId, tenantId, strIds)
	return NewCypherBuilderResult(cypher, nil, []string{"r"}), nil
}

func (c *relationCypher) FindByAggregateId(ctx context.Context, tenantId, aggregateName, aggregateId string) (result CypherResult, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) FindByGraphId(ctx context.Context, tenantId string, graphId string) (result CypherResult, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) FindAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a{tenantId:'%v'})-[r{tenantId:'%v'}]-(b{tenantId:'%v'}) RETURN r `, tenantId, tenantId, tenantId)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery) (CypherResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) Count(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) GetLabels() string {
	//TODO implement me
	panic("implement me")
}
