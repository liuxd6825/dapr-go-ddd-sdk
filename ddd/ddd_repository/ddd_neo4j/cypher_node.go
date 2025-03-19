package ddd_neo4j

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql/rsql_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"reflect"
	"strings"
)

type nodeCypher[T any] struct {
	labels string
	eb     NodeEntityBuilder[T]
	schema *dbschema.DBSchema
}

const (
	or  = " or "
	and = " and "
)

// NewNodeCypher
// @Description:
// @param labels Neo4j标签
// @return nodeCypher
func NewNodeCypher[T any](eb NodeEntityBuilder[T], schema *dbschema.DBSchema, labels ...string) Cypher[T] {
	return &nodeCypher[T]{
		eb:     eb,
		labels: getLabels(labels...),
		schema: schema,
	}
}

func (c *nodeCypher[T]) Insert(ctx context.Context, tenantId string, data T) (CypherResult, error) {
	props, dataMap, err := c.getCreateProperties(ctx, data)
	if err != nil {
		return nil, err
	}
	list := []string{}
	//list = append(list, c.labels)
	list = append(list, "case_"+c.eb.GetCaseId(data))
	list = append(list, "tenant_"+c.eb.GetTenantId(data))
	labels := c.getLabels(list...)

	cypher := fmt.Sprintf("CREATE (n%s{%s}) RETURN count(n) as rows ", labels, props)
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, dataMap, nil), nil
}

func (c *nodeCypher[T]) InsertOrUpdate(ctx context.Context, node T) (CypherResult, error) {
	props, dataMap, err := c.getUpdateProperties(ctx, node, "n")
	if err != nil {
		return nil, err
	}
	list := c.eb.GetLabels(node)
	list = append(list, "case_"+c.eb.GetCaseId(node))
	list = append(list, "tenant_"+c.eb.GetTenantId(node))
	labels := c.getLabels(list...)

	id := c.eb.GetId(node)
	cypher := fmt.Sprintf("MERGE (n%s{id:'%v'}) ON CREATE SET %v ON MATCH SET %v RETURN count(n) as rows ", labels, id, props, props)
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, dataMap, []string{"n"}), nil
}

func (c *nodeCypher[T]) InsertMany(ctx context.Context, tenantId string, list []T) (CypherResult, error) {
	/*`	CREATE (:pig{name:"猪爷爷",age:6}),
	(:pig{name:"猪奶奶",age:4}),
	(:pig{name:"猪爸爸",age:3}),
	(:pig{name:"猪妈妈",age:1})`*/
	cyphers := &strings.Builder{}
	cyphers.WriteString("UNWIND $data AS row \n ")
	cyphers.WriteString("CREATE ")
	vList := reflect.ValueOf(list)

	item := vList.Index(0).Interface()
	node := any(item).(T)
	props, _, err := c.getCreateMatchProperties(ctx, node, "row")
	if err != nil {
		return nil, err
	}
	label := c.eb.GetLabels(node)
	label = append(label, "case_"+c.eb.GetCaseId(node))
	label = append(label, "tenant_"+tenantId)
	labels := c.getLabels(label...)
	cypher := fmt.Sprintf(" (n%s{%s}) ", labels, props)
	cyphers.WriteString(cypher)
	cyphers.WriteString(" RETURN count(n) as rows ")

	createList := c.newCreateList(ctx, list)
	//logs.Debug(ctx, "", logs.Fields{"cypher": cyphers.String()})
	return NewCypherBuilderResult(cyphers.String(), map[string]any{"data": createList}, nil), nil
}

func (c *nodeCypher[T]) Update(ctx context.Context, tenantId string, data T, setFields ...string) (CypherResult, error) {
	prosNames, mapData, err := c.getUpdateProperties(ctx, data, "n", setFields...)
	if err != nil {
		return nil, err
	}
	c.eb.SetUpdatedInfo(ctx, data)
	cypher := fmt.Sprintf("MATCH (n{id:$id}) SET %s RETURN n ", prosNames)
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, mapData, []string{"n"}), nil
}

