package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql/rsql_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"strings"
)

type relationCypher[T any] struct {
	relTypes      []string
	matchTypes    string
	isEmptyLabels bool
	eb            RelationEntityBuilder[T]
	schema        *dbschema.DBSchema
}

// NewRelationCypher
// @Description:
// @param labels 关系标签，可以为空值；为空：由Relation.GetRelType()决定标签名称
// @return Cypher
func NewRelationCypher[T any](eb RelationEntityBuilder[T], schema *dbschema.DBSchema, relTypes ...string) Cypher[T] {
	matchTypes := ":" + strings.Join(relTypes, "|")

	rel := &relationCypher[T]{
		matchTypes:    matchTypes,
		relTypes:      relTypes,
		isEmptyLabels: len(relTypes) == 0,
		eb:            eb,
		schema:        schema,
	}
	return rel
}

func (c *relationCypher[T]) Insert(ctx context.Context, tenantId string, data T) (CypherResult, error) {
	list := []T{data}
	return c.InsertMany(ctx, tenantId, list)
}

func (c *relationCypher[T]) InsertOrUpdate(ctx context.Context, data T) (CypherResult, error) {
	props, dataMap, err := c.getSetFields(ctx, "r", data)
	if err != nil {
		return nil, err
	}
	tenant := c.tenant(c.eb.GetTenantId(data))
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("MATCH (s%s{id:'%v'}), (e%s{id:'%v'}) ", tenant, c.eb.GetStartId(data), tenant, c.eb.GetEndId(data)))
	sb.WriteString(fmt.Sprintf("MERGE (s)-[r%s{id:'%v'}]->(e) ", c.matchTypes, c.eb.GetId(data)))
	sb.WriteString(fmt.Sprintf("ON CREATE SET %s ", props))
	sb.WriteString(fmt.Sprintf("ON MATCH  SET %s ", props))
	logs.Debug(ctx, "", logs.Fields{"cypher": func() any { return sb.String() }})
	return NewCypherBuilderResult(sb.String(), dataMap, nil), nil
}

func (c *relationCypher[T]) InsertMany(ctx context.Context, tenantId string, list []T) (CypherResult, error) {
	/*
		UNWIND $updates AS update
		MATCH (a:Person {id: update.startId}), (b:Person {id: update.endId})
		CALL apoc.create.relationship(a, update.refType, update.properties, b)
		YIELD rel RETURN COUNT(rel) AS rows;
	*/
	if len(list) == 0 {
		return nil, errors.New("empty list")
	}
	tenant := ":tenant_" + tenantId

	cyphers := fmt.Sprintf(`UNWIND $updates AS update
	MATCH (a%s{id: update.startId}), (b%s {id: update.endId}) 
	CALL apoc.create.relationship(a, update.refType, update.properties, b) 
	YIELD rel RETURN COUNT(rel) AS rows;`, tenant, tenant)

	creates := c.newCreateList(ctx, list)

	return NewCypherBuilderResult(cyphers, map[string]any{"updates": creates}, []string{"rows"}), nil
}

func (c *relationCypher[T]) Update(ctx context.Context, tenantId string, data T, setFields ...string) (CypherResult, error) {
	return c.UpdateMany(ctx, tenantId, []T{data})
}

func (c *relationCypher[T]) UpdateByRSQL(ctx context.Context, tenantId string, rSQL string, data T, setFields ...string) (CypherResult, error) {
	process := rsql_neo4j.NewProcess(tenantId, "r")
	if err := rsql.ParseProcess(rSQL, process); err != nil {
		return nil, err
	}
	where := process.GetSQL()
	if where != "" {
		where = "WHERE " + where
	}
	tenant := ":tenant_" + tenantId
	cypher := fmt.Sprintf("MATCH (n%s)-[r]->(m%s) %s DELETE r RETURN count(r) as rows", tenant, tenant, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"rows"}), nil
}

