package tags

type Tag interface {
	GetTagType() TagType
	GetElement() Element
	SetElement(Element)
}

type TagType string

type BaseTag struct {
	TagType     TagType `json:"tagType"`
	Element     Element `json:"-"`
	ElementName string  `json:"elementName"`
}

const (
	TagTypeDomain     TagType = "@domain"
	TagTypeData       TagType = "@data"
	TagTypeAggregate  TagType = "@aggregate"
	TagTypeEnum       TagType = "@enum"
	TagTypeForm       TagType = "@form"
	TagTypeGrid       TagType = "@grid"
	TagTypeRepository TagType = "@repository"
	TagTypeService    TagType = "@service"
	TagTypeEntity     TagType = "@entity"
	TagTypeEvent      TagType = "@event"
	TagTypeApi        TagType = "@api"
	TagTypeView       TagType = "@view"
	TagTypeQuery      TagType = "@query"
	TagTypeCmd        TagType = "@cmd"
	TagTypeClass      TagType = "@class"
	TagTypeFramework  TagType = "@framework"
	TagTypeDict       TagType = "@dict"
)

func (t *BaseTag) GetTagType() TagType {
	return t.TagType
}

func (t *BaseTag) GetElement() Element {
	return t.Element
}

func (t *BaseTag) SetElement(el Element) {
	t.Element = el
	if el != nil {
		t.ElementName = el.GetName()
	}
}
