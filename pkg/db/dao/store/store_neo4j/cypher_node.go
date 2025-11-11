package store_neo4j

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql/rsql_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
)

type nodeCypher[T any] struct {
	config *Config[T]
	eb     store2.EntityBuilder[T]
	schema *store2.DBSchema
}

const (
	or  = " or "
	and = " and "
)

// NewNodeCypher
// @Description:
// @param labels Neo4j标签
// @return nodeCypher
func NewNodeCypher[T any](config *Config[T]) Cypher[T] {
	return &nodeCypher[T]{
		config: config,
		eb:     config.EntityBuilder,
		schema: config.DBSchema,
	}
}

func (c *nodeCypher[T]) NewCypher() Cypher[T] {
	return NewNodeCypher(c.config)
}

func (c *nodeCypher[T]) Insert(ctx context.Context, tenantId string, data T) (CypherResult, error) {
	props, dataMap, err := c.GetCreateProperties(ctx, data)
	if err != nil {
		return nil, err
	}
	labels := c.GetLabels(ctx, data)
	cypher := fmt.Sprintf("CREATE (n%s{%s}) RETURN count(n) as rows ", labels, props)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, dataMap, nil), nil
}

func (c *nodeCypher[T]) GetLabels(ctx context.Context, data T, labels ...string) string {
	tenantId := c.eb.GetTenantId(data)
	if tenantId == "" {
		tenantId, _ = appctx.GetTenantId(ctx)
	}
	list := labels
	list = append(list, c.config.Labels...)
	if reflectutils.IsNotEmpty[T](data) {
		list = c.eb.GetLabels(data, list...)
	}
	tags := c.getLabels(list...)
	return tags
}

func (c *nodeCypher[T]) InsertOrUpdate(ctx context.Context, node T) (CypherResult, error) {
	props, dataMap, err := c.GetUpdateProperties(ctx, node, "n")
	if err != nil {
		return nil, err
	}

	labels := c.GetLabels(ctx, node)

	id := c.eb.GetId(node)
	cypher := fmt.Sprintf("MERGE (n%s{id:'%v'}) ON CREATE SET %v ON MATCH SET %v RETURN count(n) as rows ", labels, id, props, props)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, dataMap, []string{"n"}), nil
}

func (c *nodeCypher[T]) Merge(ctx context.Context, node T, fields map[string]string) (CypherResult, error) {
	props, dataMap, err := c.GetUpdateProperties(ctx, node, "n")
	if err != nil {
		return nil, err
	}
	labels := c.GetLabels(ctx, node)

	filter := ""
	for k, v := range fields {
		filter += fmt.Sprintf("%s:'%v',", k, v)
	}
	if len(filter) > 0 {
		filter = filter[:len(filter)-1]
	}

	//id := c.eb.GetId(node)
	cypher := fmt.Sprintf("MERGE (n%s{%s}) ON CREATE SET %v ON MATCH SET %v RETURN count(n) as rows ", labels, filter, props, props)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
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
	props, _, err := c.GetCreateMatchProperties(ctx, node, "row")
	if err != nil {
		return nil, err
	}
	labels := c.GetLabels(ctx, node)

	cypher := fmt.Sprintf(" (n%s{%s}) ", labels, props)
	cyphers.WriteString(cypher)
	cyphers.WriteString(" RETURN count(n) as rows ")

	createList := c.newCreateList(ctx, list)
	//logs.Debug(ctx, "", logs.Fields{"cypher": cyphers.String()})
	return NewCypherBuilderResult(cyphers.String(), map[string]any{"data": createList}, nil), nil
}

func (c *nodeCypher[T]) Update(ctx context.Context, tenantId string, data T, setFields ...string) (CypherResult, error) {
	prosNames, mapData, err := c.GetUpdateProperties(ctx, data, "n", setFields...)
	if err != nil {
		return nil, err
	}
	c.eb.SetUpdatedInfo(ctx, data)
	cypher := fmt.Sprintf("MATCH (n{id:$id}) SET %s RETURN n ", prosNames)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
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
	labels := c.GetLabels(ctx, node)

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
		dbName := c.getFieldName(f)
		cyphers.WriteString(fmt.Sprintf("  n.%s = param.%s \n", dbName, dbName))
		if i < count {
			cyphers.WriteString(",")
		}
	}
	cyphers.WriteString("WITH n, 1 AS increment \n")
	cyphers.WriteString("RETURN sum(increment) AS rows")
	return NewCypherBuilderResult(cyphers.String(), map[string]any{"params": data}, []string{"rows"}), nil
}

func (c *nodeCypher[T]) getSchema() *store2.DBSchema {
	return c.config.DBSchema
}

