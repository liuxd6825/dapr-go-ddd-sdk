package enums

type TaskState string

const (
	TaskStateNone      TaskState = ""
	TaskStateEditing   TaskState = "编辑中"
	TaskStateGenerated TaskState = "已生成"
	TaskStateImported  TaskState = "已导入"
	TaskStateError     TaskState = "导入错误"
)

func (t TaskState) Name() string {
	return string(t)
}
