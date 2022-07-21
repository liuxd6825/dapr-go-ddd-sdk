package ddd_neo4j

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type CypherBuilder interface {
	CreateOne(ctx context.Context, data ElementEntity) (string, map[string]any, error)
	UpdateById(ctx context.Context, data ElementEntity, setFields ...string) (string, map[string]any, error)
	DeleteById(ctx context.Context, data ElementEntity) (string, map[string]any, error)
	CreateMany(ctx context.Context, data ElementEntity) (string, error)
	UpdateByIds(ctx context.Context, data ElementEntity) (string, error)
	DeleteByIds(ctx context.Context, data ElementEntity) (string, error)
	FindById(ctx context.Context, tenantId, id string) (string, error)
	FindGraphById(ctx context.Context, data ElementEntity) (string, error)
}

type ReflectBuilder struct {
}

func (r *ReflectBuilder) CreateOne(ctx context.Context, data ElementEntity) (string, map[string]any, error) {
	prosNames, mapData, err := r.getCreateProperties(ctx, data)
	if err != nil {
		return "", nil, err
	}
	labels, _ := r.getLabels(data)
	cypher := fmt.Sprintf("CREATE (n%s{%s}) RETURN n ", labels, prosNames)
	return cypher, mapData, nil
}

func (r *ReflectBuilder) UpdateById(ctx context.Context, data ElementEntity, setFields ...string) (string, map[string]any, error) {
	prosNames, mapData, err := r.getUpdateProperties(ctx, data, setFields...)
	if err != nil {
		return "", nil, err
	}
	cypher := fmt.Sprintf("MATCH (n{id:$id}) SET %s RETURN n ", prosNames)
	return cypher, mapData, nil
}

func (r *ReflectBuilder) DeleteById(ctx context.Context, data ElementEntity) (string, map[string]any, error) {
	mapData, err := r.getMap(data)
	if err != nil {
		return "", nil, err
	}

	cypher := fmt.Sprintf("MATCH (n{id:$id}) DETACH DELETE n")
	return cypher, mapData, err
}

func (r *ReflectBuilder) CreateMany(ctx context.Context, data ElementEntity) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReflectBuilder) UpdateByIds(ctx context.Context, data ElementEntity) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReflectBuilder) DeleteByIds(ctx context.Context, data ElementEntity) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReflectBuilder) FindById(ctx context.Context, tenantId, id string) (string, error) {
	return fmt.Sprintf("MATCH (n{tenantId:'%v',id:'%v'}) RETURN n", tenantId, id), nil
}

func (r *ReflectBuilder) FindGraphById(ctx context.Context, data ElementEntity) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReflectBuilder) getCreateProperties(ctx context.Context, data any) (string, map[string]any, error) {
	mapData, err := r.getMap(data)
	if err != nil {
		return "", nil, err
	}
	var properties string
	for k, _ := range mapData {
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

	for k, _ := range mapData {
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

func NewReflectBuilder() CypherBuilder {
	return &ReflectBuilder{}
}
