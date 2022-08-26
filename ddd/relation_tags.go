package ddd

import (
	"reflect"
)

const (
	dddRelTagName = "ddd-rel"
)

type Relation map[string]string

func GetRelation(data interface{}) ([]Relation, bool, error) {
	reflectValue := reflect.ValueOf(data)
	if reflectValue.Type().Kind() == reflect.Slice {
		return GetRelationByList(data)
	}

	var relations []Relation
	relation, ok, err := GetRelationByStructure(data)
	if err != nil {
		return nil, ok, err
	} else if ok {
		relations = append(relations, relation)
	}
	return relations, ok, err
}

func GetRelationByList(list interface{}) ([]Relation, bool, error) {
	reflectValue := reflect.ValueOf(list)
	var relations []Relation
	count := reflectValue.Len()
	for i := 0; i < count; i++ {
		v := reflectValue.Index(i)
		relation, ok, err := GetRelationByStructure(v.Interface())
		if err != nil {
			return nil, false, err
		}
		if ok {
			relations = append(relations, relation)
		}
	}
	return relations, len(relations) > 0, nil
}

// GetRelationByStructure 从结构中获得聚合关系
func GetRelationByStructure(structure interface{}) (Relation, bool, error) {
	reflectValue := reflect.ValueOf(structure)
	reflectType := reflectValue.Type()
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflectValue.Elem()
	}

	var relation Relation
	for i := 0; i < reflectType.NumField(); i++ {
		fieldName := reflectType.Field(i).Name
		relationName, ok := reflectType.Field(i).Tag.Lookup(dddRelTagName)
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