func (c *nodeCypher[T]) UpdateMany(ctx context.Context, tenantId string, list []T) (CypherResult, error) {
	/*
		UNWIND $params AS params
		MATCH (n:Person {id: param.id})
		SET n.name = param.name, n.age = param.age
		RETURN n
	*/
	if len(list) == 0 {
		return nil, errors.New("empty list")
	}

	node := list[0]
	label := c.eb.GetLabels(node)
	label = append(label, "case_"+c.eb.GetCaseId(node))
	label = append(label, "tenant_"+c.eb.GetTenantId(node))
	labels := c.getLabels(label...)

	cyphers := &strings.Builder{}
	cyphers.WriteString("UNWIND $params AS param \n")
	cyphers.WriteString(fmt.Sprintf("MATCH(n%s{id:param.id}) \n", labels))
	cyphers.WriteString("SET \n")

	data := c.newUpdateList(ctx, list)
	count := len(c.schema.Fields)
	i := 0
	for _, f := range c.schema.Fields {
		i++
		if f.Name == "id" || f.Updatable == false {
			continue
		}
		cyphers.WriteString(fmt.Sprintf("  n.%s = param.%s \n", f.DBName, f.DBName))
		if i < count {
			cyphers.WriteString(",")
		}
	}
	cyphers.WriteString("WITH n, 1 AS increment \n")
	cyphers.WriteString("RETURN sum(increment) AS rows")
	return NewCypherBuilderResult(cyphers.String(), map[string]any{"params": data}, []string{"rows"}), nil
}

func (c *nodeCypher[T]) UpdateByRSQL(ctx context.Context, tenantId string, rSQL string, data T, setFields ...string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", rSQL)
	if err != nil {
		return nil, err
	}
	prosNames, mapData, err := c.getUpdateProperties(ctx, data, "n", setFields...)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n{id:$id}) %s SET %s  RETURN count(n) as rows ", where, prosNames)
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, mapData, []string{"rows"}), nil
}

func (c *nodeCypher[T]) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]-(n) create (n)-[r2:%v]-(m) SET r2=r WITH r DELETE r ", tenantId, id, label)
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) UpdateLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	setLabels := getLabels(labels...)
	// 设置标签
	// match (n:CAR) set n:NEW xremove n:CAR
	cypher := fmt.Sprintf("MATCH (n{tenantId:'%v'}) %v SET n%v ", tenantId, where, setLabels)
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) Delete(ctx context.Context, node T) (CypherResult, error) {
	mapData := c.newMap(ctx, node)

	cypher := fmt.Sprintf("MATCH (n%v{tenantId:$tenantId,id:$id}) DETACH DELETE n", c.getLabels())
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (c *nodeCypher[T]) DeleteMany(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	count := len(ids)
	if count == 0 {
		return nil, errors.New("DeleteMany() ids.length is 0")
	}
	var whereIds string
	for i, id := range ids {
		whereIds = fmt.Sprintf("%v n.id='%v' ", whereIds, id)
		if i < count {
			whereIds += or
		}
	}
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) where %v DELETE n RETURN count(n) as rows", c.getLabels(), tenantId, whereIds)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *nodeCypher[T]) DeleteById(ctx context.Context, tenantId string, id string) (CypherResult, error) {
	params := make(map[string]any)
	params["id"] = id
	params["tenantId"] = tenantId
	cypher := fmt.Sprintf("MATCH (n%v{id:'%v'}) DETACH DELETE n RETURN count(n) as rows", c.getLabels("tenant_"+tenantId), id)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	count := len(ids)
	if count == 0 {
		return nil, errors.New("DeleteByIds() ids.length is 0")
	}
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	idWhere := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%s'}) WHERE n.id in [%s] DETACH DELETE n ", c.getLabels("tenant_"+tenantId), tenantId, idWhere)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteByCaseId(ctx context.Context, tenantId string, caseId string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v) DETACH DELETE n", c.getLabels("case_"+caseId, "tenant_"+tenantId))
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteByTenantId(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v) DETACH DELETE n", c.getLabels("tenant_"+tenantId))
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteByLabels(ctx context.Context, tenantId string, label ...string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v) DETACH DELETE n", c.getLabels(label...))
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v) DETACH DELETE n", c.getLabels("tenant_"+tenantId))
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteByRSQL(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v) WHERE (%v) DETACH DELETE n", c.getLabels("tenant_"+tenantId), where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// 设置标签
	// match (n:CAR) set n:NEW xremove n:CAR
	cypher := fmt.Sprintf("MATCH (n%v{id:'%v'}) REMOVE n:%v ", c.getLabels("tenant_"+tenantId), id, label)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	setLabels := getLabels(labels...)
	// 设置标签
	// match (n:CAR) set n:NEW xremove n:CAR
	cypher := fmt.Sprintf("MATCH (n%v) %v REMOVE n%v ", c.getLabels("tenant_"+tenantId), where, setLabels)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) FindById(ctx context.Context, tenantId, id string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v{id:'%v'}) RETURN n", c.getLabels("tenant_"+tenantId), id)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *nodeCypher[T]) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	idWhere := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (n%v) WHERE n.id in [%s] RETURN n ", c.getLabels("tenant_"+tenantId), idWhere)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) FindByCaseId(ctx context.Context, tenantId string, caseId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%v{caseId:'%s'}) RETURN n ", c.getLabels("tenant_"+tenantId), caseId)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher[T]) FindByAggregateId(ctx context.Context, tenantId string, aggregateName, aggregateId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s) WHERE n.%v='%s' RETURN n ", c.getLabels("tenant_"+tenantId), aggregateName, aggregateId)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher[T]) FindAll(ctx context.Context, tenantId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s{tenantId:'%s'}) RETURN n ", c.getLabels("tenant_"+tenantId))
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher[T]) GetRSQL(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v) %v RETURN n  ", c.getLabels("tenant_"+tenantId), where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *nodeCypher[T]) GetMatchCypher(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%s) ", c.getLabels("tenant_"+tenantId))
	res := NewCypherBuilderResult(cypher, nil, []string{"n"})
	return res, nil
}

