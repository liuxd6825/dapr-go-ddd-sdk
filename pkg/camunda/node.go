package camunda

import "github.com/beevik/etree"

type Node struct {
	*Element
	Outgoings []string `json:"outgoings"`
	Incomings []string `json:"incomings"`
	OutNodes  []*Node  `json:"outNodes"`
}

func NewTask(el *etree.Element) *Node {
	n := &Node{
		Element: NewElement(el),
	}
	n.Outgoings = n.getOutgoing()
	n.Incomings = n.getIncoming()
	return n
}

func (n *Node) getOutgoing() []string {
	els := n.el.FindElements("bpmn:outgoing")
	res := make([]string, len(els))
	for i, el := range els {
		res[i] = el.Text()
	}
	return res
}

func (n *Node) getIncoming() []string {
	els := n.el.FindElements("bpmn:outgoing")
	res := make([]string, len(els))
	for i, el := range els {
		res[i] = el.Text()
	}
	return res
}

func (n *Node) IsExclusiveGateway() bool {
	if n.NodeType == "exclusiveGateway" {
		return true
	}
	return false
}

func (n *Node) IsParallelGateway() bool {
	if n.NodeType == "parallelGateway" {
		return true
	}
	return false
}
