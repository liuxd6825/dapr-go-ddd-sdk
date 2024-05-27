package tags

import "github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"

type TagView struct {
	BaseTag
}

func NewTagView(text string) (*TagView, error) {
	tagData := &TagView{}
	tagData.TagType = TagTypeView
	if err := utils.TagUnmarshal(text, tagData); err != nil {
		return nil, err
	}
	return tagData, nil
}
