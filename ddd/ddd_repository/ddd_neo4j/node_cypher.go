package ddd_neo4j

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"reflect"
	"strings"
)

type nodeCypher struct {
	labels string
}

const (
	or  = " or "
	and = " and "
)

//
// NewNodeCypher
// @Description:
// @param labels Neo4j标签
// @return nodeCypher
//
func NewNodeCypher(labels ...string) Cypher {
	return &nodeCypher{
		labels: getLabels(labels...),
	}
}

func (c *nodeCypher) Insert(ctx context.Context, data interface{}) (CypherResult, error) {
	node := data.(Node)
	props, dataMap, err := c.getCreateProperties(ctx, node)
	if err != nil {
		return nil, err
	}
	list := node.GetLabels()
	list = append(list, "graph_"+node.GetGraphId())
	list = append(list, "tenant_"+node.GetTenantId())
	labels := c.getLabels(list...)

	cypher := fmt.Sprintf("CREATE (n%s{%s}) RETURN n ", labels, props)
	return NewCypherBuilderResult(cypher, dataMap, nil), nil
}

func (c *nodeCypher) InsertOrUpdate(ctx context.Context, data interface{}) (CypherResult, error) {
	node := data.(Node)
	props, dataMap, err := getUpdateProperties(ctx, node, "n")
	if err != nil {
		return nil, err
	}
	list := node.GetLabels()
	list = append(list, "graph_"+node.GetGraphId())
	list = append(list, "tenant_"+node.GetTenantId())
	labels := c.getLabels(list...)

	cypher := fmt.Sprintf("MERGE (n%s{id:'%v'}) ON CREATE SET %v ON MATCH SET %v RETURN n ", labels, node.GetId(), props, props)
	return NewCypherBuilderResult(cypher, dataMap, []string{"n"}), nil
}

func (c *nodeCypher) InsertMany(ctx context.Context, list interface{}) (CypherResult, error) {
	/*`	CREATE (:pig{name:"猪爷爷",age:6}),
	(:pig{name:"猪奶奶",age:4}),
	(:pig{name:"猪爸爸",age:3}),
	(:pig{name:"猪妈妈",age:1})`*/
	cyphers := &strings.Builder{}
	cyphers.WriteString("CREATE ")
	vList := reflect.ValueOf(list)
	count := vList.Len()
	for i := 0; i < count; i++ {
		node := vList.Index(i).Interface().(Node)
		props, _, err := c.getCreateProperties(ctx, node)
		if err != nil {
			return nil, err
		}
		label := node.GetLabels()
		label = append(label, "graph_"+node.GetGraphId())
		label = append(label, "tenant_"+node.GetTenantId())
		labels := c.getLabels(label...)
		cypher := fmt.Sprintf(" (n%s{%s}) ", labels, props)
		cyphers.WriteString(cypher)
		if i != count {
			cyphers.WriteString(",")
		}
	}
	return NewCypherBuilderResult(cyphers.String(), nil, nil), nil
}

func (c *nodeCypher) Update(ctx context.Context, data interface{}, setFields ...string) (CypherResult, error) {
	prosNames, mapData, err := getUpdateProperties(ctx, data, "n", setFields...)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n{id:$id}) SET %s RETURN n ", prosNames)
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (c *nodeCypher) UpdateMany(ctx context.Context, list interface{}) (CypherResult, error) {
	nodes := list.([]Node)
	println(nodes)

	//TODO implement me
	panic("implement me")
}

func (c *nodeCypher) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// match(n)-[r:测试]->(m) create(n)-[r2:包括]->(m) set r2=r with r delete r
	cypher := fmt.Sprintf("MATCH (n)-[r{tenantId:'%v',id:'%v'}]-(n) create (n)-[r2:%v]-(m) SET r2=r WITH r DELETE r ", tenantId, id, label)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) UpdateLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	setLabels := getLabels(labels...)
	// 设置标签
	// match (n:CAR) set n:NEW remove n:CAR
	cypher := fmt.Sprintf("MATCH (n{tenantId:'%v'}) %v SET n%v ", tenantId, where, setLabels)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) Delete(ctx context.Context, data interface{}) (CypherResult, error) {
	mapData, err := getMap(data)
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (n%v{tenantId:$tenantId,id:$id}) DETACH DELETE n", c.getLabels())
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (c *nodeCypher) DeleteMany(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
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
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) where %v DELETE n", c.getLabels(), tenantId, whereIds)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) DeleteById(ctx context.Context, tenantId string, id string) (CypherResult, error) {
	params := make(map[string]any)
	params["id"] = id
	params["tenantId"] = tenantId
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v',id:'%v'}) DETACH DELETE n", c.getLabels(), tenantId, id)
	return NewCypherBuilderResult(cypher, params, nil), nil
}