func (c *relationCypher[T]) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r
	match := c.getQueryMatch(tenantId)
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]->(n) CREATE (n)-[r2:%v]-(m) SET r2=r WITH r DELETE r ", match, id, label)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) UpdateMany(ctx context.Context, tenantId string, list []T) (CypherResult, error) {
	tenant := ":tenant_" + tenantId
	cypher := fmt.Sprintf(`UNWIND $updates AS update
		MATCH (a%s {id: update.startId})-[r]->(b%s {id: update.endId})
		WHERE type(r) = update.refType
		CALL apoc.do.when(
			r IS NOT NULL,
			'SET r += update.properties RETURN r',
			'RETURN null',
		{r: r, update: update}
		) YIELD value
		RETURN count(value) AS rows`, tenant, tenant)

	updates := c.newUpdateList(ctx, list)

	return NewCypherBuilderResult(cypher, map[string]any{"updates": updates}, nil), nil
}

func (c *relationCypher[T]) MergeMany(ctx context.Context, tenantId string, list []T) (CypherResult, error) {
	/*
		UNWIND $updates AS update
		MATCH (a:tenant_test {id: update.startId}), (b:tenant_test {id: update.endId})
		CALL apoc.merge.relationship(a, update.refType, {}, update.properties, b) YIELD rel
		RETURN count(rel) AS rows
	*/
	tenant := ":tenant_" + tenantId
	cypher := fmt.Sprintf(`
		UNWIND $updates AS update
		MATCH (a%s {id: update.startId}), (b%s {id: update.endId})
		CALL apoc.merge.relationship(a, update.refType, {}, update.properties, b) YIELD rel
		RETURN count(rel) AS rows`, tenant, tenant)

	updates := c.newUpdateList(ctx, list)

	return NewCypherBuilderResult(cypher, map[string]any{"updates": updates}, nil), nil
}

func (c *relationCypher[T]) newCreateList(ctx context.Context, list []T) []map[string]any {
	var items []map[string]any
	for _, item := range list {
		properties := c.newCreateMap(ctx, item)
		createItem := map[string]any{
			"startId":    c.eb.GetStartId(item),
			"endId":      c.eb.GetEndId(item),
			"refType":    c.eb.GetRelType(item),
			"properties": properties,
		}
		items = append(items, createItem)
	}
	return items
}

func (c *relationCypher[T]) newUpdateList(ctx context.Context, list []T) []map[string]any {
	var items []map[string]any
	for _, item := range list {
		properties := c.newUpdateMap(ctx, item)
		update := map[string]any{
			"startId":    c.eb.GetStartId(item),
			"endId":      c.eb.GetEndId(item),
			"refType":    c.eb.GetRelType(item),
			"properties": properties,
		}
		items = append(items, update)
	}
	return items
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
	cypher := fmt.Sprintf(`MATCH (a)-[r{tenantId:'%v'}]-(b) DELETE r RETURN count(r) as rows`, tenantId)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteById(ctx context.Context, tenantId string, id string) (CypherResult, error) {
	tenant := ":tenant_" + tenantId
	cypher := fmt.Sprintf(`MATCH (a%s)-[r{id:'%v'}]-(b%s) DELETE r RETURN count(r) as rows `, tenant, id, tenant)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	tenant := ":tenant_" + tenantId
	delIds := getSqlInStr(ids)
	cypher := fmt.Sprintf(`MATCH (a%s)-[r]-(b%s) WHERE r.id in [%v] DELETE r RETURN count(r) as rows `, tenant, tenant, delIds)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteAll(ctx context.Context, tenantId string) (CypherResult, error) {
	match := c.getDeleteMatch(tenantId)
	cypher := fmt.Sprintf(`%s DELETE r RETURN count(r) as rows  `, match)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) DeleteByRSQL(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "r", filter)
	if err != nil {
		return nil, err
	}
	match := c.getDeleteMatch(tenantId)
	cypher := fmt.Sprintf(`%s %v DELETE r RETURN count(r) as rows `, match, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) FindById(ctx context.Context, tenantId, id string) (CypherResult, error) {
	match := c.getQueryMatch(tenantId, fmt.Sprintf("id:'%v'", id))
	cypher := fmt.Sprintf(`%s RETURN r `, match)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) tenant(tenantId string) string {
	return ":tenant_" + tenantId
}

func (c *relationCypher[T]) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	strIds := getSqlInStr(ids)
	match := c.getQueryMatch(tenantId)
	cypher := fmt.Sprintf("%s WHERE r.id in [%v] RETURN r", match, strIds)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"r"}), nil
}

