package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"

type BusNode struct {
	Id       string `json:"id" gorm:"column:id"`
	CaseId   string `json:"caseId" gorm:"column:case_id;nodeLabel:true;nodeLabelFormat:case_%s"`
	Name     string `json:"name" gorm:"column:name"`
	TenantId string `json:"tenantId" gorm:"column:tenant_id;nodeLabel:true;nodeLabelFormat:tenant_%s"`
	Table    string `json:"table" gorm:"column:table"`
}

func NewBusNode(tableName string, data map[string]any) *BusNode {
	tenantId, _ := maputils.GetString(data, "tenantId", "")
	name, _ := maputils.GetString(data, "name", "")
	id, _ := maputils.GetString(data, "id", "")
	caseId, _ := maputils.GetString(data, "caseId", "")
	return &BusNode{
		Id:       id,
		Name:     name,
		CaseId:   caseId,
		TenantId: tenantId,
		Table:    tableName,
	}
}