func (c *nodeCypher) DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	count := len(ids)
	if count == 0 {
		return nil, errors.New("DeleteByIds() ids.length is 0")
	}
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	idWhere := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%s'}) WHERE n.id in [%s] DETACH DELETE n ", c.getLabels(), tenantId, idWhere)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) DeleteAll(ctx context.Context, tenantId string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) DETACH DELETE n", c.getLabels(), tenantId)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) DeleteByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) WHERE (%v) DETACH DELETE n", c.labels, tenantId, where)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) DeleteLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error) {
	// 设置标签
	// match (n:CAR) set n:NEW remove n:CAR
	cypher := fmt.Sprintf("MATCH (n{tenantId:'%v',id:'%v'}) REMOVE n:%v ", tenantId, id, label)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	setLabels := getLabels(labels...)
	// 设置标签
	// match (n:CAR) set n:NEW remove n:CAR
	cypher := fmt.Sprintf("MATCH (n{tenantId:'%v'}) %v REMOVE n%v ", tenantId, where, setLabels)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) FindById(ctx context.Context, tenantId, id string) (CypherResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v',id:'%v'}) RETURN n", c.getLabels(), tenantId, id)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error) {
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	idWhere := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%s'}) WHERE n.id in [%s] RETURN n ", c.getLabels(), tenantId, idWhere)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (c *nodeCypher) FindByGraphId(ctx context.Context, tenantId string, graphId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%s', graphId:'%s'}) RETURN n ", c.labels, tenantId, graphId)
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher) FindByAggregateId(ctx context.Context, tenantId string, aggregateName, aggregateId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s{tenantId:'%s'}) WHERE n.%v='%s' RETURN n ", c.getLabels(), tenantId, aggregateName, aggregateId)
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher) FindAll(ctx context.Context, tenantId string) (CypherResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s{tenantId:'%s'}) RETURN n ", c.getLabels(), tenantId)
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (c *nodeCypher) GetFilter(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) %v RETURN n  ", c.getLabels(), tenantId, where)
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (c *nodeCypher) FindByLabel(ctx context.Context, tenantId string, labels []string) (CypherResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *nodeCypher) Count(ctx context.Context, tenantId, filter string) (CypherResult, error) {
	where, err := getNeo4jWhere(tenantId, "n", filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) %v RETURN count(n) as count", c.getLabels(), tenantId, where)
	return NewCypherBuilderResult(cypher, nil, []string{"count"}), nil
}

func (c *nodeCypher) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery) (CypherResult, error) {
	where, err := getNeo4jWhere(query.GetTenantId(), "n", query.GetFilter())
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

	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) %v RETURN n %v %v SKIP %v LIMIT %v ", c.labels, query.GetTenantId(), where, count, order, skip, pageSize)
	return NewCypherBuilderResult(cypher, nil, keys), nil
}

func (c *nodeCypher) getLabels(labels ...string) string {
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

func (c *nodeCypher) getCreateProperties(ctx context.Context, data any) (string, map[string]any, error) {
	mapData, err := getMap(data)
	if err != nil {
		return "", nil, err
	}
	var properties string
	for k := range mapData {
		properties = fmt.Sprintf(`%s%s:$%s,`, properties, k, k)
	}
	if len(properties) > 0 {
		properties = properties[:len(properties)-1]
	}
	return properties, mapData, nil
}

//
//  getOrder
//  @Description: 返回排序bson.D
//  @receiver r
//  @param sort  排序语句 "name:desc,id:asc"
//  @return bson.D
//  @return error
//
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

func getUpdateProperties(ctx context.Context, data any, dataKey string, setFields ...string) (string, map[string]any, error) {
	mapData, err := getMap(data)
	if err != nil {
		return "", nil, err
	}
	return getUpdatePropertiesByMap(ctx, mapData, dataKey, setFields...)
}

func getUpdatePropertiesByMap(ctx context.Context, mapData map[string]any, dataKey string, setFields ...string) (string, map[string]any, error) {
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

func getMap(data any) (map[string]interface{}, error) {
	mapData := make(map[string]any)
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}
