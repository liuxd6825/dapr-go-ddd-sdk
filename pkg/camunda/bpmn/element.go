package bpmn

import (
	"github.com/beevik/etree"
	"strings"
)

type Element struct {
	el       *etree.Element
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Text     string   `json:"text"`
	NodeType NodeType `json:"nodeType"`
}

func NewElement(el *etree.Element) *Element {
	e := &Element{el: el}
	e.Id = e.getId()
	e.Name = e.getName()
	e.Text = e.getText()
	e.NodeType = e.getType()
	return e
}

func (n *Element) getId() string {
	return n.el.SelectAttrValue("id", "")
}

func (n *Element) getName() string {
	return n.el.SelectAttrValue("name", "")
}

func (n *Element) getType() NodeType {
	return NodeType(n.el.Tag)
}

func (n *Element) getText() string {
	text := n.el.Text()
	text = strings.ReplaceAll(text, "\n", "")
	return strings.TrimSpace(text)
}
