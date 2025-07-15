package enums

type MasterType string

const (
	MasterTypeRecord   MasterType = "银行流水"
	MasterTypeContract            = "合同"
	MasterTypeProduct             = "产品"
	MasterTypeHuman               = "人员"
	MasterTypeCompany             = "公司"
	MasterTypeAccount             = "账号"
)

func (t MasterType) Name() string {
	return string(t)
}
