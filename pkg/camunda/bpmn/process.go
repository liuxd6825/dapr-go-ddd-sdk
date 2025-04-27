package bpmn

import (
	"fmt"
	"github.com/beevik/etree"
)

type NodeType string

const (
	NodeType_StartEvent       NodeType = "startEvent"
	NodeType_EndEvent         NodeType = "endTask"
	NodeType_UserTask         NodeType = "userTask"
	NodeType_ServiceTask      NodeType = "serviceTask"
	NodeType_ExclusiveGateway NodeType = "exclusiveGateway"
	NodeType_ParallelGateway  NodeType = "parallelGateway"
)

type Process struct {
	doc     *etree.Document
	process *etree.Element
}

func NewProcessWidthXml(bpmnXml string) *Process {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(bpmnXml); err != nil {
		panic(err)
	}
	return NewProcess(doc)
}

func NewProcess(doc *etree.Document) *Process {
	process := doc.FindElement("//process")
	return &Process{
		doc:     doc,
		process: process,
	}
}

func (b *Process) GetNode(id string) *Node {
	nodeId := fmt.Sprintf("//*[@id='%s']", id)
	el := b.process.FindElement(nodeId)
	return NewTask(el)
}

func (b *Process) GetNextNodes(id string) []*Node {
	seqList := b.GetSequenceByAttr(id)
	result := make([]*Node, 0)
	for _, seq := range seqList {
		targetId := seq.GetTargetRef()
		if targetId != "" {
			node := b.GetNode(targetId)
			if node != nil {
				switch node.NodeType {
				case NodeType_ExclusiveGateway, NodeType_ParallelGateway:
					node.OutNodes = b.GetNextNodes(node.Id)
				}
				result = append(result, node)
			}
		}
	}
	return result
}

func (b *Process) GetSequenceByAttr(sourceRef string) []*Sequence {
	find := fmt.Sprintf("bpmn:sequenceFlow[@sourceRef='%s']", sourceRef)
	list := b.process.FindElements(find)
	res := make([]*Sequence, len(list))
	for i, item := range list {
		res[i] = NewSequence(item)
	}
	return res
}
