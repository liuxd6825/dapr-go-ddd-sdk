package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetData(ctx context.Context, result neo4j.Result) map[string][]any {
	dataSet := make(map[string][]any)
	init := false

	for result.Next() {
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

	return dataSet
}
