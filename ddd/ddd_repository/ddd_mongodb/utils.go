package ddd_mongodb

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/pkg/errors"
)

func getMongoFieldName(s string) string {
	id := stringutils.ToLower(s)
	if id == "id" || id == "_id" {
		return ConstIdField
	}
	return stringutils.SnakeString(s)
}

func asDocument(doc map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range doc {
		key := getMongoFieldName(k)
		res[key] = v
	}
	return res
}

func getDocumentId(doc interface{}) (id string, err error) {
	switch doc.(type) {
	case ddd.Entity:
		e, _ := doc.(ddd.Entity)
		id = e.GetId()
		break
	case map[string]interface{}:
		m, _ := doc.(map[string]interface{})
		if v, ok := m[ConstIdField]; !ok {
			err = errors.New("ddd_mongodb.getDocumentId(doc) err: doc is not have \"_id\" key. ")
		} else {
			if s, ok := v.(string); !ok {
				err = errors.New("ddd_mongodb.getDocumentId(doc) err: doc id field is not string .")
			} else {
				id = s
			}
		}
		break
	default:
		err = errors.New("ddd_mongodb.getDocumentId(doc) err: doc is ddd.Entity or map[string]interface.")
	}
	return
}
