package ddd_neo4j

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"reflect"
	"strconv"
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
	var dataList []interface{}
	var ok bool = false
	if len(key) == 0 {
		for _, v := range r.data {
			dataList = v
			ok = true
			break
		}
	} else {
		dataList, ok = r.data[key]
	}

	if !ok {
		return nil
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

//
// GetInteger
// @Description: 获取整数值，如count查询结果；当neo4j结果是列表时，只取第一条；当neo4j没有结果时，返回defaultValue值
// @receiver r
// @param  key   数据集Key名称
// @param  defaultValue 默认值
// @return int64 total汇总数据量
// @return error 错误
//
func (r *Neo4jResult) GetInteger(key string, defaultValue int64) (int64, error) {
	var total int64 = 0
	dataList, ok := r.data[key]
	if !ok {
		return defaultValue, nil
	}
	if len(dataList) > 0 {
		v := dataList[0]
		s := fmt.Sprintf("%v", v)
		if count, err := strconv.ParseInt(s, 10, 64); err != nil {
			return 0, err
		} else {
			total = count
		}
	}
	return total, nil
}

//
// GetOne
// @Description:
// @receiver r
// @param key
// @param entity
// @return bool
// @return error
//
func (r *Neo4jResult) GetOne(key string, entity interface{}) (bool, error) {
	var list []any
	if len(key) == 0 {
		for _, v := range r.data {
			list = v
			break
		}
	} else {
		neo4jList, ok := r.data[key]
		if !ok {
			return false, fmt.Errorf("GetOne(key, entity) key \"%s\" not exist ", key)
		}
		list = neo4jList
	}

	count := len(list)
	if count == 0 {
		return false, nil
	} else if count > 1 {
		return false, fmt.Errorf("GetList(key, list) key \"%s\" entity length != 1  not exist ", key)
	}
	node := list[0].(neo4j.Node)
	err := reflectutils.MappingStruct(node, entity, func(source reflect.Value, target reflect.Value) error {
		return r.setEntity(source, target)
	})
	if e, ok := entity.(ElementEntity); ok {
		e.SetNid(node.Id)
		e.SetLabels(node.Labels)
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func decode(input interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Squash:           true,
		Result:           out,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	decoder.Decode(input)
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
	if err := maputils.Decode(node.Props, data); err != nil {
		return err
	}
	id, hasId := node.Props["id"]
	tenantId, hasTenantId := node.Props["tenantId"]
	v := reflectutils.GetValuePointer(data)
	element := v.Interface()
	switch element.(type) {
	case ElementEntity:
		n := element.(ElementEntity)
		n.SetNid(node.Id)
		n.SetLabels(node.Labels)
		n.SetNid(node.Id)
		if hasId {
			n.SetId(id.(string))
		}
		if hasTenantId {
			n.SetTenantId(tenantId.(string))
		}
	}
	return nil
}

func setRelationship(data interface{}, rel neo4j.Relationship) error {
	if err := decode(rel.Props, data); err != nil {
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
