package model

import (
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type MasterNode struct {
	Id          string `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	CaseId      string `json:"caseId" gorm:"column:case_id;nodeLabel:true;nodeLabelFormat:case_%s"`
	TenantId    string `json:"tenantId" gorm:"column:tenant_id;nodeLabel:true;nodeLabelFormat:tenant_%s"`
	SourceIds   string `json:"sourceIds" gorm:"column:source_ids"`
	SourceType  string `json:"sourceType" gorm:"column:source_type"`
	Description string `json:"description" gorm:"column:description"`
	Type        string `json:"type" gorm:"column:type"`
	Table       string `json:"table" gorm:"column:table"`
}

func NewMasterNode(data map[string]any, dbSch *dbschema.DBSchema) *MasterNode {
	tenantId, _ := maputils.GetString(data, "tenant_id", "")
	meta := schema.GetMetaExtension(dbSch.JsonSchema)
	var nameField string = "name"
	var name string = ""

	if meta != nil {
		if meta.Graph != nil && meta.Graph.IsEnable {
			if len(meta.Graph.Name) > 0 {
				if field := dbSch.LookedField(meta.Graph.Name); field != nil {
					nameField = field.DBName
				}
			}
		}
	}
	if val, ok := data[nameField]; ok {
		name = stringutils.AnyToString(val)
	}

	caseId, _ := maputils.GetString(data, "case_id", "")
	id, _ := maputils.GetString(data, "id", "")
	desc := utils.GetNodeDescription(data, dbSch)
	typeName := utils.GetNodeType(dbSch)
	if list := stringutils.MacroValues(typeName); len(list) > 0 {
		for _, key := range list {
			dbField := stringutils.AsFieldName(key)
			value, err := maputils.GetString(data, dbField, "")
			if err != nil {
				continue
			}
			typeName = strings.Replace(typeName, "${"+key+"}", value, 1)
		}
	}
	return &MasterNode{
		//Kid:         kid,
		Id:          id,
		Name:        name,
		CaseId:      caseId,
		TenantId:    tenantId,
		SourceIds:   id,
		SourceType:  SourceType,
		Type:        typeName,
		Table:       dbSch.TableName,
		Description: desc,
	}
}
