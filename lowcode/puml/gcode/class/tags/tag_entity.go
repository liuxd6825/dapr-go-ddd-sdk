package tags

import "github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"

type TagEntity struct {
	BaseTag
}

func NewTagEntity(text string) (*TagEntity, error) {
	tagData := &TagEntity{}
	tagData.TagType = TagTypeEntity
	if err := utils.TagUnmarshal(text, tagData); err != nil {
		return nil, err
	}
	return tagData, nil
}
