package command

type CreateCommand struct {
	CommandId string `json:"commandId"`
	Data      *Data  `json:"data"`
}

type DeleteByIdsCommand struct {
	CommandId string   `json:"commandId"`
	Data      []string `json:"data"`
}

type SubmitCommand struct {
	CommandId string `json:"commandId"`
	Data      *Data  `json:"data"`
}

type UpdateCommand struct {
	CommandId string `json:"commandId"`
	Data      *Data  `json:"data"`
}

type DeleteCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		Id string `json:"id" gorm:"type:varchar(50);primary_key"  `
	} `json:"data"`
}

type Data struct {
	Id     string `json:"id" gorm:"type:varchar(50);primary_key"`
	Name   string `json:"name" gorm:"name"`
	CaseId string `json:"caseId" gorm:"type:varchar(50);case_id"`
	Status string `json:"status" gorm:"status"`
	Remark string `json:"remark" gorm:"remark"`
}
