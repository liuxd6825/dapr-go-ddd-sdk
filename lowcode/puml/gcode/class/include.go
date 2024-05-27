package class

import "context"

type Include struct {
	NameElement
	FileName string `json:"fileName"`
}

func NewInclude(line string, namespaceName string, notes []*Comment) *Include {
	include := &Include{}
	include.InitBase(include, line, "include", namespaceName, notes)
	return include
}

func (b *Include) Parse(ctx context.Context, reader ParseReader) (err error) {
	return nil
}
