package tags

import "github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"

type TagQuery struct {
	BaseTag
}

func NewTagQuery(text string) (*TagView, error) {
	tagData := &TagView{}
	tagData.TagType = TagTypeQuery
	if err := utils.TagUnmarshal(text, tagData); err != nil {
		return nil, err
	}
	return tagData, nil
}
