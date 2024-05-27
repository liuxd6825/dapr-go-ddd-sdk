package tags

import "github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"

type TagEvent struct {
	BaseTag
	Name   string `json:"name"`
	Fields Fields `json:"fields"`
}

func NewTagEvent(text string) (*TagEvent, error) {
	tagData := &TagEvent{}
	tagData.TagType = TagTypeEvent
	if err := utils.TagUnmarshal(text, tagData); err != nil {
		return nil, err
	}
	return tagData, nil
}
