package code_generator

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"text/template"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
)

// FieldData 传递给模板的字段信息
type FieldData struct {
	Name       string // PascalCase (Go Struct Field)
	JsonName   string // camelCase
	DbName     string // snake_case
	Type       string // Go Type
	Title      string // Comment/Title
	Tag        string // Full struct tag
	Comment    string // Field comment
	IsBaseTime bool   // 是否是时间类型，用于添加 imports
}

// ModelData 传递给模板的整体数据
type ModelData struct {
	PackageName string
	Imports     []string
	StructName  string
	TableName   string
	Fields      []FieldData
	HasBase     bool // 是否继承 BaseModel
	Values      any
}

type Generator struct {
	TemplateStr string
}

func NewGenerator() *Generator {
	return &Generator{}
}

// Generate 处理 JsonSchema 并生成代码
func (g *Generator) Generate(tmpl *template.Template, schema *jsonschema.Schema, pkgName string, values any) ([]byte, error) {
	data := g.parseToModelData(schema, pkgName, values)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template error: %v", err)
	}

	// 格式化代码 (go fmt)
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// 如果格式化失败，返回原始内容以便调试
		return buf.Bytes(), fmt.Errorf("format source error: %v", err)
	}

	return formatted, nil
}

func (g *Generator) parseToModelData(sch *jsonschema.Schema, pkgName string, values any) ModelData {
	modelData := ModelData{
		PackageName: pkgName,
		StructName:  ToPascalCase(sch.Name()),
		TableName:   ToSnakeCase(sch.Name()),
		Imports:     []string{},
		HasBase:     false,
		Values:      values,
	}

	// 检查是否继承 Base (根据 allOf 或 meta 判断)
	if len(sch.AllOf) > 0 {
		modelData.HasBase = true
		modelData.Imports = append(modelData.Imports, "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase")
	}
	tableName := sch.Name()

	schMeta := schema.GetMetaExtension(sch)
	if schMeta != nil && schMeta.DBTable != nil {
		tableName = schMeta.DBTable.Name
	}
	// 如果 Meta 中定义了表名
	modelData.TableName = tableName

	// 收集字段
	// Map 遍历是无序的，为了保证生成代码一致性，我们需要对 Key 排序
	var keys []string
	for k := range sch.Properties {
		keys = append(keys, k)
	}

	// 根据 Order 排序，如果 Schema 中有 order 字段最好，这里先按 Name 排序，
	// 实际应用中可以读取 property 的 "order" 字段进行 sort
	type sortField struct {
		Key   string
		Order int
	}
	var sortedFields []sortField
	for _, k := range keys {
		prop := sch.Properties[k]
		if prop.Order != nil {
			order := *prop.Order
			sortedFields = append(sortedFields, sortField{Key: k, Order: order})
		}

	}

	sort.Slice(sortedFields, func(i, j int) bool {
		if sortedFields[i].Order != sortedFields[j].Order {
			return sortedFields[i].Order < sortedFields[j].Order
		}
		return sortedFields[i].Key < sortedFields[j].Key
	})

	hasTime := false

	for _, sf := range sortedFields {
		key := sf.Key
		prop := sch.Properties[key]

		goType := GetGoType(prop)
		if goType == "time.Time" || prop.Types.String() == "date" {
			goType = "*time.Time" // 通常 Model 里日期用指针方便处理 null
			hasTime = true
		}

		// 构造 Tags
		// json: camelCase
		// gorm/bson: snake_case
		jsonName := key
		snakeName := ToSnakeCase(key)
		title := prop.Title

		// 构造 Tag 字符串
		// `json:"homeId" gorm:"home_id" bson:"home_id" title:"主页ID"`
		tag := fmt.Sprintf(`json:"%s" gorm:"%s" bson:"%s" title:"%s"`, jsonName, snakeName, snakeName, title)

		field := FieldData{
			Name:     ToPascalCase(key),
			JsonName: jsonName,
			DbName:   snakeName,
			Type:     goType,
			Title:    title,
			Tag:      tag,
			Comment:  title,
		}

		modelData.Fields = append(modelData.Fields, field)
	}

	if hasTime {
		modelData.Imports = append(modelData.Imports, "time")
	}

	return modelData
}
