package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"strings"
)

type relationCypher struct {
	labels        string
	isEmptyLabels bool
}

// NewRelationCypher
// @Description:
// @param labels 关系标签，可以为空值；为空：由Relation.GetRelType()决定标签名称
// @return Cypher
func NewRelationCypher(labels string) Cypher {
	return &relationCypher{
		labels:        getLabels(labels),
		isEmptyLabels: len(labels) == 0,
	}
}

func (c *relationCypher) Insert(ctx context.Context, data interface{}) (CypherResult, error) {
	rel, ok := data.(Relation)
	if !ok {
		return nil, errors.ErrorOf(" parameter data is not ddd_neo4j.Relation Type")
	}
	props, dataMap, err := c.getCreateProperties(ctx, data)
	if err != nil {
		return nil, err
	}
	labels := c.getLabels(rel.GetRelType())
	cypher := fmt.Sprintf(`
	MATCH (a{tenantId:'%v'}),(b{tenantId:'%v'})
	WHERE a.id = '%v' AND b.id = '%v'
	CREATE (a)-[r%v{%v}]->(b)
	RETURN r`, rel.GetTenantId(), rel.GetTenantId(), rel.GetStartId(), rel.GetEndId(), labels, props)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, dataMap, nil), nil
}

func (c *relationCypher) InsertOrUpdate(ctx context.Context, data interface{}) (CypherResult, error) {
	rel, ok := data.(Relation)
	if !ok {
		return nil, errors.ErrorOf(" parameter data is not ddd_neo4j.Relation Type")
	}
	props, dataMap, err := c.getSetFields(ctx, "r", data)
	if err != nil {
		return nil, err
	}
	labels := c.getLabels(rel.GetRelType())
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("MATCH (s{id:'%v'}), (e{id:'%v'}) ", rel.GetStartId(), rel.GetEndId()))
	sb.WriteString(fmt.Sprintf("MERGE (s)-[r%v{id:'%v'}]->(e) ", labels, rel.GetId()))
	sb.WriteString(fmt.Sprintf("ON CREATE SET %s ", props))
	sb.WriteString(fmt.Sprintf("ON MATCH  SET %s ", props))
	logs.Debug(ctx, func() any { return sb.String() })
	return NewCypherBuilderResult(sb.String(), dataMap, nil), nil
}

func (c *relationCypher) InsertMany(ctx context.Context, list interface{}) (CypherResult, error) {
	rels := list.([]Relation)
	println(rels)

	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) Update(ctx context.Context, data interface{}, setFields ...string) (CypherResult, error) {

	// 只更新关系标签
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r

	// 更新关系的标签与属性
	// match(n)-[r:relation{id:'cbc4d7be-43fa-427e-956d-e812b335bc12'}]->(m) create (n)-[r2:relation]->(m) set r2=r, r2.title='title' with r delete r

	rel := data.(Relation)
	prosNames, mapData, err := getUpdateProperties(ctx, data, "r2", setFields...)
	if err != nil {
		return nil, err
	}

	labels := c.getLabels(rel.GetRelType())
	if len(labels) == 0 {
		return nil, errors.New("neo4j relation.Type is nil")
	}

	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]->(m) CREATE (n)-[r2%s]->(m) SET r2=r, %s  WITH r DELETE r ", rel.GetTenantId(), rel.GetId(), labels, prosNames)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (c *relationCypher) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]->(n) CREATE (n)-[r2:%v]-(m) SET r2=r WITH r DELETE r ", tenantId, id, label)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) UpdateMany(ctx context.Context, list interface{}) (CypherResult, error) {
	rels := list.([]Relation)
	println(rels)

	//TODO implement me
	panic("implement me")
}

func (c *relationCypher) UpdateLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v'}]-(n) create (n)-[r2:{tenantId:'%v'}]-(m)  %s SET r2=r WITH r DELETE r ", tenantId, tenantId, where)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) DeleteLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// neo4j 不支持删除关系标签
	return nil, nil
}

func (c *relationCypher) DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	// neo4j 不支持删除关系标签
	return nil, nil
}

func (c *relationCypher) DeleteByLabels(ctx context.Context, tenantId string, label ...string) (CypherResult, error) {
	// neo4j 不支持删除关系标签
	return nil, nil
}

func (c *relationCypher) DeleteByTenantId(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v',id:'%v'}]-(b) delete r `, tenantId)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) DeleteById(ctx context.Context, tenantId string, id string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v',id:'%v'}]-(b) delete r `, tenantId, id)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	strIds := getSqlInStr(ids)
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) WHERE r.id in [%v] delete r `, tenantId, strIds)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) DeleteAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) delete r `, tenantId)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) DeleteByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "r", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) %v delete r `, tenantId, where)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) FindById(ctx context.Context, tenantId, id string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v',id:'%v'}]->(b) RETURN r `, tenantId, id)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	strIds := getSqlInStr(ids)
	cypher := fmt.Sprintf("MATCH (a)-[r{tenantId:'%v'}]-(b) where r.id in [%v] RETURN r", tenantId, strIds)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, []string{"r"}), nil
}

