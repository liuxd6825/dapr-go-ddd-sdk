package xbase

type Command[T any] struct {
	CommandId   string   `json:"commandId"  validate:"required" title:"命令ID"` // 命令ID
	IsValidOnly bool     `json:"isValidOnly"  validate:"-" title:"仅验证"`       // 是否仅验证，不执行
	UpdateMask  []string `json:"updateMask"`
	Data        T        `json:"data"  validate:"required" title:"命令数据" `
}

type DeleteByIdCommand = Command[IdField]

type IdField struct {
	Id string `json:"id" validate:"required"`
}

type IdsField struct {
	Id []string `json:"ids" validate:"required"`
}
