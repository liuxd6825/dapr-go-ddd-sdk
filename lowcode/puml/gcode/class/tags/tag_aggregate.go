package tags

import (
	"github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"
)

type TagAggregate struct {
	BaseTag
}

func NewTagAggregate(text string) (*TagAggregate, error) {
	tagData := &TagAggregate{}
	tagData.TagType = TagTypeAggregate
	if err := utils.TagUnmarshal(text, tagData); err != nil {
		return nil, err
	}
	return tagData, nil
}
