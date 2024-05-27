package tags

import "github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"

type TagDomain struct {
	BaseTag
}

func NewTagDomain(text string) (*TagDomain, error) {
	tagData := &TagDomain{}
	tagData.TagType = TagTypeDomain
	if err := utils.TagUnmarshal(text, tagData); err != nil {
		return nil, err
	}
	return tagData, nil
}
