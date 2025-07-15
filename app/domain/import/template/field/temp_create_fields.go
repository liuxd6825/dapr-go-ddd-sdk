package field

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/enums"
)

type TempCreateField struct {
	Id         string           `json:"id"`
	CaseId     string           `json:"caseId"`
	TenantId   string           `json:"tenantId"`
	Name       string           `json:"name"`
	BankName   string           `json:"bankName"`
	FileId     string           `json:"fileId"`
	FileName   string           `json:"fileName"`
	SheetName  string           `json:"sheetName"`
	MasterType enums.MasterType `json:"masterType"`
	MapHeads   []*MapHead       `json:"mapHeads"`
	Fields     []*TemplateField `json:"fields"`
	Remark     string           `json:"remark,omitempty" desc:"备注"` // 备注
}

/*
func NewRecordTemplateFields() []*TemplateField {
	var fields []*TemplateField
	fields = append(fields, &TemplateField{Key: "iden", Name: ""})
	fields = append(fields, &TemplateField{Key: "name", Name: ""})
	fields = append(fields, &TemplateField{Key: "acct", Name: ""})
	fields = append(fields, &TemplateField{Key: "acctType", Name: ""})
	fields = append(fields, &TemplateField{Key: "category", Name: ""})
	fields = append(fields, &TemplateField{Key: "bankName", Name: ""})
	fields = append(fields, &TemplateField{Key: "balance", Name: ""})

	fields = append(fields, &TemplateField{Key: "oppIden", Name: ""})
	fields = append(fields, &TemplateField{Key: "oppName", Name: ""})
	fields = append(fields, &TemplateField{Key: "oppAcct", Name: ""})
	fields = append(fields, &TemplateField{Key: "oppAcctType", Name: ""})
	fields = append(fields, &TemplateField{Key: "oppCategory", Name: ""})
	fields = append(fields, &TemplateField{Key: "oppBankName", Name: ""})

	fields = append(fields, &TemplateField{Key: "serial", Name: ""})
	fields = append(fields, &TemplateField{Key: "income", Name: ""})
	fields = append(fields, &TemplateField{Key: "payout", Name: ""})
	fields = append(fields, &TemplateField{Key: "amount", Name: ""})
	fields = append(fields, &TemplateField{Key: "date", Name: ""})

	fields = append(fields, &TemplateField{Key: "type", Name: ""})
	fields = append(fields, &TemplateField{Key: "ccy", Name: ""})
	fields = append(fields, &TemplateField{Key: "place", Name: ""})
	fields = append(fields, &TemplateField{Key: "summary", Name: ""})
	fields = append(fields, &TemplateField{Key: "notes", Name: ""})

	return fields
}

*/
