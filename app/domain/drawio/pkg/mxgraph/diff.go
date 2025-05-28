package mxgraph

import (
	"encoding/json"
	"github.com/antchfx/xmlquery"
	"strings"
)

type BaseCell struct {
	Id     string  `json:"id"`
	Edge   *int    `json:"edge"`   //是否为关系
	Source *string `json:"source"` // 关系来源
	Target *string `json:"target"` // 关系目录
	style  map[string]string
}

type Extend struct {
	Type     string   `json:"type"`
	Label    string   `json:"label"`
	Labels   []*Label `json:"labels"`
	ParentId string   `json:"parentId"`
	SourceId string   `json:"sourceId"`
	TargetId string   `json:"targetId"`
	OldValue string   `json:"oldValue"`
}

type Label struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type DiffCell struct {
	Geometry *string `json:"geometry"`
	Id       string  `json:"id"`
	Parent   *string `json:"parent"`
	Previous *string `json:"previous"`
	Style    *string `json:"style"`
	Vertex   *int64  `json:"vertex"`
	XmlValue *string `json:"xmlValue"`
	Value    *string `json:"value"`
	Edge     *int    `json:"edge"`   //是否为关系
	Source   *string `json:"source"` // 关系来源
	Target   *string `json:"target"` // 关系目录
	Extend   Extend  `json:"extend"`
	State    DiffState
	newXml   *xmlquery.Node
	oldXml   *xmlquery.Node
}

type DiffState int

const (
	DiffState_Remove DiffState = 0
	DiffState_Insert DiffState = 1
	DiffState_Update DiffState = 2
)

type DiffCells struct {
	I []*DiffCell          `json:"i"`
	R []*DiffCell          `json:"r"`
	U map[string]*DiffCell `json:"u"`
}

type UpdatePage struct {
	Name          *string    `json:"name"`
	PageProperty2 *string    `json:"page_property2"`
	Cells         *DiffCells `json:"cells"`
}

type InsertPage struct {
	Id       *string `json:"id"`
	Data     *string `json:"data"`
	Previous *string `json:"previous"`
}

type FileDiff struct {
	I []*InsertPage          `json:"i"`
	R []string               `json:"r"`
	U map[string]*UpdatePage `json:"u"`
}

func NewFileDiff(jsonText string) *FileDiff {
	var fileDiff FileDiff
	err := json.Unmarshal([]byte(jsonText), &fileDiff)
	if err != nil {
		panic(err)
	}
	return &fileDiff
}

func (c *DiffCell) GetNodeName() string {
	if c.Extend.OldValue != "" {
		oldXml := c.GetOldXml()
		if oldXml != nil && oldXml.Data == "object" {
			return oldXml.SelectAttr("label")
		}
	} else if c.XmlValue != nil {
		newXml := c.GetNewXml()
		if newXml != nil && newXml.Data == "object" {
			return newXml.SelectAttr("label")
		}
	}
	return ""
}

func (c *DiffCell) GetNodeLabels() []string {
	var labels []string

	if c.Extend.OldValue != "" {
		oldXml := c.GetOldXml()
		if oldXml != nil && oldXml.Type == xmlquery.ElementNode {
			return labels
		}
	} else if c.XmlValue != nil {
		newXml := c.GetNewXml()
		return []string{c.GetObjectLabel(newXml)}
	}
	return labels
}

func (c *DiffCell) GetObjectLabel(xml *xmlquery.Node) string {
	if xml != nil && xml.Data == "object" {
		return xml.SelectAttr("label")
	}
	return ""
}

func (c *DiffCell) GetObjectCellType(xml *xmlquery.Node) string {
	if xml != nil && xml.Data == "object" {
		return xml.SelectAttr("cellType")
	}
	return ""
}

func (c *DiffCell) GetNewXml() *xmlquery.Node {
	if c.XmlValue == nil || *c.XmlValue == "" {
		return nil
	}
	if c.newXml != nil {
		return c.newXml
	}
	newXml, err := xmlquery.Parse(strings.NewReader(*c.XmlValue))
	if err != nil {
		panic(err)
	}
	c.newXml = newXml.LastChild
	return c.newXml
}

func (c *DiffCell) GetOldXml() *xmlquery.Node {
	if c.Extend.OldValue == "" {
		return nil
	}
	if c.oldXml != nil {
		return c.oldXml
	}
	oldXml, err := xmlquery.Parse(strings.NewReader(c.Extend.OldValue))
	if err != nil {
		panic(err)
	}
	c.oldXml = oldXml
	return c.oldXml
}

func (c *DiffCell) GetTagType() string {
	xml := c.Extend.OldValue
	i := strings.Index(xml, " ")
	if i < 2 {
		return ""
	}
	return xml[1 : i-1]
}

func (c *DiffCell) GetRelType() string {
	return c.Extend.Label
}

func (c *DiffCell) GetTargetId() string {
	return c.Extend.TargetId
}

func (c *DiffCell) GetSourceId() string {
	return c.Extend.SourceId
}

func (c *DiffCell) IsEdge() bool {
	return c.Extend.Type == "edge" || c.Extend.Type == "edgeLabel"
}

func (c *DiffCell) IsNode() bool {
	return c.Extend.Type == "node"
}
