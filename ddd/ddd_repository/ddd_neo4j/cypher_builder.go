package ddd_neo4j

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"strings"
)

type CypherBuilderResult interface {
	Cypher() string
	Params() map[string]any
	ResultKeys() []string
	ResultOneKey() string
}

type cypherBuilderResult struct {
	cypher    string
	params    map[string]any
	resultKey []string
}

type CypherBuilder interface {
	Insert(ctx context.Context, data ElementEntity) (CypherBuilderResult, error)
	InsertMany(ctx context.Context, data []ElementEntity) (CypherBuilderResult, error)

	Update(ctx context.Context, data ElementEntity, setFields ...string) (CypherBuilderResult, error)
	UpdateMany(ctx context.Context, list []ElementEntity) (CypherBuilderResult, error)

	DeleteById(ctx context.Context, tenantId string, id string) (CypherBuilderResult, error)
	DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error)
	DeleteAll(ctx context.Context, tenantId string) (CypherBuilderResult, error)
	DeleteByFilter(ctx context.Context, tenantId string, filter string) (CypherBuilderResult, error)

	GetFilter(ctx context.Context, tenantId, filter string) (CypherBuilderResult, error)
	FindById(ctx context.Context, tenantId, id string) (CypherBuilderResult, error)
	FindByIds(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error)
	FindByAggregateId(ctx context.Context, tenantId, aggregateName, aggregateId string) (result CypherBuilderResult, err error)
	FindByGraphId(ctx context.Context, tenantId string, graphId string) (result CypherBuilderResult, err error)
	FindAll(ctx context.Context, tenantId string) (CypherBuilderResult, error)
	FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery) (CypherBuilderResult, error)
	Count(ctx context.Context, tenantId, filter string) (CypherBuilderResult, error)

	GetLabels() string
}

type ReflectBuilder struct {
	labels string
}

const (
	or  = " or "
	and = " and "
)

//
// NewReflectBuilder
// @Description:
// @param labels Neo4j标签
// @return CypherBuilder
//
func NewReflectBuilder(labels ...string) CypherBuilder {
	var s string
	for _, l := range labels {
		s = s + ":" + l
	}
	return &ReflectBuilder{labels: s}
}

func NewCypherBuilderResult(cypher string, params map[string]any, resultKey []string) CypherBuilderResult {
	return &cypherBuilderResult{
		cypher:    cypher,
		params:    params,
		resultKey: resultKey,
	}
}

func (c *cypherBuilderResult) Cypher() string {
	return c.cypher
}

func (c *cypherBuilderResult) Params() map[string]any {
	return c.params
}

func (c *cypherBuilderResult) ResultKeys() []string {
	return c.resultKey
}

func (c *cypherBuilderResult) ResultOneKey() string {
	if len(c.resultKey) > 0 {
		return c.resultKey[0]
	}
	return ""
}

func (r *ReflectBuilder) Insert(ctx context.Context, data ElementEntity) (CypherBuilderResult, error) {
	props, dataMap, err := r.getCreateProperties(ctx, data)
	if err != nil {
		return nil, err
	}
	labels := r.GetLabels()
	cypher := fmt.Sprintf("CREATE (n%s{%s}) RETURN n ", labels, props)
	return NewCypherBuilderResult(cypher, dataMap, nil), nil
}

func (r *ReflectBuilder) InsertMany(ctx context.Context, list []ElementEntity) (CypherBuilderResult, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReflectBuilder) Update(ctx context.Context, data ElementEntity, setFields ...string) (CypherBuilderResult, error) {
	prosNames, mapData, err := r.getUpdateProperties(ctx, data, setFields...)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n{id:$id}) SET %s RETURN n ", prosNames)
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (r *ReflectBuilder) UpdateMany(ctx context.Context, list []ElementEntity) (CypherBuilderResult, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReflectBuilder) Delete(ctx context.Context, data ElementEntity) (CypherBuilderResult, error) {
	mapData, err := r.getMap(data)
	if err != nil {
		return nil, err
	}

	cypher := fmt.Sprintf("MATCH (n%v{tenantId:$tenantId,id:$id}) DETACH DELETE n", r.GetLabels())
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (r *ReflectBuilder) DeleteMany(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error) {
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
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) where %v DELETE n", r.GetLabels(), tenantId, whereIds)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) DeleteById(ctx context.Context, tenantId string, id string) (CypherBuilderResult, error) {
	params := make(map[string]any)
	params["id"] = id
	params["tenantId"] = tenantId
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v',id:'%v'}) DETACH DELETE n", r.GetLabels(), tenantId, id)
	return NewCypherBuilderResult(cypher, params, nil), nil
}