func (c *relationCypher) FindByAggregateId(ctx context.Context, tenantId, aggregateName, aggregateId string) (result CypherResult, err error) {
	return c.DeleteByFilter(ctx, tenantId, fmt.Sprintf("%v=='%v'", aggregateName, aggregateId))
}

func (c *relationCypher) FindByGraphId(ctx context.Context, tenantId string, graphId string) (result CypherResult, err error) {
	return c.FindByFilter(ctx, tenantId, fmt.Sprintf("graphId=='%v'", graphId))
}

func (c *relationCypher) FindByCaseId(ctx context.Context, tenantId string, caseId string) (result CypherResult, err error) {
	return c.FindByFilter(ctx, tenantId, fmt.Sprintf("caseId=='%v'", caseId))
}

func (c *relationCypher) FindAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r%v{tenantId:'%v'}]->(b) RETURN r `, c.getLabels(""), tenantId)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher) FindPaging(ctx context.Context, qry ddd_repository.FindPagingQuery) (CypherResult, error) {
	where, err := getNeo4jWhere(qry.GetTenantId(), "r", qry.GetFilter())
	if err != nil {
		return nil, err
	}
	skip := qry.GetPageNum() * qry.GetPageSize()
	pageSize := qry.GetPageSize()
	keys := []string{"r", "count"}
	order, err := getOrder(qry.GetSort())
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH ()-[r%v{tenantId:'%v'}]->() %v RETURN r %s SKIP %v LIMIT %v ", c.getLabels(c.labels), qry.GetTenantId(), where, order, skip, pageSize)
	countCypher := fmt.Sprintf("MATCH ()-[r%v{tenantId:'%v'}]->() %v RETURN count(r) as count  ", c.getLabels(c.labels), qry.GetTenantId(), where)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, keys, NewCypherResultOptions().SetCountCypher(countCypher)), nil
}

func (c *relationCypher) FindByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf(`MATCH (a{tenantId:'%v'})-[n%v{tenantId:'%v'}]-(b{tenantId:'%v'}) %v return n `, tenantId, c.getLabels(""), tenantId, tenantId, where)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher) Count(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (a{tenantId:'%v'})-[n%v]->(b{tenantId:'%v'}) %v RETURN count(n) as n  ", tenantId, c.labels, tenantId, where)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher) GetFilter(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (a{tenantId:'%v'})-[n%v]->(b{tenantId:'%v'}) %v RETURN n  ", tenantId, c.getLabels(""), tenantId, where)
	logs.Debug(ctx, cypher)
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher) getLabels(labels string) string {
	if c.isEmptyLabels && len(labels) == 0 {
		return ""
	} else if c.isEmptyLabels {
		return fmt.Sprintf(":`%v`", labels)
	}
	return c.labels
}

func (c *relationCypher) getCreateProperties(ctx context.Context, data interface{}) (string, map[string]any, error) {
	mapData, err := getMap(data)
	if err != nil {
		return "", nil, err
	}

	strBuilder := strings.Builder{}
	for k := range mapData {
		if k != "properties" {
			strBuilder.WriteString(fmt.Sprintf(`%s:$%s,`, k, k))
		}
	}

	if relation, ok := data.(Relation); ok {
		for k := range relation.GetProperties() {
			strBuilder.WriteString(fmt.Sprintf(`%s:$properties.%s,`, k, k))
		}
	}

	res := strBuilder.String()
	if len(res) > 0 {
		res = res[:len(res)-1]
	}
	return res, mapData, nil
}

func (c *relationCypher) getSetFields(ctx context.Context, resName string, data interface{}) (string, map[string]any, error) {
	mapData, err := getMap(data)
	if err != nil {
		return "", nil, err
	}

	strBuilder := strings.Builder{}
	for k := range mapData {
		if k != "properties" {
			strBuilder.WriteString(fmt.Sprintf(`%s.%s=$%s,`, resName, k, k))
		}
	}

	if relation, ok := data.(Relation); ok {
		for k := range relation.GetProperties() {
			strBuilder.WriteString(fmt.Sprintf(`%s.%s=$properties.%s,`, resName, k, k))
		}
	}

	res := strBuilder.String()
	if len(res) > 0 {
		res = res[:len(res)-1]
	}
	return res, mapData, nil
}

// getIds 将id数组，转换成sql形式。如：'111','222'。
func getSqlInStr(ids []string) string {
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	strIds := strings.Join(ids, ",")
	return strIds
}
