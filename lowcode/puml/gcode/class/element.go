package class

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/tags"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"strings"
)

type Element interface {
	GetLine() string
	GetTypeName() string
	GetComments() []*Comment
	GetNamespaceName() string
	GetName() string
	GetPackage() string
}

type NodeElement interface {
	GetTags() []tags.Tag
	FindTags(typeName tags.TagType) []tags.Tag
}

type BaseElement struct {
	Line          string     `json:"-"`                       // 行文本
	TypeName      string     `json:"typeName,omitempty"`      // 元素类型
	Comments      []*Comment `json:"comments,omitempty"`      // 注释
	Tags          []tags.Tag `json:"tags,omitempty"`          // 标签
	NamespaceName string     `json:"namespaceName,omitempty"` // 命名空间
	Package       string     `json:"package,omitempty"`       // 包名
	DataType      string     `json:"dataType,omitempty"`      // 数据类型

}

type NameElement struct {
	BaseElement
	Name  string `json:"name,omitempty"`  // 名称
	Alias string `json:"alias,omitempty"` // 别名
}

func (n *NameElement) GetDataType() string {
	return n.DataType
}

func (n *BaseElement) InitBase(parentElement tags.Element, line string, typeName string, namespaceName string, comments []*Comment) (err error) {
	list := strings.Split(namespaceName, ".")
	n.Package = list[len(list)-1]
	n.Line = line
	n.TypeName = typeName
	n.NamespaceName = namespaceName
	n.Tags, err = tags.ParseComments(parentElement, commentsToStrList(comments))
	n.Comments = initComments(comments)
	return
}

func (n *BaseElement) InitDataTag(elementName string) {
	title := ""
	if len(n.Comments) > 0 {
		title = n.Comments[0].Text
	} else {
		title = elementName
	}

	dataTags := n.FindTags(tags.TagTypeData)
	if len(n.Comments) == 1 {
		for _, dataTag := range dataTags {
			data := dataTag.(*tags.TagData)
			if data != nil && data.Titles.CN == "" {
				data.Titles.CN = title
			}
		}
	}
	for _, dataTag := range dataTags {
		data := dataTag.(*tags.TagData)
		if data != nil {
			if data.Titles.CN == "" {
				data.Titles.CN = title
			}
			if data.Titles.EN == "" {
				data.Titles.EN = elementName
			}
			if data.Titles.DE == "" {
				data.Titles.DE = elementName
			}
		}
	}
}

func (n *BaseElement) FindTags(tagType tags.TagType) (tags []tags.Tag) {
	for _, tag := range n.Tags {
		if tagType == tag.GetTagType() {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (n *BaseElement) GetPackage() string {
	return n.Package
}

func (n *BaseElement) GetLine() string {
	return n.Line
}

func (n *BaseElement) SetLine(val string) {
	n.Line = val
}

func (n *BaseElement) GetTypeName() string {
	return n.TypeName
}

func (n *BaseElement) SetTypeName(val string) {
	n.TypeName = val
}

func (n *BaseElement) GetNamespaceName() string {
	return n.NamespaceName
}

func (n *BaseElement) SetNamespaceName(val string) {
	n.NamespaceName = val
}

func (n *BaseElement) GetComments() []*Comment {
	return n.Comments
}

func (n *BaseElement) SetComments(val []*Comment) {
	n.Comments = val
}

func (n *BaseElement) GetTags() []tags.Tag {
	return n.Tags
}

func (n *BaseElement) SetTags(val []tags.Tag) {
	n.Tags = val
}

func (n *BaseElement) CommentsText() string {
	return GetNotesText(n.Comments)
}

func (n *BaseElement) AddComment(commentText string) {
	n.Comments = append(n.Comments, &Comment{Text: commentText})
}

func (n *BaseElement) Register(ctx context.Context, node Element) {
	pctx := GetParseContext(ctx)
	if pctx != nil {
		pctx.AddNode(node)
	}
}

func (n *BaseElement) DataTypeIsArray() bool {
	if strings.HasPrefix(n.DataType, "[]") {
		return true
	}
	return false
}

func (n *NameElement) GetName() string {
	return n.Name
}

func (n *NameElement) FUpperName() string {
	return utils.FirstUpperName(n.Name)
}

func (n *NameElement) FLowerName() string {
	return utils.FirstLowerName(n.Name)
}

func (n *NameElement) SetName(val string) {
	n.Name = val
}

func (n *NameElement) GetAlias() string {
	return n.Alias
}

func (n *NameElement) SetAlias(val string) {
	n.Alias = val
}

// SN
// @Description: 蛇形名称
func (n *NameElement) SN() string {
	if len(n.Name) > 0 {
		return strings.ToLower(n.Name[0:1])
	}
	return "n"
}

func (n *NameElement) InitName(name string, alias string) {
	n.Name = name
	n.Alias = alias
}

func (n *NameElement) Init(line string, typeName string, namespace string, notes []*Comment, elementName string, elementAlias string) {
	n.BaseElement.InitBase(n, line, typeName, namespace, notes)
	n.BaseElement.InitDataTag(elementName)
	n.InitName(elementName, elementAlias)
}

func commentsToStrList(notes []*Comment) []string {
	var list []string
	for _, note := range notes {
		list = append(list, note.Text)
	}
	return list
}

func initComments(comments []*Comment) []*Comment {
	var newComments []*Comment
	for _, note := range comments {
		str := utils.RemoveTags(note.Text)
		if str != "" {
			note.Text = str
			newComments = append(newComments, note)
		}
	}
	return newComments
}
