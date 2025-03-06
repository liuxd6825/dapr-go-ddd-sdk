package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/jsonschema/v6"
	"strings"
)

type relationCypher[T any] struct {
	labels        string
	isEmptyLabels bool
	eb            RelationEntityBuilder[T]
	schema        *jsonschema.Schema
}

// NewRelationCypher
// @Description:
// @param labels 关系标签，可以为空值；为空：由Relation.GetRelType()决定标签名称
// @return Cypher
func NewRelationCypher[T any](eb RelationEntityBuilder[T], schema *jsonschema.Schema, labels ...string) Cypher[T] {
	return &relationCypher[T]{
		labels:        getLabels(labels...),
		isEmptyLabels: len(labels) == 0,
		eb:            eb,
		schema:        schema,
	}
}

func (c *relationCypher[T]) Insert(ctx context.Context, data T) (CypherResult, error) {
	props, dataMap, err := c.getCreateProperties(ctx, data)
	if err != nil {
		return nil, err
	}
	labels := c.getLabels(c.labels)
	cypher := fmt.Sprintf(`
	MATCH (a{tenantId:'%v'}),(b{tenantId:'%v'})
	WHERE a.id = '%v' AND b.id = '%v'
	CREATE (a)-[r%v{%v}]->(b)
	RETURN r`, c.eb.GetTenantId(data), c.eb.GetTenantId(data), c.eb.GetStartId(data), c.eb.GetEndId(data), labels, props)
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, dataMap, nil), nil
}

func (c *relationCypher[T]) InsertOrUpdate(ctx context.Context, data T) (CypherResult, error) {
	props, dataMap, err := c.getSetFields(ctx, "r", data)
	if err != nil {
		return nil, err
	}
	labels := c.getLabels(c.eb.GetRelType(data))
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("MATCH (s{id:'%v'}), (e{id:'%v'}) ", c.eb.GetStartId(data), c.eb.GetEndId(data)))
	sb.WriteString(fmt.Sprintf("MERGE (s)-[r%v{id:'%v'}]->(e) ", labels, c.eb.GetId(data)))
	sb.WriteString(fmt.Sprintf("ON CREATE SET %s ", props))
	sb.WriteString(fmt.Sprintf("ON MATCH  SET %s ", props))
	logs.Debug(ctx, "", logs.Fields{"cypher": func() any { return sb.String() }})
	return NewCypherBuilderResult(sb.String(), dataMap, nil), nil
}

func (c *relationCypher[T]) InsertMany(ctx context.Context, list []T) (CypherResult, error) {
	return nil, nil
}

func (c *relationCypher[T]) Update(ctx context.Context, data T, setFields ...string) (CypherResult, error) {

	// 只更新关系标签
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r

	// 更新关系的标签与属性
	// match(n)-[r:relation{id:'cbc4d7be-43fa-427e-956d-e812b335bc12'}]->(m) create (n)-[r2:relation]->(m) set r2=r, r2.title='title' with r delete r

	prosNames, mapData, err := getUpdateProperties(ctx, data, "r2", setFields...)
	if err != nil {
		return nil, err
	}
	relType := c.eb.GetRelType(data)
	labels := c.getLabels(relType)
	if len(labels) == 0 {
		return nil, errors.New("neo4j relation.Type is nil")
	}

	tenantId := c.eb.GetTenantId(data)
	id := c.eb.GetId(data)

	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]->(m) CREATE (n)-[r2%s]->(m) SET r2=r, %s  WITH r DELETE r ", tenantId, id, labels, prosNames)
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (c *relationCypher[T]) UpdateByRSQL(ctx context.Context, tenantId string, rSQL string, data T, setFields ...string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]->(n) CREATE (n)-[r2:%v]-(m) SET r2=r WITH r DELETE r ")
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]->(n) CREATE (n)-[r2:%v]-(m) SET r2=r WITH r DELETE r ", tenantId, id, label)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) UpdateMany(ctx context.Context, list []T) (CypherResult, error) {
	panic("implement me")
}

func (c *relationCypher[T]) UpdateLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v'}]-(n) create (n)-[r2:{tenantId:'%v'}]-(m)  %s SET r2=r WITH r DELETE r ", tenantId, tenantId, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// neo4j 不支持删除关系标签
	return nil, nil
}

func (c *relationCypher[T]) DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	// neo4j 不支持删除关系标签
	return nil, nil
}

func (c *relationCypher[T]) DeleteByLabels(ctx context.Context, tenantId string, label ...string) (CypherResult, error) {
	// neo4j 不支持删除关系标签
	return nil, nil
}

func (c *relationCypher[T]) DeleteByTenantId(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) delete r `, tenantId)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteById(ctx context.Context, tenantId string, id string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v',id:'%v'}]-(b) delete r `, tenantId, id)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	strIds := getSqlInStr(ids)
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) WHERE r.id in [%v] delete r `, tenantId, strIds)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) delete r `, tenantId)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "r", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) %v delete r `, tenantId, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) FindById(ctx context.Context, tenantId, id string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v',id:'%v'}]->(b) RETURN r `, tenantId, id)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	strIds := getSqlInStr(ids)
	cypher := fmt.Sprintf("MATCH (a)-[r{tenantId:'%v'}]-(b) where r.id in [%v] RETURN r", tenantId, strIds)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"r"}), nil
}

func (c *relationCypher[T]) FindByAggregateId(ctx context.Context, tenantId, aggregateName, aggregateId string) (result CypherResult, err error) {
	return c.DeleteByFilter(ctx, tenantId, fmt.Sprintf("%v=='%v'", aggregateName, aggregateId))
}

func (c *relationCypher[T]) FindByGraphId(ctx context.Context, tenantId string, graphId string) (result CypherResult, err error) {
	return c.FindByFilter(ctx, tenantId, fmt.Sprintf("graphId=='%v'", graphId))
}

func (c *relationCypher[T]) FindByCaseId(ctx context.Context, tenantId string, caseId string) (result CypherResult, err error) {
	return c.FindByFilter(ctx, tenantId, fmt.Sprintf("caseId=='%v'", caseId))
}

func (c *relationCypher[T]) FindAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf(`MATCH (a)-[r%v{tenantId:'%v'}]->(b) RETURN r `, c.getLabels(""), tenantId)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) FindPaging(ctx context.Context, qry ddd_repository.FindPagingQuery) (CypherResult, error) {
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
	logs.Debug(ctx, qry.GetTenantId(), logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, keys, NewCypherResultOptions().SetCountCypher(countCypher)), nil
}

func (c *relationCypher[T]) FindByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf(`MATCH (a{tenantId:'%v'})-[n%v{tenantId:'%v'}]-(b{tenantId:'%v'}) %v return n `, tenantId, c.getLabels(""), tenantId, tenantId, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher[T]) Count(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (a{tenantId:'%v'})-[n%v]->(b{tenantId:'%v'}) %v RETURN count(n) as n  ", tenantId, c.labels, tenantId, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher[T]) GetFilter(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (a{tenantId:'%v'})-[n%v]->(b{tenantId:'%v'}) %v RETURN n  ", tenantId, c.getLabels(""), tenantId, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher[T]) getLabels(labels string) string {
	if c.isEmptyLabels && len(labels) == 0 {
		return ""
	} else if c.isEmptyLabels {
		return fmt.Sprintf(":`%v`", labels)
	}
	return c.labels
}

func (c *relationCypher[T]) getCreateProperties(ctx context.Context, data interface{}) (string, map[string]any, error) {
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

func (c *relationCypher[T]) getSetFields(ctx context.Context, resName string, data interface{}) (string, map[string]any, error) {
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
