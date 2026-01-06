package schema

import (
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
)

type Graph struct {
	IsEnable  bool     `json:"isEnable"`
	GraphType []string `json:"graphType"` // 图中的类型
	RelStart  string   `json:"relStart"`  // 是图关系中的开始节点字段
	RelEnd    string   `json:"relEnd"`    // 是图关系中的结束节点字段
	RelType   string   `json:"relType"`   // 是图关系中的类型字段
	RelName   string   `json:"relName"`   // 图关系名称字段
	Labels    []string `json:"labels"`    // 是图节点的标签
	Name      string   `json:"name"`      // 名称字段名
}

type GraphType string

const (
	GraphType_Node GraphType = "node"
	GraphType_Rel  GraphType = "rel"
)

func NewGraph() *Graph {
	return &Graph{
		IsEnable: false,
	}
}

func (g *Graph) Valid() error {
	if g.IsEnable == false {
		return nil
	}
	err := errors.NewVerifyError()
	if g.IsNodeType() {
		if len(g.Labels) == 0 {
			err.AppendField("labels", "must have labels")
		}
		if len(g.Name) == 0 {
			err.AppendField("Name", "must have name")
		}
	}
	if g.IsRelType() {
		if len(g.RelEnd) == 0 {
			err.AppendField("relEnd", "must have relEnd")
		}
		if len(g.RelStart) == 0 {
			err.AppendField("relStart", "must have relStart")
		}
	}
	return nil
}

func (g *Graph) IsRelTypeField() bool {
	return strings.Contains(g.RelType, "${")
}

func (g *Graph) GetRelTypeField() string {
	list := stringutils.MacroValues(g.RelType)
	if len(list) > 0 {
		return list[0]
	}
	panic("invalid graph relType")
}

func (g *Graph) IsRelStartField() bool {
	return strings.Contains(g.RelStart, "${")
}

func (g *Graph) IsRelEndField() bool {
	return strings.Contains(g.RelEnd, "${")
}

func (g *Graph) IsNodeType() bool {
	for _, v := range g.GraphType {
		if v == string(GraphType_Node) {
			return true
		}
	}
	return false
}

func (g *Graph) IsRelType() bool {
	for _, v := range g.GraphType {
		if v == string(GraphType_Rel) {
			return true
		}
	}
	return false
}

func (g *Graph) init(ctx *jsonschema.CompilerContext, values map[string]any) error {
	for k, v := range values {
		switch k {
		case "isEnable":
			g.IsEnable = v.(bool)
		case "name":
			g.Name = v.(string)
		case "relStart":
			g.RelStart = v.(string)
		case "relEnd":
			g.RelEnd = v.(string)
		case "relType":
			g.RelType = v.(string)
		case "labels":
			labels, err := GetStrings(v)
			if err != nil {
				return errors.New("invalid graph labels " + err.Error())
			}
			g.Labels = labels
		case "graphType":
			graphType, err := GetStrings(v)
			if err != nil {
				return errors.New("invalid graph graphType " + err.Error())
			}
			g.GraphType = graphType
		}
	}
	return g.Valid()
}
