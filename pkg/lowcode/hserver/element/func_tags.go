package element

import (
	"golang.org/x/net/html"
)

type FuncTagType = string

const (
	FuncTag_DomainEvent FuncTagType = "event"
	FuncTag_RestApi     FuncTagType = "rest-api"
)

type FuncTag struct {
	Type  string
	Attrs map[string]string
}

func NewFuncTag(node *html.Node) *FuncTag {
	tag := &FuncTag{
		Type:  node.Data,
		Attrs: make(map[string]string),
	}
	for _, a := range node.Attr {
		tag.Attrs[a.Key] = a.Val
	}
	return tag
}

func (tag *FuncTag) GetType() string {
	return tag.Type
}

func (tag *FuncTag) HAttr(key string) (string, bool) {
	val, ok := tag.Attrs[key]

	return val, ok
}

func (tag *FuncTag) AttrOr(key string, def string) string {
	val, ok := tag.Attrs[key]
	if !ok {
		return def
	}
	return val
}

func (tag *FuncTag) Attr(key string) string {
	val, ok := tag.Attrs[key]
	if !ok {
		return ""
	}
	return val
}
