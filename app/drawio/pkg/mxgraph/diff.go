package mxgraph

import "strings"

type Cell struct {
	Geometry *string `json:"geometry"`
	Id       *string `json:"id"`
	Parent   *string `json:"parent"`
	Previous *string `json:"previous"`
	Style    *string `json:"style"`
	Vertex   *int64  `json:"vertex"`
	XmlValue *string `json:"xmlValue"`
	Value    *string `json:"value"`
	Edge     *int    `json:"edge"`   //是否为关系
	Source   *string `json:"source"` // 关系来源
	Target   *string `json:"target"` // 关系目录
	style    map[string]string
	element  Element
}

func (c *Cell) GetElement() Element {
	return c.element
}

func (c *Cell) IsNode() bool {
	if c.Edge != nil {
		if c.Style != nil {
			style := c.GetStyle()
			if style != nil {
			}
		}
	}
	return false
}

func (c *Cell) IsEdge() bool {
	if c.Edge != nil {
		if c.Style != nil {
			style := c.GetStyle()
			if style != nil {
			}
		}
	}
	return true
}

func (c *Cell) GetStyle() map[string]string {
	if c.style != nil {
		return c.style
	}
	if c.Style == nil {
		return nil
	}
	result := make(map[string]string)

	// 使用分号分隔每个键值对
	pairs := strings.Split(*c.Style, ";")
	for _, pair := range pairs {
		// 跳过空字符串
		if pair == "" {
			continue
		}

		// 使用等号分隔键和值
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			// 存储键值对
			result[kv[0]] = kv[1]
		} else {
			// 如果没有值，则将值设置为空字符串
			result[kv[0]] = ""
		}
	}

	c.style = result
	return result
}

type Cells struct {
	I []*Cell          `json:"i"`
	R []string         `json:"r"`
	U map[string]*Cell `json:"u"`
}

type UpdatePage struct {
	Name          *string `json:"name"`
	PageProperty2 *string `json:"page_property2"`
	Cells         *Cells  `json:"cells"`
}

type InsertPage struct {
	Id       *string `json:"id"`
	Data     *string `json:"data"`
	Previous *string `json:"previous"`
}

type FileDiff struct {
	I []*InsertPage          `json:"i"`
	R []string               `json:"r"`
	U map[string]*UpdatePage `json:"u"`
}
