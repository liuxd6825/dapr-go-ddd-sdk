package tags

import "github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"

type TagClass struct {
	BaseTag
	Name   string `json:"name"`
	Fields Fields `json:"fields"`
}

func NewTagClass(text string) (*TagClass, error) {
	tagData := &TagClass{}
	tagData.TagType = TagTypeClass
	if err := utils.TagUnmarshal(text, tagData); err != nil {
		return nil, err
	}
	return tagData, nil
}
