package field

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/enums"
)

type TempUpdateField struct {
	Id         string           `json:"id"`
	CaseId     string           `json:"caseId"`
	TenantId   string           `json:"tenantId"`
	Name       string           `json:"name"`
	BankName   string           `json:"bankName"`
	SheetName  string           `json:"sheetName"`
	MasterType enums.MasterType `json:"masterType"`
	MapHeads   []*MapHead       `json:"mapHeads"`
	Fields     []*TemplateField `json:"fields"`
	Remarks    string           `json:"remarks,omitempty" desc:"备注"` // 备注
}
