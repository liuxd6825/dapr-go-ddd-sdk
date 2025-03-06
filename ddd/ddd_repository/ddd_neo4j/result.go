package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Neo4jResult struct {
	dataSet map[string][]any
}

type KeyResult[T interface{}] struct {
	result     *Neo4jResult
	key        string
	newEntity  func() T
	list       []T
	isInitList bool
}

type Mapping func(sourceValue reflect.Value, targetValue reflect.Value) error

type MappingOptions struct {
	mapping Mapping
}

const timeLayout = "2006-01-02 15:04:05Z07:00"

var jsonTimeType = reflect.TypeOf(times.NewTime())

func NewMappingOptions() *MappingOptions {
	return &MappingOptions{
		mapping: defaultMapping,
	}
}

func NewNeo4jResult(ctx context.Context, result neo4j.ResultWithContext, keys ...*KeyResult[interface{}]) *Neo4jResult {
	dataSet := make(map[string][]any)
	init := false

	for result.Next(ctx) {
		record := result.Record()
		if !init {
			init = true
			for _, key := range record.Keys {
				dataSet[key] = make([]interface{}, 0)
			}
		}
		for _, key := range record.Keys {
			list := dataSet[key]
			if value, ok := record.Get(key); ok {
				if items, ok := value.([]interface{}); ok {
					for _, item := range items {
						list = append(list, item)
					}
					dataSet[key] = list
				} else {
					dataSet[key] = append(list, value)
				}

			}
		}
	}

	return &Neo4jResult{
		dataSet: dataSet,
	}
}

func (r *Neo4jResult) Data() map[string][]any {
	return r.dataSet
}

func (r *Neo4jResult) GetData(key string) ([]any, bool) {
	v, ok := r.dataSet[key]
	return v, ok
}

func (r *Neo4jResult) GetList(ctx context.Context, key string, resList any, opts ...*MappingOptions) error {
	// 检查 res 是否为 *[]map[string]any 类型
	resType := reflect.TypeOf(resList)
	if resType.Kind() != reflect.Ptr || resType.Elem().Kind() != reflect.Slice || resType.Elem().Elem().Kind() != reflect.Map {
		return errors.New("res must be a pointer to a slice of map[string]any")
	}

	// 解引用 res
	resValue := reflect.ValueOf(resList).Elem()

	// 遍历 Neo4j 结果集

	items, found := r.dataSet[key]
	if !found {
		return fmt.Errorf("dataKey '%s' not found in record", key)
	}

	for _, item := range items {
		if node, ok := item.(dbtype.Node); ok {
			m := reflect.ValueOf(node.Props)
			resValue.Set(reflect.Append(resValue, m))
		} else if rel, ok := item.(neo4j.Relationship); ok {
			m := reflect.ValueOf(rel.Props)
			resValue.Set(reflect.Append(resValue, m))
		}
	}

	return nil
}

