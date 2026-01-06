package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/pkg/mxgraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type DrawSaveCommand struct {
	xbase.Command[DrawSaveCommandData]
}

type IDrawSaveCommand interface {
	GetCaseId() string
	GetDrawId() string
	GetFileName() string
	GetXML() string
	GetDiff() *mxgraph.FileDiff
}

type DrawSaveCommandData struct {
	CaseId   string            `json:"caseId" param:"case-id"`
	DrawId   string            `json:"drawId" param:"draw-id"`
	FileName string            `json:"fileName" `
	XML      string            `json:"xml"`
	Diff     *mxgraph.FileDiff `json:"diff"`
}

func (c *DrawSaveCommand) GetCaseId() string {
	return c.Data.CaseId
}

func (c *DrawSaveCommand) GetDrawId() string {
	return c.Data.DrawId
}

func (c *DrawSaveCommand) GetFileName() string {
	return c.Data.FileName
}

func (c *DrawSaveCommand) GetXML() string {
	return c.Data.XML
}

func (c *DrawSaveCommand) GetDiff() *mxgraph.FileDiff {
	return c.Data.Diff
}
