package tags

import "github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"

type TagUiSchema struct {
	BaseTag
	Order      []string           `json:"order"` // 控件顺序
	Options    TagUiSchemaOptions `json:"options"`
	ClassNames string             `json:"classNames"`
}

type TagUiSchemaOptions struct {
	Expandable bool
}

func NewUiSchema(text string) (*TagUiSchema, error) {
	tagForm := &TagUiSchema{}
	tagForm.TagType = TagTypeForm
	if err := utils.TagUnmarshal(text, tagForm); err != nil {
		return nil, err
	}
	return tagForm, nil
}
