package hugegraph

import (
	"fmt"
	"reflect"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/gremlin"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
)

type Result struct {
	Nodes   map[string]*Node
	Edges   map[string]*Edge
	Paths   []*Path
	Records []Record
}

// Path 路径结构
type Path struct {
	Labels  [][]string `json:"labels"`
	Objects []any      `json:"objects"`
}

// Record 普通聚合或动态返回记录
type Record any

// Node 顶点结构
type Node struct {
	ID         string         `json:"id"`
	Label      string         `json:"label"`
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
}

// Edge 边结构
type Edge struct {
	ID         string         `json:"id"`
	Label      string         `json:"label"`
	Type       string         `json:"type"`
	OutV       string         `json:"outV"`
	InV        string         `json:"inV"`
	Properties map[string]any `json:"properties"`
}

type ResultOptions struct {
	Data any
}

func NewRecord(data any) Record {
	return data
}

func NewResult(data any) *Result {
	res := &Result{
		Paths:   make([]*Path, 0),
		Records: make([]Record, 0),
		Nodes:   make(map[string]*Node),
		Edges:   make(map[string]*Edge),
	}
	res.init(data)
	return res
}

func NewResultWithResponse(response *gremlin.PostResponseData) *Result {
	if response == nil {
		panic("NewResultWithResponse(response) response is nil")
	}
	return NewResult(response.Result.Data)
}

func (r *Result) init(data any) {
	if data == nil {
		return
	}
	switch data.(type) {
	case []any:
		list := data.([]any)
		for _, v := range list {
			r.addAny(v)
		}
		break
	case map[string]any:
		r.addAny(data)
		break
	case string:
		r.addRecord(data)
		break
	default:
		r.addAny(data)
	}
}

func (r *Result) addAny(data any) {
	if data == nil {
		return
	}
	typeVal := reflect.TypeOf(data)
	fmt.Println(typeVal.Name())

	switch data.(type) {
	case map[string]interface{}:
		vData := data.(map[string]interface{})
		dataType := r.getDataType(vData)
		switch dataType {
		case VERTEX:
			node := r.newNode(vData)
			r.Nodes[node.ID] = node
			break
		case EDGE:
			edge := r.newEdge(vData)
			r.Edges[edge.ID] = edge
			break
		case PATH:
			path := r.newPath(vData)
			r.Paths = append(r.Paths, path)
			break
		}
	default:
		r.addRecord(data)
		break
	}
}

func (r *Result) addRecord(data any) {
	if data == nil {
		return
	}
	r.Records = append(r.Records, Record(data))
}

func (r *Result) newEdge(data map[string]any) *Edge {
	edge := &Edge{}
	edge.ID = maputils.GetStringErr(data, "id", "")
	edge.Label = maputils.GetStringErr(data, "label", "")
	edge.Type = maputils.GetStringErr(data, "type", "")
	edge.OutV = maputils.GetStringErr(data, "outV", "")
	edge.InV = maputils.GetStringErr(data, "inV", "")
	edge.Properties = maputils.GetMapErr(data, "properties", nil)
	return edge
}

func (r *Result) newNode(data map[string]any) *Node {
	node := &Node{}
	node.ID = maputils.GetStringErr(data, "id", "")
	node.Label = maputils.GetStringErr(data, "label", "")
	node.Type = maputils.GetStringErr(data, "type", "")
	node.Properties = maputils.GetMapErr(data, "properties", nil)
	return node
}

func (r *Result) newPath(data map[string]any) *Path {
	path := &Path{}
	path.Labels = r.getPathLabels(data)
	path.Objects = r.getPathObjects(data)
	return path
}

func (r *Result) isPathData(data map[string]any) bool {
	if data == nil {
		return false
	}
	_, isLabels := data["labels"]
	_, isObjects := data["objects"]
	if isLabels && isObjects {
		return true
	}
	return false
}

type GraphType string

const (
	VERTEX GraphType = "vertex"
	EDGE   GraphType = "edge"
	PATH   GraphType = "path"
	RECORD GraphType = "record"
	NONE   GraphType = "none"
)

func (g GraphType) Name() string {
	return string(g)
}

func (r *Result) getDataType(data map[string]any) GraphType {
	if data == nil {
		return NONE
	}
	typeVal, hasType := data["type"]
	if hasType {
		typeVal = typeVal.(string)
		switch typeVal {
		case VERTEX.Name():
			return VERTEX
		case EDGE.Name():
			return EDGE
		}
	} else if r.isPathData(data) {
		return PATH
	}
	return RECORD
}

func (r *Result) getPathLabels(data map[string]any) [][]string {
	val, isLabels := data["labels"]
	res := make([][]string, 0)
	if isLabels {
		listVal := val.([]any)
		labels := make([]string, 0)
		for _, v := range listVal {
			sList := v.([]any)
			for _, l := range sList {
				label := l.(string)
				labels = append(labels, label)
			}
			res = append(res, labels)
		}

	}
	return res
}

func (r *Result) getPathObjects(data map[string]any) []any {
	res := []any{}
	objects, isHas := data["objects"]
	if !isHas {
		return []any{}
	}
	list := objects.([]any)
	for _, v := range list {
		switch v.(type) {
		case map[string]interface{}:
			vData := v.(map[string]interface{})
			switch r.getDataType(vData) {
			case VERTEX:
				node := r.newNode(vData)
				res = append(res, node)
				break
			case EDGE:
				edge := r.newEdge(vData)
				res = append(res, edge)
				break
			case PATH:
				if r.isPathData(vData) {
					path := r.newPath(vData)
					res = append(res, path)
				}
			}
		}
	}
	return res
}

func (r *Result) GetGraphView() *graph.GraphView {
	graphView := graph.NewGraphView()
	var nodes []*graph.Node
	for _, i := range r.Nodes {
		node := &graph.Node{
			Nid:   i.ID,
			Id:    i.ID,
			Label: i.Label,
			Props: i.Properties,
		}
		nodes = append(nodes, node)
	}

	var edges []*graph.Edge
	for _, i := range r.Edges {
		edge := &graph.Edge{
			Nid:   i.ID,
			NFrom: i.InV,
			NTo:   i.OutV,
			Label: i.Label,
			Props: i.Properties,
		}
		edges = append(edges, edge)
	}

	for _, p := range r.Paths {
		for _, v := range p.Objects {
			switch v.(type) {
			case *Node:
				i := v.(*Node)
				node := &graph.Node{
					Nid:   i.ID,
					Id:    i.ID,
					Label: i.Label,
					Props: i.Properties,
				}
				nodes = append(nodes, node)
				break
			case *Edge:
				i := v.(*Edge)
				edge := &graph.Edge{
					Nid:   i.ID,
					NFrom: i.InV,
					NTo:   i.OutV,
					Label: i.Label,
					Props: i.Properties,
				}
				edges = append(edges, edge)
				break
			}
		}
	}

	graphView.AddNodes("nodes", nodes)
	graphView.AddEdges("edges", edges)
	return graphView
}
