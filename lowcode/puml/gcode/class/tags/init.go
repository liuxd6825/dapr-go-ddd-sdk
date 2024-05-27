package tags

import (
	"github.com/jfeliu007/goplantuml/pkg/tocode/classdiagram/utils"
	"strings"
)

type NewTagFunc func(note string) (Tag, error)
type Note interface {
	GetText() string
}

var parseTags map[TagType]NewTagFunc

func init() {
	parseTags = make(map[TagType]NewTagFunc)
	parseTags[TagTypeData] = func(note string) (Tag, error) {
		return NewTagData(note)
	}
	parseTags[TagTypeGrid] = func(note string) (Tag, error) {
		return NewTagGrid(note)
	}
	parseTags[TagTypeEnum] = func(note string) (Tag, error) {
		return NewTagEnum(note)
	}
	parseTags[TagTypeForm] = func(note string) (Tag, error) {
		return NewTagForm(note)
	}
	parseTags[TagTypeAggregate] = func(note string) (Tag, error) {
		return NewTagAggregate(note)
	}
	parseTags[TagTypeRepository] = func(note string) (Tag, error) {
		return NewTagRepository(note)
	}
	parseTags[TagTypeEntity] = func(note string) (Tag, error) {
		return NewTagEntity(note)
	}
	parseTags[TagTypeEvent] = func(note string) (Tag, error) {
		return NewTagEvent(note)
	}
	parseTags[TagTypeService] = func(note string) (Tag, error) {
		return NewTagService(note)
	}
	parseTags[TagTypeView] = func(note string) (Tag, error) {
		return NewTagView(note)
	}
	parseTags[TagTypeQuery] = func(note string) (Tag, error) {
		return NewTagQuery(note)
	}
	parseTags[TagTypeCmd] = func(note string) (Tag, error) {
		return NewTagCmd(note)
	}
	parseTags[TagTypeApi] = func(note string) (Tag, error) {
		return NewTagApi(note)
	}
	parseTags[TagTypeClass] = func(note string) (Tag, error) {
		return NewTagClass(note)
	}
	parseTags[TagTypeDict] = func(note string) (Tag, error) {
		return NewTagDict(note)
	}
	parseTags[TagTypeDomain] = func(note string) (Tag, error) {
		return NewTagDomain(note)
	}
}

func ParseComments(element Element, comments []string) ([]Tag, error) {
	var list []Tag
	for _, text := range comments {
		tags, ok, err := ParseTags(element, text)
		if err != nil {
			return nil, err
		}
		if ok {
			list = append(list, tags...)
		}
	}
	return list, nil
}

func ParseTags(element Element, text string) ([]Tag, bool, error) {
	if !strings.Contains(text, "@") {
		return nil, false, nil
	}
	var list []Tag
	for tagType, fun := range parseTags {
		tag, ok, err := newTag(text, tagType, fun)
		if err != nil {
			return nil, false, err
		} else if ok {
			tag.SetElement(element)
			list = append(list, tag)
		}
	}
	return list, len(list) > 0, nil
}

func GetTagText(text string, tagType TagType) (string, bool) {
	isHas, res, _, _, _ := utils.StringBetween(text, string(tagType), ")")
	if isHas {
		return res[1:], true
	}
	return "", false
}

func newTag(text string, tagType TagType, newTagFun NewTagFunc) (tag Tag, ok bool, err error) {
	dataText, dataHas := GetTagText(text, tagType)
	if dataHas {
		tag, err = newTagFun(dataText)
	}
	if tag != nil {
		ok = true
	}
	return
}
