package view

import (
	"fmt"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
)

type GraphDataView struct {
	Nodes     map[int64]*GraphNode `json:"nodes"`
	Edges     map[int64]*GraphEdge `json:"lines"`
	PageNum   *int64               `json:"pageNum"`
	PageSize  *int64               `json:"pageSize"`
	Total     *int64               `json:"total"`
	PageTotal *int64               `json:"pageTotal"`
}

type NodeType string

const (
	NodeTypeAccount  NodeType = "account"
	NodeTypeCompany  NodeType = "company"
	NodeTypeHuman    NodeType = "human"
	NodeTypeProduct  NodeType = "product"
	NodeTypePayment  NodeType = "payment"
	NodeTypeContract NodeType = "contract"
)

func (e NodeType) Name() string {
	return string(e)
}

func (e NodeType) Color() string {
	var color = ""
	switch e {
	case NodeTypeAccount: // 橙色
		color = "#F75000"
	case NodeTypeCompany: // 绿色
		color = "#73BF00"
	case NodeTypeHuman: // 黄色
		color = "#FFDC35"
	case NodeTypeProduct: // 棕色
		color = "#A23400"
	case NodeTypePayment: // 蓝色
		color = "#484891"
	case NodeTypeContract: // 紫色
		color = "#8F4586"
	}
	return color
}

func (e NodeType) FontColor() string {
	var color = "#FFFFFF"
	switch e {
	case NodeTypeAccount:
		color = "#FFFFFF"
	case NodeTypeCompany:
		color = "#FFFFFF"
	case NodeTypeHuman: // 灰色
		color = "#272727"
	case NodeTypeProduct:
		color = "#FFFFFF"
	case NodeTypePayment:
		color = "#FFFFFF"
	case NodeTypeContract:
		color = "#FFFFFF"
	}
	return color
}

type GraphNode struct {
	Identity int64                  `json:"identity"`           // Identity
	Labels   []string               `json:"labels"`             // 标签
	Props    map[string]interface{} `json:"props"`              // 属性
	SubNodes map[string]*GraphNode  `json:"subNodes,omitempty"` // 子节点
	Hidden   bool                   `json:"hidden"`             // 是否隐藏
}
type EdgeType string

const (
	EdgeTypeAccountHolder EdgeType = "开户人"
	EdgeTypeRecord        EdgeType = "银行流水"
)

func (e EdgeType) Name() string {
	return string(e)
}

func (e EdgeType) Color() string {
	var color = ""
	switch e {
	case EdgeTypeRecord:
		color = "#FF2D2D"
	case EdgeTypeAccountHolder:
		color = "#0000C6"
	}
	return color
}

type GraphEdge struct {
	Identity int64                  `json:"identity"`           // neo4j id
	Type     string                 `json:"type"`               // 标签
	Start    int64                  `json:"start"`              // 开始节点Id
	End      int64                  `json:"end"`                // 结束节点Id
	Props    map[string]interface{} `json:"props"`              // 属性
	SubEdges map[string]*GraphEdge  `json:"subEdges,omitempty"` // 子关系
	Hidden   bool                   `json:"hidden"`             // 是否隐藏
	Virtual  bool                   `json:"virtual"`            // 是虚拟
}

type MergeEdge struct {
	Key   *MergeEdgeKey         `json:"key"`
	Items map[string]*GraphEdge `json:"subItems"` // 子关系
}

type MergeEdgeKey struct {
	RelType string `json:"relType"` // 关系类型
	Sid     int64  `json:"sid"`     // 开始节点Id
	Eid     int64  `json:"eid"`     // 结束节点Id
	Key     string `json:"key"`
}

func NewMergeEdges(key *MergeEdgeKey) *MergeEdge {
	ml := &MergeEdge{
		Key:   key,
		Items: make(map[string]*GraphEdge),
	}
	return ml
}

type GraphDataViewOptions struct {
	MergeRelTypes *[]string `json:"mergeRelTypes"`
	IsMergeRel    *bool     `json:"isMergeRel"`
}

