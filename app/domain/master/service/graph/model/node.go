package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
)

type Node struct {
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

func NewNode(data map[string]any, dbSch *dbschema.DBSchema) *Node {
	tenantId, _ := maputils.GetString(data, "tenant_id", "")
	name, _ := maputils.GetString(data, "name", "")
	caseId, _ := maputils.GetString(data, "case_id", "")
	id, _ := maputils.GetString(data, "id", "")
	desc := utils.GetDescription(data, dbSch)
	typeName := utils.GetNodeType(dbSch)
	return &Node{
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
