package tags

import "github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"

type TagApi struct {
	BaseTag
}

func NewTagApi(text string) (*TagApi, error) {
	tagData := &TagApi{}
	tagData.TagType = TagTypeApi
	if err := utils.TagUnmarshal(text, tagData); err != nil {
		return nil, err
	}
	return tagData, nil
}
