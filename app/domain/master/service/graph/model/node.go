package model

import (
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
)

type Node struct {
	Id          string `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	CaseId      string `json:"caseId" gorm:"column:case_id;nodeLabel:true;nodeLabelFormat:case_%s"`
	TenantId    string `json:"tenantId" gorm:"column:tenant_id;nodeLabel:true;nodeLabelFormat:tenant_%s"`
	SourceIds   string `json:"sourceIds" gorm:"column:source_ids"`
	Description string `json:"description" gorm:"column:description"`
	Type        string `json:"type" gorm:"column:type"`
}

func NewNode(tableName string, data map[string]any) *Node {
	tenantId, _ := maputils.GetString(data, "tenantId", "")
	name, _ := maputils.GetString(data, "name", "")
	id, _ := maputils.GetString(data, "id", "")
	caseId, _ := maputils.GetString(data, "caseId", "")
	desc := description(data)
	return &Node{
		Id:          id,
		Name:        name,
		CaseId:      caseId,
		TenantId:    tenantId,
		SourceIds:   tableName,
		Type:        tableName,
		Description: desc,
	}
}

func description(data map[string]any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
