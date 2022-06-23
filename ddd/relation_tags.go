package ddd

import (
	"reflect"
)

const (
	relationTagName = "ddd-rel"
)

type Relation map[string]string

func GetRelation(data interface{}) (Relation, bool, error) {
	reflectValue := reflect.ValueOf(data)
	reflectType := reflectValue.Type()
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflectValue.Elem()
	}

	var relation Relation
	for i := 0; i < reflectType.NumField(); i++ {
		fieldName := reflectType.Field(i).Name
		relationName, ok := reflectType.Field(i).Tag.Lookup(relationTagName)
		relationId := "null"
		if ok {
			if relation == nil {
				relation = Relation{}
			}
			if len(relationName) == 0 || relationName == "-" {
				relationName = fieldName
			}
			relationId = reflectValue.Field(i).String()
			relation[relationName] = relationId
		}
	}
	return relation, relation != nil, nil
}
