package tags

import (
	"fmt"
	"github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"
	"strings"
)

type TagDict struct {
	BaseTag
	Items []*DictItem `json:"items"`
}

type DictItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Index int    `json:"index"`
}

func NewTagDict(text string) (*TagDict, error) {
	tagDict := &TagDict{}
	tagDict.TagType = TagTypeEnum
	list := strings.Split(text, ",")
	for i, str := range list {
		item := &DictItem{Index: i}
		s := strings.Split(str, ":")
		count := len(s)
		if count == 1 {
			item.Name = utils.StringTrim(s[0])
			item.Title = fmt.Sprintf("%d", i)
		} else if count > 1 {
			item.Name = utils.StringTrim(s[0])
			item.Title = utils.StringTrim(s[1], `"`, `'`)
		} else {
			continue
		}
		tagDict.Items = append(tagDict.Items, item)
	}
	return tagDict, nil
}