func (c *nodeCypher[T]) UpdateByRSQL(ctx context.Context, tenantId string, rSQL string, data T, setFields ...string) (CypherResult, error) {
	where, err := GetNeo4jWhere(tenantId, "n", rSQL)
	if err != nil {
		return nil, err
	}
	prosNames, mapData, err := c.GetUpdateProperties(ctx, data, "n", setFields...)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n{id:$id}) %s SET %s  RETURN count(n) as rows ", where, prosNames)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, mapData, []string{"rows"}), nil
}

func (c *nodeCypher[T]) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]-(n) create (n)-[r2:%v]-(m) SET r2=r WITH r DELETE r ", tenantId, id, label)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) UpdateLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	where, err := GetNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	setLabels := GetLabels(labels...)
	// 设置标签
	// match (n:CAR) set n:NEW xremove n:CAR
	cypher := fmt.Sprintf("MATCH (n{tenantId:'%v'}) %v SET n%v ", tenantId, where, setLabels)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) Delete(ctx context.Context, node T) (CypherResult, error) {
	mapData := c.newMap(ctx, node)

	cypher := fmt.Sprintf("MATCH (n%v{tenantId:$tenantId,id:$id}) DETACH DELETE n", c.getLabels())
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
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
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *nodeCypher[T]) DeleteById(ctx context.Context, tenantId string, id string) (CypherResult, error) {
	params := make(map[string]any)
	params["id"] = id
	params["tenantId"] = tenantId
	cypher := fmt.Sprintf("MATCH (n%v{id:'%v'}) DETACH DELETE n RETURN count(n) as rows", c.getLabels("tenant_"+tenantId), id)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
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
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteByCaseId(ctx context.Context, tenantId string, caseId string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v) DETACH DELETE n", c.getLabels("case_"+caseId, "tenant_"+tenantId))
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteByTenantId(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v) DETACH DELETE n", c.getLabels("tenant_"+tenantId))
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteByLabels(ctx context.Context, tenantId string, label ...string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v) DETACH DELETE n", c.getLabels(label...))
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v) DETACH DELETE n", c.getLabels("tenant_"+tenantId))
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteByRSQL(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := GetNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v) WHERE (%v) DETACH DELETE n", c.getLabels("tenant_"+tenantId), where)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// 设置标签
	// match (n:CAR) set n:NEW xremove n:CAR
	cypher := fmt.Sprintf("MATCH (n%v{id:'%v'}) REMOVE n:%v ", c.getLabels("tenant_"+tenantId), id, label)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	where, err := GetNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	setLabels := GetLabels(labels...)
	// 设置标签
	// match (n:CAR) set n:NEW xremove n:CAR
	cypher := fmt.Sprintf("MATCH (n%v) %v REMOVE n%v ", c.getLabels("tenant_"+tenantId), where, setLabels)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) FindById(ctx context.Context, tenantId, id string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v{id:'%v'}) RETURN n", c.getLabels("tenant_"+tenantId), id)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *nodeCypher[T]) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	idWhere := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (n%v) WHERE n.id in [%s] RETURN n ", c.getLabels("tenant_"+tenantId), idWhere)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher[T]) FindByCaseId(ctx context.Context, tenantId string, caseId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%v{caseId:'%s'}) RETURN n ", c.getLabels("tenant_"+tenantId), caseId)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher[T]) FindByAggregateId(ctx context.Context, tenantId string, aggregateName, aggregateId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s) WHERE n.%v='%s' RETURN n ", c.getLabels("tenant_"+tenantId), aggregateName, aggregateId)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher[T]) FindAll(ctx context.Context, tenantId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s{tenantId:'%s'}) RETURN n ", c.getLabels("tenant_"+tenantId))
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher[T]) GetRSQL(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := GetNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v) %v RETURN n  ", c.getLabels("tenant_"+tenantId), where)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
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
	where, err := GetNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v) %v RETURN count(n) as count", c.getLabels("tenant_"+tenantId), where)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, []string{"count"}), nil
}

func (c *nodeCypher[T]) FindPaging(ctx context.Context, query store2.FindPagingQuery) (CypherResult, error) {
	tenantId := appctx.GetTenantId2(ctx)
	where, err := GetNeo4jWhere(tenantId, "n", query.GetFilter())
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
	cypher := fmt.Sprintf("MATCH (n%v) %v RETURN n %v SKIP %v LIMIT %v ", c.getLabels("tenant_"+tenantId), where, order, skip, pageSize)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	return NewCypherBuilderResult(cypher, nil, keys), nil
}

