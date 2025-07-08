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

type TaskState string

const (
	TaskStateNone      TaskState = ""
	TaskStateEditing   TaskState = "编辑中"
	TaskStateGenerated TaskState = "已生成"
	TaskStateImported  TaskState = "已导入"
)

func (t TaskState) Name() string {
	return string(t)
}

func (t MasterType) Name() string {
	return string(t)
}
