package model

import (
	"errors"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type MasterRelation struct {
	Id          string   `json:"id" gorm:"column:id"`
	CaseId      string   `json:"caseId" gorm:"column:case_id;nodeLabel:true;nodeLabelFormat:case_%"`
	TenantId    string   `json:"tenantId" gorm:"column:tenant_id;nodeLabel:true;nodeLabelFormat:tenant_%"`
	Source      string   `json:"source" gorm:"column:source;relStartId:true"`
	Target      string   `json:"target" gorm:"column:target;relEndId:true"`
	SourceIds   string   `json:"sourceIds" gorm:"column:source_ids"`
	SourceType  string   `json:"sourceType" gorm:"column:source_type"`
	RelType     string   `json:"relType" gorm:"column:rel_type;relType:true"`
	Keywords    []string `json:"keywords" gorm:"column:keywords;type:text;serializer:json"`
	Description string   `json:"description" gorm:"column:description"`
	Table       string   `json:"table" gorm:"column:table"`
}

func NewMasterRelation(data map[string]any, dbSchema *dbschema.DBSchema) *MasterRelation {
	if dbSchema == nil {
		panic(errors.New("dbSch is nil"))
	}
	tenantId, _ := maputils.GetString(data, "tenant_id", "")
	caseId, _ := maputils.GetString(data, "case_id", "")
	id, _ := maputils.GetString(data, "id", "")
	relStartId := getRelStartId(dbSchema, data)
	relEndId := getRelEndId(dbSchema, data)
	relType := getRelType(dbSchema, data)
	desc := utils.GetDescription(data, dbSchema)
	return &MasterRelation{
		Id:          id,
		CaseId:      caseId,
		TenantId:    tenantId,
		Keywords:    []string{relType},
		RelType:     relType,
		Source:      relStartId,
		Target:      relEndId,
		SourceIds:   id,
		Table:       dbSchema.TableName,
		SourceType:  SourceType,
		Description: desc,
	}
}

func getRelStartId(dbSch *dbschema.DBSchema, data map[string]any) string {
	return getGraphMetaValue(dbSch, data, func(graph *schema.Graph) string {
		return graph.RelStart
	})
}

func getRelEndId(dbSch *dbschema.DBSchema, data map[string]any) string {
	return getGraphMetaValue(dbSch, data, func(graph *schema.Graph) string {
		return graph.RelEnd
	})
}

func getRelType(dbSch *dbschema.DBSchema, data map[string]any) string {
	return getGraphMetaValue(dbSch, data, func(graph *schema.Graph) string {
		return graph.RelType
	})
}

// getGraphMetaValue
// @Description: 取实体数据值，通过jsonSchema中的meta.graph配置内容
// @param dbSch
// @param data
// @param getFieldName
// @return string
func getGraphMetaValue(dbSch *dbschema.DBSchema, data map[string]any, getFieldName func(graph *schema.Graph) string) string {
	meta := schema.GetMetaExtension(dbSch.JsonSchema)
	if meta != nil && meta.Graph != nil {
		graphMeta := meta.GetGraph()
		if graphMeta.IsEnable {
			value := getFieldName(meta.Graph)
			names := stringutils.MacroValues(value)
			for _, name := range names {
				dbName := stringutils.AsFieldName(name)
				val, err := maputils.GetString(data, dbName, "")
				if err != nil {
					panic(err)
				}
				value = strings.ReplaceAll(value, "${"+name+"}", val)
			}
			return value
		}
	}
	panic(errors.New("getGraphMetaValue() schema.meta.Graph is nil"))
}
