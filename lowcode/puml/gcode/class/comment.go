package class

import (
	"context"
	"strings"
)

type Comment struct {
	Text string `json:"text"`
}

func NewComment(ctx context.Context, line string) *Comment {
	text := strings.Trim(line[1:], " ")
	return &Comment{
		Text: text,
	}
}
func (n *Comment) GetType() string {
	return "note"
}
func (n *Comment) ParseStart() bool {
	return false
}
func (n *Comment) ParseEnd() bool {
	return false
}
func (n *Comment) StartTag() string {
	return "'"
}
func (n *Comment) EndTag() string {
	return ""
}
func (n *Comment) AddLine(line string) {

}

func (n *Comment) GetText() string {
	return n.Text
}