func (c *relationCypher[T]) getQueryMatch(tenantId string, properties ...string) string {
	tenant := c.tenant(tenantId)
	props := strings.Join(properties, "")
	cypher := fmt.Sprintf("MATCH (a%s)-[r%s{%s}]->(b%s)", tenant, c.matchTypes, props, tenant)
	return cypher
}

func (c *relationCypher[T]) getDeleteMatch(tenantId string, properties ...string) string {
	tenant := c.tenant(tenantId)
	props := strings.Join(properties, "")
	cypher := fmt.Sprintf("MATCH (a%s)-[r%s{%s}]->(b%s)", tenant, c.matchTypes, props, tenant)
	return cypher
}

func (c *relationCypher[T]) Sum(ctx context.Context, tenantId, rSQL string, valueCols []*ddd_repository.ValueCol) (CypherResult, error) {
	// 解析rsql
	process := rsql_neo4j.NewProcess(tenantId, "r")
	if err := rsql.ParseProcess(rSQL, process); err != nil {
		return nil, err
	}

	// 取match
	match := c.getQueryMatch(tenantId)
	sb := strings.Builder{}
	sb.WriteString(match)

	// 取得where条件
	where, err := getNeo4jWhere(tenantId, "r", rSQL)
	if err != nil {
		return nil, err
	}
	sb.WriteString(where)

	// 取得return sum 条件
	sumFields := make([]string, 0)
	for _, col := range valueCols {
		sumFields = append(sumFields, fmt.Sprintf("sum(r.%s) AS %s", col.Field, col.Field))
	}
	sum := strings.Join(sumFields, ", ")
	sb.WriteString(" RETURN ")
	sb.WriteString(sum)

	return NewCypherBuilderResult(sb.String(), nil, []string{"r"}), nil
}

func (c *relationCypher[T]) FindByAggregateId(ctx context.Context, tenantId, aggregateName, aggregateId string) (result CypherResult, err error) {
	return c.DeleteByRSQL(ctx, tenantId, fmt.Sprintf("%v=='%v'", aggregateName, aggregateId))
}

func (c *relationCypher[T]) FindByGraphId(ctx context.Context, tenantId string, graphId string) (result CypherResult, err error) {
	return c.FindByRSQL(ctx, tenantId, fmt.Sprintf("graphId=='%v'", graphId))
}

func (c *relationCypher[T]) FindByCaseId(ctx context.Context, tenantId string, caseId string) (result CypherResult, err error) {
	return c.FindByRSQL(ctx, tenantId, fmt.Sprintf("caseId=='%v'", caseId))
}

func (c *relationCypher[T]) FindAll(ctx context.Context, tenantId string) (CypherResult, error) {
	match := c.getQueryMatch(tenantId)
	cypher := fmt.Sprintf(`%s RETURN r, count(r) AS rows `, match)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *relationCypher[T]) FindPaging(ctx context.Context, qry ddd_repository.FindPagingQuery) (CypherResult, error) {
	where, err := getNeo4jWhere(qry.GetTenantId(), "r", qry.GetFilter())
	if err != nil {
		return nil, err
	}

	skip := qry.GetPageNum() * qry.GetPageSize()
	pageSize := qry.GetPageSize()
	order, err := getOrder(qry.GetSort())
	if err != nil {
		return nil, err
	}
	match := c.getQueryMatch(qry.GetTenantId())
	cypher := fmt.Sprintf("%s %s RETURN r %s SKIP %v LIMIT %v ", match, where, order, skip, pageSize)
	countCypher := ""
	if qry.GetIsTotalRows() {
		countCypher = fmt.Sprintf("%s %s RETURN count(r) AS rows  ", match, where)
	}
	logs.Debug(ctx, qry.GetTenantId(), logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"r", "rows"}, NewCypherResultOptions().SetCountCypher(countCypher)), nil
}