func NewGraphDataViewEmpty() *GraphDataView {
	return &GraphDataView{
		Nodes: map[int64]*GraphNode{},
		Edges: map[int64]*GraphEdge{},
	}
}

func NewMergeEdgeKey(relType string, sid int64, eid int64) *MergeEdgeKey {
	return &MergeEdgeKey{
		relType, sid, eid, fmt.Sprintf("%d-%d-%v", sid, eid, relType),
	}
}

func NewGraphDataView(nodes map[int64]*GraphNode, lines map[int64]*GraphEdge, pageNum *int64, pageSize *int64, total *int64, opts ...*GraphDataViewOptions) *GraphDataView {
	opt := &GraphDataViewOptions{}
	opt.Merge(opts...)
	res := &GraphDataView{}
	res.Nodes = nodes
	if opt.IsMergeRel != nil && *opt.IsMergeRel == true {
		res.Edges = res.mergeLines(lines)
	} else {
		res.Edges = lines
	}

	res.Total = total
	res.PageNum = pageNum
	res.PageSize = pageSize

	var pageTotal *int64
	if pageSize != nil && pageNum != nil {
		var vPageSize int64 = *pageSize
		var vTotal int64 = *total / vPageSize
		if (*total % vPageSize) > 0 {
			vTotal = vTotal + 1
		}
		pageTotal = &vTotal
	}
	res.PageTotal = pageTotal

	return res
}

func (g *GraphDataView) mergeLines(lines map[int64]*GraphEdge) map[int64]*GraphEdge {
	merge := make(map[string]*MergeEdge)
	for _, item := range lines {
		key := NewMergeEdgeKey(item.Type, item.Start, item.End)
		if m, ok := merge[key.Key]; ok {
			m.AddItem(item)
		} else {
			ms := NewMergeEdges(key)
			ms.AddItem(item)
			merge[key.Key] = ms
		}
	}

	res := make(map[int64]*GraphEdge)
	for _, item := range merge {
		item.AddTo(res)
	}
	return res
}

func (n *GraphNode) Id() (string, error) {
	if v, ok := n.Props["id"].(string); ok {
		return v, nil
	}
	return "", errors.New("node property 'id' is null, neo4j node id = %v ", n.Identity)
}

func (n *GraphNode) AddSubNode(node *GraphNode) {
	if n.SubNodes == nil {
		n.SubNodes = make(map[string]*GraphNode)
	}
	if id, err := node.Id(); err == nil {
		n.SubNodes[id] = node
	}

}

func (n *GraphNode) GetPropertyStr(propName string, doValue ...func(value string)) (string, bool) {
	defer func() {
		if e := errors.GetRecoverError(nil, recover()); e != nil {
			fmt.Println(e)
		}
	}()
	v, ok := n.Props[propName].(string)
	if ok && len(doValue) == 1 {
		doValue[0](v)
	}
	return v, ok
}

func (n *GraphNode) GetPropertyInt(propName string, doValue ...func(value int64)) (int64, bool) {
	v, ok := n.Props[propName].(int64)
	if ok && len(doValue) == 1 {
		doValue[0](v)
	}
	return v, ok
}

func (n *GraphNode) GetPropertyTime(propName string, doValue ...func(value *time.Time)) (*time.Time, bool) {
	var res *time.Time
	v, ok := n.Props[propName].(string)
	if ok && len(doValue) == 1 {
		if v, err := timeutils.StrToDateTime(v); err != nil {
			return nil, false
		} else {
			res = &v
		}
		doValue[0](res)
	}
	return res, ok
}

func (n *GraphNode) IsLabel(label string) bool {
	for _, item := range n.Labels {
		if item == label {
			return true
		}
	}
	return false
}

func (n *GraphNode) IsHuman() bool {
	return n.IsLabel("human")
}

func (n *GraphNode) IsPayment() bool {
	return n.IsLabel("payment")
}

func (n *GraphNode) IsCompany() bool {
	return n.IsLabel("company")
}

func (n *GraphNode) IsContract() bool {
	return n.IsLabel("contract")
}

