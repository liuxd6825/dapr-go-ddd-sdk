package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"strings"
)

type relationCypher struct {
	labels        string
	isEmptyLabels bool
}

//
//  NewRelationCypher
//  @Description:
//  @param labels 关系标签，可以为空值；为空：由Relation.GetRelType()决定标签名称
//  @return Cypher
//
func NewRelationCypher(labels string) Cypher {
	return &relationCypher{
		labels:        getLabels(labels),
		isEmptyLabels: len(labels) == 0,
	}
}

func (c *relationCypher) Insert(ctx context.Context, data interface{}) (CypherResult, error) {
	rel := data.(Relation)
	props, dataMap, err := getCreateProperties(ctx, data)
	if err != nil {
		return nil, err
	}
	labels := c.getLabels(rel.GetRelType())
	cypher := fmt.Sprintf(`
	MATCH (a{tenantId:'%v'}),(b{tenantId:'%v'})
	WHERE a.id = '%v' AND b.id = '%v'
	CREATE (a)-[r%v{%v}]->(b)
	RETURN r`, rel.GetTenantId(), rel.GetTenantId(), rel.GetStartId(), rel.GetEndId(), labels, props)

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

	labels := c.getLabels(rel.GetRelType())

	cypher := fmt.Sprintf("MATCH (a)-[r%v{tenantId:'%v',id:'%v'}]-(b) SET %s ", labels, rel.GetTenantId(), rel.GetId(), prosNames)
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
	strIds := getSqlInStr(ids)
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) WHERE r.id in [%v] delete r `, tenantId, strIds)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) DeleteAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) delete r `, tenantId)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) DeleteByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getSqlWhere(tenantId, filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) %v delete r `, tenantId, where)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) FindById(ctx context.Context, tenantId, id string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v',id:'%v'}]->(b) RETURN r `, tenantId, id)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	strIds := getSqlInStr(ids)
	cypher := fmt.Sprintf("MATCH (a)-[r{tenantId:'%v'}]-(b) where r.id in [%v] RETURN r", tenantId, strIds)
	return NewCypherBuilderResult(cypher, nil, []string{"r"}), nil
}

func (c *relationCypher) FindByAggregateId(ctx context.Context, tenantId, aggregateName, aggregateId string) (result CypherResult, err error) {
	return c.DeleteByFilter(ctx, tenantId, fmt.Sprintf("%v=='%v'", aggregateName, aggregateId))
}

func (c *relationCypher) FindByGraphId(ctx context.Context, tenantId string, graphId string) (result CypherResult, err error) {
	return c.FindByFilter(ctx, tenantId, fmt.Sprintf("graphId=='%v'", graphId))
}

func (c *relationCypher) FindAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r%v{tenantId:'%v'}]->(b) RETURN r `, c.getLabels(""), tenantId)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery) (CypherResult, error) {
	where, err := getSqlWhere(query.GetTenantId(), query.GetFilter())
	if err != nil {
		return nil, err
	}
	skip := query.GetPageNum() * query.GetPageSize()
	pageSize := query.GetPageSize()

	keys := []string{"n"}
	count := ""
	if query.GetIsTotalRows() {
		count = ", count(n) as count "
		keys = append(keys, "count")
	}

	order, err := getOrder(query.GetSort())
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (a)-(r%v{tenantId:'%v'})->(b) %v RETURN r %v %v SKIP %v LIMIT %v ", c.labels, query.GetTenantId(), where, count, order, skip, pageSize)
	return NewCypherBuilderResult(cypher, nil, keys), nil
}

func (c *relationCypher) FindByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getSqlWhere(tenantId, filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf(`MATCH (a{tenantId:'%v'})-[n%v{tenantId:'%v'}]-(b{tenantId:'%v'}) %v return n `, tenantId, c.getLabels(""), tenantId, tenantId, where)
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher) Count(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getSqlWhere(tenantId, filter)
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (a{tenantId:'%v'})-[n%v]->(b{tenantId:'%v'}) %v RETURN count(n) as n  ", tenantId, c.labels, tenantId, where)
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher) GetFilter(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getSqlWhere(tenantId, filter)
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (a{tenantId:'%v'})-[n%v]->(b{tenantId:'%v'}) %v RETURN n  ", tenantId, c.labels, tenantId, where)
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher) getLabels(labels string) string {
	if c.isEmptyLabels && len(labels) == 0 {
		return ""
	} else if c.isEmptyLabels {
		return ":" + labels
	}
	return c.labels
}

// getIds 将id数组，转换成sql形式。如：'111','222'。
func getSqlInStr(ids []string) string {
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	strIds := strings.Join(ids, ",")
	return strIds
}
