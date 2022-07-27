package ddd_neo4j

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"reflect"
)

type Neo4jResult struct {
	data map[string][]interface{}
}

type KeyResult[T interface{}] struct {
	result     *Neo4jResult
	key        string
	newEntity  func() T
	list       []T
	isInitList bool
}

func NewNeo4jResult(result neo4j.Result, keys ...*KeyResult[interface{}]) *Neo4jResult {
	mapList := make(map[string][]interface{})
	init := false
	for result.Next() {
		record := result.Record()
		if !init {
			init = true
			for _, key := range record.Keys {
				mapList[key] = make([]interface{}, 0)
			}
		}
		for _, key := range record.Keys {
			list := mapList[key]
			if value, ok := record.Get(key); ok {
				if items, ok := value.([]interface{}); ok {
					for _, item := range items {
						list = append(list, item)
					}
					mapList[key] = list
				} else {
					mapList[key] = append(list, value)
				}

			}
		}
	}
	return &Neo4jResult{
		data: mapList,
	}
}

func (r *Neo4jResult) GetLists(keys []string, resultList ...interface{}) error {
	if len(keys) != len(resultList) {
		return fmt.Errorf("GetLists(keys, list...) keys.length != list.length")
	}
	for i, key := range keys {
		list := resultList[i]
		if err := r.GetList(key, list); err != nil {
			return fmt.Errorf("error: GetList(key) by key \"%s\"", err.Error())
		}
	}
	return nil
}

func (r *Neo4jResult) GetList(key string, list interface{}) error {
	dataList, ok := r.data[key]
	if !ok {
		return fmt.Errorf("GetList(key, list) key \"%s\" not exist ", key)
	}
	err := reflectutils.MappingSlice(dataList, list, func(i int, source reflect.Value, target reflect.Value) error {
		s := source
		t := target
		return r.setEntity(s, t)
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Neo4jResult) GetOne(key string, entity interface{}) error {
	var list []any
	if len(key) == 0 {
		for _, v := range r.data {
			list = v
			break
		}
	} else {
		neo4jList, ok := r.data[key]
		if !ok {
			return fmt.Errorf("GetList(key, list) key \"%s\" not exist ", key)
		}
		list = neo4jList
	}

	if len(list) != 1 {
		return fmt.Errorf("GetList(key, list) key \"%s\" entity length != 1  not exist ", key)
	}
	err := reflectutils.MappingStruct(list[0], entity, func(source reflect.Value, target reflect.Value) error {
		return r.setEntity(source, target)
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Neo4jResult) AddEntity(key string, value interface{}) []interface{} {
	var list []interface{}
	if v, ok := r.data[key]; ok {
		list = v
	} else {
		list = make([]interface{}, 0)
		r.data[key] = list
	}
	list = append(list, value)
	return list
}

func (r *Neo4jResult) setEntity(sourceValue reflect.Value, targetValue reflect.Value) error {
	source := sourceValue.Interface()
	target := targetValue.Interface()
	switch source.(type) {
	case neo4j.Node:
		node := source.(neo4j.Node)
		if err := setNode(target, node); err != nil {
			return err
		}
		break
	case neo4j.Relationship:
		rel := source.(neo4j.Relationship)
		if err := setRelationship(target, rel); err != nil {
			return err
		}
		break
	}
	return nil
}

func setNode(data interface{}, node neo4j.Node) error {
	if err := mapstructure.Decode(node.Props, data); err != nil {
		return err
	}
	id, isId := node.Props["id"]
	tenantId, isTenantId := node.Props["tenantId"]
	v := reflectutils.GetValuePointer(data)
	element := v.Interface()
	switch element.(type) {
	case ElementEntity:
		n := element.(ElementEntity)
		n.SetNid(node.Id)
		n.SetLabels(node.Labels)
		n.SetNid(node.Id)
		if isId {
			n.SetId(id.(string))
		}
		if isTenantId {
			n.SetTenantId(tenantId.(string))
		}
	}
	return nil
}

func setRelationship(data interface{}, rel neo4j.Relationship) error {
	if err := mapstructure.Decode(rel.Props, data); err != nil {
		return err
	}
	if r, ok := data.(RelationshipEntity); ok {
		r.SetNid(rel.Id)
		r.SetType(rel.Type)
		r.SetEndId(rel.EndId)
		r.SetStartId(rel.StartId)
	} else if e, ok := data.(ElementEntity); ok {
		e.SetNid(rel.Id)
	}
	return nil
}