func (c *nodeCypher[T]) FindByLabel(ctx context.Context, tenantId string, labels []string) (CypherResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *nodeCypher[T]) Count(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v) %v RETURN count(n) as count", c.getLabels("tenant_"+tenantId), where)
	logs.Debug(ctx, tenantId, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"count"}), nil
}

func (c *nodeCypher[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery) (CypherResult, error) {
	where, err := getNeo4jWhere(query.GetTenantId(), "n", query.GetFilter())
	if err != nil {
		return nil, err
	}
	skip := query.GetPageNum() * query.GetPageSize()
	pageSize := query.GetPageSize()

	keys := []string{"n"}

	order, err := getOrder(query.GetSort())
	if err != nil {
		return nil, err
	}
	tenantId := query.GetTenantId()
	cypher := fmt.Sprintf("MATCH (n%v) %v RETURN n %v SKIP %v LIMIT %v ", c.getLabels("tenant_"+tenantId), where, order, skip, pageSize)
	logs.Debug(ctx, query.GetTenantId(), logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, keys), nil
}

func (c *nodeCypher[T]) getNodeLabels(node T) string {
	label := c.eb.GetLabels(node)
	label = append(label, "case_"+c.eb.GetCaseId(node))
	label = append(label, "tenant_"+c.eb.GetTenantId(node))
	return getLabels(label...)
}

func (c *nodeCypher[T]) getLabels(labels ...string) string {
	s := c.labels
	for _, l := range labels {
		if len(l) > 0 {
			s = fmt.Sprintf("%v:`%v`", s, l)
		}
	}
	if strings.HasSuffix(s, ":") {
		s = s[:len(s)-1]
	}
	return strings.ToLower(s)
}

func (c *nodeCypher[T]) getCreateProperties(ctx context.Context, data any) (string, map[string]any, error) {
	mapData := c.newCreateMap(ctx, data)

	var properties string
	for _, f := range c.schema.Fields {
		properties = fmt.Sprintf(`%s%s:$%s,`, properties, f.DBName, f.DBName)
	}

	if len(properties) > 0 {
		properties = properties[:len(properties)-1]
	}
	return properties, mapData, nil
}

func (c *nodeCypher[T]) getCreateMatchProperties(ctx context.Context, data any, asName string) (string, map[string]any, error) {
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

func (c *nodeCypher[T]) Sum(ctx context.Context, tenantId string, rSQL string, valueCols []*ddd_repository.ValueCol) (CypherResult, error) {
	// 解析rsql
	process := rsql_neo4j.NewProcess(tenantId, "n")
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
		sumFields = append(sumFields, fmt.Sprintf("sum(%s) AS %s", col.Field, col.Field))
	}
	sum := strings.Join(sumFields, ", ")
	sb.WriteString(" RETURN ")
	sb.WriteString(sum)

	return NewCypherBuilderResult(sb.String(), nil, []string{"r"}), nil
}

func (c *nodeCypher[T]) tenant(tenantId string) string {
	return ":tenant_" + tenantId
}

func (c *nodeCypher[T]) getQueryMatch(tenantId string, properties ...string) string {
	props := strings.Join(properties, "")
	cypher := fmt.Sprintf("MATCH (a%s{%s})", c.getLabels(c.labels), props)
	return cypher
}

// getOrder
// @Description: 返回排序bson.D
// @receiver r
// @param sort  排序语句 "name:desc,id:asc"
// @return bson.D
// @return error
func getOrder(sort string) (string, error) {
	if len(sort) == 0 {
		return "", nil
	}
	// 输入
	// name:desc,id:asc
	// 输出
	// order by n.name desc , n.id asc
	res := " order by "
	list := strings.Split(sort, ",")
	for _, s := range list {
		sortItem := strings.Split(s, ":")
		orderName := sortItem[0]
		orderName = strings.Trim(orderName, " ")
		if orderName == "id" {
			orderName = "id"
		}
		order := "asc"
		if len(sortItem) > 1 {
			order = sortItem[1]
			order = strings.ToLower(order)
			order = strings.Trim(order, " ")
		}

		// 其中 1 为升序排列，而-1是用于降序排列.
		orderVal := "asc"
		var oerr error
		switch order {
		case "asc":
			orderVal = "asc"
		case "desc":
			orderVal = "desc"
		default:
			oerr = errors.New("order " + order + " is error")
		}
		if oerr != nil {
			return "", oerr
		}

		res = fmt.Sprintf("%v n.%v %v,", res, orderName, orderVal)
	}
	res = res[0 : len(res)-1]
	return res, nil
}

func (c *nodeCypher[T]) getUpdateProperties(ctx context.Context, data any, dataKey string, setFields ...string) (string, map[string]any, error) {
	mapData := c.newUpdateMap(ctx, data)
	return c.getUpdatePropertiesByMap(ctx, mapData, dataKey, setFields...)
}

func (c *nodeCypher[T]) getUpdatePropertiesByMap(ctx context.Context, mapData map[string]any, dataKey string, setFields ...string) (string, map[string]any, error) {
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

func (c *nodeCypher[T]) newUpdateMap(ctx context.Context, data any) map[string]any {
	res := c.newMap(ctx, data, func(m map[string]any) {
		c.eb.SetUpdatedInfo(ctx, m)
	})
	return res
}

func (c *nodeCypher[T]) newCreateMap(ctx context.Context, data any) map[string]any {
	mapData := c.newMap(ctx, data, func(m map[string]any) {
		c.eb.SetCreatedInfo(ctx, m)
	})
	return mapData
}

func (c *nodeCypher[T]) newMap(ctx context.Context, data any, opts ...func(map[string]any)) map[string]any {
	v, err := c.schema.NewMap(ctx, data, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *nodeCypher[T]) newCreateList(ctx context.Context, list []T) []map[string]any {
	var items []map[string]any
	for _, node := range list {
		item := c.newCreateMap(ctx, node)
		items = append(items, item)
	}
	return items
}

func (c *nodeCypher[T]) newUpdateList(ctx context.Context, list []T) []map[string]any {
	var items []map[string]any
	for _, node := range list {
		item := c.newUpdateMap(ctx, node)
		items = append(items, item)
	}
	return items
}
