package tags

import (
	"github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"
)

type TagForm struct {
	BaseTag
	Required bool     `json:"required"` //是否必填
	Widget   string   `json:"widget"`   // 控件
	Help     string   `json:"help"`     // 帮助说明
	Order    []string `json:"order"`    // 控件顺序
}

func NewTagForm(text string) (*TagForm, error) {
	tagForm := &TagForm{}
	tagForm.TagType = TagTypeForm
	if err := utils.TagUnmarshal(text, tagForm); err != nil {
		return nil, err
	}
	return tagForm, nil
}
