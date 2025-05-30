package vis_network

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/graph/view"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"strconv"
	"strings"
	"time"
)

func getMinMaxTime(value any, minTime time.Time, maxTime time.Time) (time.Time, time.Time) {
	if vDate, err := timeutils.StrToDateTime(value.(string)); err == nil {
		if minTime.After(vDate) {
			minTime = vDate
		}
		if maxTime.Before(vDate) {
			maxTime = vDate
		}
	}
	return minTime, maxTime
}

/*
func (v *VisData) AddTagType(tagId string) bool {
	if strings.HasPrefix(tagId, "tenant_") {
		return false
	}
	if strings.HasPrefix(tagId, "graph_") {
		return false
	}
	if _, ok := v.TagTypes[tagId]; ok {
		return false
	}
	tagType := &TagType{
		Id:    tagId,
		Title: tagId,
	}

	v.TagTypes[tagId] = tagType
	switch tagId {
	case "human":
		{
			tagType.Title = "人员"
			tagType.Color = "#FF9900"
		}
	case "product":
		{
			tagType.Title = "产品"
			tagType.Color = "#2B7CE9"
		}
	case "company":
		{
			tagType.Title = "公司"
			tagType.Color = "#5A1E5C"
		}
	case "contract":
		{
			tagType.Title = "合同"
			tagType.Color = "#C5000B"
		}
	case "account":
		{
			tagType.Title = "帐号"
			tagType.Color = "#109618"
		}
	case "payment":
		{
			tagType.Title = "支付"
			tagType.Color = "#F09618"
		}
	default:
		{
			tagType.Color = getColor()
		}
	}
	return true
}

*/

func (v *Node) String(s string) *string {
	return &s
}

func (v *Edge) String(s string) *string {
	return &s
}

func newVisEdge(e *view.GraphEdge) *Edge {
	edge := &Edge{}
	edge.Id = strconv.FormatInt(e.Identity, 10)
	edge.From = strconv.FormatInt(e.Start, 10)
	edge.To = strconv.FormatInt(e.End, 10)
	edge.Arrows = edge.String("to")
	edge.Label = e.Type
	edge.Props = e.Props
	if e.Hidden {
		edge.Hidden = &e.Hidden
	}
	if edge.Label == view.EdgeTypeRecord.Name() {
		var color = view.EdgeTypeRecord.Color()
		edge.Color = &color
		edge.Label = ""
	}
	if len(e.SubEdges) > 0 {
		subEdges := make([]any, 0)
		for _, item := range e.SubEdges {
			subEdges = append(subEdges, &item.Props)
		}
		edge.SubEdges = subEdges
	}
	return edge
}

// newVisNode
// @Description: 创建新的节点
// @param n
// @return visNode
// @return err
func newVisNode(n *view.GraphNode) (visNode *Node, err error) {
	node := &Node{}
	node.Id = strconv.FormatInt(n.Identity, 10)
	node.Props = n.Props
	if n.Hidden {
		node.Hidden = &n.Hidden
	}
	for _, v := range n.SubNodes {
		if newNode, err := newVisNode(v); err != nil {
			return nil, err
		} else {
			node.SubNodes = append(node.SubNodes, newNode.Id)
		}
	}
	if node.Font == nil {
		node.Font = &Font{
			Color: pstr("#ffffff"),
		}
	}
	return node, nil
}

func newFont(color string) *Font {
	return &Font{
		Color: pstr(color),
	}
}

func newTags(labels []string) map[string]string {
	res := make(map[string]string)
	for _, v := range labels {
		res[v] = v
	}
	return res
}

/**
 * 颜色代码转换为RGB
 * input int
 * output int red, green, blue
 **/
func colorToRGB(color int) (red, green, blue int) {
	red = color >> 16
	green = (color & 0x00FF00) >> 8
	blue = color & 0x0000FF
	return
}

func getColor() string {
	color_str := "0xAD99C0"                         //颜色代码的值
	color_str = strings.TrimPrefix(color_str, "0x") //过滤掉16进制前缀

	color64, err := strconv.ParseInt(color_str, 16, 32) //字串到数据整型
	if err != nil {
		panic(err)
	}
	color32 := int(color64) //类型强转
	r, g, b := colorToRGB(color32)
	rgb := fmt.Sprintf("#%d%d%d", r, g, b)
	return rgb
}

func pstr(v string) *string {
	return &v
}

func pint(v int64) *int64 {
	return &v
}

func pbool(v bool) *bool {
	return &v
}