// GetOne
// @Description:
// @receiver r
// @param key
// @param entity
// @return bool
// @return error
func (r *Neo4jResult) GetOne(dataKey string, entity interface{}, opts ...*MappingOptions) (bool, error) {
	options := NewMappingOptions()
	options.Merge(opts...)

	var list []any
	key := dataKey
	if len(key) == 0 {
		for k, v := range r.dataSet {
			list = v
			key = k
			break
		}
	} else {
		neo4jList, ok := r.dataSet[key]
		if !ok {
			return false, fmt.Errorf("GetOne(dataKey, entity) dataKey \"%s\" not exist ", key)
		}
		list = neo4jList
	}

	count := len(list)
	if count == 0 {
		return false, nil
	} else if count > 1 {
		return false, fmt.Errorf("GetOne(dataKey, data) dataKey \"%s\" neo4j result %v > 1  not exist ", key, count)
	}
	var err error
	item := list[0]
	if m, ok := entity.(map[string]any); ok {
		if node, ok := item.(dbtype.Node); ok {
			for k, v := range node.GetProperties() {
				m[k] = v
			}
		}
	} else {
		err = reflectutils.MappingStruct(item, entity, func(source reflect.Value, target reflect.Value) error {
			return options.mapping(source, target)
		})
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

// GetInteger
// @Description: 获取整数值，如count查询结果；当neo4j结果是列表时，只取第一条；当neo4j没有结果时，返回defaultValue值
// @receiver r
// @param  key   数据集Key名称
// @param  defaultValue 默认值
// @return int64 total汇总数据量
// @return error 错误
func (r *Neo4jResult) GetInteger(key string, defaultValue int64) (int64, error) {
	var total int64 = 0
	dataList, ok := r.dataSet[key]
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

func (r *Neo4jResult) GetRowsAffected() int64 {
	count, err := r.GetInteger("rows", 0)
	if err != nil {
		panic(err)
	}
	return count
}

func (r *Neo4jResult) AddEntity(key string, value interface{}) []interface{} {
	var list []interface{}
	if v, ok := r.dataSet[key]; ok {
		list = v
	} else {
		list = make([]interface{}, 0)
		r.dataSet[key] = list
	}
	list = append(list, value)
	return list
}

func (r *Neo4jResult) GetCount(dataKey string) int64 {
	var total int64 = 0
	dataList, ok := r.dataSet[dataKey]
	if !ok {
		return total
	}
	total = int64(len(dataList))
	return total
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

func (o *MappingOptions) SetMapping(v Mapping) *MappingOptions {
	o.mapping = v
	return o
}

func (o *MappingOptions) GetMapping() Mapping {
	return o.mapping
}

func (o *MappingOptions) Merge(options ...*MappingOptions) {
	for _, i := range options {
		if i.mapping != nil {
			o.mapping = i.mapping
		}
	}
}

func defaultMapping(sourceValue reflect.Value, targetValue reflect.Value) error {
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
	if err := maputils.Decode(node.Props, &data); err != nil {
		return err
	}
	v := reflectutils.GetValuePointer(data)
	element := v.Interface()
	switch element.(type) {
	case Node:
		n := element.(Node)
		n.SetNid(node.Id)
		n.SetLabels(node.Labels)
		if id, ok := node.Props["id"]; ok {
			n.SetId(id.(string))
		}
		if tenantId, ok := node.Props["tenantId"]; ok {
			n.SetTenantId(tenantId.(string))
		}
	}
	return nil
}

func setRelationship(data interface{}, rel neo4j.Relationship) error {
	if err := decode(rel.Props, &data); err != nil {
		return err
	}
	v := reflectutils.GetValuePointer(data)
	element := v.Interface()
	if pr, ok := element.(Relation); ok {
		r := pr
		r.SetNid(rel.Id)
		r.SetRelType(rel.Type)
		r.SetEid(rel.EndId)
		r.SetSid(rel.StartId)
		if id, ok := rel.Props["id"]; ok {
			r.SetId(id.(string))
		}
		if tenantId, ok := rel.Props["tenantId"]; ok {
			r.SetTenantId(tenantId.(string))
		}
		r.SetProperties(rel.Props)
	} else if n, ok := data.(Node); ok {
		n.SetNid(rel.Id)
	} else {
		v := reflect.ValueOf(data)
		if r, ok := v.Elem().Interface().(Relation); ok {
			println(r)
		}
	}
	return nil
}

func decode(input interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook:       decodeHook,
		WeaklyTypedInput: true,
		Squash:           true,
		Result:           out,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(input)
	return err
}

func decodeHook(fromType reflect.Type, toType reflect.Type, v interface{}) (interface{}, error) {
	if toType == jsonTimeType {
		switch fromType.Name() {
		case "string":
			sTime := v.(string)
			format := timeLayout
			if strings.Contains(sTime, "T") {
				format = time.RFC3339
			}
			res, err := time.Parse(format, sTime)
			return times.GetTime(&res), err
		}
	} else if fromType.Kind() == reflect.String {
		switch toType.Name() {
		case "Time":
			sTime := v.(string)
			format := timeLayout
			if strings.Contains(sTime, "T") {
				format = time.RFC3339
			}
			res, err := time.Parse(format, sTime)
			return res, err
		}
	}
	return v, nil
}

// getList 从 Neo4j 返回的数据集中提取 dataKey 对应的内容，并解析到 res 中
func getList(ctx context.Context, dataKey string, res any, dataSet map[string][]any) error {
	// 检查 res 是否为 *[]map[string]any 类型
	resType := reflect.TypeOf(res)
	if resType.Kind() != reflect.Ptr || resType.Elem().Kind() != reflect.Slice || resType.Elem().Elem().Kind() != reflect.Map {
		return errors.New("res must be a pointer to a slice of map[string]any")
	}

	// 解引用 res
	resValue := reflect.ValueOf(res).Elem()

	// 遍历 Neo4j 结果集

	items, found := dataSet[dataKey]
	if !found {
		return fmt.Errorf("dataKey '%s' not found in record", dataKey)
	}

	for _, item := range items {
		if node, ok := item.(dbtype.Node); ok {
			m := reflect.ValueOf(node.Props)
			resValue.Set(reflect.Append(resValue, m))
		} else if rel, ok := item.(neo4j.Relationship); ok {
			m := reflect.ValueOf(rel.Props)
			resValue.Set(reflect.Append(resValue, m))
		}
	}

	return nil
}
