package ddd_neo4j

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type CypherBuilderResult interface {
	Cypher() string
	Params() map[string]any
	ResultKeys() []string
	ResultOneKey() string
}

func NewCypherBuilderResult(cypher string, params map[string]any, resultKey []string) CypherBuilderResult {
	return &cypherBuilderResult{
		cypher:    cypher,
		params:    params,
		resultKey: resultKey,
	}
}

type cypherBuilderResult struct {
	cypher    string
	params    map[string]any
	resultKey []string
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

type CypherBuilder interface {
	Insert(ctx context.Context, data ElementEntity) (CypherBuilderResult, error)
	InsertMany(ctx context.Context, data []ElementEntity) (CypherBuilderResult, error)

	Update(ctx context.Context, data ElementEntity, setFields ...string) (CypherBuilderResult, error)
	UpdateMany(ctx context.Context, list []ElementEntity) (CypherBuilderResult, error)

	DeleteById(ctx context.Context, tenantId string, id string) (CypherBuilderResult, error)
	DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error)
	DeleteAll(ctx context.Context, tenantId string) (CypherBuilderResult, error)

	FindById(ctx context.Context, tenantId, id string) (CypherBuilderResult, error)
	FindByIds(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error)
	FindByGraphId(ctx context.Context, tenantId, graphId string) (result CypherBuilderResult, err error)
	FindAll(ctx context.Context, tenantId string) (CypherBuilderResult, error)

	GetLabels() string
}

type ReflectBuilder struct {
	labels string
}

func (r *ReflectBuilder) Insert(ctx context.Context, data ElementEntity) (CypherBuilderResult, error) {
	props, dataMap, err := r.getCreateProperties(ctx, data)
	if err != nil {
		return nil, err
	}
	labels, _ := r.getLabels(data)
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

	cypher := fmt.Sprintf("MATCH (n{id:$id}) DETACH DELETE n")
	return NewCypherBuilderResult(cypher, mapData, nil), nil
}

func (r *ReflectBuilder) DeleteMany(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error) {
	count := len(ids)
	if count == 0 {
		return nil, errors.New("DeleteManyById() ids.length is 0")
	}
	var whereIds string
	for i, id := range ids {
		whereIds = fmt.Sprintf("%v n.id='%v' ", whereIds, id)
		if i < count {
			whereIds += " or "
		}
	}
	cypher := fmt.Sprintf("MATCH (n {tenantId:'%v'}) where %v DELETE n", tenantId, whereIds)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) DeleteById(ctx context.Context, tenantId string, id string) (CypherBuilderResult, error) {
	params := make(map[string]any)
	params["id"] = id
	params["tenantId"] = tenantId
	cypher := fmt.Sprintf("MATCH (n{id:$id, tenantId:$id}) DETACH DELETE n")
	return NewCypherBuilderResult(cypher, params, nil), nil
}

func (r *ReflectBuilder) DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error) {
	count := len(ids)
	if count == 0 {
		return nil, errors.New("DeleteByIds() ids.length is 0")
	}
	var whereIds string
	for i, id := range ids {
		whereIds = fmt.Sprintf("%v n.id='%v' ", whereIds, id)
		if i < count {
			whereIds += " or "
		}
	}
	cypher := fmt.Sprintf("MATCH (n %v {tenantId:'%v'}) where %v DELETE n", r.GetLabels(), tenantId, whereIds)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) DeleteAll(ctx context.Context, tenantId string) (CypherBuilderResult, error) {
	cypher := fmt.Sprintf("MATCH (n {tenantId:'%v'}) DELETE n", tenantId)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) FindById(ctx context.Context, tenantId, id string) (CypherBuilderResult, error) {
	cypher := fmt.Sprintf("MATCH (n {tenantId:'%v',id:'%v'}) RETURN n", tenantId, id)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) FindByIds(ctx context.Context, tenantId string, ids []string) (CypherBuilderResult, error) {
	for i, id := range ids {
		ids[i] = fmt.Sprintf(`'%v'`, id)
	}
	idWhere := strings.Join(ids, ",")
	cypher := fmt.Sprintf("MATCH (n%s) WHERE  n.tenantId = '%s' and n.id in [%s] RETURN n ", r.labels, tenantId, idWhere)
	return NewCypherBuilderResult(cypher, nil, nil), nil
}

func (r *ReflectBuilder) FindByGraphId(ctx context.Context, tenantId string, graphId string) (CypherBuilderResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s) WHERE  n.tenantId = '%s' and n.graphId= '%s'  RETURN n ", r.labels, tenantId, graphId)
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
}

func (r *ReflectBuilder) FindAll(ctx context.Context, tenantId string) (CypherBuilderResult, error) {
	var params map[string]any
	cypher := fmt.Sprintf("MATCH (n%s) WHERE  n.tenantId = '%s' RETURN n ", r.labels, tenantId)
	return NewCypherBuilderResult(cypher, params, []string{"n"}), nil
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
