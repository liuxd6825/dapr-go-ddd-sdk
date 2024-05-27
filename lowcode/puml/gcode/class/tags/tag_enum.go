package tags

import (
	"fmt"
	"github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"
	"strings"
)

type TagEnum struct {
	BaseTag
	Items []*TagItem `json:"items"`
}

type TagItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Index int    `json:"index"`
}

func NewTagEnum(text string) (*TagEnum, error) {
	tagEnum := &TagEnum{}
	tagEnum.TagType = TagTypeEnum
	list := strings.Split(text, ",")
	for i, str := range list {
		item := &TagItem{Index: i}
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
		tagEnum.Items = append(tagEnum.Items, item)
	}
	return tagEnum, nil
}
