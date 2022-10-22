package metadata

import (
	"fmt"
	"reflect"
	"testing"
)

type BankBase struct {
	AccountCode string `json:"accountCode" validate:"-" description:"账号"`            // 账号
	BankName    string `json:"bankName" validate:"-" description:"银行名称"`             // 银行名称
	CaseId      string `json:"caseId" validate:"required"  description:"案件ID"`       // 案件ID
	Columns     string `json:"columns" copier:"-" validate:"-"  description:"扫描列信息"` // 扫描列信息
}

type BankBaseMetadata struct {
	AccountCode Property `json:"accountCode"` // 账号
	BankName    Property `json:"bankName"`    // 银行名称
	CaseId      Property `json:"caseId"`      // 案件ID
	Columns     Property `json:"columns"`     // 扫描列信息
	Properties  Properties
}

type BankTable struct {
	BankBase
	EndTime   string `json:"endTime" validate:"-"  description:"流水结束时间"`                     // 流水结束时间
	FileId    string `json:"fileId" validate:"required" description:"扫描文件ID"`                // 扫描文件ID
	Id        string `json:"id" validate:"required"  description:"主键"`                       // 主键
	IsDeleted bool   `json:"isDeleted" validate:"-"  description:"已删除"`                      // 已删除
	Pages     string `json:"pages" copier:"-" validate:"-"  description:"扫描页,用于管理表中的页数和页顺序"` // 扫描页,用于管理表中的页数和页顺序
	Remarks   string `json:"remarks" validate:"-"  description:"备注"`                         // 备注
	StartTime string `json:"startTime" validate:"-" description:"流水开始时间"`                    // 流水开始时间
	TenantId  string `json:"tenantId" validate:"required" description:"租户ID"`                // 租户ID
	Title     string `json:"title" validate:"-"  description:"标题"`                           // 标题
}

type BankTableMetadata struct {
	BankBaseMetadata
	Properties  Properties
	AccountCode Property `json:"accountCode"` // 账号
	/*
		BankName    Property `json:"bankName"`    // 银行名称
		CaseId      Property `json:"caseId"`      // 案件ID
		Columns     Property `json:"columns"`     // 扫描列信息
		EndTime     Property `json:"endTime"`     // 流水结束时间
		FileId      Property `json:"fileId"`      // 扫描文件ID
		Id          Property `json:"id"`          // 主键
		IsDeleted   Property `json:"isDeleted"`   // 已删除
		Pages       Property `json:"pages"`       // 扫描页,用于管理表中的页数和页顺序
		Remarks     Property `json:"remarks"`     // 备注
		StartTime   Property `json:"startTime"`   // 流水开始时间
		TenantId    Property `json:"tenantId"`    // 租户ID
		Title       Property `json:"title"`       // 标题
	*/

}

func TestNewProperties(t *testing.T) {
	m := &BankTableMetadata{}
	entity := &BankTable{}
	rValue := reflect.ValueOf(m)
	rType := reflect.TypeOf(m)
	for i := 0; i < rType.Elem().NumField(); i++ {
		fieldName := rType.Elem().Field(i).Name
		fv := rValue.Elem().FieldByName(fieldName)
		t.Log(fv.Interface())
	}

	if prop, err := NewProperties(m, entity, t); err != nil {
		t.Error(err)
	} else {
		m.Properties = prop
		fmt.Println(m.BankBaseMetadata.BankName)
		for _, p := range m.Properties.Values() {
			t.Logf("name=%v; jsonName=%v; type=%v; description=%v", p.Name(), p.JsonName(), p.TypeName(), p.Description())
		}
	}
}
