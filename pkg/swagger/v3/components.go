package swagger3

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

var components = map[string]*Schema{
	"string":     &Schema{Type: []string{TypeString}},
	"number":     &Schema{Type: []string{TypeNumber}},
	"boolean":    &Schema{Type: []string{TypeBoolean}},
	"findPaging": findPaging(),
}

func GetComponents() map[string]*Schema {
	return components
}

func findPaging() *Schema {
	return &Schema{
		Type: []string{TypeObject},
		Properties: map[string]*Property{
			"filter": &Property{
				Description: store.FilterDescription,
				Type:        []string{TypeString},
			},
			"sort": &Property{
				Description: store.SortDescription,
				Type:        []string{TypeString},
			},
			"fields": &Property{
				Description: store.FieldsDescription,
				Type:        []string{TypeString},
			},
			"pageNum": &Property{
				Type: []string{TypeNumber},
			},
			"pageSize": &Property{
				Type: []string{TypeNumber},
			},
			"groupCols": &Property{
				Description: store.GroupColsDescription,
				Type:        []string{TypeString},
			},
			"groupKeys": &Property{
				Description: store.GroupKeysDescription,
				Type:        []string{TypeString},
			},
			"valueCols": &Property{
				Description: store.ValueColsDescription,
				Type:        []string{TypeString},
			},
		},
	}
}