func (n *GraphNode) IsProduct() bool {
	return n.IsLabel("product")
}

func (n *GraphNode) IsAccount() bool {
	return n.IsLabel("account")
}

// FindNodesByRelType
// @Description: 根据关系类型查找节点
// @receiver n
// @param findLabels  标签
// @param relType 关系类型
// @param data
// @return []*GraphNode
// @return bool
func (n *GraphNode) FindNodesByRelType(findLabels []string, relType string, data *GraphDataView) ([]*GraphNode, bool) {
	nodes := make([]*GraphNode, 0)
	for _, line := range data.Edges {
		if line.Type == relType {
			//打到对方节点的Id
			var nodeId int64 = -1
			if line.Start == n.Identity {
				nodeId = line.End
			} else if line.End == n.Identity {
				nodeId = line.Start
			} else {
				continue
			}

			// 找到对方节点
			if node, ok := data.Nodes[nodeId]; ok {
				for _, label := range findLabels {
					if node.IsLabel(label) {
						nodes = append(nodes, node)
						break
					}
				}
			}
		}
	}
	return nodes, len(nodes) != 0
}

func (n *GraphNode) FindEdges(data *GraphDataView, start *bool, relTypes ...string) (map[int64]*GraphEdge, bool) {
	lines := make(map[int64]*GraphEdge)
	for _, item := range data.Edges {
		if len(relTypes) == 0 {
			if start == nil && (item.Start == n.Identity || item.End == n.Identity) {
				lines[item.Identity] = item
			} else if start != nil && *start == true && item.Start == n.Identity {
				lines[item.Identity] = item
			} else if start != nil && *start == false && item.End == n.Identity {
				lines[item.Identity] = item
			}
		} else {
			for _, relType := range relTypes {
				if item.Type == relType {
					if start == nil && (item.Start == n.Identity || item.End == n.Identity) {
						lines[item.Identity] = item
					} else if start != nil && *start == true && item.Start == n.Identity {
						lines[item.Identity] = item
					} else if start != nil && *start == false && item.End == n.Identity {
						lines[item.Identity] = item
					}
				}
			}
		}
	}
	return lines, len(lines) != 0
}

func (n *GraphEdge) Id() (string, error) {
	if v, ok := n.Props["id"].(string); ok {
		return v, nil
	}
	return "", errors.New("edge property 'id' is null, neo4j edge id = %v ", n.Identity)
}

func (n *GraphEdge) StartId() (string, error) {
	if v, ok := n.Props["startId"].(string); ok {
		return v, nil
	}
	return "", errors.New("edge property 'startId' is null, neo4j edge id = %v ", n.Identity)
}

func (n *GraphEdge) EndId() (string, error) {
	if v, ok := n.Props["endId"].(string); ok {
		return v, nil
	}
	return "", errors.New("edge property 'endId' is null, neo4j edge id = %v ", n.Identity)
}

func (n *GraphEdge) GetPropertyValue(propName string, doValue ...func(value string)) (string, bool) {
	v, ok := n.Props[propName].(string)
	if ok && len(doValue) == 1 {
		doValue[0](v)
	}
	return v, ok
}

func (o *GraphDataViewOptions) Merge(opts ...*GraphDataViewOptions) *GraphDataViewOptions {
	for _, opt := range opts {
		if opt.IsMergeRel != nil {
			o.IsMergeRel = opt.IsMergeRel
		}
	}
	return o
}

func (m *MergeEdge) AddItem(item *GraphEdge) {
	if id, err := item.Id(); err == nil {
		m.Items[id] = item
	}
}

func (m *MergeEdge) AddTo(items map[int64]*GraphEdge) {
	fi := m.firstItem()
	if fi == nil {
		return
	}
	line := &GraphEdge{}
	line.Identity = fi.Identity
	line.End = fi.End
	line.Start = fi.Start
	line.Type = fi.Type
	line.SubEdges = m.Items

	items[line.Identity] = line
}

func (m *MergeEdge) firstItem() *GraphEdge {
	for _, i := range m.Items {
		return i
	}
	return nil
}