func (r *ReflectBuilder) DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error) {
	count := len(ids)
	if count == 0 {
		return nil, errors.New("DeleteByIds() ids.length is 0")
	}
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	idWhere := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%s'}) WHERE n.id in [%s] DETACH DELETE n ", r.GetLabels(), tenantId, idWhere)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) DeleteAll(ctx context.Context, tenantId string) (CypherBuilderResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) DETACH DELETE n", r.GetLabels(), tenantId)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) DeleteByFilter(ctx context.Context, tenantId string, filter string) (CypherBuilderResult, error) {
	where, err := getSqlWhere(tenantId, filter)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) WHERE (%v) DETACH DELETE n", r.labels, tenantId, where)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) FindById(ctx context.Context, tenantId, id string) (CypherBuilderResult, error) {
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v',id:'%v'}) RETURN n", r.GetLabels(), tenantId, id)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error) {
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	idWhere := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%s'}) WHERE n.id in [%s] RETURN n ", r.GetLabels(), tenantId, idWhere)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) FindByGraphId(ctx context.Context, tenantId string, graphId string) (CypherBuilderResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%s', graphId:'%s'}) RETURN n ", r.labels, tenantId, graphId)
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (r *ReflectBuilder) FindByAggregateId(ctx context.Context, tenantId string, aggregateName, aggregateId string) (CypherBuilderResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s{tenantId:'%s'}) WHERE n.%v='%s' RETURN n ", r.GetLabels(), tenantId, aggregateName, aggregateId)
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (r *ReflectBuilder) FindAll(ctx context.Context, tenantId string) (CypherBuilderResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s{tenantId:'%s'}) RETURN n ", r.GetLabels(), tenantId)
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (r *ReflectBuilder) GetFilter(ctx context.Context, tenantId, filter string) (CypherBuilderResult, error) {
	where, err := getSqlWhere(tenantId, filter)
	if err != nil {
		return nil, err
	}
	if len(where) > 0 {
		where = "WHERE " + where
	}
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) %v RETURN n  ", r.labels, tenantId, where)
	return NewCypherBuilderResult(cypher, nil, []string{"n"}), nil
}

func (r *ReflectBuilder) Count(ctx context.Context, tenantId, fitler string) (CypherBuilderResult, error) {
	where, err := getSqlWhere(tenantId, fitler)
	if err != nil {
		return nil, err
	}
	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) %v RETURN count(n) as count", r.GetLabels(), tenantId, where)
	return NewCypherBuilderResult(cypher, nil, []string{"count"}), nil
}

func (r *ReflectBuilder) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery) (CypherBuilderResult, error) {
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

	order, err := r.getOrder(query.GetSort())
	if err != nil {
		return nil, err
	}

	if len(where) > 0 {
		where = fmt.Sprintf("WHERE %v", where)
	}

	cypher := fmt.Sprintf("MATCH (n%v{tenantId:'%v'}) %v RETURN n %v %v SKIP %v LIMIT %v ", r.labels, query.GetTenantId(), where, count, order, skip, pageSize)
	return NewCypherBuilderResult(cypher, nil, keys), nil
}

func (r *ReflectBuilder) getCreateProperties(ctx context.Context, data any) (string, map[string]any, error) {
	mapData, err := r.getMap(data)
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

func (r *ReflectBuilder) getUpdateProperties(ctx context.Context, data any, setFields ...string) (string, map[string]any, error) {
	mapData, err := r.getMap(data)
	if err != nil {
		return "", nil, err
	}
	return r.getUpdatePropertiesByMap(ctx, mapData, setFields...)
}

func (r *ReflectBuilder) getUpdatePropertiesByMap(ctx context.Context, mapData map[string]any, setFields ...string) (string, map[string]any, error) {
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
				properties = fmt.Sprintf(`%sn.%s=$%s,`, properties, k, k)
			}
		} else {
			properties = fmt.Sprintf(`%sn.%s=$%s,`, properties, k, k)
		}
	}

	if len(properties) > 0 {
		properties = properties[:len(properties)-1]
	}

	return properties, mapData, nil
}

func (r *ReflectBuilder) getLabels(data any) (string, error) {
	var labels string
	if element, ok := data.(ElementEntity); ok {
		for _, l := range element.GetLabels() {
			labels = fmt.Sprintf(":%s ", l)
		}
	}
	return labels, nil
}

func (r *ReflectBuilder) getMap(data any) (map[string]interface{}, error) {
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

func (r *ReflectBuilder) GetLabels() string {
	return r.labels
}

//
//  getOrder
//  @Description: 返回排序bson.D
//  @receiver r
//  @param sort  排序语句 "name:desc,id:asc"
//  @return bson.D
//  @return error
//
func (r *ReflectBuilder) getOrder(sort string) (string, error) {
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
