package schema

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

var components = map[string]*Schema{
	"string":     &Schema{Type: TypeString},
	"number":     &Schema{Type: TypeNumber},
	"boolean":    &Schema{Type: TypeBoolean},
	"findPaging": findPaging(),
}

func Components() map[string]*Schema {
	return components
}

func findPaging() *Schema {
	return &Schema{
		Type: TypeObject,
		Properties: map[string]*Property{
			"filter": &Property{
				Description: ddd_repository.FilterDescription,
				Type:        TypeString,
			},
			"sort": &Property{
				Description: ddd_repository.SortDescription,
				Type:        TypeString,
			},
			"fields": &Property{
				Description: ddd_repository.FieldsDescription,
				Type:        TypeString,
			},
			"pageNum": &Property{
				Type: TypeNumber,
			},
			"pageSize": &Property{
				Type: TypeNumber,
			},
			"groupCols": &Property{
				Description: ddd_repository.GroupColsDescription,
				Type:        TypeString,
			},
			"groupKeys": &Property{
				Description: ddd_repository.GroupKeysDescription,
				Type:        TypeString,
			},
			"valueCols": &Property{
				Description: ddd_repository.ValueColsDescription,
				Type:        TypeString,
			},
		},
	}
}
