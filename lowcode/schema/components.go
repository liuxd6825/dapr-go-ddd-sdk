package schema

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

var components = map[string]*Schema{
	"string":     &Schema{Type: []string{TypeString}},
	"number":     &Schema{Type: []string{TypeNumber}},
	"boolean":    &Schema{Type: []string{TypeBoolean}},
	"findPaging": findPaging(),
}

func Components() map[string]*Schema {
	return components
}

func findPaging() *Schema {
	return &Schema{
		Type: []string{TypeObject},
		Properties: map[string]*Property{
			"filter": &Property{
				Description: ddd_repository.FilterDescription,
				Type:        []string{TypeString},
			},
			"sort": &Property{
				Description: ddd_repository.SortDescription,
				Type:        []string{TypeString},
			},
			"fields": &Property{
				Description: ddd_repository.FieldsDescription,
				Type:        []string{TypeString},
			},
			"pageNum": &Property{
				Type: []string{TypeNumber},
			},
			"pageSize": &Property{
				Type: []string{TypeNumber},
			},
			"groupCols": &Property{
				Description: ddd_repository.GroupColsDescription,
				Type:        []string{TypeString},
			},
			"groupKeys": &Property{
				Description: ddd_repository.GroupKeysDescription,
				Type:        []string{TypeString},
			},
			"valueCols": &Property{
				Description: ddd_repository.ValueColsDescription,
				Type:        []string{TypeString},
			},
		},
	}
}
