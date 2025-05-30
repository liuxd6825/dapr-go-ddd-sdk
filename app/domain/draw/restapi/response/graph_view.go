package response

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/graph/vis_network"
	"strings"
)

var nodeType = map[string]string{
	"human":    "hm",
	"company":  "cp",
	"product":  "pd",
	"contract": "ct",
	"account":  "ac",
	"bank":     "bk",
	"record":   "rc",
}

func NewResultWithGraphView(graphView *graph.GraphView) *vis_network.Result {
	isMerge := true
	return vis_network.NewResultWithGraphView(graphView, &vis_network.ResultWithGraphViewOption{
		Merge: &isMerge,
		OnNode: func(result *vis_network.Result, n *vis_network.Node) {
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
		OnEdge: func(result *vis_network.Result, edge *vis_network.Edge) {

		},
	})
}
