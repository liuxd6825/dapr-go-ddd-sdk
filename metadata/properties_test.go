package metadata

import "testing"

type BankTableAggregateMetadata struct {
	AccountCode Property `json:"accountCode" validate:"-"`        // 账号
	BankName    Property `json:"bankName" validate:"-"`           // 银行名称
	CaseId      Property `json:"caseId" validate:"required"`      // 案件ID
	Columns     Property `json:"columns" copier:"-" validate:"-"` // 扫描列信息
	EndTime     Property `json:"endTime" validate:"-"`            // 流水结束时间
	FileId      Property `json:"fileId" validate:"required"`      // 扫描文件ID
	Id          Property `json:"id" validate:"required"`          // 主键
	IsDeleted   Property `json:"isDeleted" validate:"-"`          // 已删除
	Pages       Property `json:"pages" copier:"-" validate:"-"`   // 扫描页,用于管理表中的页数和页顺序
	Remarks     Property `json:"remarks" validate:"-"`            // 备注
	StartTime   Property `json:"startTime" validate:"-"`          // 流水开始时间
	TenantId    Property `json:"tenantId" validate:"required"`    // 租户ID
	Title       Property `json:"title" validate:"-"`              // 标题
	Properties  Properties
}

func TestNewProperties(t *testing.T) {
	m := &BankTableAggregateMetadata{
		AccountCode: NewProperty("AccountCode", DataTypeStr, "账号"),
		BankName:    NewProperty("BankName", DataTypeStr, "银行名称"),
		CaseId:      NewProperty("CaseId", DataTypeStr, "案件ID"),
		Columns:     NewProperty("Columns", DataTypeStruct, "扫描列信息"),
		EndTime:     NewProperty("EndTime", DataTypeStr, "流水结束时间"),
		FileId:      NewProperty("FileId", DataTypeStr, "扫描文件ID"),
		Id:          NewProperty("Id", DataTypeStr, "主键"),
		IsDeleted:   NewProperty("IsDeleted", DataTypeBool, "已删除"),
		Pages:       NewProperty("Pages", DataTypeStruct, "扫描页,用于管理表中的页数和页顺序"),
		Remarks:     NewProperty("Remarks", DataTypeStr, "备注"),
		StartTime:   NewProperty("StartTime", DataTypeTime, "流水开始时间"),
		TenantId:    NewProperty("TenantId", DataTypeStr, "租户ID"),
		Title:       NewProperty("Title", DataTypeStr, "标题"),
	}
	m.Properties = NewProperties(*m)
}
