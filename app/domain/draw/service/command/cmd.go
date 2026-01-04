package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/pkg/mxgraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type CreateCommand struct {
	CommandId string `json:"commandId"`
	Data      *Data  `json:"data"`
}

type DeleteByIdsCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		Ids    []string `json:"ids"`
		CaseId string   `json:"caseId"`
	} `json:"data"`
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

type DrawSaveFileCommand struct {
	xbase.Command[DrawSaveFileCommandData]
}

type DrawSaveFileCommandData struct {
	CaseId   string            `json:"caseId" param:"case-id" validate:"required"`
	DrawId   string            `json:"drawId" param:"draw-id" validate:"required"`
	FileName string            `json:"fileName" validate:"required"`
	XML      string            `json:"xml" validate:"required"`
	Diff     *mxgraph.FileDiff `json:"diff" validate:"required"`
}

func (c *DrawSaveFileCommand) GetCaseId() string {
	return c.Data.CaseId
}

func (c *DrawSaveFileCommand) GetDrawId() string {
	return c.Data.DrawId
}

func (c *DrawSaveFileCommand) GetFileName() string {
	return c.Data.FileName
}

func (c *DrawSaveFileCommand) GetXML() string {
	return c.Data.XML
}

func (c *DrawSaveFileCommand) GetDiff() *mxgraph.FileDiff {
	return c.Data.Diff
}
