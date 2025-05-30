package response

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/graph/vis_network"
	"strings"
)

var nodeType = map[string]string{
	"human":    "hm", // 人员
	"company":  "cp", // 公司
	"product":  "pd", // 产品
	"contract": "ct", // 合同
	"account":  "ac", // 账号
	"bank":     "bk", // 银行
	"record":   "rc", // 流水
}

func NewResultMap(graphView *graph.GraphView) *vis_network.ResultMap {
	isMerge := true
	return vis_network.NewResultMapWithGraphView(graphView, &vis_network.ResultMapWithGraphViewOption{
		Merge: &isMerge,
		OnNode: func(result *vis_network.ResultMap, n *vis_network.Node) {
			tags := n.Tags
			for key, _ := range tags {
				if strings.Contains(key, "_") {
					delete(tags, key)
				}
			}
			for key, group := range nodeType {
				if _, ok := tags[key]; ok {
					n.Group = &group
					return
				}
			}
		},
		OnEdge: func(result *vis_network.ResultMap, edge *vis_network.Edge) {

		},
	})
}

func NewResultList(graphView *graph.GraphView) *vis_network.ResultList {
	return vis_network.NewResultListWithGraphView(graphView, &vis_network.ResultListWithGraphViewOption{
		OnNode: func(result *vis_network.ResultList, n *vis_network.Node) {
			tags := n.Tags
			for key, _ := range tags {
				if strings.Contains(key, "_") {
					delete(tags, key)
				}
			}
			for key, group := range nodeType {
				if _, ok := tags[key]; ok {
					n.Group = &group
					return
				}
			}
		},
		OnEdge: func(result *vis_network.ResultList, edge *vis_network.Edge) {

		},
	})
}