func (c *relationCypher[T]) FindByRSQL(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "r", filter)
	if err != nil {
		return nil, err
	}
	match := c.getQueryMatch(tenantId)
	cypher := fmt.Sprintf(`%s %s return r, count(r) AS rows `, match, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *relationCypher[T]) Count(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "r", filter)
	if err != nil {
		return nil, err
	}

	match := c.getQueryMatch(tenantId)
	cypher := fmt.Sprintf("%s %s RETURN count(r) AS rows ", match, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"rows"}), nil
}

func (c *relationCypher[T]) GetRSQL(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "r", filter)
	if err != nil {
		return nil, err
	}
	match := c.getQueryMatch(tenantId)
	cypher := fmt.Sprintf("%s %s RETURN r, count(r) AS rows  ", match, where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"r"}), nil
}

func (c *relationCypher[T]) getSetFields(ctx context.Context, resName string, data interface{}) (string, map[string]any, error) {
	mapData := c.newMap(ctx, data)

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

func (c *relationCypher[T]) getUpdatePropertiesByMap(ctx context.Context, mapData map[string]any, dataKey string, setFields ...string) (string, map[string]any, error) {
	var properties string
	isSetFields := len(setFields) > 0
	var keyFields map[string]string
	if isSetFields {
		keyFields = make(map[string]string)
		for _, k := range setFields {
			keyFields[strings.ToLower(k)] = k
		}
	}

	for k := range mapData {
		if isSetFields {
			if _, ok := keyFields[strings.ToLower(k)]; ok {
				properties = fmt.Sprintf(`%s%s.%s=$%s,`, properties, dataKey, k, k)
			}
		} else {
			properties = fmt.Sprintf(`%s%s.%s=$%s,`, properties, dataKey, k, k)
		}
	}

	if len(properties) > 0 {
		properties = properties[:len(properties)-1]
	}

	return properties, mapData, nil
}

func (c *relationCypher[T]) newUpdateMap(ctx context.Context, data any) map[string]any {
	mapData := c.newMap(ctx, data)
	return mapData
}

func (c *relationCypher[T]) newCreateMap(ctx context.Context, data any) map[string]any {
	mapData := c.newMap(ctx, data, func(m map[string]any) {
		c.eb.SetCreatedInfo(ctx, m)
	})
	return mapData
}

func (c *relationCypher[T]) newMap(ctx context.Context, data any, opts ...func(map[string]any)) map[string]any {
	v, err := c.schema.NewMap(ctx, data, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *relationCypher[T]) getUpdateProperties(ctx context.Context, data any, dataKey string, setFields ...string) (string, map[string]any, error) {
	mapData := c.newMap(ctx, data, func(m map[string]any) {
		c.eb.SetUpdatedInfo(ctx, m)
	})
	return c.getUpdatePropertiesByMap(ctx, mapData, dataKey, setFields...)
}

func (c *relationCypher[T]) getCreateProperties(ctx context.Context, data interface{}) (string, map[string]any, error) {
	mapData := c.newMap(ctx, data, func(m map[string]any) {
		c.eb.SetCreatedInfo(ctx, data)
	})

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

func (c *relationCypher[T]) getCreateMatchProperties(ctx context.Context, data any, asName string) (string, map[string]any, error) {
	mapData := c.newCreateMap(ctx, data)

	var properties string
	for _, f := range c.schema.Fields {
		dbName := f.DBName
		if f.DataType == dbschema.Time {
			properties = fmt.Sprintf(`%s%s:dateTime(%s.%s),`, properties, dbName, asName, dbName)
		} else if f.DataType == dbschema.Date {
			properties = fmt.Sprintf(`%s%s:date(%s.%s),`, properties, dbName, asName, dbName)
		} else {
			properties = fmt.Sprintf(`%s%s:%s.%s,`, properties, dbName, asName, dbName)
		}
	}

	if len(properties) > 0 {
		properties = properties[:len(properties)-1]
	}
	return properties, mapData, nil
}
