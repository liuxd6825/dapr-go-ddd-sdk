package store_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
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

type Neo4jResult[T any] struct {
	dataSet map[string][]any
	eb      NodeEntityBuilder[T]
}

type KeyResult[T any] struct {
	result     *Neo4jResult[T]
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

func NewNeo4jResult[T any](ctx context.Context, eb store.EntityBuilder[T], result neo4j.ResultWithContext, keys ...*KeyResult[interface{}]) *Neo4jResult[T] {
	dataSet := make(map[string][]any)
	init := false

	for result.Next(ctx) {
		record := result.Record()
		if !init {
			init = true
			for _, key := range record.Keys {
				dataSet[key] = make([]any, 0)
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

	return &Neo4jResult[T]{
		dataSet: dataSet,
		eb:      eb,
	}
}

func (r *Neo4jResult[T]) SetMapEntity(schema *store.DBSchema, neo4jData map[string]any, entity map[string]any) {
	for _, field := range schema.Fields {
		if !field.Readable {
			continue
		}
		if val, ok := neo4jData[field.DBName]; ok {
			entity[field.Name] = val
		}
	}
}

func (r *Neo4jResult[T]) Data() map[string][]any {
	return r.dataSet
}

func (r *Neo4jResult[T]) GetData(key string) ([]any, bool) {
	v, ok := r.dataSet[key]
	return v, ok
}

func (r *Neo4jResult[T]) GetList(ctx context.Context, key string, resList any, schema *store.DBSchema, opts ...*MappingOptions) error {
	// 检查 res 是否为 *[]map[string]any 类型
	resType := reflect.TypeOf(resList)
	println(resType.Kind().String())
	if resType.Kind() != reflect.Ptr || resType.Elem().Kind() != reflect.Slice || resType.Elem().Elem().Kind() != reflect.Map {
		return errors.New("resList must be a pointer to a slice of map[string]any")
	}

	// 解引用 res
	resValue := reflect.ValueOf(resList).Elem()

	// 遍历 Neo4j 结果集
	items, found := r.dataSet[key]
	if !found {
		return nil
	}
	isMap := r.eb.GetConfig().IsMap
	if isMap {
		for _, item := range items {
			if node, ok := item.(dbtype.Node); ok {
				mapEntity := map[string]any{}
				r.SetMapEntity(schema, node.Props, mapEntity)
				m := reflect.ValueOf(mapEntity)
				resValue.Set(reflect.Append(resValue, m))

			} else if rel, ok := item.(neo4j.Relationship); ok {
				mapEntity := map[string]any{}
				r.SetMapEntity(schema, rel.Props, mapEntity)
				m := reflect.ValueOf(mapEntity)
				resValue.Set(reflect.Append(resValue, m))
			}
		}
	} else {
		for _, item := range items {
			if node, ok := item.(dbtype.Node); ok {
				m := reflect.ValueOf(node.Props)
				resValue.Set(reflect.Append(resValue, m))
			} else if rel, ok := item.(neo4j.Relationship); ok {
				m := reflect.ValueOf(rel.Props)
				resValue.Set(reflect.Append(resValue, m))
			}
		}
	}

	return nil
}

func (r *Neo4jResult[T]) GetAnyList(key string) []any {
	var resList []any
	items, found := r.dataSet[key]
	if !found {
		return nil
	}
	for _, item := range items {
		resList = append(resList, item)
	}
	return resList
}

func (r *Neo4jResult[T]) GetIntList(key string) []int64 {
	var resList []int64
	items, found := r.dataSet[key]
	if !found {
		return nil
	}
	for _, item := range items {
		if iVal, ok := item.(int64); ok {
			resList = append(resList, iVal)
		}
	}
	return resList
}

func (r *Neo4jResult[T]) GetSum(data any) error {
	if m, ok := data.(map[string]any); ok {
		for k, list := range r.dataSet {
			m[k] = list[0]
		}
		return nil
	} else if v, ok := data.(*map[string]any); ok {
		m := *v
		for k, list := range r.dataSet {
			m[k] = list[0]
		}
		return nil
	}
	return errors.New("data is not a map[string]any")
}

// GetOne
// @Description:
// @receiver r
// @param key
// @param entity
// @return bool
// @return error
func (r *Neo4jResult[T]) GetOne(dataKey string, entity interface{}, schema *store.DBSchema, opts ...*MappingOptions) (bool, error) {
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
	if entityMap, ok := entity.(map[string]any); ok {
		if node, ok := item.(dbtype.Node); ok {
			r.SetMapEntity(schema, node.GetProperties(), entityMap)
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
func (r *Neo4jResult[T]) GetInteger(key string, defaultValue int64) (int64, error) {
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

func (r *Neo4jResult[T]) GetRowsAffected() int64 {
	count, err := r.GetInteger("rows", 0)
	if err != nil {
		panic(err)
	}
	return count
}

/*
func (r *Neo4jResult[T]) AddEntity(key string, value interface{}) []T {
	var list []T
	if v, ok := r.dataSet[key]; ok {
		list = v
	} else {
		list = make([]interface{}, 0)
		r.dataSet[key] = list
	}
	list = append(list, value)
	return list
}
*/

func (r *Neo4jResult[T]) GetLength(dataKey string) int64 {
	var total int64 = 0
	dataList, ok := r.dataSet[dataKey]
	if !ok {
		return total
	}
	total = int64(len(dataList))
	return total
}

func (r *Neo4jResult[T]) GetInt(dataKey string) int64 {
	var total int64 = 0
	dataList, ok := r.dataSet[dataKey]
	if !ok {
		return total
	}
	val := dataList[0]
	if v, ok := val.(int64); ok {
		return v
	} else if v, ok := val.(int); ok {
		return int64(v)
	}
	return total
}

func (r *Neo4jResult[T]) setEntity(sourceValue reflect.Value, targetValue reflect.Value) error {
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