func (c *nodeCypher[T]) getNodeLabels(node T) string {
	label := c.eb.GetLabels(node)
	label = append(label, c.config.Labels...)
	label = append(label, "case_"+c.eb.GetCaseId(node))
	label = append(label, "tenant_"+c.eb.GetTenantId(node))
	return GetLabels(label...)
}

func (c *nodeCypher[T]) getLabels(labels ...string) string {
	s := ""
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

func (c *nodeCypher[T]) GetCreateProperties(ctx context.Context, data any) (string, map[string]any, error) {
	mapData := c.NewCreateMap(ctx, data)

	var properties string
	for _, f := range c.schema.Fields {
		if f.Creatable {
			dbName := c.getFieldName(f)
			propName := c.getPropertyName(f)
			properties = fmt.Sprintf(`%s%s:$%s,`, properties, dbName, propName)
		}
	}

	if len(properties) > 0 {
		properties = properties[:len(properties)-1]
	}
	return properties, mapData, nil
}

func (c *nodeCypher[T]) GetCreateMatchProperties(ctx context.Context, data any, asName string) (string, map[string]any, error) {
	mapData := c.NewCreateMap(ctx, data)

	var properties string
	for _, f := range c.schema.Fields {
		dbName := c.getFieldName(f)
		if f.DataType == store2.DataType_Time {
			properties = fmt.Sprintf(`%s%s:dateTime(%s.%s),`, properties, dbName, asName, dbName)
		} else if f.DataType == store2.DataType_Date {
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

func (c *nodeCypher[T]) Sum(ctx context.Context, tenantId string, rSQL string, valueCols []*store2.ValueCol) (CypherResult, error) {
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
	where, err := GetNeo4jWhere(tenantId, "r", rSQL)
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
	cypher := fmt.Sprintf("MATCH (a%s{%s})", c.getLabels(), props)
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

func (c *nodeCypher[T]) GetUpdateProperties(ctx context.Context, data any, dataKey string, setFields ...string) (string, map[string]any, error) {
	mapData := c.NewUpdateMap(ctx, data)
	return c.GetUpdatePropertiesByMap(ctx, mapData, dataKey, setFields...)
}

func (c *nodeCypher[T]) getFieldName(field *store2.Field) string {
	return field.DBName
}

func (c *nodeCypher[T]) getPropertyName(field *store2.Field) string {
	return field.Name
}

func (c *nodeCypher[T]) GetUpdatePropertiesByMap(ctx context.Context, mapData map[string]any, dataKey string, updateProperties ...string) (string, map[string]any, error) {
	var properties string

	if updateProperties != nil {
		for _, propName := range updateProperties {
			field := c.schema.LookedField(propName)
			if field != nil && field.Updatable {
				dbName := c.getFieldName(field)
				propName := c.getPropertyName(field)
				properties = fmt.Sprintf(`%s%s.%s=$%s,`, properties, dataKey, dbName, propName)
			}
		}
	} else {
		for _, field := range c.schema.Fields {
			if field.Updatable {
				dbName := c.getFieldName(field)
				propName := c.getPropertyName(field)
				properties = fmt.Sprintf(`%s%s.%s=$%s,`, properties, dataKey, dbName, propName)
			}
		}
	}

	if len(properties) > 0 {
		properties = properties[:len(properties)-1]
	}

	return properties, mapData, nil
}

func (c *nodeCypher[T]) NewUpdateMap(ctx context.Context, data any) map[string]any {
	res := c.newMap(ctx, data, func(m map[string]any) {
		c.eb.SetUpdatedInfo(ctx, m)
	})
	return res
}

func (c *nodeCypher[T]) NewCreateMap(ctx context.Context, data any) map[string]any {

	mapData := c.newMap(ctx, data, func(m map[string]any) {
		c.eb.SetCreatedInfo(ctx, m)
	})
	return mapData
}

func (c *nodeCypher[T]) newMap(ctx context.Context, data any, opts ...func(map[string]any)) map[string]any {
	dataMap, err := c.schema.NewMap(ctx, data, opts...)
	if err != nil {
		panic(err)
	}
	for key, val := range dataMap {
		if timeVal, ok := val.(*time.Time); ok {
			if timeVal != nil {
				dataMap[key] = timeVal.Unix()
			} else {
				dataMap[key] = nil
			}
		}
	}
	return dataMap
}

func (c *nodeCypher[T]) newCreateList(ctx context.Context, list []T) []map[string]any {
	var items []map[string]any
	for _, node := range list {
		item := c.NewCreateMap(ctx, node)
		items = append(items, item)
	}
	return items
}

func (c *nodeCypher[T]) newUpdateList(ctx context.Context, list []T) []map[string]any {
	var items []map[string]any
	for _, node := range list {
		item := c.NewUpdateMap(ctx, node)
		items = append(items, item)
	}
	return items
}
